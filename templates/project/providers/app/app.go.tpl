package app

import (
	"git.ipao.vip/rogeecn/atom/container"
	"git.ipao.vip/rogeecn/atom/utils/opt"
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
