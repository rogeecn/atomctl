package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"{{.ModuleName}}/pkg/atom/container"
	"{{.ModuleName}}/pkg/atom/opt"
)

func Provide(opts ...opt.Option) error {
	o := opt.New(opts...)
	var config Config
	if err := o.UnmarshalConfig(&config); err != nil {
		return err
	}
	config.format()
	return container.Container.Provide(func() (redis.UniversalClient, error) {
		rdb := redis.NewClient(&redis.Options{
			Addr:     config.Addr(),
			Password: config.Password,
			DB:       int(config.DB),
		})

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		if _, err := rdb.Ping(ctx).Result(); err != nil {
			return nil, err
		}

		return rdb, nil
	}, o.DiOptions()...)
}
