package cmd

import (
	"github.com/spf13/cobra"
)

func CommandInit(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "new [project|module]",
		Short: "new project/module",
	}

	cmd.PersistentFlags().BoolP("force", "f", false, "Force init project if exists")

	cmds := []func(*cobra.Command){
		CommandNewProject,
		CommandNewModule,
		CommandNewProvider,
		CommandNewEvent,
	}

	for _, c := range cmds {
		c(cmd)
	}

	root.AddCommand(cmd)
}
