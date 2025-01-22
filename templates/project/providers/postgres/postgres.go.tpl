package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"{{.ModuleName}}/pkg/atom/container"
	"{{.ModuleName}}/pkg/atom/opt"
)

func Provide(opts ...opt.Option) error {
	o := opt.New(opts...)
	var conf Config
	if err := o.UnmarshalConfig(&conf); err != nil {
		return err
	}

	return container.Container.Provide(func() (*sql.DB, *Config, error) {
		log.Debugf("connect postgres with dsn: '%s'", conf.DSN())
		db, err := sql.Open("postgres", conf.DSN())
		if err != nil {
			return nil, nil, errors.Wrap(err, "connect database")
		}

		if err := db.Ping(); err != nil {
			db.Close()
			return nil, nil, errors.Wrap(err, "ping database")
		}

		return db, &conf, err
	}, o.DiOptions()...)
}
