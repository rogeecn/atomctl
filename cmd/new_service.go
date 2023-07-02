package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rogeecn/atomctl/templates/service"
	"github.com/rogeecn/atomctl/utils"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var addService = &ServiceGenerator{}

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:     "service",
	Short:   "create service in target module path",
	Long:    `new service file. support chain module moduleA.moduleB.moduleC`,
	Example: "atomctl new service [module] [name]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("invalid params")
		}
		// a => modules/a/service
		// a.b => modules/a/modules/b/services

		file := fmt.Sprintf("modules/%s/service", strings.ReplaceAll(args[0], ".", "/modules/"))
		if !utils.IsDir(file) {
			return errors.New("module not exists")
		}

		addService.Path = file
		addService.Package = getPackage()
		addService.Name = args[1]
		addService.PascalName = strcase.ToCamel(addService.Name)
		addService.CamelName = strcase.ToLowerCamel(addService.Name)

		generateFiles, err := addService.prepareFiles(service.Files)
		if err != nil {
			return err
		}

		if err := utils.Generate(generateFiles, service.Templates, addService); err != nil {
			return err
		}

		// register controller provider
		providerFile := filepath.Join(addService.Path, "provider.go")
		if !utils.IsFile(providerFile) {
			log.Println("[Warn] " + providerFile + " not exists, please add new provider manually")
			return err
		}

		providerContent, err := os.ReadFile(providerFile)
		if err != nil {
			log.Println("[Warn] read " + providerFile + " failed, please add new provider manually")
			return err
		}

		providerFunc := fmt.Sprintf("New%sService", addService.PascalName)

		content := string(providerContent)
		if !strings.Contains(content, providerFunc) {
			provider := fmt.Sprintf("\n\tif err := container.Container.Provide(%s); err!=nil {\n\t\treturn err\n\t}\n\treturn nil", providerFunc)
			content = strings.Replace(content, "return nil", provider, 1)
		}
		containerPackage := `"github.com/rogeecn/atom/container"`
		if !strings.Contains(content, containerPackage) {
			content = strings.Replace(content, "import (", "import (\n\t"+containerPackage, 1)
		}

		return os.WriteFile(providerFile, []byte(content), os.ModePerm)
	},
}

func init() {
	newCmd.AddCommand(serviceCmd)
}

type ServiceGenerator struct {
	Name       string
	Package    string
	Path       string
	PascalName string
	CamelName  string

	PkgName string
}

func (m *ServiceGenerator) prepareFiles(files map[string]string) (map[string]string, error) {
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
