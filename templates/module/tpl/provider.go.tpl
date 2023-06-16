package {{.Name}}

import (
	"{{.Package}}/{{.Path}}/dao"
	"{{.Package}}/{{.Path}}/service"
	"{{.Package}}/{{.Path}}/controller"
	"{{.Package}}/{{.Path}}/routes"

	"github.com/rogeecn/atom/container"
)

func Providers() container.Providers {
	return container.Providers{
		{ Provider: dao.Provide },
		{ Provider: service.Provide },
		{ Provider: controller.Provide },
		{ Provider: routes.Provide },
	}
}
