package main

import (
	"{{ .Pkg }}/database/seeders"

	{{- if eq .Driver "mysql" }}
	dbProvider "github.com/atom-providers/database-mysql"
	{{- else if eq .Driver "postgres" }}
	dbProvider "github.com/atom-providers/database-postgres"
	{{- else if eq .Driver "sqlite" }}
	dbProvider "github.com/atom-providers/database-sqlite"
	{{- end }}
	"github.com/atom-providers/faker"
	"github.com/atom-providers/app"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/rogeecn/atom"
	"github.com/atom-providers/log"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
	"gorm.io/gorm"

	"github.com/rogeecn/atom/container"
	"github.com/rogeecn/atom/contracts"
)

type SeedersContainer struct {
	dig.In

	DB      *gorm.DB
	Faker   *gofakeit.Faker
	Seeders []contracts.Seeder `group:"seeders"`
}

func main() {
	for _, seeder := range seeders.Seeders {
		if err := container.Container.Provide(seeder, dig.Group("seeders")); err != nil {
			log.Fatal(err)
		}
	}

	providers := container.Providers{
		log.DefaultProvider(),
		app.DefaultProvider(),
		faker.DefaultProvider(),
		dbProvider.DefaultProvider(),
	}

	opts := []atom.Option{
		atom.Name("{{.AppName}}"),
		atom.RunE(func(cmd *cobra.Command, args []string) error {
			return container.Container.Invoke(func(c SeedersContainer) error {
				if len(c.Seeders) == 0 {
					log.Info("no seeder exists")
					return nil
				}

				for _, seeder := range c.Seeders {
					seeder.Run(c.Faker, c.DB)
				}
				return nil
			})
		}),
	}
	if err := atom.Serve(providers, opts...); err != nil {
		log.Fatal(err)
	}
}
