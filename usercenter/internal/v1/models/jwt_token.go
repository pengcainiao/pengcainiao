package models

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pengcainiao/pengcainiao/usercenter/internal/v1/constant"
	"github.com/pengcainiao/zero/core/logx"
	"github.com/pengcainiao/zero/core/timex"
	"github.com/pengcainiao/zero/rest/httprouter"
	"github.com/pengcainiao/zero/tools"
	"github.com/pengcainiao/zero/tools/syncer"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ExpiringTokens struct {
	sync.Map
}

type (
	expireUserToken struct {
		cancelFunc     context.CancelFunc
		userID         string
		plaform        string
		onTokenExpired func()
	}
)

var (
	hmacSampleSecret               = []byte("mk4YzTnR$W2pT83mK/Lhi30UTk9wRA1")
	tokenExp                       = &ExpiringTokens{Map: sync.Map{}}
	redisClient                    = syncer.Redis()
	TokenVersion2MiniClientVersion = 15
)

const (
	ExpireSoonTokenKey = "_o" //一分钟后过期的token
)

type PlatformDefine string

//TokenObject JWT 凭证内容
type TokenObject struct {
	*jwt.StandardClaims
	SetNewToken         bool                    `json:"-"` //是否创建新token
	MaxAge              int64                   `json:"-"` //过期时间
	Token               string                  `json:"-"` //凭证
	ClientVersionNumber int                     `json:"-"` //token版本号
	UserID              string                  //用户ID
	DeviceID            string                  //设备ID
	Platform            constant.PlatformDefine //凭证所属平台
	ClientVersion       string                  //客户端版本
	Phone               string
	NickName            string
	Avatar              string
}

const (
	PC PlatformDefine = "pc"

	//Wechat     PlatformDefine = "wechat"
	//Mobile     PlatformDefine = "mobile"
	//Web        PlatformDefine = "web"
	//PCWechat   PlatformDefine = "pc_wechat"
	//H5         PlatformDefine = "h5"
	//CorpWechat PlatformDefine = "corp_wechat"
)

func NewTokenObject() *TokenObject {
	return &TokenObject{}
}

//VerifyToken 验证token
func VerifyToken(tokenString string) *TokenObject {
	var (
		claims     TokenObject
		tokenValid bool
	)
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})
	if err != nil {
		return nil
	}
	if !claims.IsValid() {
		logx.NewTraceLogger(context.Background()).Debug().Interface("token", tokenString).Msg("token验证无效")
		return nil
	}
	var (
		redisKey       = fmt.Sprintf(RedisKeyForOnlineUsers, claims.UserID)
		platform       = string(claims.Platform)
		allowPlatforms = []string{string(PC)}
		//allowPlatforms = []string{string(Wechat), string(PC), string(Mobile), string(Web), string(PCWechat), string(H5), string(CorpWechat)}
	)
	if !tools.InArray(platform, allowPlatforms) {
		logx.NewTraceLogger(context.Background()).Debug().Str("platform", platform).Interface("token", token).Msg("不支持的平台类型")
		return nil
	}

	var fields = []string{platform, getPlatformExpireKey(platform)}
	//if platform == "wechat" {
	//	fields = append(fields, "pc_wechat", getPlatformExpireKey("pc_wechat"))
	//}

	storageToken, err := syncer.Redis().HMGet(context.Background(), redisKey, fields...).Result()
	if err != nil {
		logx.NewTraceLogger(context.Background()).Err(err).Str("rdsKey", redisKey).Msg("查询用户token出错")
	}
	for _, tk := range storageToken {
		if tk == nil {
			continue
		}
		if tk.(string) == tokenString {
			tokenValid = true
			break
		}
	}
	if !tokenValid {
		return nil
	}

	if token != nil && token.Valid && claims.IsValid() {
		claims.Token = tokenString
		return &claims
	}
	return nil
}

//IsValid 验证token携带的用户是否合法
func (token TokenObject) IsValid() bool {
	return len(token.UserID) > 10 && BloomCheckUserExists(token.UserID)
}

//BloomAddNewUser 往布隆过滤器中添加新的用户ID
func BloomAddNewUser(userID string) {
	syncer.BloomClient(syncer.HoldingUsers).Add(userID)
}

func BloomCheckUserExists(userID string) bool {
	return syncer.UserExists(userID).Code == 0
}

func getPlatformExpireKey(platform string) string {
	return fmt.Sprintf("%s%s", platform, ExpireSoonTokenKey)
}

func (token *TokenObject) parseVersion() {
	var (
		versionStr string
		verArr     []string
	)

	if token.ClientVersion == "" {
		//logx.NewTraceLogger(context.Background()).Err(errors.New("client_version not set")).
		//	Str("stack", string(debug.Stack())).Msg("WARN")
		return
	}
	var versionByte0 = token.ClientVersion[0]
	if versionByte0 > 48 && versionByte0 <= 57 { //
		verArr = strings.Split(token.ClientVersion, ".")
	} else if versionByte0 == 86 || versionByte0 == 118 { //v V
		verArr = strings.Split(token.ClientVersion[1:], ".")
	} else {
		return
	}
	versionStr = verArr[0] + verArr[1]
	for i := 0; i < len(versionStr); i++ {
		var cur = versionStr[i]
		if cur < 48 || cur > 57 {
			return
		}
	}
	v, _ := strconv.Atoi(versionStr)
	token.ClientVersionNumber = v
}

//New 创建新token
func (token *TokenObject) New(ctx *httprouter.Context, currentUser UserOnline, oldToken string) error {
	token.UserID = currentUser.UserID
	token.Platform = "pc"

	token.DeviceID = currentUser.DeviceID
	token.SetNewToken = true
	token.ClientVersion = currentUser.ClientVersion

	token.parseVersion()
	var tokenMaxAge time.Duration
	if token.ClientVersionNumber < TokenVersion2MiniClientVersion {
		tokenMaxAge = timex.RandomExpireSeconds(time.Hour * 24 * 10) //time.Hour*24*2
	} else {
		tokenMaxAge = timex.RandomExpireSeconds(time.Hour * 2) //time.Hour*24*2
	}

	token.MaxAge = int64(tokenMaxAge.Seconds())
	token.StandardClaims = &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(tokenMaxAge).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "api.aaa.net",
	}
	tokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, token)
	if s, err := tokenClaim.SignedString(hmacSampleSecret); err != nil {
		logx.NewTraceLogger(ctx).Err(err).Msg("tokenClaim SignedString err")
		return err
	} else {
		token.Token = s
	}

	return token.refreshToken(oldToken, "pc")
}

func (token *TokenObject) refreshToken(oldToken string, params constant.PlatformDefine) error {
	var (
		redisKey = fmt.Sprintf(RedisKeyForOnlineUsers, token.UserID)
		platform = "pc"
	)
	var tokenMap = map[string]interface{}{
		platform: token.Token,
	}
	if oldToken != "" {
		tokenMap[getPlatformExpireKey(platform)] = oldToken
		go tokenExp.expireToken(token.UserID, platform)
	}
	logx.NewTraceLogger(context.Background()).Error().Fields(tokenMap).Str("user_id", token.UserID).Msg("！！刷新token")
	if err := redisClient.HMSet(context.Background(), redisKey, tokenMap).Err(); err != nil {
		logx.NewTraceLogger(context.Background()).Err(err).Str("rdsKey", redisKey).Msg("无法写入用户token")
		return nil
	}
	if err := redisClient.Expire(context.Background(), redisKey, time.Hour*24*30).Err(); err != nil {
		logx.NewTraceLogger(context.Background()).Err(err).Str("rdsKey", redisKey).Msg("无法设置过期时间")
	}
	return nil
}

func (exp *ExpiringTokens) expireToken(userID, platform string) {
	if v, ok := exp.Load(userID); ok {
		if m, ok := v.(*expireUserToken); ok {
			if m.cancelFunc != nil {
				m.cancelFunc()
			}
		}
	}
	_, cancelFunc := context.WithCancel(context.Background())
	var expObj = &expireUserToken{
		cancelFunc: cancelFunc,
		userID:     userID,
		plaform:    platform,
		onTokenExpired: func() {
			var redisKey = fmt.Sprintf(RedisKeyForOnlineUsers, userID)
			redisClient.HDel(context.Background(), redisKey, getPlatformExpireKey(platform))
			exp.Delete(userID)
		},
	}
	go expObj.exp(context.Background())
	exp.Store(userID, expObj)

}

func (exp *expireUserToken) exp(ctx context.Context) {
	select {
	case <-ctx.Done():
		break
	case <-time.After(time.Minute):
		exp.onTokenExpired()
	}
}
