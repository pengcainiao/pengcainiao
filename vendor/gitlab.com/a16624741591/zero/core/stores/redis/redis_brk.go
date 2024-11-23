package redis

import (
	"context"
	"encoding/json"

	red "github.com/go-redis/redis/v8"

	"gitlab.com/a16624741591/zero/core/env"
)

func Client() RedisNode {
	if env.RedisMasterName != "" {
		var sentinelAddrs []string
		_ = json.Unmarshal([]byte(env.RedisAddr), &sentinelAddrs)
		conn, _ := getRedis(&Redis{
			MasterName:    env.RedisMasterName,
			SentinelAddrs: sentinelAddrs,
			Pass:          env.RedisPwd,
			Type:          FailoverType,
		})
		return conn
	}

	conn, _ := getRedis(&Redis{
		Addr: env.RedisAddr,
		Pass: env.RedisPwd,
		Type: NodeType,
	})
	return conn
}

func RedisClient() *red.Client {
	if env.RedisMasterName != "" {
		var sentinelAddrs []string
		_ = json.Unmarshal([]byte(env.RedisAddr), &sentinelAddrs)
		conn, _ := getFailover(&Redis{
			MasterName:    env.RedisMasterName,
			SentinelAddrs: sentinelAddrs,
			Pass:          env.RedisPwd,
			Type:          FailoverType,
		})
		return conn
	}

	conn, _ := getClient(&Redis{
		Addr: env.RedisAddr,
		Pass: env.RedisPwd,
		Type: NodeType,
	})
	return conn
}

func x() {
	Client().SRandMember(context.Background(), "xxx").Val()
}
