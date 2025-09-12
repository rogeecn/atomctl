package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	v := viper.New()
	v.SetConfigType("toml")
	v.SetConfigFile(cfgFile)
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		return errors.Wrap(err, "read config")
	}

	var dbc modelDBConfig
	if err := v.UnmarshalKey("Database", &dbc); err != nil {
		return errors.Wrap(err, "unmarshal Database config")
	}
	dsn := dbc.DSN()

	log.Infof("parsed DSN: %s (schema=%s)", dsn, dbc.Schema)

	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn}))
	if err != nil {
		return errors.Wrapf(err, "open database with dsn: %s", dsn)
	}

	// 默认同包同目录生成到 ./database
	gen.GenerateWithDefault(db, "./database/.transform.yaml")

	return nil
}

// local config and helpers

type modelDBConfig struct {
	Username string
	Password string
	Database string
	Schema   string
	Host     string
	Port     uint
	SslMode  string
	TimeZone string
}

func (c *modelDBConfig) applyDefaults() {
	if c.Username == "" {
		c.Username = "postgres"
	}
	if c.SslMode == "" {
		c.SslMode = "disable"
	}
	if c.TimeZone == "" {
		c.TimeZone = "Asia/Shanghai"
	}
	if c.Port == 0 {
		c.Port = 5432
	}
	if c.Schema == "" {
		c.Schema = "public"
	}
}

func (c *modelDBConfig) DSN() string {
	c.applyDefaults()
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		c.Host,
		c.Username,
		c.Password,
		c.Database,
		c.Port,
		c.SslMode,
		c.TimeZone,
	)
}
