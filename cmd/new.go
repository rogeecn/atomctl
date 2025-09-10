package cmd

import (
	"github.com/spf13/cobra"
)

func CommandInit(root *cobra.Command) {
    cmd := &cobra.Command{
        Use:   "new [project|module]",
        Short: "new project/module",
    }

    cmd.PersistentFlags().BoolP("force", "f", false, "Force overwrite existing files or directories")
    cmd.PersistentFlags().Bool("dry-run", false, "Preview actions without writing files")
    cmd.PersistentFlags().String("dir", ".", "Base directory for outputs")

    cmds := []func(*cobra.Command){
        CommandNewProject,
        // deprecate CommandNewModule,
        CommandNewProvider,
        CommandNewEvent,
        CommandNewJob,
    }

	for _, c := range cmds {
		c(cmd)
	}

	root.AddCommand(cmd)
}
