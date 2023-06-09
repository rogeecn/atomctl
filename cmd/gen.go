package cmd

import (
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "A brief description of your command",
}

func init() {
	rootCmd.AddCommand(genCmd)
}
