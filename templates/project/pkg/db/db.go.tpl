package db

import (
	"context"
	"database/sql"
	"fmt"

	"{{.ModuleName}}/pkg/consts"

	"github.com/go-jet/jet/v2/qrm"
)

func FromContext(ctx context.Context, db *sql.DB) qrm.DB {
	if tx, ok := ctx.Value(consts.CtxKeyTx).(*sql.Tx); ok {
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
