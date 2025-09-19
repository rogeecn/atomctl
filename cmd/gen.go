package cmd

import "github.com/spf13/cobra"

func CommandGen(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate code",
		Long: `代码生成命令组：包含 route、provider、model、enum、service 等。

持久化参数：
- -c, --config  数据库配置文件（默认 config.toml），供 gen model 使用

说明：
- 子命令执行完成后会自动运行 atomctl fmt 进行格式化`,
		PersistentPostRunE: commandFmtE,
	}
	cmd.PersistentFlags().StringP("config", "c", "config.toml", "database config file")

	cmds := []func(*cobra.Command){
		CommandGenProvider,
		CommandGenRoute,
		CommandGenModel,
		CommandGenEnum,
		CommandGenService,
	}

	for _, c := range cmds {
		c(cmd)
	}

	root.AddCommand(cmd)
}
