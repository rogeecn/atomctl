package service

import (
	"{{.ModuleName}}/providers/app"
	"{{.ModuleName}}/providers/events"

	"git.ipao.vip/rogeecn/atom/container"
)

func Default(providers ...container.ProviderContainer) container.Providers {
	return append(container.Providers{
		app.DefaultProvider(),
		events.DefaultProvider(),
	}, providers...)
}
