package dao

import (
	"context"

	"{{ .PkgName }}/common"
	"{{ .PkgName }}/database/models"
	q "{{ .PkgName }}/database/query"
	"{{ .PkgName }}/{{ .Module }}/dto"

	"github.com/jinzhu/copier"
)

type {{ .Model.Name }}Dao struct {
	query *q.Query
}

func New{{ .Model.Name }}Dao(query *q.Query) *{{ .Model.Name }}Dao {
	return &{{ .Model.Name }}Dao{query: query}
}

func (dao *{{ .Model.Name }}Dao) Update(ctx context.Context, id int32, model *models.{{ .Model.Name }}) error {
	oldModel, err := dao.GetByID(ctx, id)
	if err != nil {
		return err
	}
	_ = copier.Copy(oldModel, model)

	query := dao.query.{{ .Model.Name }}
	_, err = query.WithContext(ctx).Where(query.ID.Eq(id)).Updates(model)

	return err
}

func (dao *{{ .Model.Name }}Dao) Delete(ctx context.Context, id int32) error {
	query := dao.query.{{ .Model.Name }}
	_, err := query.WithContext(ctx).Where(query.ID.Eq(id)).Delete()
	return err
}

func (dao *{{ .Model.Name }}Dao) Create(ctx context.Context, model *models.{{ .Model.Name }}) error {
	return dao.query.{{ .Model.Name }}.WithContext(ctx).Create(model)
}

func (dao *{{ .Model.Name }}Dao) GetByID(ctx context.Context, id int32) (*models.{{ .Model.Name }}, error) {
	query := dao.query.{{ .Model.Name }}
	return query.WithContext(ctx).Where(query.ID.Eq(id)).First()
}

func (dao *{{ .Model.Name }}Dao) PageByQueryFilter(
	ctx context.Context, 
	queryFilter *dto.{{ .Model.Name }}ListQueryFilter,
	pageFilter *common.PageQueryFilter, 
	sortFilter *common.SortQueryFilter,
) ([]*models.{{ .Model.Name }}, int64, error) {
	query := dao.query.{{ .Model.Name }}
	{{ .Model.CamelName }}Query := query.WithContext(ctx)
	if queryFilter != nil {
	{{- range $index, $item := .Model.Fields }}
		if queryFilter.{{ $item.Name }} != nil {
			{{ $.Model.CamelName }}Query = {{ $.Model.CamelName }}Query.Where(query.{{ $item.Name }}.Eq(*queryFilter.{{ $item.Name }}))
		}
	{{- end }}
	}

	if sortFilter != nil {
		orderExprs := []field.Expr{}
		for _, v := range sortFilter.AscField() {
			if expr, ok := query.GetFieldByName(v); ok {
				orderExprs = append(orderExprs, expr)
			}
		}
		for _, v := range sortFilter.DescField() {
			if expr, ok := query.GetFieldByName(v); ok {
				orderExprs = append(orderExprs, expr.Desc())
			}
		}
		{{ .Model.CamelName }}Query = {{ .Model.CamelName }}Query.Order(orderExprs...)
	}

	return {{ .Model.CamelName }}Query.FindByPage(pageFilter.Offset(), pageFilter.Limit)
}


func (dao *{{ .Model.Name }}Dao) FindByQueryFilter(
	ctx context.Context, 
	queryFilter *dto.{{ .Model.Name }}ListQueryFilter,
	sortFilter *common.SortQueryFilter,
) ([]*models.{{ .Model.Name }}, error) {
	query := dao.query.{{ .Model.Name }}
	{{ .Model.CamelName }}Query := query.WithContext(ctx)
	if queryFilter != nil {
	{{- range $index, $item := .Model.Fields }}
		if queryFilter.{{ $item.Name }} != nil {
			{{ $.Model.CamelName }}Query = {{ $.Model.CamelName }}Query.Where(query.{{ $item.Name }}.Eq(*queryFilter.{{ $item.Name }}))
		}
	{{- end }}
	}

	if sortFilter != nil {
		orderExprs := []field.Expr{}
		for _, v := range sortFilter.AscField() {
			if expr, ok := query.GetFieldByName(v); ok {
				orderExprs = append(orderExprs, expr)
			}
		}
		for _, v := range sortFilter.DescField() {
			if expr, ok := query.GetFieldByName(v); ok {
				orderExprs = append(orderExprs, expr.Desc())
			}
		}
		{{ .Model.CamelName }}Query = {{ .Model.CamelName }}Query.Order(orderExprs...)
	}

	return {{ .Model.CamelName }}Query.Find()
}
