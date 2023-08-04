package test

import (
	"github.com/rogeecn/atomctl/modules/test/dao"
	"github.com/rogeecn/atomctl/modules/test/service"
	"github.com/rogeecn/atomctl/modules/test/controller"
	"github.com/rogeecn/atomctl/modules/test/routes"

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
