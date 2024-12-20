package cmd

import (
	"github.com/spf13/cobra"
	"github.com/swaggo/swag/format"
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
