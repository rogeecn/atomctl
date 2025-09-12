package cmd

import (
	"github.com/rogeecn/swag/format"
	"github.com/spf13/cobra"
)

func CommandSwagFmt(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "fmt",
		Aliases: []string{"f"},
		Short:   "swag format",
		Long: `格式化接口注释并协助生成更规范的 Swagger 注释。

参数：
- --dir   扫描目录（默认 ./app/http）
- --main  主入口文件（默认 main.go）`,
		RunE:    commandSwagFmtE,
	}

    cmd.Flags().String("dir", "./app/http", "SearchDir for swag format")
    cmd.Flags().String("main", "main.go", "MainFile for swag format")

	root.AddCommand(cmd)
}

func commandSwagFmtE(cmd *cobra.Command, args []string) error {
    dir := cmd.Flag("dir").Value.String()
    main := cmd.Flag("main").Value.String()
    return format.New().Build(&format.Config{
        SearchDir: dir,
        Excludes:  "",
        MainFile:  main,
    })
}
