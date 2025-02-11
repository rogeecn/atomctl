package app

import (
	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/opt"
)

func Provide(opts ...opt.Option) error {
	o := opt.New(opts...)
	var config Config
	if err := o.UnmarshalConfig(&config); err != nil {
		return err
	}

	return container.Container.Provide(func() (*Config, error) {
		return &config, nil
	}, o.DiOptions()...)
}
