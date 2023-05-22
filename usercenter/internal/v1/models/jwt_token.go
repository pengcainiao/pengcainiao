package models

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pengcainiao/zero/core/logx"
	"github.com/pengcainiao/zero/tools"
	"github.com/pengcainiao/zero/tools/syncer"
)

var (
	hmacSampleSecret = []byte("mk4YzTnR$W2pT83mK/Lhi30UTk9wRA1")
)

const (
	ExpireSoonTokenKey = "_o" //一分钟后过期的token
)

type PlatformDefine string

//TokenObject JWT 凭证内容
type TokenObject struct {
	*jwt.StandardClaims
	UserID        string         //用户ID
	Token         string         //凭证
	DeviceID      string         //设备ID
	Platform      PlatformDefine //凭证所属平台
	ClientVersion string         //客户端版本
	//DeviceID    int64  //设备ID
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
