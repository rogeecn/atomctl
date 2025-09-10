package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
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
