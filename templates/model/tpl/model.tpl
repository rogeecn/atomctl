package main

import (
	"errors"
	"os"
	"strings"

	{{- if eq .Driver "mysql" }}
	dbProvider "github.com/atom-providers/database-mysql"
	{{- else if eq .Driver "postgres" }}
	dbProvider "github.com/atom-providers/database-postgres"
	{{- else if eq .Driver "sqlite" }}
	dbProvider "github.com/atom-providers/database-sqlite"
	{{- end }}
	"github.com/atom-providers/log"
	"github.com/atom-providers/app"

	"github.com/rogeecn/atom"
	"github.com/rogeecn/atom/container"
	"github.com/rogeecn/atom/utils/fs"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
	"gorm.io/gen"
	"gorm.io/gorm"
)

type ModelGenerator struct {
	Driver string
}

// MigrationInfo http service container
type GenQueryGenerator struct {
	dig.In
	DB *gorm.DB
}

func main() {
	providers := container.Providers{
		log.DefaultProvider(),
		app.DefaultProvider(),
		dbProvider.DefaultProvider(),
	}

	opts := []atom.Option{
		atom.Name("{{ .AppName }}"),
		atom.RunE(func(cmd *cobra.Command, args []string) error {
			return container.Container.Invoke(func(gq GenQueryGenerator) error {
				var tables []string

				{{- if eq .Driver "mysql" }}
				err := gq.DB.Raw("show tables").Scan(&tables).Error
				if err != nil {
					return err
				}
				{{- else if eq .Driver "postgres" }}
				err := gq.DB.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tables).Error
				if err != nil {
					return err
				}
				{{- else if eq .Driver "sqlite" }}
				err := gq.DB.Raw("SELECT name FROM sqlite_master WHERE type='table'").Scan(&tables).Error
				if err != nil {
					return err
				}
				{{- end }}

				if len(tables) == 0 {
					return errors.New("no tables in database, run migrate first")
				}

				g := gen.NewGenerator(gen.Config{
					OutPath:          "database/query",
					OutFile:          "query.gen.go",
					ModelPkgPath:     "database/models",
					FieldSignable:    true,
					FieldWithTypeTag: true,
					Mode:             gen.WithDefaultQuery | gen.WithQueryInterface,
				})
				tags := map[string]string{
					"swaggertype": "string",
				}
				g.WithOpts(
					gen.WithMethod(gen.DefaultMethodTableWithNamer),
					gen.FieldNewTag("deleted_at", tags),
				)
				g.WithImportPkgPath(
					"gorm.io/datatypes",
					"github.com/lib/pq",
					"{{ .PkgName }}/common/consts",
				)

				g.UseDB(gq.DB) // reuse your gorm db

				transforms := getConvertModelFields()
				log.Infof("Transforms: %+v", transforms)

				models := []interface{}{}
				for _, table := range tables {
					opts := []gen.ModelOpt{}
					if transform, ok := transforms[table]; ok {
						for from, to := range transform {
							opts = append(opts, gen.FieldType(from, to))
						}
					}
					models = append(models, g.GenerateModel(table, opts...))
				}

				// Generate basic type-safe DAO API for struct `model.User` following conventions
				g.ApplyBasic(models...)

				// Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User` and `model.Company`
				// g.ApplyInterface(func(Querier) {}, model.User{}, model.Company{})

				// Generate the code
				g.Execute()
				return nil
			})
		}),
	}
	if err := atom.Serve(providers, opts...); err != nil {
		log.Fatal(err)
	}
}

func getConvertModelFields() map[string]map[string]string {
	convert := make(map[string]map[string]string)

	filePath := "database/.transform"
	if !fs.FileExist(filePath) {
		return convert
	}

	b, err := os.ReadFile(filePath)
	if err != nil {
		log.Error(err)
		return convert
	}
	lines := strings.Split(string(b), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "//") {
			continue
		}
		if strings.HasPrefix(line, "--") {
			continue
		}
		if line == "" {
			continue
		}

		items := strings.Split(line, "=>")
		if len(items) != 2 {
			log.Warn("invalid line: ", line)
			continue
		}
		tableMark := strings.Split(strings.TrimSpace(items[0]), ".")
		if len(tableMark) != 2 {
			log.Warn("invalid line: ", line)
			continue
		}

		toStruct := strings.TrimSpace(items[1])
		if _, ok := convert[tableMark[0]]; !ok {
			convert[tableMark[0]] = make(map[string]string)
		}
		convert[tableMark[0]][tableMark[1]] = toStruct
	}
	return convert
}
