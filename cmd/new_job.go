package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"go.ipao.vip/atomctl/v2/pkg/utils/gomod"
	"go.ipao.vip/atomctl/v2/templates"
)

// CommandNewProvider 注册 new_provider 命令
func CommandNewJob(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "job",
		Short: "创建新的 job",
		Long: `在 app/jobs 下渲染创建任务模板文件。

行为：
- 名称转换：输入名的 Snake 与 Pascal 形式分别用于文件名与导出名
- 输出到 app/jobs/<snake>.go；指定 --cron 时输出到 app/jobs/cron_<snake>.go
- --dry-run 仅打印渲染与写入动作；--dir 指定输出基目录（默认 .）

参数：
- --cron  生成定时任务（渲染 cronjob.go.tpl，输出 cron_<snake>.go）

示例：
  atomctl new job SendDailyReport`,
		Args: cobra.ExactArgs(1),
		RunE: commandNewJobE,
	}

	cmd.Flags().Bool("cron", false, "创建 cron 任务（渲染 cronjob.go.tpl，输出 cron_<name>.go）")
	root.AddCommand(cmd)
}

func commandNewJobE(cmd *cobra.Command, args []string) error {
	snakeName := lo.SnakeCase(args[0])
	camelName := lo.PascalCase(args[0])

	// shared flags
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	baseDir, _ := cmd.Flags().GetString("dir")
	cron, _ := cmd.Flags().GetBool("cron")

	basePath := filepath.Join(baseDir, "app/jobs")

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	path, _ = filepath.Abs(path)
	err = gomod.Parse(filepath.Join(path, "go.mod"))
	if err != nil {
		return err
	}

	if dryRun {
		fmt.Printf("[dry-run] mkdir -p %s\n", basePath)
	} else {
		if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
			return err
		}
	}

	// always render job file
	jobOut := filepath.Join(basePath, snakeName+".go")
	jobTpl, err := template.ParseFS(templates.Jobs, "jobs/job.go.tpl")
	if err != nil {
		return err
	}
	if dryRun {
		fmt.Printf("[dry-run] render > %s\n", jobOut)
	} else {
		fd, err := os.Create(jobOut)
		if err != nil {
			return err
		}
		if err := jobTpl.Execute(fd, map[string]string{
			"Name":       camelName,
			"ModuleName": gomod.GetModuleName(),
		}); err != nil {
			fd.Close()
			return err
		}
		fd.Close()
	}

	// optionally render cron job file when --cron is set
	if cron {
		cronOut := filepath.Join(basePath, "cron_"+snakeName+".go")
		cronTpl, err := template.ParseFS(templates.Jobs, "jobs/cronjob.go.tpl")
		if err != nil {
			return err
		}
		if dryRun {
			fmt.Printf("[dry-run] render > %s\n", cronOut)
		} else {
			fd, err := os.Create(cronOut)
			if err != nil {
				return err
			}
			if err := cronTpl.Execute(fd, map[string]string{
				"Name":       camelName,
				"ModuleName": gomod.GetModuleName(),
			}); err != nil {
				fd.Close()
				return err
			}
			fd.Close()
		}
	}

	fmt.Printf("job 已创建: %s\n", snakeName)
	return nil
}
