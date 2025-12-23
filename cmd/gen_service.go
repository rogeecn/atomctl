package cmd

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.ipao.vip/atomctl/v2/pkg/ast/provider"
	"go.ipao.vip/atomctl/v2/pkg/utils/gomod"
	"go.ipao.vip/atomctl/v2/templates"
)

func CommandGenService(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "generate services",
		Long: `扫描 --path 指定目录（默认 ./app/services）下的 Go 文件，汇总服务名并渲染生成 services.gen.go。

规则：
- 扫描目录中带有 @provider 注释的结构体
- 以结构体名称作为服务名：
  - StructName 作为 ServiceName，用于变量/字段类型
  - PascalCase(StructName) 作为 CamelName，用于导出变量名
- 使用内置模板 services/services.go.tpl 渲染
- 生成完成后会自动运行 gen provider 以补全注入`,
		RunE:     commandGenServiceE,
		PostRunE: commandGenProviderE,
	}

	cmd.Flags().String("path", "./app/services", "base path to scan")

	root.AddCommand(cmd)
}

func commandGenServiceE(cmd *cobra.Command, args []string) error {
	path := cmd.Flag("path").Value.String()
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	// Try to parse go.mod from CWD or target path to ensure parser context
	wd, _ := os.Getwd()
	if err := gomod.Parse(filepath.Join(wd, "go.mod")); err != nil {
		// fallback to check if go.mod is in the target path
		if err := gomod.Parse(filepath.Join(absPath, "go.mod")); err != nil {
			// If both fail, we might still proceed, but parser might lack module info.
			// However, for just getting struct names, it might be fine.
			// Logging warning could be good but we stick to error if critical.
			// provider.ParseDir might depend on it.
		}
	}

	log := log.WithField("path", absPath)
	log.Info("finding service providers...")

	parser := provider.NewGoParser()
	providers, err := parser.ParseDir(absPath)
	if err != nil {
		return err
	}
	log.Infof("found %d providers", len(providers))

	type srv struct {
		CamelName   string
		ServiceName string
	}

	// get services from providers
	var services []srv
	for _, p := range providers {
		name := filepath.Base(p.Location.File)
		if strings.HasSuffix(name, "_test.go") || strings.HasSuffix(name, ".gen.go") ||
			!strings.HasSuffix(name, ".go") {
			log.Warnf("ignore file %s provider, %+v", p.Location.File, p)
			continue
		}

		log.Infof("found service %s", p.StructName)

		services = append(services, srv{
			CamelName:   lo.PascalCase(p.StructName),
			ServiceName: p.StructName,
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
