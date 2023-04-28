package cmd

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rogeecn/atomctl/templates/dao"
	"github.com/rogeecn/atomctl/utils"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var addDao = &DaoGenerator{}

// daoCmd represents the dao command
var daoCmd = &cobra.Command{
	Use:     "dao",
	Short:   "create dao in target module path",
	Long:    `new dao file. support chain module moduleA.moduleB.moduleC`,
	Example: "atomctl new dao [module] [name]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("invalid params")
		}
		// a => modules/a/dao
		// a.b => modules/a/modules/b/daos

		file := fmt.Sprintf("modules/%s/dao", strings.ReplaceAll(args[0], ".", "/modules/"))
		if !utils.IsDir(file) {
			return errors.New("module not exists")
		}

		addDao.Path = file
		addDao.Package = getPackage()
		addDao.Name = args[1]
		addDao.PascalName = strcase.ToCamel(addDao.Name)
		addDao.CamelName = strcase.ToLowerCamel(addDao.Name)

		generateFiles, err := addDao.prepareFiles(dao.Files)
		if err != nil {
			return err
		}

		if err := utils.Generate(generateFiles, dao.Templates, addDao); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	newCmd.AddCommand(daoCmd)
}

type DaoGenerator struct {
	Package    string
	Name       string
	Path       string
	PascalName string
	CamelName  string

	PkgName string
}

func (m *DaoGenerator) prepareFiles(files map[string]string) (map[string]string, error) {
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
