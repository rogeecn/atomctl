package cmd

import (
	"atomctl/templates/migration"
	"atomctl/utils"
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// newMigrationCmd represents the newMigration command
var newMigrationCmd = &cobra.Command{
	Use:     "migration",
	Example: "atomctl new migration create_posts_table",
	Short:   "create new migration file",
	Long:    `create new migration file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			log.Fatal("need migration name")
		}

		addMigration.ID = time.Now().Format("20060102_150405")

		addMigration.MigrationName = strings.TrimSpace(args[0])
		addMigration.SnakeMigrationName = strcase.ToSnake(addMigration.MigrationName)
		addMigration.PascalMigrationName = strcase.ToCamel(addMigration.MigrationName)

		files, err := addMigration.prepareFiles()
		if err != nil {
			return errors.Wrap(err, "file prepare failed")
		}

		if err := addMigration.prepare(); err != nil {
			return errors.Wrap(err, "module prepare failed")
		}

		if err := utils.Generate(files, migration.Templates, addMigration); err != nil {
			return err
		}
		return nil
	},
}
var addMigration = &MigrationGenerator{}

func init() {
	newCmd.AddCommand(newMigrationCmd)
}

type MigrationGenerator struct {
	ID                  string
	MigrationName       string
	SnakeMigrationName  string
	PascalMigrationName string
}

func (m *MigrationGenerator) prepare() error {
	if !utils.IsDir(filepath.Dir(m.migrationPath())) {
		return errors.New("not found Migrations director")
	}

	if err := os.MkdirAll(m.migrationPath(), os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (m *MigrationGenerator) migrationPath() string {
	return "./database/migrations/"
}

func (m *MigrationGenerator) prepareFiles() (map[string]string, error) {
	files := migration.Files
	result := make(map[string]string)
	for tpl, target := range files {
		// get target file name
		if len(target) == 0 {
			if utils.IsTplFile(tpl) {
				target = utils.TplToGo(tpl)
			} else {
				target = tpl
			}
		} else {
			newName := bytes.NewBuffer(nil)
			t, err := template.New("name").Parse(target)
			if err != nil {
				return nil, errors.Wrap(err, "init template failed")
			}
			if err := t.Execute(newName, m); err != nil {
				return nil, errors.Wrapf(err, "generate target file failed, tpl: %s", tpl)
			}
			target = newName.String()
		}
		tplFilePath := "tpl/" + tpl

		result[tplFilePath] = target
	}

	return result, nil
}
