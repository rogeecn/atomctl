package grpc

import (
	"fmt"
	"net"

	"{{.ModuleName}}/pkg/atom/container"
	"{{.ModuleName}}/pkg/atom/opt"

	"google.golang.org/grpc"
)

const DefaultPrefix = "Grpc"

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

type Grpc struct {
	Server *grpc.Server
	config *Config
}

// Serve
func (g *Grpc) Serve() error {
	l, err := net.Listen("tcp", g.config.Address())
	if err != nil {
		return err
	}

	return g.Server.Serve(l)
}

func (g *Grpc) ServeWithListener(ln net.Listener) error {
	return g.Server.Serve(ln)
}
