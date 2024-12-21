package cmd

import "github.com/spf13/cobra"

func CommandGen(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:                "gen",
		Short:              "Generate code",
		PersistentPostRunE: commandFmtE,
	}
	cmd.PersistentFlags().StringP("config", "c", "config.toml", "database config file")

	cmds := []func(*cobra.Command){
		CommandGenProvider,
		CommandGenRoute,
		CommandGenModel,
		CommandGenEnum,
	}

	for _, c := range cmds {
		c(cmd)
	}

	root.AddCommand(cmd)
}
