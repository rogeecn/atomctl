package model

import (
	"database/sql"
	"fmt"
	"strings"

	db "{{.ModuleName}}/providers/postgres"

	"git.ipao.vip/rogeecn/atom"
	"git.ipao.vip/rogeecn/atom/container"
	"github.com/go-jet/jet/v2/generator/metadata"
	"github.com/go-jet/jet/v2/generator/postgres"
	"github.com/go-jet/jet/v2/generator/template"
	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/gofiber/fiber/v3/log"
	_ "github.com/lib/pq"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/dig"
)

func Default(providers ...container.ProviderContainer) container.Providers {
	return append(container.Providers{
		db.DefaultProvider(),
	}, providers...)
}

func Options() []atom.Option {
	return []atom.Option{
		atom.Name("model"),
		atom.Short("run model generator"),
		atom.RunE(Serve),
		atom.Providers(Default()),
		atom.Arguments(func(cmd *cobra.Command) {
			cmd.Flags().String("path", "./database/models", "generate to path")
			cmd.Flags().String("transform", "./database/.transform.yaml", "transform config")
		}),
	}
}

func Command() atom.Option {
	return atom.Command(Options()...)
}

type Migrate struct {
	dig.In
	DB     *sql.DB
	Config *db.Config
}

type Transformer struct {
	Ignores []string                     `mapstructure:"ignores"`
	Types   map[string]map[string]string `mapstructure:"types"`
}

func Serve(cmd *cobra.Command, args []string) error {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(cmd.Flag("transform").Value.String())

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	var conf Transformer
	if err := v.Unmarshal(&conf); err != nil {
		return err
	}

	return container.Container.Invoke(func(migrate Migrate) error {
		return postgres.GenerateDSN(
			migrate.Config.DSN(),
			migrate.Config.Schema,
			cmd.Flag("path").Value.String(),
			template.Default(pg.Dialect).
				UseSchema(func(schema metadata.Schema) template.Schema {
					return template.
						DefaultSchema(schema).
						UseModel(
							template.
								DefaultModel().
								UseTable(func(table metadata.Table) template.TableModel {
									if lo.Contains(conf.Ignores, table.Name) {
										table := template.DefaultTableModel(table)
										table.Skip = true
										return table
									}

									return template.DefaultTableModel(table).UseField(func(column metadata.Column) template.TableModelField {
										defaultTableModelField := template.DefaultTableModelField(column)
										defaultTableModelField = defaultTableModelField.UseTags(fmt.Sprintf(`json:"%s"`, column.Name))

										if schema.Name != migrate.Config.Schema {
											return defaultTableModelField
										}

										fields, ok := conf.Types[table.Name]
										if !ok {
											return defaultTableModelField
										}

										toType, ok := fields[column.Name]
										if !ok {
											return defaultTableModelField
										}

										splits := strings.Split(toType, ".")
										typeName := splits[len(splits)-1]

										pkgSplits := strings.Split(splits[0], "/")
										typePkg := pkgSplits[len(pkgSplits)-1]

										defaultTableModelField = defaultTableModelField.
											UseType(template.Type{
												Name:       fmt.Sprintf("%s.%s", typePkg, typeName),
												ImportPath: splits[0],
											})

										log.Infof("Convert table %s field %s type to : %s", table.Name, column.Name, toType)
										return defaultTableModelField
									})
								}),
						)
				}),
		)
	})
}
