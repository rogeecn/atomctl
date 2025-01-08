package cmux

import (
	"net"

	"{{.ModuleName}}/providers/grpc"
	"{{.ModuleName}}/providers/http"

	"git.ipao.vip/rogeecn/atom/container"
	"git.ipao.vip/rogeecn/atom/utils/opt"
	"github.com/soheilhy/cmux"
)

func Provide(opts ...opt.Option) error {
	o := opt.New(opts...)
	var config Config
	if err := o.UnmarshalConfig(&config); err != nil {
		return err
	}
	return container.Container.Provide(func(http *http.Service, grpc *grpc.Grpc) (*CMux, error) {
		l, err := net.Listen("tcp", config.Address())
		if err != nil {
			return nil, err
		}

		return &CMux{
			Http: http,
			Grpc: grpc,
			Mux:  cmux.New(l),
		}, nil
	}, o.DiOptions()...)
}
