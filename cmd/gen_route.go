package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.ipao.vip/atomctl/v2/pkg/ast/route"
	"go.ipao.vip/atomctl/v2/pkg/utils/gomod"
)

func CommandGenRoute(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "route",
		Short: "generate routes",
		Long: `扫描项目控制器，解析注释生成 routes.gen.go。

用法与规则：
- 扫描根目录通过 --path 指定（默认 CWD），会在 <path>/app/http 下递归搜索。
- 使用注释定义路由与参数绑定：
  - @Router <path> [<method>]  例如：@Router /users/:id [get]
  - @Bind <name> (<position>) key(<key>) [model(<field[:field_type]>)]
    - position 枚举：path|query|body|header|cookie|local|file
    - model() 形式：model() 默认 id:int；model(id) 默认 int；model(id:int) 显式类型

参数位置与类型建议：
- path：标量，或结合 model() 从路径值加载模型
- query/header：标量或结构体
- cookie：标量；string 有快捷写法
- body：结构体（或标量）
- file：固定 multipart.FileHeader
- local：任意类型（上下文本地值）

说明：生成完成后会自动运行 gen provider 以补全依赖注入。`,
		RunE:     commandGenRouteE,
		PostRunE: commandGenProviderE,
	}

	cmd.Flags().String("path", ".", "Base path to scan (defaults to CWD)")

	root.AddCommand(cmd)
}

// https://go.ipao.vip/atomctl/v2/pkg/swag?tab=readme-ov-file#api-operation
func commandGenRouteE(cmd *cobra.Command, args []string) error {
	var err error
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	// allow overriding via --path flag
	if f := cmd.Flag("path"); f != nil && f.Value.String() != "." && f.Value.String() != "" {
		path = f.Value.String()
	}

	path, _ = filepath.Abs(path)

	err = gomod.Parse(filepath.Join(path, "go.mod"))
	if err != nil {
		return err
	}

	routes := []route.RouteDefinition{}

	modulePath := filepath.Join(path, "app/http")
	if _, err := os.Stat(modulePath); os.IsNotExist(err) {
		return fmt.Errorf("routes directory not found: %s (set --path to your project root)", modulePath)
	}

	// controllerPattern := regexp.MustCompile(`controller(_?\w+)?\.go`)
	err = filepath.WalkDir(modulePath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		// if !controllerPattern.MatchString(d.Name()) {
		// 	return nil
		// }
		if strings.HasSuffix(path, ".gen.go") {
			return nil
		}

		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		routes = append(routes, route.ParseFile(path)...)
		return nil
	})
	if err != nil {
		return err
	}

	routeGroups := lo.GroupBy(routes, func(item route.RouteDefinition) string {
		return filepath.Dir(item.Path)
	})

	for path, routes := range routeGroups {
		if err := route.Render(path, routes); err != nil {
			log.WithError(err).WithField("path", path).Error("render route failed")
		}
	}

	return nil
}
