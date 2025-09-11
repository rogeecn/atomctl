package cmux

import (
	"fmt"
	"net"
	"time"

	"{{.ModuleName}}/providers/grpc"
	"{{.ModuleName}}/providers/http"

	log "github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/opt"
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
	Base net.Listener
}

func (c *CMux) Serve() error {
	// Protect against slowloris connections when sniffing protocol
	// Safe even if SetReadTimeout is a no-op in the cmux version in use
	c.Mux.SetReadTimeout(1 * time.Second)

	addr := ""
	if c.Base != nil && c.Base.Addr() != nil {
		addr = c.Base.Addr().String()
	}
	log.WithFields(log.Fields{
		"addr": addr,
	}).Info("cmux starting")

	// Route classic HTTP/1.x traffic to the HTTP service
	httpL := c.Mux.Match(cmux.HTTP1Fast())

	// Route gRPC (HTTP/2 with content-type application/grpc) to the gRPC service.
	// Additionally, send other HTTP/2 traffic to gRPC since Fiber (HTTP) does not serve HTTP/2.
	grpcL := c.Mux.Match(
		cmux.HTTP2HeaderField("content-type", "application/grpc"),
		cmux.HTTP2(),
	)

	var eg errgroup.Group
	eg.Go(func() error {
		log.WithField("addr", addr).Info("grpc serving via cmux")
		err := c.Grpc.ServeWithListener(grpcL)
		if err != nil {
			log.WithError(err).Error("grpc server exited with error")
		} else {
			log.Info("grpc server exited")
		}
		return err
	})

	eg.Go(func() error {
		log.WithField("addr", addr).Info("http serving via cmux")
		err := c.Http.Listener(httpL)
		if err != nil {
			log.WithError(err).Error("http server exited with error")
		} else {
			log.Info("http server exited")
		}
		return err
	})

	// Run cmux dispatcher; wait for the first error from any goroutine
	eg.Go(func() error {
		err := c.Mux.Serve()
		if err != nil {
			log.WithError(err).Error("cmux exited with error")
		} else {
			log.Info("cmux exited")
		}
		return err
	})
	err := eg.Wait()
	if err == nil {
		log.Info("cmux and sub-servers exited cleanly")
	}
	return err
}
