package main

import (
	"errors"

	{{- if eq .Driver "mysql" }}
	mysqlProvider "github.com/rogeecn/atom-addons/providers/database/mysql"
	{{- else if eq .Driver "postgres" }}
	postgresProvider "github.com/rogeecn/atom-addons/providers/database/postgres"
	{{- else if eq .Driver "sqlite" }}
	sqliteProvider "github.com/rogeecn/atom-addons/providers/database/sqlite"
	{{- end }}
	"github.com/rogeecn/atom-addons/providers/log"

	"github.com/rogeecn/atom"
	"github.com/rogeecn/atom/container"
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
	}

	{{- if eq .Driver "mysql" }}
	providers = append(providers, mysqlProvider.DefaultProvider())
	{{- else if eq .Driver "postgres" }}
	providers = append(providers, postgresProvider.DefaultProvider())
	{{- else if eq .Driver "sqlite" }}
	providers = append(providers, sqliteProvider.DefaultProvider())
	{{- end }}

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
				g.WithOpts(gen.WithMethod(gen.DefaultMethodTableWithNamer))
				g.WithImportPkgPath(
					"gorm.io/datatypes",
					"{{ .PkgName }}/common",
				)

				transforms := getConvertModelFields()
				for from, to := range transforms {
					g.WithOpts(gen.FieldType(from, to))
					g.WithOpts(gen.FieldGenType(from, "Field"))
				}


				g.UseDB(gq.DB) // reuse your gorm db

				models := []interface{}{}
				for _, table := range tables {
					models = append(models, g.GenerateModel(table))
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



func getConvertModelFields() map[string]string {
	convert := make(map[string]string)

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
		if strings.TrimSpace(line) == "" {
			continue
		}

		items := strings.Split(line, "=>")
		if len(items) != 2 {
			log.Println("invalid line: ", line)
			continue
		}
		toStruct := strings.TrimSpace(items[1])
		convert[strings.TrimSpace(items[0])] = toStruct
	}
	return convert
}