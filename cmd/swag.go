package cmd

import "github.com/spf13/cobra"

func CommandSwag(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "swag",
		Short: "Generate swag docs",
	}
	cmds := []func(*cobra.Command){
		CommandSwagInit,
		CommandSwagFmt,
	}

	for _, c := range cmds {
		c(cmd)
	}

	root.AddCommand(cmd)
}
