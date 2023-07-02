package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rogeecn/atomctl/templates/controller"
	"github.com/rogeecn/atomctl/utils"

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

		modulePath, _ := dotToModule(args[0])
		file := fmt.Sprintf("%s/controller", modulePath)
		if !utils.IsDir(file) {
			return errors.New("module not exists")
		}

		addController.Path = file
		addController.Package = getPackage()
		addController.Name = args[1]
		addController.PascalName = strcase.ToCamel(addController.Name)
		addController.CamelName = strcase.ToLowerCamel(addController.Name)

		generateFiles, err := addController.prepareFiles(controller.Files)
		if err != nil {
			return err
		}

		if err := utils.Generate(generateFiles, controller.Templates, addController); err != nil {
			return err
		}

		// register controller provider
		providerFile := filepath.Join(addController.Path, "provider.go")
		if !utils.IsFile(providerFile) {
			log.Println("[Warn] " + providerFile + " not exists, please add new provider manually")
			return err
		}

		providerContent, err := os.ReadFile(providerFile)
		if err != nil {
			log.Println("[Warn] read " + providerFile + " failed, please add new provider manually")
			return err
		}

		providerFunc := fmt.Sprintf("New%sController", addController.PascalName)

		content := string(providerContent)
		if !strings.Contains(content, providerFunc) {
			provider := fmt.Sprintf("if err := container.Container.Provide(%s); err!=nil {\n\treturn err\n\t}\n\treturn nil", providerFunc)
			content = strings.Replace(content, "return nil", provider, 1)
		}

		return os.WriteFile(providerFile, []byte(content), os.ModePerm)
	},
}

func init() {
	newCmd.AddCommand(controllerCmd)
}

type ControllerGenerator struct {
	Package    string
	Name       string
	Path       string
	PascalName string
	CamelName  string

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

		target = filepath.Join(m.Path, target)
		result[tplFilePath] = target
		if utils.IsFile(target) {
			return nil, errors.New(target + " file exists")
		}
	}

	return result, nil
}
