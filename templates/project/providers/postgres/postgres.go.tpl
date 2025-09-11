package postgres

import (
	"context"
	"time"

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
		dbConfig := postgres.Config{DSN: conf.DSN()}

		// 安全日志：不打印密码，仅输出关键连接信息
		logrus.
			WithFields(
				logrus.Fields{
					"host":   conf.Host,
					"port":   conf.Port,
					"db":     conf.Database,
					"schema": conf.Schema,
					"ssl":    conf.SslMode,
				},
			).
			Info("opening PostgreSQL connection")

		// 映射日志等级
		lvl := conf.GormLogLevel()
		slow := conf.GormSlowThreshold()

		gormConfig := gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   conf.Prefix,
				SingularTable: conf.Singular,
			},
			DisableForeignKeyConstraintWhenMigrating: true,
			PrepareStmt:                              conf.PrepareStmt,
			SkipDefaultTransaction:                   conf.SkipDefaultTransaction,
			Logger: logger.New(logrus.StandardLogger(), logger.Config{
				SlowThreshold:             slow,
				LogLevel:                  lvl,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
				ParameterizedQueries:      conf.ParameterizedQueries,
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
		if conf.ConnMaxLifetimeSeconds > 0 {
			sqlDB.SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetimeSeconds) * time.Second)
		}
		if conf.ConnMaxIdleTimeSeconds > 0 {
			sqlDB.SetConnMaxIdleTime(time.Duration(conf.ConnMaxIdleTimeSeconds) * time.Second)
		}

		// Ping 校验
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := sqlDB.PingContext(ctx); err != nil {
			return nil, nil, err
		}

		// 关闭钩子
		container.AddCloseAble(func() { _ = sqlDB.Close() })

		return db, &conf, nil
	}, o.DiOptions()...)
}
