package cmd

import (
    "github.com/pkg/errors"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    apg "go.ipao.vip/atomctl/v2/pkg/postgres"
    "go.ipao.vip/atomctl/v2/pkg/utils/gomod"
    "go.ipao.vip/gen"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func CommandGenModel(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "model",
		Aliases: []string{"m"},
		Short:   "Generate models",
		Long: `根据数据库连接配置生成模型代码，输出到 ./database，并支持基于 ./database/.transform.yaml 的类型/命名转换。

配置：通过 -c/--config 指定配置文件（默认 config.toml），从 [Database] 段读取：
  Username, Password, Database, Host, Port, Schema, SslMode, TimeZone

行为：
- 解析 go.mod 以识别项目模块名
- 连接 PostgreSQL（gorm + postgres driver），校验连通性
- 调用 go.ipao.vip/gen 生成模型与相关文件到 ./database
- 根据 .transform.yaml 应用转换规则

示例：
  atomctl gen -c config.toml model`,
		RunE:    commandGenModelE,
	}

	root.AddCommand(cmd)
}

func commandGenModelE(cmd *cobra.Command, args []string) error {
    if err := gomod.Parse("go.mod"); err != nil {
        return errors.Wrap(err, "parse go.mod")
    }

    cfgFile := cmd.Flag("config").Value.String()
    if cfgFile == "" {
        cfgFile = "config.toml"
    }

    sqlDB, conf, err := apg.GetDB(cfgFile)
    if err != nil {
        return errors.Wrap(err, "load database config")
    }
    defer sqlDB.Close()

    dsn := conf.DSN()
    log.Infof("parsed DSN: %s (schema=%s)", dsn, conf.Schema)

    db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn}))
    if err != nil {
        return errors.Wrapf(err, "open database with dsn: %s", dsn)
    }

	// 默认同包同目录生成到 ./database
	gen.GenerateWithDefault(db, "./database/.transform.yaml")

    return nil
}
