package cmd

import (
	"context"

	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.ipao.vip/atomctl/v2/pkg/postgres"
)

// migrate
func CommandMigrate(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "migrate [up|up-by-one|up-to|create|down|down-to|fix|redo|reset|status|version]",
		Aliases: []string{"m"},
		RunE:    commandMigrate,
		Long: `基于 goose 的数据库迁移工具。

action：
  up|up-by-one|up-to|create|down|down-to|fix|redo|reset|status|version

参数：
  -c, --config   数据库配置文件（默认 config.toml，读取 [Database]）
  --dir          迁移目录（默认 database/migrations）
  --table        迁移版本表名（默认 migrations）

说明：
- 执行 create 时会在缺省情况下追加 sql 类型
- 使用配置连接数据库，并通过 goose.Run 执行对应 action
- 可在代码中通过 goose.SetTableName 自定义版本表名`,
	}
	cmd.Flags().StringP("config", "c", "config.toml", "database config file")
	cmd.Flags().String("dir", "database/migrations", "migrations directory")
	cmd.Flags().String("table", "migrations", "migrations table name")

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

	dir := cmd.Flag("dir").Value.String()
	table := cmd.Flag("table").Value.String()

	goose.SetTableName(table)

	return goose.RunContext(context.Background(), action, db, dir, args...)
}
