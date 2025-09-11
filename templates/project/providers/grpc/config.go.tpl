package grpc

import (
	"fmt"
	"net"
	"time"

	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/opt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
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
	// EnableReflection enables grpc/reflection registration when true
	EnableReflection *bool
	// EnableHealth enables gRPC health service registration when true
	EnableHealth *bool
	// ShutdownTimeoutSeconds controls graceful stop timeout; 0 uses default
	ShutdownTimeoutSeconds uint
}

func (h *Config) Address() string {
	if h.Port == 0 {
		h.Port = 8081
	}

	if h.Host == nil {
		return fmt.Sprintf(":%d", h.Port)
	}
	return fmt.Sprintf("%s:%d", *h.Host, h.Port)
}

type Grpc struct {
	Server *grpc.Server
	config *Config

	options            []grpc.ServerOption
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
}

func (g *Grpc) Init() error {
	// merge options and build interceptor chains if provided
	var srvOpts []grpc.ServerOption
	if len(g.unaryInterceptors) > 0 {
		srvOpts = append(srvOpts, grpc.ChainUnaryInterceptor(g.unaryInterceptors...))
	}
	if len(g.streamInterceptors) > 0 {
		srvOpts = append(srvOpts, grpc.ChainStreamInterceptor(g.streamInterceptors...))
	}
	srvOpts = append(srvOpts, g.options...)

	g.Server = grpc.NewServer(srvOpts...)

	// optional reflection and health
	if g.config.EnableReflection != nil && *g.config.EnableReflection {
		reflection.Register(g.Server)
	}
	if g.config.EnableHealth != nil && *g.config.EnableHealth {
		hs := health.NewServer()
		grpc_health_v1.RegisterHealthServer(g.Server, hs)
	}

	// graceful stop with timeout fallback to Stop()
	container.AddCloseAble(func() {
		timeout := g.config.ShutdownTimeoutSeconds
		if timeout == 0 {
			timeout = 10
		}
		done := make(chan struct{})
		go func() {
			g.Server.GracefulStop()
			close(done)
		}()
		select {
		case <-done:
			// graceful stop finished
		case <-time.After(time.Duration(timeout) * time.Second):
			// timeout, force stop
			g.Server.Stop()
		}
	})

	return nil
}

// Serve
func (g *Grpc) Serve() error {
	if g.Server == nil {
		if err := g.Init(); err != nil {
			return err
		}
	}

	l, err := net.Listen("tcp", g.config.Address())
	if err != nil {
		return err
	}

	return g.Server.Serve(l)
}

func (g *Grpc) ServeWithListener(ln net.Listener) error {
	return g.Server.Serve(ln)
}

// UseOptions appends gRPC ServerOptions to be applied when constructing the server.
func (g *Grpc) UseOptions(opts ...grpc.ServerOption) {
	g.options = append(g.options, opts...)
}

// UseUnaryInterceptors appends unary interceptors to be chained.
func (g *Grpc) UseUnaryInterceptors(inters ...grpc.UnaryServerInterceptor) {
	g.unaryInterceptors = append(g.unaryInterceptors, inters...)
}

// UseStreamInterceptors appends stream interceptors to be chained.
func (g *Grpc) UseStreamInterceptors(inters ...grpc.StreamServerInterceptor) {
	g.streamInterceptors = append(g.streamInterceptors, inters...)
}

// Reset clears all configured options and interceptors.
// Useful in tests to ensure isolation.
func (g *Grpc) Reset() {
	g.options = nil
	g.unaryInterceptors = nil
	g.streamInterceptors = nil
}
