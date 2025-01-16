package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	"git.ipao.vip/rogeecn/atomctl/pkg/utils/gomod"
	"git.ipao.vip/rogeecn/atomctl/templates"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

// CommandNewProvider 注册 new_provider 命令
func CommandNewJob(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "job",
		Short: "创建新的 job",
		Args:  cobra.ExactArgs(1),
		RunE:  commandNewJobE,
	}

	root.AddCommand(cmd)
}

func commandNewJobE(cmd *cobra.Command, args []string) error {
	snakeName := lo.SnakeCase(args[0])
	camelName := lo.PascalCase(args[0])

	destPath := "app/jobs"

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	path, _ = filepath.Abs(path)
	err = gomod.Parse(filepath.Join(path, "go.mod"))
	if err != nil {
		return err
	}

	if err := os.MkdirAll(destPath, os.ModePerm); err != nil {
		return err
	}

	err = fs.WalkDir(templates.Jobs, "jobs", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		destPath := filepath.Join(destPath, snakeName+".go")
		tmpl, err := template.ParseFS(templates.Jobs, path)
		if err != nil {
			return err
		}

		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		return tmpl.Execute(destFile, map[string]string{
			"Name":       camelName,
			"ModuleName": gomod.GetModuleName(),
		})
	})
	if err != nil {
		return err
	}

	fmt.Printf("job 已创建: %s\n", snakeName)
	return nil
}
