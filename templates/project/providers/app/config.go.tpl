package app

import (
	"{{.ModuleName}}/pkg/atom/container"
	"{{.ModuleName}}/pkg/atom/opt"
)

const DefaultPrefix = "App"

func DefaultProvider() container.ProviderContainer {
	return container.ProviderContainer{
		Provider: Provide,
		Options: []opt.Option{
			opt.Prefix(DefaultPrefix),
		},
	}
}

// swagger:enum AppMode
// ENUM(development, release, test)
type AppMode string

type Config struct {
	Mode    AppMode
	Cert    *Cert
	BaseURI *string
}

func (c *Config) IsDevMode() bool {
	return c.Mode == AppModeDevelopment
}

func (c *Config) IsReleaseMode() bool {
	return c.Mode == AppModeRelease
}

func (c *Config) IsTestMode() bool {
	return c.Mode == AppModeTest
}

type Cert struct {
	CA   string
	Cert string
	Key  string
}
