package dao

import (
	"github.com/rogeecn/atom/database/query"
	"context"
)

type {{.PascalName}}Dao struct {
	query *query.Query
}

func New{{.PascalName}}Dao(query *query.Query) *{{.PascalName}}Dao {
	return &{{.PascalName}}Dao{query:query}
}

func (dao *{{.PascalName}}Dao) GetByID(ctx context.Context, id int) error {
	return nil
}
