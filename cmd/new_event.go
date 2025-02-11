package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	"go.ipao.vip/atomctl/pkg/utils/gomod"
	"go.ipao.vip/atomctl/templates"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

// CommandNewProvider 注册 new_provider 命令
func CommandNewEvent(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "event",
		Aliases: []string{"e"},
		Short:   "创建新的 event publish & subscriber",
		Args:    cobra.ExactArgs(1),
		RunE:    commandNewEventE,
	}

	root.AddCommand(cmd)
}

func commandNewEventE(cmd *cobra.Command, args []string) error {
	snakeName := lo.SnakeCase(args[0])
	camelName := lo.PascalCase(args[0])

	publisherPath := "app/events/publishers"
	subscriberPath := "app/events/subscribers"

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	path, _ = filepath.Abs(path)
	err = gomod.Parse(filepath.Join(path, "go.mod"))
	if err != nil {
		return err
	}

	if err := os.MkdirAll(publisherPath, os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll(subscriberPath, os.ModePerm); err != nil {
		return err
	}

	err = fs.WalkDir(templates.Events, "events", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel("events", path)
		if err != nil {
			return err
		}

		var destPath string
		if relPath == "publisher.go.tpl" {
			destPath = filepath.Join(publisherPath, snakeName+".go")
		} else if relPath == "subscriber.go.tpl" {
			destPath = filepath.Join(subscriberPath, snakeName+".go")
		}

		tmpl, err := template.ParseFS(templates.Events, path)
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

	topicStr := fmt.Sprintf("const Topic%s = %q\n", camelName, snakeName)
	// 写入到 app/events/topic.go
	topicFile, err := os.OpenFile("app/events/topics.go", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer topicFile.Close()

	_, err = topicFile.WriteString(topicStr)
	if err != nil {
		return err
	}

	fmt.Printf("event 已创建: %s\n", snakeName)
	return nil
}
