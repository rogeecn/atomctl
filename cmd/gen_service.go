package cmd

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"go.ipao.vip/atomctl/v2/templates"
)

func CommandGenService(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:      "service",
		Short:    "generate services",
		RunE:     commandGenServiceE,
		PostRunE: commandGenProviderE,
	}

	cmd.Flags().String("path", "./app/services", "base path to scan")

	root.AddCommand(cmd)
}

func commandGenServiceE(cmd *cobra.Command, args []string) error {
	path := cmd.Flag("path").Value.String()

	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	type srv struct {
		CamelName   string
		ServiceName string
	}

	// get services from files
	var services []srv
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		if strings.HasSuffix(name, "_test.go") || strings.HasSuffix(name, ".gen.go") ||
			!strings.HasSuffix(name, ".go") {
			continue
		}

		name = strings.TrimSuffix(name, ".go")

		services = append(services, srv{
			CamelName:   lo.PascalCase(name),
			ServiceName: lo.CamelCase(name),
		})
	}

	// 生成 services.gen.go 文件，使用 text/template 渲染 services.go.tpl 模板
	if err := renderTemplateFS(templates.Services, "services/services.go.tpl", services, filepath.Join(path, "services.gen.go"), true); err != nil {
		return err
	}

	return nil
}

// renderTemplateFS 使用 text/template 渲染模板并写入目标文件
func renderTemplateFS(fs embed.FS, tplName string, data interface{}, targetPath string, overwrite bool) error {
	tplContent, err := fs.ReadFile(tplName)
	if err != nil {
		return err
	}
	tpl, err := template.New(tplName).Parse(string(tplContent))
	if err != nil {
		return err
	}
	if !overwrite {
		if _, err := os.Stat(targetPath); err == nil {
			return nil // 文件已存在且不覆盖
		}
	}
	f, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return tpl.Execute(f, data)
}
