package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-jet/jet/v2/qrm"
)

const CtxDB = "__db__tx:"

func FromContext(ctx context.Context, db *sql.DB) qrm.DB {
	if tx, ok := ctx.Value(CtxDB).(*sql.Tx); ok {
		return tx
	}
	return db
}

func TruncateAllTables(ctx context.Context, db *sql.DB, tableName ...string) error {
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
