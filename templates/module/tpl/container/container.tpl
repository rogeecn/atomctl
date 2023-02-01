package container

import (
	"atom/container"
	"atom/{{.Path}}/routes"
	"log"

	"go.uber.org/dig"
)

func init() {
	if err := container.Container.Provide(routes.NewRoute, dig.Group("route")); err != nil {
		log.Fatal(err)
	}
}
