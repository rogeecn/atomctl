package grpc

import (
	"git.ipao.vip/rogeecn/atom/container"
	"git.ipao.vip/rogeecn/atom/utils/opt"
	"google.golang.org/grpc"
)

func Provide(opts ...opt.Option) error {
	o := opt.New(opts...)
	var config Config
	if err := o.UnmarshalConfig(&config); err != nil {
		return err
	}
	return container.Container.Provide(func() (*Grpc, error) {
		server := grpc.NewServer()

		grpc := &Grpc{
			Server: server,
			config: &config,
		}
		container.AddCloseAble(grpc.Server.GracefulStop)

		return grpc, nil
	}, o.DiOptions()...)
}
