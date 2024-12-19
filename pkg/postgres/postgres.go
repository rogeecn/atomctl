package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func GetDB(cfgFile string) (*sql.DB, *Config, error) {
	var conf Config
	v := viper.New()
	v.SetConfigType("toml")
	v.SetConfigFile(cfgFile)
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return nil, &conf, errors.Wrap(err, "read config file")
	}

	if err := v.UnmarshalKey(DefaultPrefix, &conf); err != nil {
		return nil, &conf, errors.Wrap(err, "unmarshal config")
	}

	log.Debugf("connect postgres with dsn: '%s'", conf.DSN())
	db, err := sql.Open("postgres", conf.DSN())
	if err != nil {
		return nil, &conf, errors.Wrap(err, "connect database")
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, &conf, errors.Wrap(err, "ping database")
	}

	return db, &conf, err
}
