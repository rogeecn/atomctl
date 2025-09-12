package cmd

import "github.com/spf13/cobra"

func CommandSwag(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "swag",
		Short: "Generate swag docs",
		Long:  `Swagger 文档相关命令组：包含 init（生成文档）与 fmt（格式化注释）。`,
	}
	cmds := []func(*cobra.Command){
		CommandSwagInit,
		CommandSwagFmt,
	}

	for _, c := range cmds {
		c(cmd)
	}

	root.AddCommand(cmd)
}
