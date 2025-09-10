package cmd

import (
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CommandFmt(root *cobra.Command) {
    cmd := &cobra.Command{
        Use:   "fmt",
        Short: "fmt codes",
        RunE:  commandFmtE,
    }

    cmd.Flags().Bool("check", false, "Check formatting without writing changes")
    cmd.Flags().String("path", ".", "Path to format (default .)")

	root.AddCommand(cmd)
}

func commandFmtE(cmd *cobra.Command, args []string) error {
    log.Info("开始格式化代码")
	if _, err := exec.LookPath("gofumpt"); err != nil {
		log.Info("gofumpt 不存在，正在安装...")
		installCmd := exec.Command("go", "install", "mvdan.cc/gofumpt@latest")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			return fmt.Errorf("安装 gofumpt 失败: %v", err)
		}
		if _, err := exec.LookPath("gofumpt"); err != nil {
			return fmt.Errorf("gofumpt 已经安装，但是执行失败")
		}
	}

    check, _ := cmd.Flags().GetBool("check")
    path, _ := cmd.Flags().GetString("path")

    if check {
        log.Info("运行 gofumpt 检查模式...")
        out, err := exec.Command("gofumpt", "-l", "-extra", path).CombinedOutput()
        if err != nil {
            return fmt.Errorf("运行 gofumpt 失败: %v", err)
        }
        if len(out) > 0 {
            fmt.Fprintln(os.Stdout, string(out))
            return fmt.Errorf("发现未格式化文件，请运行: gofumpt -l -extra -w %s", path)
        }
        log.Info("代码格式良好")
        return nil
    }

    log.Info("运行 gofumpt...")
    gofumptCmd := exec.Command("gofumpt", "-l", "-extra", "-w", path)
    gofumptCmd.Stdout = os.Stdout
    gofumptCmd.Stderr = os.Stderr
    if err := gofumptCmd.Run(); err != nil {
        return fmt.Errorf("运行 gofumpt 失败: %v", err)
    }

    log.Info("格式化代码完成")
    return nil
}
