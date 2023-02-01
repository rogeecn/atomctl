package dao

import (
	"atom/providers/config"
	"context"
	"errors"

	"gorm.io/gorm"
)

type {{.PascalName}}Dao interface {
	Release(context.Context, int, string) error
}

type {{.CamelName}}DaoImpl struct {
	db   *gorm.DB
}

func New{{.PascalName}}Dao(db *gorm.DB) {{.PascalName}}Dao {
	return &{{.CamelName}}DaoImpl{db: db}
}

func (c *{{.CamelName}}DaoImpl) Release(ctx context.Context, a int, b string) error {
	return nil
}
