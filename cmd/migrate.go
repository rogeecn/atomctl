package cmd

import (
	"context"

	"go.ipao.vip/atomctl/pkg/postgres"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// migrate
func CommandMigrate(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "migrate [up|up-by-one|up-to|create|down|down-to|fix|redo|reset|status|version]",
		Aliases: []string{"m"},
		RunE:    commandMigrate,
	}
	cmd.Flags().StringP("config", "c", "config.toml", "database config file")

	root.AddCommand(cmd)
}

func commandMigrate(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		args = append(args, "up")
	}

	if args[0] == "create" {
		if !lo.Contains(args, "sql") {
			args = append(args, "sql")
		}
	}

	db, _, err := postgres.GetDB(cmd.Flag("config").Value.String())
	if err != nil {
		return errors.Wrap(err, "get db")
	}

	action, args := args[0], args[1:]
	log.Infof("migration action: %s args: %+v", action, args)

	goose.SetTableName("migrations")

	return goose.RunContext(context.Background(), action, db, "database/migrations", args...)
}
