package {{.Name}}

import (
	"{{.Package}}/{{.Path}}/service"
	"{{.Package}}/{{.Path}}/controller"
	"{{.Package}}/{{.Path}}/routes"

	"github.com/rogeecn/atom/container"
)

func Providers() container.Providers {
	return container.Providers{
		{ Provider: service.Provide },
		{ Provider: controller.Provide },
		{ Provider: routes.Provide },
	}
}
