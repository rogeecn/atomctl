package http

import (
	"errors"
	"fmt"
	"net"
	"runtime/debug"
	"time"

	log "github.com/sirupsen/logrus"
	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/opt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
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
		OnShutdownSuccess: func() {
			log.Info("http server shutdown success")
		},
		OnShutdownError: func(err error) {
			log.Error("http server shutdown error: ", err)
		},

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
		svc.Engine.Shutdown()
	})
	return listenConfig
}

func (svc *Service) Listener(ln net.Listener) error {
	return svc.Engine.Listener(ln, svc.listenerConfig())
}

func (svc *Service) Serve() error {
	return svc.Engine.Listen(svc.conf.Address(), svc.listenerConfig())
}

func Provide(opts ...opt.Option) error {
	o := opt.New(opts...)
	var config Config
	if err := o.UnmarshalConfig(&config); err != nil {
		return err
	}

	return container.Container.Provide(func() (*Service, error) {
		engine := fiber.New(fiber.Config{
			StrictRouting: true,
		})
		engine.Use(recover.New(recover.Config{
			EnableStackTrace: true,
			StackTraceHandler: func(c fiber.Ctx, e any) {
				log.Error(fmt.Sprintf("panic: %v\n%s\n", e, debug.Stack()))
			},
		}))

		if config.StaticRoute != nil && config.StaticPath != nil {
			engine.Use(config.StaticRoute, config.StaticPath)
		}

		engine.Use(logger.New(logger.Config{
			Format:     `[${ip}:${port}] - [${time}] - ${method} - ${status} - ${path} ${latency} "${ua}"` + "\n",
			TimeFormat: time.RFC1123,
			TimeZone:   "Asia/Shanghai",
		}))

		return &Service{
			Engine: engine,
			conf:   &config,
		}, nil
	}, o.DiOptions()...)
}
