package dao

import (
	"atom/database/query"
	"context"
	"errors"
)

type {{.PascalName}}Dao interface {
	GetByID(context.Context, uint64) error
}

type {{.CamelName}}DaoImpl struct {
	query *query.Query
}

func New{{.PascalName}}Dao(query *query.Query) {{.PascalName}}Dao {
	return &{{.CamelName}}DaoImpl{query:query}
}

func (dao *{{.CamelName}}DaoImpl) GetByID(ctx context.Context, id int) error {
	return nil
}
