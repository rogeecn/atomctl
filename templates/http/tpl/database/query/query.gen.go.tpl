package query

import (
	"gorm.io/gen"
	"gorm.io/gorm"
)

var (
	Q = new(Query)
)

func SetDefault(db *gorm.DB, opts ...gen.DOOption) {
	*Q = *Use(db, opts...)
}

func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db: db,
	}
}

type Query struct {
	db *gorm.DB
}
