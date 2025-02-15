package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/martian/log"
	"gitlab.com/a16624741591/zero/core/logx"
	"time"

	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
	"github.com/scylladb/go-set/strset"
	"gopkg.in/square/go-jose.v2"
)

type RedisConfig struct {
	Host         string        `json:"host" yaml:"host"`
	Port         int           `json:"port" yaml:"port"`
	Protocol     string        `json:"protocol" yaml:"protocol"`
	PingInterval time.Duration `json:"ping_interval" yaml:"ping_interval"`
	PoolSize     int           `json:"pool_size" yaml:"pool_size"`
	Password     string        `json:"password" yaml:"password"`
	Database     int           `json:"database" yaml:"database"`
}

func (cfg *RedisConfig) Options() *redis.Options {
	return &redis.Options{
		Network:            cfg.Protocol,
		Addr:               cfg.Addr(),
		Password:           cfg.Password,
		DB:                 cfg.Database,
		IdleCheckFrequency: cfg.PingInterval,
		PoolSize:           cfg.PoolSize,
	}
}

func (cfg *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}

type RefreshTokenConfig struct {
	Redis      RedisConfig
	TimeToLive time.Duration
}

type _tokenInfo struct {
	*TokenInfo
}

func (t *_tokenInfo) Reset() {
	//TODO implement me
	panic("implement me")
}

func (t *_tokenInfo) String() string {
	//TODO implement me
	panic("implement me")
}

func (t *_tokenInfo) ProtoMessage() {
	//TODO implement me
	panic("implement me")
}

func (t *_tokenInfo) MarshalBinary() ([]byte, error) {
	return proto.Marshal(t)
}

func (t *_tokenInfo) UnmarshalBinary(b []byte) error {
	return proto.Unmarshal(b, t)
}

type refreshTokenProvider struct {
	rc *redis.Client
}

func NewRefreshTokenProvider(conf *RefreshTokenConfig) (TokenProvider, error) {

	opts := conf.Redis.Options()
	opts.OnConnect = func(c *redis.Conn) error {
		log.Infof("OnConnect: %v", c)
		return nil
	}

	rc := redis.NewClient(opts)
	err := rc.Ping().Err()
	if err != nil {
		return nil, err
	}

	return &refreshTokenProvider{
		rc: rc,
	}, nil
}

func (p *refreshTokenProvider) GetPublicJWKSet(ctx context.Context, keyID string) jose.JSONWebKeySet {
	return jose.JSONWebKeySet{}
}

func (p *refreshTokenProvider) GrantToken(ctx context.Context, tokenInfo *TokenInfo, ttl time.Duration) (*GrantedToken, error) {
	logx.NewTraceLogger(ctx).Info().Msg(fmt.Sprintf("222grant token: %d %s %s", tokenInfo.Uid, tokenInfo.DeviceId, ttl))

	token := uuid.NewV4().String()
	v, _ := json.Marshal(tokenInfo)

	userTokensKey := p.userTokensKey(tokenInfo.Uid)

	_, err := p.rc.TxPipelined(func(pl redis.Pipeliner) error {
		pl.Set(token, string(v), ttl)
		pl.HSet(userTokensKey, token, tokenInfo.DeviceId)
		pl.Expire(userTokensKey, ttl)

		return nil
	})

	if err != nil {
		log.Errorf("GrantToken: %v", err)
		return nil, err
	}

	log.Infof("grant token %d %s %s:%s", tokenInfo.Uid, tokenInfo.DeviceId, ttl, token)

	return &GrantedToken{
		Type:      TokenType_REFRESH_TOKEN,
		Token:     token,
		TokenId:   token,
		ExpiresIn: int64(ttl.Seconds()),
	}, nil
}

func (p *refreshTokenProvider) ValidateToken(ctx context.Context, token string) (*TokenInfo, error) {
	v, err := p.rc.Get(token).Bytes()
	switch err {
	case nil:
		var t TokenInfo
		err = json.Unmarshal(v, &t)
		if err != nil {
			return nil, invalidTokenError(err)
		}
		log.Debugf("validate %s %v", token, &t)
		return &t, nil
	case redis.Nil:
		return nil, invalidTokenError(redis.Nil)
	default:
		return nil, err
	}
}

func (p *refreshTokenProvider) RevokeToken(ctx context.Context, token string, def time.Duration) error {
	info, _ := p.ValidateToken(ctx, token)
	_, err := p.rc.TxPipelined(func(pl redis.Pipeliner) error {
		pl.Del(token)
		if info != nil {
			userTokensKey := p.userTokensKey(info.Uid)
			pl.HDel(userTokensKey, token)
		}

		return nil
	})
	return err
}

func (p *refreshTokenProvider) RevokeAllTokens(ctx context.Context, uid uint64, deviceId string, excludes ...string) ([]string, error) {
	log.Infof("RevokeAllTokens: revoke with uid %d and deviceId %s", uid, deviceId)
	userTokensKey := p.userTokensKey(uid)
	tokenToDeviceMap, err := p.rc.HGetAll(userTokensKey).Result()
	if err != nil {
		log.Errorf("RevokeAllTokens: revoke uid %d and deviceId %s", uid, deviceId)
		return nil, err
	}
	var (
		tokensToExpire = make([]string, 0, len(tokenToDeviceMap))
		tokensToDelete = make([]string, 0, len(tokenToDeviceMap))
		excludeSet     = strset.New(excludes...)
	)
	for token, device := range tokenToDeviceMap {
		if excludeSet.Has(token) {
			continue
		}
		if device == deviceId {
			tokensToExpire = append(tokensToExpire, token)
			continue
		}
		tokensToDelete = append(tokensToDelete, token)
	}
	if len(tokensToExpire) != 0 {
		boolCmdMap := make(map[string]*redis.BoolCmd, len(tokensToExpire))
		_, err := p.rc.TxPipelined(func(pl redis.Pipeliner) error {
			for _, token := range tokensToExpire {
				// expire 空 key没有问题 , 空的说明过期 从hmap中删除 // todo
				boolCmdMap[token] = pl.Expire(token, time.Hour*24)
			}
			return nil
		})
		if err != nil && err != redis.Nil {
			log.Errorf("RevokeAllTokens: expire with uid %d and deviceId %s : %s", uid, deviceId, err)
			return nil, err
		}
		for token, cmd := range boolCmdMap {
			ok, err := cmd.Result()
			if err != nil {
				log.Errorf("RevokeAllTokens: check expire with uid %d and deviceId %s : %s", uid, deviceId, err)
				tokensToDelete = append(tokensToDelete, token)
				continue
			}
			if !ok {
				tokensToDelete = append(tokensToDelete, token)
			}
		}

	}

	if len(tokensToDelete) != 0 {
		_, err = p.rc.TxPipelined(func(pl redis.Pipeliner) error {
			pl.Del(tokensToDelete...)
			pl.HDel(userTokensKey, tokensToDelete...)
			return nil
		})
		if err != nil {
			log.Errorf("RevokeAllTokens: delete token by tx with uid %d and deviceId %s : %s", uid, deviceId, err)
			return nil, err
		}
	}
	return tokensToDelete, nil
}

func (p *refreshTokenProvider) userTokensKey(uid uint64) string {
	return fmt.Sprintf("utk-%d", uid)
}
