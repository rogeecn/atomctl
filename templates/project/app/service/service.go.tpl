package service

import (
	"{{.ModuleName}}/pkg/atom/container"
	"{{.ModuleName}}/providers/app"
	"{{.ModuleName}}/providers/event"
)

func Default(providers ...container.ProviderContainer) container.Providers {
	return append(container.Providers{
		app.DefaultProvider(),
		event.DefaultProvider(),
	}, providers...)
}
