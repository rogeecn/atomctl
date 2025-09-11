package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/opt"
)

func Provide(opts ...opt.Option) error {
	o := opt.New(opts...)
	var config Config
	if err := o.UnmarshalConfig(&config); err != nil {
		return err
	}
	config.format()
	return container.Container.Provide(func() (redis.UniversalClient, error) {
		// Build options via config helper (encapsulates defaults/decisions)
		ro := config.Options()

		// Safe structured log (no password)
		log.WithFields(log.Fields{
			"addr":        ro.Addr,
			"db":          ro.DB,
			"user":        ro.Username,
			"client_name": ro.ClientName,
			"pool_size":   ro.PoolSize,
			"min_idle":    ro.MinIdleConns,
			"retries":     ro.MaxRetries,
		}).Info("opening Redis connection")

		rdb := redis.NewClient(ro)

		// ping to verify connectivity
		ctx, cancel := context.WithTimeout(context.Background(), config.PingTimeout())
		defer cancel()
		if _, err := rdb.Ping(ctx).Result(); err != nil {
			return nil, err
		}

		// close hook
		container.AddCloseAble(func() { _ = rdb.Close() })

		return rdb, nil
	}, o.DiOptions()...)
}
