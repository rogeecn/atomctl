package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rogeecn/atomctl/templates/suite"
	"github.com/rogeecn/atomctl/utils"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var addSuite = &SuiteGenerator{}

// suiteCmd represents the suite command
var suiteCmd = &cobra.Command{
	Use:   "suite",
	Short: "create testing suit file for target",
	Long:  `new suite file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("invalid params")
		}

		file := args[0]
		if !utils.IsFile(file) {
			return errors.New("target is not a valid file")
		}

		if filepath.Ext(file) != ".go" {
			return errors.New("target is not a valid go file ")
		}

		addSuite.Path = filepath.Dir(file)
		addSuite.Name = strings.ReplaceAll(filepath.Base(file), ".go", "")
		addSuite.PascalName = strcase.ToCamel(addSuite.Name)
		addSuite.Package = getPackage()

		addSuite.PkgName = filepath.Base(addSuite.Path)

		generateFiles, err := addSuite.prepareFiles(suite.Files)
		if err != nil {
			return err
		}

		if err := utils.Generate(generateFiles, suite.Templates, addSuite); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	newCmd.AddCommand(suiteCmd)
}

type SuiteGenerator struct {
	Name       string
	Package    string
	Path       string
	PascalName string

	PkgName string
}

func (m *SuiteGenerator) prepareFiles(files map[string]string) (map[string]string, error) {
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

		target = filepath.Join(m.Path, target)
		result[tplFilePath] = target
		if utils.IsFile(target) {
			return nil, errors.New(target + " file exists")
		}
	}

	return result, nil
}
