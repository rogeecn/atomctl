package main

import (
	"sort"

	"{{ .Pkg }}/database/migrations"

	"github.com/go-gormigrate/gormigrate/v2"
	{{- if eq .Driver "mysql" }}
	dbProvider "github.com/atom-providers/database-mysql"
	{{- else if eq .Driver "postgres" }}
	dbProvider "github.com/atom-providers/database-postgres"
	{{- else if eq .Driver "sqlite" }}
	dbProvider "github.com/atom-providers/database-sqlite"
	{{- end }}
	"github.com/rogeecn/atom"
	"github.com/atom-providers/log"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
	"gorm.io/gorm"
	"github.com/rogeecn/atom/container"
	"github.com/rogeecn/atom/contracts"
)

type MigrationInfo struct {
	dig.In
	DB         *gorm.DB
	Migrations []contracts.Migration `group:"migrations"`
}

func main() {
	for _, migration := range migrations.Migrations {
		if err := container.Container.Provide(migration, dig.Group("migrations")); err != nil {
			log.Fatal(err)
		}
	}

	providers := container.Providers{
		log.DefaultProvider(),
		dbProvider.DefaultProvider(),
	}

	opts := []atom.Option{
		atom.Name("{{ .AppName }}"),
		atom.RunE(func(cmd *cobra.Command, args []string) error {
			return container.Container.Invoke(func(mi MigrationInfo) error {
				migrateToId := "{{ .MigrateToId }}"
				m := gormigrate.New(mi.DB, gormigrate.DefaultOptions, sortedMigrations(mi.Migrations))
				if len(migrateToId) > 0 {
					log.Infof("migrate down to [%s]\n", migrateToId)
					m.RollbackTo(migrateToId)
				}
				return m.RollbackLast()
			})
		}),
	}
	if err := atom.Serve(providers, opts...); err != nil {
		log.Fatal(err)
	}
}

func sortedMigrations(ms []contracts.Migration) []*gormigrate.Migration {
	migrationKeys := []string{}
	migrationMaps := make(map[string]*gormigrate.Migration)
	for _, m := range ms {
		migrationKeys = append(migrationKeys, m.ID())
		migrationMaps[m.ID()] = &gormigrate.Migration{
			ID:       m.ID(),
			Migrate:  m.Up,
			Rollback: m.Down,
		}
	}
	sort.Strings(migrationKeys)

	migrations := []*gormigrate.Migration{}
	for _, key := range migrationKeys {
		migrations = append(migrations, migrationMaps[key])
	}

	return migrations
}
