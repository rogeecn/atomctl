package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"go.ipao.vip/atomctl/v2/templates"
)

// CommandNewProvider 注册 new_provider 命令
func CommandNewProvider(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "provider",
		Short: "创建新的 provider",
		Long: `在 providers/<name> 目录下渲染创建 Provider 模板。

行为：
- 当 name 与内置预置目录同名时，渲染该目录；否则回退渲染 providers/default
- 从内置模板 templates/providers 渲染相关文件
- 使用 name 的 CamelCase 作为导出名
- --dry-run 仅打印渲染与写入动作；--dir 指定输出基目录（默认 .）

不带名称直接运行时，将列出可用预置 provider 列表供选择。

示例：
  atomctl new provider email
  atomctl new --dry-run --dir ./demo provider cache`,
		Args: cobra.MaximumNArgs(1),
		RunE: commandNewProviderE,
	}

	root.AddCommand(cmd)
}

func commandNewProviderE(cmd *cobra.Command, args []string) error {
	// no-arg: list available preset providers
	if len(args) == 0 {
		entries, err := templates.Providers.ReadDir("providers")
		if err != nil {
			return err
		}
		var names []string
		for _, e := range entries {
			if e.IsDir() {
				names = append(names, e.Name())
			}
		}
		sort.Strings(names)
		fmt.Println("可用预置 providers:")
		for _, n := range names {
			fmt.Printf("  - %s\n", n)
		}
		return nil
	}

	providerName := args[0]
	// shared flags
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	baseDir, _ := cmd.Flags().GetString("dir")

	targetPath := filepath.Join(baseDir, "providers", providerName)

	if _, err := os.Stat(targetPath); err == nil {
		return fmt.Errorf("目录 %s 已存在", targetPath)
	}

	if dryRun {
		fmt.Printf("[dry-run] mkdir -p %s\n", targetPath)
	} else {
		if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
			return err
		}
	}

	// choose template source: providers/<name> or providers/default
	srcDir := filepath.Join("providers", providerName)
	if _, err := templates.Providers.ReadDir(srcDir); err != nil {
		srcDir = filepath.Join("providers", "default")
	}

	err := fs.WalkDir(templates.Providers, srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(targetPath, strings.TrimSuffix(relPath, ".tpl"))
		if dryRun {
			fmt.Printf("[dry-run] mkdir -p %s\n", filepath.Dir(destPath))
		} else {
			if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
				return err
			}
		}

		tmpl, err := template.ParseFS(templates.Providers, path)
		if err != nil {
			return err
		}

		if dryRun {
			fmt.Printf("[dry-run] render > %s\n", destPath)
			return nil
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
