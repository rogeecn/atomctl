package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"{{.ModuleName}}/database/models"

	"go.ipao.vip/atom"
	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/contracts"
	"go.ipao.vip/atom/opt"
	"gorm.io/gorm"
)

//go:embed migrations/*
var MigrationFS embed.FS

func Truncate(ctx context.Context, db *sql.DB, tableName ...string) error {
	for _, name := range tableName {
		sql := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY", name)
		if _, err := db.ExecContext(ctx, sql); err != nil {
			return err
		}
	}
	return nil
}

func WrapLike(v string) string {
	return "%" + v + "%"
}

func WrapLikeLeft(v string) string {
	return "%" + v
}

func WrapLikeRight(v string) string {
	return "%" + v
}

func Provide(...opt.Option) error {
	return container.Container.Provide(func(db *gorm.DB) contracts.Initial {
		models.SetDefault(db)
		return models.Q
	}, atom.GroupInitial)
}

func DefaultProvider() container.ProviderContainer {
	return container.ProviderContainer{
		Provider: Provide,
		Options:  []opt.Option{},
	}
}
