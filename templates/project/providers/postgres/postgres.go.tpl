package postgres

import (
	"github.com/sirupsen/logrus"
	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/opt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func Provide(opts ...opt.Option) error {
	o := opt.New(opts...)
	var conf Config
	if err := o.UnmarshalConfig(&conf); err != nil {
		return err
	}

	return container.Container.Provide(func() (*gorm.DB, *Config, error) {
		dbConfig := postgres.Config{
			DSN: conf.DSN(), // DSN data source name
		}
		logrus.Info("Open PostgreSQL:", conf.DSN())

		gormConfig := gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   conf.Prefix,
				SingularTable: conf.Singular,
			},
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger: logger.New(logrus.StandardLogger(), logger.Config{
				SlowThreshold:             200, // 慢 SQL 阈值
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true, // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  false,
			}),
		}

		db, err := gorm.Open(postgres.New(dbConfig), &gormConfig)
		if err != nil {
			return nil, nil, err
		}

		sqlDB, err := db.DB()
		if err != nil {
			return nil, nil, err
		}
		sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
		sqlDB.SetMaxOpenConns(conf.MaxOpenConns)

		return db, &conf, err
	}, o.DiOptions()...)
}
