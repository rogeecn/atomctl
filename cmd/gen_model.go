package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	astModel "go.ipao.vip/atomctl/pkg/ast/model"

	"github.com/go-jet/jet/v2/generator/metadata"
	"github.com/go-jet/jet/v2/generator/postgres"
	"github.com/go-jet/jet/v2/generator/template"
	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/jackc/pgconn"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.ipao.vip/atomctl/pkg/ast/model"
	pgDatabase "go.ipao.vip/atomctl/pkg/postgres"
	"go.ipao.vip/atomctl/pkg/utils/gomod"
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

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile("database/transform.yaml")

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	var transformer model.Transformer
	if err := v.Unmarshal(&transformer); err != nil {
		return err
	}

	jsonReg := regexp.MustCompile(`Json\[\[?\]?(\w+)\]`)
	builtinTypes := []string{
		"string",
		"int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64",
		"bool",
	}

	generatedTables := []string{}

	err = postgres.GenerateDSN(
		dbConf.DSN(),
		dbConf.Schema,
		"database",
		template.Default(pg.Dialect).
			UseSchema(func(schema metadata.Schema) template.Schema {
				return template.
					DefaultSchema(schema).
					UseSQLBuilder(
						template.
							DefaultSQLBuilder().
							UseTable(func(table metadata.Table) template.TableSQLBuilder {
								tbl := template.DefaultTableSQLBuilder(table)

								if lo.Contains(transformer.Ignores.Jet, table.Name) {
									tbl.Skip = true
									log.Infof("Skip table %s", table.Name)
								}
								return tbl
							}).
							UseEnum(func(meta metadata.Enum) template.EnumSQLBuilder {
								enum := template.DefaultEnumSQLBuilder(meta)
								if lo.Contains(transformer.Ignores.Jet, meta.Name) {
									enum.Skip = true
									log.Infof("Skip enum %s", meta.Name)
								}
								return enum
							}),
					).
					UseModel(
						template.
							DefaultModel().
							UseEnum(func(meta metadata.Enum) template.EnumModel {
								enum := template.DefaultEnumModel(meta)
								if lo.Contains(transformer.Ignores.Jet, meta.Name) {
									enum.Skip = true
									log.Infof("Skip enum  %s", meta.Name)
								}
								return enum
							}).
							UseTable(func(table metadata.Table) template.TableModel {
								tbl := template.DefaultTableModel(table)
								if lo.Contains(transformer.Ignores.Jet, table.Name) {
									tbl.Skip = true
									log.Infof("Skip table %s", table.Name)
									return tbl
								}

								generatedTables = append(generatedTables, table.Name)

								return tbl.UseField(func(column metadata.Column) template.TableModelField {
									defaultTableModelField := template.DefaultTableModelField(column)
									defaultTableModelField = defaultTableModelField.UseTags(fmt.Sprintf(`json:"%s"`, column.Name))

									if schema.Name != dbConf.Schema {
										return defaultTableModelField
									}

									fields, ok := transformer.Types[table.Name]
									if !ok {
										return defaultTableModelField
									}

									toType, ok := fields[column.Name]
									if !ok {
										return defaultTableModelField
									}

									// toType = jsonReg.ReplaceAllString(toType, "fields.$1")
									if jsonReg.MatchString(toType) {
										matches := jsonReg.FindStringSubmatch(toType)
										if len(matches) == 2 && !lo.Contains(builtinTypes, matches[1]) {
											toType = strings.Replace(toType, matches[1], "fields."+matches[1], 1)
										}
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
	if err != nil {
		return err
	}

	cfg, err := pgconn.ParseConfig(dbConf.DSN())
	if err != nil {
		return err
	}

	if err := os.RemoveAll("database/schemas"); err != nil {
		return err
	}

	dataPath := fmt.Sprintf("database/%s", cfg.Database)
	if err := os.Rename(dataPath, "database/schemas"); err != nil {
		return err
	}

	if err := astModel.Generate(generatedTables, transformer); err != nil {
		return err
	}
	return nil
}
