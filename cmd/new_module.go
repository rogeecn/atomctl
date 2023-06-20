package cmd

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/rogeecn/atomctl/templates/module"
	"github.com/rogeecn/atomctl/utils"

	pluralize "github.com/gertd/go-pluralize"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var addModule = &ModuleGenerator{}

// moduleCmd represents the module command
var moduleCmd = &cobra.Command{
	Use:     "module",
	Short:   "create module in target module path",
	Long:    `new module file. support chain module moduleA.moduleB.moduleC`,
	Example: "atomctl new module [name]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("invalid params")
		}
		// a => modules/a
		// a.b => modules/a/modules/b

		file, moduleName := dotToModule(args[0])
		if utils.IsDir(file) {
			return errors.New("module already exists")
		}
		addModule.Name = moduleName
		addModule.NamePlural = moduleName

		pcli := pluralize.NewClient()
		if !pcli.IsPlural(moduleName) {
			addModule.NamePlural = pcli.Plural(moduleName)
		}

		addModule.Path = file
		addModule.Package = getPackage()

		generateFiles, err := addModule.prepareFiles(module.Files)
		if err != nil {
			return err
		}

		if err := utils.Generate(generateFiles, module.Templates, addModule); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	newCmd.AddCommand(moduleCmd)
}

type ModuleGenerator struct {
	Package    string
	Name       string
	NamePlural string
	Path       string

	PkgName string
}

func (m *ModuleGenerator) prepareFiles(files map[string]string) (map[string]string, error) {
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
