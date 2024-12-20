package cmd

import (
	"fmt"
	"strings"

	pgDatabase "git.ipao.vip/rogeecn/atomctl/pkg/postgres"
	"github.com/go-jet/jet/v2/generator/metadata"
	"github.com/go-jet/jet/v2/generator/postgres"
	"github.com/go-jet/jet/v2/generator/template"
	pg "github.com/go-jet/jet/v2/postgres"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CommandGenModel(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "model",
		Short: "Generate jet models",
		RunE:  commandGenModelE,
	}

	root.AddCommand(cmd)
}

func commandGenModelE(cmd *cobra.Command, args []string) error {
	_, dbConf, err := pgDatabase.GetDB(cmd.Flag("config").Value.String())
	if err != nil {
		return errors.Wrap(err, "get db")
	}

	type Transformer struct {
		Ignores []string                     `mapstructure:"ignores"`
		Types   map[string]map[string]string `mapstructure:"types"`
	}

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile("database/transform.yaml")

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	var conf Transformer
	if err := v.Unmarshal(&conf); err != nil {
		return err
	}

	return postgres.GenerateDSN(
		dbConf.DSN(),
		dbConf.Schema,
		"database/models",
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

									if schema.Name != dbConf.Schema {
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
}
