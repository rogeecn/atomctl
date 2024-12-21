package main

import (
	"git.ipao.vip/rogeecn/atomctl/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "atomctl",
		Short: "atom framework command line tool",
	}

	cmds := []func(*cobra.Command){
		cmd.CommandInit,
		cmd.CommandMigrate,
		cmd.CommandGen,
		cmd.CommandSwag,
		cmd.CommandFmt,
	}

	for _, c := range cmds {
		c(rootCmd)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
