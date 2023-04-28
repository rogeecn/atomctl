package cmd

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rogeecn/atomctl/templates/seeder"
	"github.com/rogeecn/atomctl/utils"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// newSeederCmd represents the newSeeder command
var newSeederCmd = &cobra.Command{
	Use:     "seeder",
	Example: "atomctl new seeder [table_name]",
	Short:   "create new seeder file",
	Long:    `create new seeder file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			log.Fatal("need Seeder name")
		}

		addSeeder.Package = ModPath
		addSeeder.SeederName = strings.TrimSpace(args[0])
		addSeeder.SnakeSeederName = strcase.ToSnake(addSeeder.SeederName)
		addSeeder.PascalSeederName = strcase.ToCamel(addSeeder.SeederName)

		files, err := addSeeder.prepareFiles()
		if err != nil {
			return errors.Wrap(err, "file prepare failed")
		}

		if err := addSeeder.prepare(); err != nil {
			return errors.Wrap(err, "module prepare failed")
		}

		if err := utils.Generate(files, seeder.Templates, addSeeder); err != nil {
			return err
		}
		return nil
	},
}
var addSeeder = &SeederGenerator{}

func init() {
	newCmd.AddCommand(newSeederCmd)
}

type SeederGenerator struct {
	Package          string
	SeederName       string
	SnakeSeederName  string
	PascalSeederName string
}

func (m *SeederGenerator) prepare() error {
	if !utils.IsDir(filepath.Dir(m.SeederPath())) {
		return errors.New("not found Seeders director")
	}

	if err := os.MkdirAll(m.SeederPath(), os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (m *SeederGenerator) SeederPath() string {
	return "./database/seeders/"
}

func (m *SeederGenerator) prepareFiles() (map[string]string, error) {
	files := seeder.Files
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
