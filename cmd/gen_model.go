package cmd

import (
	"fmt"
	"strings"

	pgDatabase "git.ipao.vip/rogeecn/atomctl/pkg/postgres"
	"git.ipao.vip/rogeecn/atomctl/pkg/utils/gomod"
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
		Use:     "model",
		Aliases: []string{"m"},
		Short:   "Generate jet models",
		RunE:    commandGenModelE,
	}

	root.AddCommand(cmd)
}

func commandGenModelE(cmd *cobra.Command, args []string) error {
	if err := gomod.Parse("go.mod"); err != nil {
		return errors.Wrap(err, "parse go.mod")
	}

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

									if strings.Contains(toType, "[") && strings.HasSuffix(toType, "]") {
										toType = strings.Replace(toType, "[", "[fields.", 1)
									}

									defaultTableModelField = defaultTableModelField.
										UseType(template.Type{
											Name:       fmt.Sprintf("fields.%s", toType),
											ImportPath: fmt.Sprintf("%s/database/fields", gomod.GetModuleName()),
										})

									log.Infof("Convert table %s field %s type to : %s", table.Name, column.Name, toType)
									return defaultTableModelField
								})
							}),
					)
			}),
	)
}
