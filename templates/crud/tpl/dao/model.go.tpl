package dao

import (
	"context"

	"{{ .PkgName }}/pkg/common"
	"{{ .PkgName }}/database/models"
	"{{ .PkgName }}/database/query"
	"{{ .PkgName }}/{{ .Module }}/dto"

	"gorm.io/gen/field"
)

type {{ .Model.Name }}Dao struct {
	query *query.Query
}

func New{{ .Model.Name }}Dao(query *query.Query) *{{ .Model.Name }}Dao {
	return &{{ .Model.Name }}Dao{query: query}
}

func (dao *{{ .Model.Name }}Dao) Context(ctx context.Context) query.I{{ .Model.Name }}Do {
	return dao.query.{{ .Model.Name }}.WithContext(ctx)
}

func (dao *{{ .Model.Name }}Dao) decorateSortQueryFilter(query query.I{{ .Model.Name }}Do, sortFilter *common.SortQueryFilter) query.I{{ .Model.Name }}Do {
	if sortFilter == nil {
		return query
	}

	orderExprs := []field.Expr{}
	for _, v := range sortFilter.AscFields() {
		if expr, ok := dao.query.{{ .Model.Name }}.GetFieldByName(v); ok {
			orderExprs = append(orderExprs, expr)
		}
	}
	for _, v := range sortFilter.DescFields() {
		if expr, ok := dao.query.{{ .Model.Name }}.GetFieldByName(v); ok {
			orderExprs = append(orderExprs, expr.Desc())
		}
	}
	return query.Order(orderExprs...)
}

func (dao *{{ .Model.Name }}Dao) decorateQueryFilter(query query.I{{ .Model.Name }}Do, queryFilter *dto.{{ .Model.Name }}ListQueryFilter) query.I{{ .Model.Name }}Do {
	if queryFilter == nil {
		return query
	}

	{{- range $index, $item := .Model.Fields }}
	if queryFilter.{{ $item.Name }} != nil {
		{{- if eq $item.Type "*bool" }}
		query = query.Where(dao.query.{{ $.Model.Name }}.{{ $item.Name }}.Is(*queryFilter.{{ $item.Name }}))
		{{- else }}
		query = query.Where(dao.query.{{ $.Model.Name }}.{{ $item.Name }}.Eq(*queryFilter.{{ $item.Name }}))
		{{- end }}
	}
	{{- end }}

	return query
}

func (dao *{{ .Model.Name }}Dao) UpdateColumn(ctx context.Context, id int32, field field.Expr, value interface{}) error {
	_, err := dao.Context(ctx).Where(dao.query.{{ .Model.Name }}.ID.Eq(id)).Update(field, value)
	return err
}

func (dao *{{ .Model.Name }}Dao) Update(ctx context.Context, model *models.{{ .Model.Name }}) error {
	_, err := dao.Context(ctx).Where(dao.query.{{ .Model.Name }}.ID.Eq(model.ID)).Updates(model)
	return err
}

func (dao *{{ .Model.Name }}Dao) Delete(ctx context.Context, id {{ .Model.IntType }}) error {
	_, err := dao.Context(ctx).Where(dao.query.{{ .Model.Name }}.ID.Eq(id)).Delete()
	return err
}

func (dao *{{ .Model.Name }}Dao) Create(ctx context.Context, model *models.{{ .Model.Name }}) error {
	return dao.Context(ctx).Create(model)
}

func (dao *{{ .Model.Name }}Dao) GetByID(ctx context.Context, id {{ .Model.IntType }}) (*models.{{ .Model.Name }}, error) {
	return dao.Context(ctx).Where(dao.query.{{ .Model.Name }}.ID.Eq(id)).First()
}

func (dao *{{ .Model.Name }}Dao) GetByIDs(ctx context.Context, ids []{{ .Model.IntType }}) ([]*models.{{ .Model.Name }}, error) {
	return dao.Context(ctx).Where(dao.query.{{ .Model.Name }}.ID.In(ids...)).Find()
}

func (dao *{{ .Model.Name }}Dao) PageByQueryFilter(
	ctx context.Context, 
	queryFilter *dto.{{ .Model.Name }}ListQueryFilter,
	pageFilter *common.PageQueryFilter, 
	sortFilter *common.SortQueryFilter,
) ([]*models.{{ .Model.Name }}, int64, error) {
	query := dao.query.{{ .Model.Name }}
	{{ .Model.CamelName }}Query := query.WithContext(ctx)
	{{ .Model.CamelName }}Query = dao.decorateQueryFilter({{ .Model.CamelName }}Query, queryFilter)
	{{ .Model.CamelName }}Query = dao.decorateSortQueryFilter({{ .Model.CamelName }}Query, sortFilter)
	return {{ .Model.CamelName }}Query.FindByPage(pageFilter.Offset(), pageFilter.Limit)
}


func (dao *{{ .Model.Name }}Dao) FindByQueryFilter(
	ctx context.Context, 
	queryFilter *dto.{{ .Model.Name }}ListQueryFilter,
	sortFilter *common.SortQueryFilter,
) ([]*models.{{ .Model.Name }}, error) {
	query := dao.query.{{ .Model.Name }}
	{{ .Model.CamelName }}Query := query.WithContext(ctx)
	{{ .Model.CamelName }}Query = dao.decorateQueryFilter({{ .Model.CamelName }}Query, queryFilter)
	{{ .Model.CamelName }}Query = dao.decorateSortQueryFilter({{ .Model.CamelName }}Query, sortFilter)
	return {{ .Model.CamelName }}Query.Find()
}
