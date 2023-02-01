package cmd

import (
	"atomctl/templates/controller"
	"atomctl/utils"
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var addController = &ControllerGenerator{}

// controllerCmd represents the controller command
var controllerCmd = &cobra.Command{
	Use:     "controller",
	Short:   "create controller in target module path",
	Long:    `new controller file. support chain module moduleA.moduleB.moduleC`,
	Example: "atomctl new controller [module] [name]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("invalid params")
		}
		// a => modules/a/controller
		// a.b => modules/a/modules/b/controllers

		file := fmt.Sprintf("modules/%s/controller", strings.ReplaceAll(args[0], ".", "/modules/"))
		if !utils.IsDir(file) {
			return errors.New("module not exists")
		}

		addController.Path = file
		addController.Name = args[1]
		addController.PascalName = strcase.ToCamel(addController.Name)

		generateFiles, err := addController.prepareFiles(controller.Files)
		if err != nil {
			return err
		}

		if err := utils.Generate(generateFiles, controller.Templates, addController); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	newCmd.AddCommand(controllerCmd)
}

type ControllerGenerator struct {
	Name       string
	Path       string
	PascalName string

	PkgName string
}

func (m *ControllerGenerator) prepareFiles(files map[string]string) (map[string]string, error) {
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

		result[tplFilePath] = filepath.Join(m.Path, target)
		if utils.IsFile(target) {
			return nil, errors.New(target + " file exists")
		}
	}

	return result, nil
}
