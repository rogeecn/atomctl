package container

import (
	"{{.Package}}/{{.Path}}/routes"

	"github.com/rogeecn/atom/container"
)

func Providers() container.Providers {
	return container.Providers{
		{
			Provider: routes.Provide,
		},
	}
}
