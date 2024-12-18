package migrate

import (
	"context"
	"database/sql"

	"{{.ModuleName}}/database"
	"{{.ModuleName}}/providers/postgres"

	"git.ipao.vip/rogeecn/atom"
	"git.ipao.vip/rogeecn/atom/container"
	"github.com/pressly/goose/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func Default(providers ...container.ProviderContainer) container.Providers {
	return append(container.Providers{
		postgres.DefaultProvider(),
	}, providers...)
}

func Command() atom.Option {
	return atom.Command(
		atom.Name("migrate"),
		atom.Short("run migrations"),
		atom.RunE(Serve),
		atom.Providers(Default()),
	)
}

type Migrate struct {
	dig.In
	DB *sql.DB
}

func Serve(cmd *cobra.Command, args []string) error {
	return container.Container.Invoke(func(migrate Migrate) error {
		if len(args) == 0 {
			args = append(args, "up")
		}

		action, args := args[0], args[1:]
		log.Infof("migration action: %s args: %+v", action, args)

		goose.SetBaseFS(database.MigrationFS)
		goose.SetTableName("migrations")

		return goose.RunContext(context.Background(), action, migrate.DB, "migrations", args...)
	})
}
