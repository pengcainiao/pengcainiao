package redis

import (
	"crypto/tls"
	"io"

	red "github.com/go-redis/redis/v8"
	"gitlab.com/a16624741591/zero/core/syncx"
)

const (
	defaultDatabase = 0
	MaxRetries      = 3
	IdleConns       = 8
)

var clientManager = syncx.NewResourceManager()

func getClient(r *Redis) (*red.Client, error) {
	val, err := clientManager.GetResource(r.Addr, func() (io.Closer, error) {
		var tlsConfig *tls.Config
		if r.tls {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		store := red.NewClient(&red.Options{
			Addr:         r.Addr,
			Password:     r.Pass,
			DB:           defaultDatabase,
			MaxRetries:   MaxRetries,
			MinIdleConns: IdleConns,
			TLSConfig:    tlsConfig,
		})
		store.AddHook(durationHook)
		return store, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*red.Client), nil
}
