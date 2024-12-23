package cmd

import (
	"git.ipao.vip/rogeecn/atomctl/pkg/swag/format"
	"github.com/spf13/cobra"
)

func CommandSwagFmt(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "fmt",
		Aliases: []string{"f"},
		Short:   "swag format",
		RunE:    commandSwagFmtE,
	}

	root.AddCommand(cmd)
}

func commandSwagFmtE(cmd *cobra.Command, args []string) error {
	return format.New().Build(&format.Config{
		SearchDir: "./app/http",
		Excludes:  "",
		MainFile:  "main.go",
	})
}
