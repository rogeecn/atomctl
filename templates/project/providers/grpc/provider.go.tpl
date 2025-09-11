package grpc

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

	return container.Container.Provide(func() (*Grpc, error) {
		return &Grpc{config: &config}, nil
	}, o.DiOptions()...)
}
