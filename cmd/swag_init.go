package cmd

import (
	"os"
	"path/filepath"

	"github.com/rogeecn/swag"
	"github.com/rogeecn/swag/gen"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CommandSwagInit(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "init",
		Short:   "swag init",
		Aliases: []string{"i"},
		Long: `生成 Swagger 文档（go/json/yaml）。

参数：
- --dir   项目根目录（默认 .）
- --out   输出目录（默认 docs）
- --main  主入口文件（默认 main.go）

说明：基于 rogeecn/swag 的 gen 构建器，支持模板分隔符定制、依赖解析等配置。`,
		RunE:    commandSwagInitE,
	}

    cmd.Flags().String("dir", ".", "SearchDir (project root)")
    cmd.Flags().String("out", "docs", "Output dir for generated docs")
    cmd.Flags().String("main", "main.go", "Main API file path")

	root.AddCommand(cmd)
}

func commandSwagInitE(cmd *cobra.Command, args []string) error {
    root := cmd.Flag("dir").Value.String()
    if root == "" {
        var err error
        root, err = os.Getwd()
        if err != nil { return err }
    }

	leftDelim, rightDelim := "{{", "}}"

    outDir := cmd.Flag("out").Value.String()
    mainFile := cmd.Flag("main").Value.String()

    return gen.New().Build(&gen.Config{
        SearchDir:           root,
        Excludes:            "",
        ParseExtension:      "",
        MainAPIFile:         mainFile,
        PropNamingStrategy:  swag.CamelCase,
        OutputDir:           filepath.Join(root, outDir),
        OutputTypes:         []string{"go", "json", "yaml"},
        ParseVendor:         false,
        ParseDependency:     0,
        MarkdownFilesDir:    "",
        ParseInternal:       false,
		Strict:              false,
		GeneratedTime:       false,
		RequiredByDefault:   false,
		CodeExampleFilesDir: "",
		ParseDepth:          100,
		InstanceName:        "",
		OverridesFile:       ".swaggo",
		ParseGoList:         true,
		Tags:                "",
		LeftTemplateDelim:   leftDelim,
		RightTemplateDelim:  rightDelim,
		PackageName:         "",
		Debugger:            log.WithField("module", "swag.init"),
		CollectionFormat:    "csv",
		PackagePrefix:       "",
		State:               "",
		ParseFuncBody:       false,
	})
}
