package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	"github.com/rogeecn/atomctl/templates/http"
	"github.com/rogeecn/atomctl/utils"
	"github.com/samber/lo"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var addHttp = &HttpGenerator{}

// moduleCmd represents the module command
var httpCmd = &cobra.Command{
	Use:     "http",
	Short:   "create http project",
	Example: "atomctl new http [pkg]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("invalid params")
		}
		addHttp.Package = args[0]
		addHttp.Name = filepath.Base(args[0])
		addHttp.Path = filepath.Join(lo.Must1(os.Getwd()), addHttp.Name)
		if utils.IsDir(addHttp.Path) {
			return errors.New("project already exists")
		}

		generateFiles, err := addHttp.prepareFiles(http.Files)
		if err != nil {
			return err
		}

		if err := utils.Generate(generateFiles, http.Templates, addHttp); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	newCmd.AddCommand(httpCmd)
}

type HttpGenerator struct {
	Package string
	Name    string
	Path    string

	PkgName string
}

func (m *HttpGenerator) prepareFiles(files map[string]string) (map[string]string, error) {
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
