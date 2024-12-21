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

	log.Info("运行 gofumpt...")
	gofumptCmd := exec.Command("gofumpt", "-l", "-extra", "-w", ".")
	gofumptCmd.Stdout = os.Stdout
	gofumptCmd.Stderr = os.Stderr
	if err := gofumptCmd.Run(); err != nil {
		return fmt.Errorf("运行 gofumpt 失败: %v", err)
	}

	log.Info("格式化代码完成")
	return nil
}
