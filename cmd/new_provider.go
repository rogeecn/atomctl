package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"git.ipao.vip/rogeecn/atomctl/templates"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

// CommandNewProvider 注册 new_provider 命令
func CommandNewProvider(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "provider",
		Short: "创建新的 provider",
		Args:  cobra.ExactArgs(1),
		RunE:  commandNewProviderE,
	}

	root.AddCommand(cmd)
}

func commandNewProviderE(cmd *cobra.Command, args []string) error {
	providerName := args[0]
	targetPath := filepath.Join("providers", providerName)

	if _, err := os.Stat(targetPath); err == nil {
		return fmt.Errorf("目录 %s 已存在", targetPath)
	}

	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		return err
	}

	err := fs.WalkDir(templates.Provider, "provider", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel("provider", path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(targetPath, strings.TrimSuffix(relPath, ".tpl"))
		if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
			return err
		}

		tmpl, err := template.ParseFS(templates.Provider, path)
		if err != nil {
			return err
		}

		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		return tmpl.Execute(destFile, map[string]string{
			"Name":      providerName,
			"CamelName": strcase.ToCamel(providerName),
		})
	})
	if err != nil {
		return errors.New("渲染 provider 模板失败")
	}

	fmt.Printf("Provider 已创建: %s\n", targetPath)
	return nil
}
