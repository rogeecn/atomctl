package cmux

import (
	"fmt"
	"qq/providers/grpc"
	"qq/providers/http"

	"git.ipao.vip/rogeecn/atom/container"
	"git.ipao.vip/rogeecn/atom/utils/opt"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
)

const DefaultPrefix = "Cmux"

func DefaultProvider() container.ProviderContainer {
	return container.ProviderContainer{
		Provider: Provide,
		Options: []opt.Option{
			opt.Prefix(DefaultPrefix),
		},
	}
}

type Config struct {
	Host *string
	Port uint
}

func (h *Config) Address() string {
	if h.Host == nil {
		return fmt.Sprintf(":%d", h.Port)
	}
	return fmt.Sprintf("%s:%d", *h.Host, h.Port)
}

type CMux struct {
	Http *http.Service
	Grpc *grpc.Grpc
	Mux  cmux.CMux
}

func (c *CMux) Serve() error {
	// grpcL := c.Mux.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	// httpL := c.Mux.Match(cmux.HTTP1Fast())
	// httpL := c.Mux.Match(cmux.Any())
	httpL := c.Mux.Match(cmux.HTTP1Fast())
	grpcL := c.Mux.Match(cmux.Any())

	var eg errgroup.Group
	eg.Go(func() error {
		return c.Grpc.ServeWithListener(grpcL)
	})

	eg.Go(func() error {
		return c.Http.Listener(httpL)
	})

	return c.Mux.Serve()
}
