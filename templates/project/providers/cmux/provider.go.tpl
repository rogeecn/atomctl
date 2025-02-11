package cmux

import (
	"net"

	"{{.ModuleName}}/providers/grpc"
	"{{.ModuleName}}/providers/http"

	"github.com/soheilhy/cmux"
	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/opt"
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
