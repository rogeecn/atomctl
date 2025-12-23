package migrate

import (
	"context"
	"database/sql"

	"{{.ModuleName}}/app/commands"
	"{{.ModuleName}}/database"
	"{{.ModuleName}}/providers/postgres"

	"github.com/pressly/goose/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.ipao.vip/atom"
	"go.ipao.vip/atom/container"
	"go.uber.org/dig"

	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/rivermigrate"
)

func defaultProviders() container.Providers {
	return commands.Default(container.Providers{
		postgres.DefaultProvider(),
	}...)
}

func Command() atom.Option {
	return atom.Command(
		atom.Name("migrate"),
		atom.Short("run migrations"),
		atom.RunE(Serve),
		atom.Providers(defaultProviders()),
		atom.Example("migrate [up|up-by-one|up-to|create|down|down-to|fix|redo|reset|status|version]"),
	)
}

type Service struct {
	dig.In

	DB *sql.DB
}

// migrate
func Serve(cmd *cobra.Command, args []string) error {
	return container.Container.Invoke(func(ctx context.Context, svc Service) error {
		if len(args) == 0 {
			args = append(args, "up")
		}

		if args[0] == "create" {
			return nil
		}

		action, args := args[0], args[1:]
		log.Infof("migration action: %s args: %+v", action, args)

		goose.SetBaseFS(database.MigrationFS)
		goose.SetTableName("migrations")
		goose.AddNamedMigrationNoTxContext(
			"10000000000001_river_job.go",
			func(ctx context.Context, db *sql.DB) error {
				migrator, err := rivermigrate.New(riverdatabasesql.New(db), nil)
				if err != nil {
					return err
				}

				_, err = migrator.Migrate(ctx, rivermigrate.DirectionUp, &rivermigrate.MigrateOpts{TargetVersion: -1})
				return err
			},
			func(ctx context.Context, db *sql.DB) error {
				migrator, err := rivermigrate.New(riverdatabasesql.New(db), nil)
				if err != nil {
					return err
				}

				_, err = migrator.Migrate(ctx, rivermigrate.DirectionDown, &rivermigrate.MigrateOpts{TargetVersion: -1})
				return err
			})

		return goose.RunContext(context.Background(), action, svc.DB, "migrations", args...)
	})
}
