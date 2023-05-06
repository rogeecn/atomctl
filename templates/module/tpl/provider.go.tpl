package {{.Name}}

import (
	"{{.Package}}/{{.Path}}/controller"
	"{{.Package}}/{{.Path}}/service"
	"{{.Package}}/{{.Path}}/routes"

	"github.com/rogeecn/atom/container"
)

func Providers() container.Providers {
	return container.Providers{
		{ Provider: controller.Provide },
		{ Provider: service.Provide },
		{ Provider: routes.Provide },
	}
}
