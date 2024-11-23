package redis

import (
	"crypto/tls"
	"io"

	red "github.com/go-redis/redis/v8"
	"gitlab.com/a16624741591/zero/core/syncx"
)

var failoverManager = syncx.NewResourceManager()

func getFailover(r *Redis) (*red.Client, error) {
	val, err := failoverManager.GetResource(r.Addr, func() (io.Closer, error) {
		var tlsConfig *tls.Config
		if r.tls {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		store := red.NewFailoverClient(&red.FailoverOptions{
			MasterName:    r.MasterName,
			SentinelAddrs: r.SentinelAddrs,
			Password:      r.Pass,
			MaxRetries:    MaxRetries,
			MinIdleConns:  IdleConns,
			TLSConfig:     tlsConfig,
		})

		store.AddHook(durationHook)
		return store, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*red.Client), nil
}
