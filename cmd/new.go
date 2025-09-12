package cmd

import (
	"github.com/spf13/cobra"
)

func CommandInit(root *cobra.Command) {
    cmd := &cobra.Command{
        Use:   "new [project|module]",
        Short: "new project/module",
        Long: `脚手架命令组：创建项目与常用组件模板。

持久化参数（所有子命令通用）：
- --force, -f  覆盖已存在文件/目录
- --dry-run    仅预览渲染与写入，不落盘
- --dir        指定输出基目录（默认 .）

子命令：project、provider、event、job（module 已弃用）`,
    }

    cmd.PersistentFlags().BoolP("force", "f", false, "Force overwrite existing files or directories")
    cmd.PersistentFlags().Bool("dry-run", false, "Preview actions without writing files")
    cmd.PersistentFlags().String("dir", ".", "Base directory for outputs")

    cmds := []func(*cobra.Command){
        CommandNewProject,
        // deprecate CommandNewModule,
        CommandNewProvider,
        CommandNewEvent,
        CommandNewJob,
    }

	for _, c := range cmds {
		c(cmd)
	}

	root.AddCommand(cmd)
}
