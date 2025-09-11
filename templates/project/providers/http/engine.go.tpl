package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"runtime/debug"
	"time"

	log "github.com/sirupsen/logrus"
	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/opt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/samber/lo"
)

func DefaultProvider() container.ProviderContainer {
	return container.ProviderContainer{
		Provider: Provide,
		Options: []opt.Option{
			opt.Prefix(DefaultPrefix),
		},
	}
}

type Service struct {
	conf   *Config
	Engine *fiber.App
}

func (svc *Service) listenerConfig() fiber.ListenConfig {
	listenConfig := fiber.ListenConfig{
		EnablePrintRoutes: true,
		// DisableStartupMessage: true,
	}

	if svc.conf.Tls != nil {
		if svc.conf.Tls.Cert == "" || svc.conf.Tls.Key == "" {
			panic(errors.New("tls cert and key must be set"))
		}
		listenConfig.CertFile = svc.conf.Tls.Cert
		listenConfig.CertKeyFile = svc.conf.Tls.Key
	}
	container.AddCloseAble(func() {
		svc.Engine.ShutdownWithTimeout(time.Second * 10)
	})
	return listenConfig
}

func (svc *Service) Listener(ln net.Listener) error {
	return svc.Engine.Listener(ln, svc.listenerConfig())
}

func (svc *Service) Serve(ctx context.Context) error {
	ln, err := net.Listen("tcp", svc.conf.Address())
	if err != nil {
		return err
	}

	// Run the server in a goroutine so we can listen for context cancellation
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- svc.Engine.Listener(ln, svc.listenerConfig())
	}()

	select {
	case <-ctx.Done():
		// Shutdown the server gracefully
		if shutdownErr := svc.Engine.Shutdown(); shutdownErr != nil {
			return shutdownErr
		}
		// treat context cancellation as graceful shutdown
		return nil
	case err := <-serverErr:
		return err
	}
}

func Provide(opts ...opt.Option) error {
	o := opt.New(opts...)
	var config Config
	if err := o.UnmarshalConfig(&config); err != nil {
		return err
	}

	return container.Container.Provide(func() (*Service, error) {
		engine := fiber.New(fiber.Config{
			StrictRouting:      true,
			CaseSensitive:      true,
			BodyLimit:          10 * 1024 * 1024, // 10 MiB
			ReadTimeout:        10 * time.Second,
			WriteTimeout:       10 * time.Second,
			IdleTimeout:        60 * time.Second,
			ProxyHeader:        fiber.HeaderXForwardedFor,
			EnableIPValidation: true,
		})

		// request id first for correlation
		engine.Use(requestid.New())

		// recover with stack + request id
		engine.Use(recover.New(recover.Config{
			EnableStackTrace: true,
			StackTraceHandler: func(c fiber.Ctx, e any) {
				rid := c.Get(fiber.HeaderXRequestID)
				log.WithField("request_id", rid).Error(fmt.Sprintf("panic: %v\n%s\n", e, debug.Stack()))
			},
		}))

		// basic security + compression
		engine.Use(helmet.New())
		engine.Use(compress.New(compress.Config{Level: compress.LevelDefault}))

		// optional CORS based on config
		if config.Cors != nil {
			corsCfg := buildCORSConfig(config.Cors)
			if corsCfg != nil {
				engine.Use(cors.New(*corsCfg))
			}
		}

		// logging with request id and latency
		engine.Use(logger.New(logger.Config{
			// requestid middleware stores ctx.Locals("requestid")
			Format:     `${time} [${ip}] ${method} ${status} ${path} ${latency} rid=${locals:requestid} "${ua}"\n`,
			TimeFormat: time.RFC3339,
			TimeZone:   "Asia/Shanghai",
		}))

		// rate limit (enable standard headers; adjust Max via config if needed)
		engine.Use(limiter.New(limiter.Config{Max: 0}))

		// static files (Fiber v3 Static helper moved; enable via filesystem middleware later)
		// if config.StaticRoute != nil && config.StaticPath != nil { ... }

		// health endpoints
		engine.Get("/healthz", func(c fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })
		engine.Get("/readyz", func(c fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })

		engine.Hooks().OnPostShutdown(func(err error) error {
			if err != nil {
				log.Error("http server shutdown error: ", err)
			}
			log.Info("http server has shutdown success")
			return nil
		})

		return &Service{
			Engine: engine,
			conf:   &config,
		}, nil
	}, o.DiOptions()...)
}

// buildCORSConfig converts provider Cors config into fiber cors.Config
func buildCORSConfig(c *Cors) *cors.Config {
	if c == nil {
		return nil
	}
	if c.Mode == "disabled" {
		return nil
	}
	var (
		origins    []string
		headers    []string
		methods    []string
		exposes    []string
		allowCreds bool
	)
	for _, w := range c.Whitelist {
		if w.AllowOrigin != "" {
			origins = append(origins, w.AllowOrigin)
		}
		if w.AllowHeaders != "" {
			headers = append(headers, w.AllowHeaders)
		}
		if w.AllowMethods != "" {
			methods = append(methods, w.AllowMethods)
		}
		if w.ExposeHeaders != "" {
			exposes = append(exposes, w.ExposeHeaders)
		}
		allowCreds = allowCreds || w.AllowCredentials
	}

	cfg := cors.Config{
		AllowOrigins:     lo.Uniq(origins),
		AllowHeaders:     lo.Uniq(headers),
		AllowMethods:     lo.Uniq(methods),
		ExposeHeaders:    lo.Uniq(exposes),
		AllowCredentials: allowCreds,
	}
	return &cfg
}
