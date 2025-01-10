package service

import (
	"{{.ModuleName}}/providers/app"
	"{{.ModuleName}}/providers/event"

	"git.ipao.vip/rogeecn/atom/container"
)

func Default(providers ...container.ProviderContainer) container.Providers {
	return append(container.Providers{
		app.DefaultProvider(),
		event.DefaultProvider(),
	}, providers...)
}
