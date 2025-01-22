package {{.Name}}

import (
	"{{.ModuleName}}/pkg/atom/container"
	"{{.ModuleName}}/pkg/atom/opt"
)

const DefaultPrefix = "{{.CamelName}}"

func DefaultProvider() container.ProviderContainer {
	return container.ProviderContainer{
		Provider: Provide,
		Options: []opt.Option{
			opt.Prefix(DefaultPrefix),
		},
	}
}

type Config struct {
}