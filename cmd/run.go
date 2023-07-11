package cmd

import (
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run commands",
}

func init() {
	rootCmd.AddCommand(runCmd)
}
