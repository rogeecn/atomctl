package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"git.ipao.vip/rogeecn/atomctl/pkg/utils/gomod"
	"git.ipao.vip/rogeecn/atomctl/templates"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CommandNewModule(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "module",
		Short: "new module",
		Args:  cobra.ExactArgs(1),
		RunE:  commandNewModuleE,
	}

	root.AddCommand(cmd)
}

func commandNewModuleE(cmd *cobra.Command, args []string) error {
	module := lo.Filter(strings.Split(args[0], "."), func(s string, _ int) bool {
		return s != ""
	})
	module = append([]string{"app/http"}, module...)

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	modFile := filepath.Join(pwd, "go.mod")
	if err := gomod.Parse(modFile); err != nil {
		return errors.New("parse go.mod file failed")
	}

	moduleName := module[len(module)-1]
	modulePath := filepath.Join(module...)
	log.Infof("new module: %s", modulePath)

	force, _ := cmd.Flags().GetBool("force")
	if _, err := os.Stat(modulePath); err == nil {
		if !force {
			return fmt.Errorf("module directory %s already exists", modulePath)
		}
		log.Warnf("强制删除已存在的目录: %s", modulePath)
		if err := os.RemoveAll(modulePath); err != nil {
			return fmt.Errorf("failed to remove existing directory: %v", err)
		}
	}

	err = os.MkdirAll(modulePath, os.ModePerm)
	if err != nil {
		return err
	}

	err = fs.WalkDir(templates.Module, "module", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel("module", path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(modulePath, strings.TrimSuffix(relPath, ".tpl"))
		destDir := filepath.Dir(destPath)
		err = os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			return err
		}

		tmpl, err := template.ParseFS(templates.Module, path)
		if err != nil {
			return err
		}

		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		return tmpl.Execute(destFile, map[string]string{
			"ModuleName":    moduleName,
			"ProjectModule": gomod.GetModuleName(),
		})
	})

	return err
}
