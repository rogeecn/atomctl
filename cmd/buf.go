package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CommandBuf(root *cobra.Command) {
    cmd := &cobra.Command{
        Use:   "buf",
        Short: "run buf commands",
        RunE:  commandBufE,
    }

    cmd.Flags().String("dir", ".", "Directory to run buf from")
    cmd.Flags().Bool("dry-run", false, "Preview buf command without executing")

	root.AddCommand(cmd)
}

func commandBufE(cmd *cobra.Command, args []string) error {
    dir := cmd.Flag("dir").Value.String()
    dryRun, _ := cmd.Flags().GetBool("dry-run")
	if _, err := exec.LookPath("buf"); err != nil {
		log.Warn("buf 命令不存在，正在安装 buf...")
		log.Info("go install github.com/bufbuild/buf/cmd/buf@v1.48.0")
		installCmd := exec.Command("go", "install", "github.com/bufbuild/buf/cmd/buf@v1.48.0")
		if err := installCmd.Run(); err != nil {
			return fmt.Errorf("安装 buf 失败: %v", err)
		}
		log.Info("buf 安装成功")

		if _, err := exec.LookPath("buf"); err != nil {
			return fmt.Errorf("buf 命令不存在，请检查 $PATH")
		}
	}

    // preflight: ensure buf.yaml exists
    if _, err := os.Stat(filepath.Join(dir, "buf.yaml")); err != nil {
        log.Warnf("未找到 %s，buf generate 可能失败", filepath.Join(dir, "buf.yaml"))
    }

    log.Info("buf 命令已存在，正在运行 buf generate...")
    log.Info("PROTOBUF GUIDE: https://buf.build/docs/best-practices/style-guide/")
    if dryRun {
        log.Infof("[dry-run] (cd %s && buf generate)", dir)
        return nil
    }
    generateCmd := exec.Command("buf", "generate")
    generateCmd.Dir = dir
	if err := generateCmd.Run(); err != nil {
		return fmt.Errorf("运行 buf generate 失败: %v", err)
	}
	log.Info("buf generate 运行成功")
	return nil
}
