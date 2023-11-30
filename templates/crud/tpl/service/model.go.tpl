package service

import (
	"context"

	"{{ .PkgName }}/common"
	"{{ .PkgName }}/database/query"
	"{{ .PkgName }}/database/models"
	"{{ .PkgName }}/{{ .Module }}/dto"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/copier"
	"gorm.io/gen/field"
)

// @provider
type {{ .Model.Name }}Service struct { }

func (svc *{{ .Model.Name }}Service) DecorateItem(model *models.{{ .Model.Name }}, id int) *dto.{{ .Model.Name }}Item {
	return &dto.{{ .Model.Name }}Item{
	{{- range .Model.Fields }}
	{{- if eq .Name "DeletedAt" }}
	{{- else }}
		{{ .Name }}: model.{{ .Name }},
	{{- end }}
	{{- end }}
	}
}

func (svc *{{ .Model.Name }}Service) FirstByID(ctx context.Context, id uint64) (*models.{{ .Model.Name }}, error) {
	t, q := query.{{ .Model.Name }}, query.{{ .Model.Name }}.WithContext(ctx)
	return q.Where(t.ID.Eq(id)).First()
}

func (svc *{{ .Model.Name }}Service) FindByIDs(ctx context.Context, id []uint64) ([]*models.{{ .Model.Name }}, error) {
	t, q := query.{{ .Model.Name }}, query.{{ .Model.Name }}.WithContext(ctx)
	return q.Where(t.ID.In(id...)).Find()
}

// Create
func (svc *{{ .Model.Name }}Service) Create(ctx context.Context, body *dto.{{ .Model.Name }}Form) error {
	model := &models.{{ .Model.Name }}{}
	_ = copier.Copy(model, body)

	return svc.CreateFromModel(ctx, model)
}

// CreateFromModel
func (svc *{{ .Model.Name }}Service) CreateFromModel(ctx context.Context, model *models.{{ .Model.Name }}) error {
	_, q := query.{{ .Model.Name }}, query.{{ .Model.Name }}.WithContext(ctx)
	return q.Create(model)
}

// Update
func (svc *{{ .Model.Name }}Service) Update(ctx context.Context, id uint64, body *dto.{{ .Model.Name }}Form) error {
	model, err := svc.FirstByID(ctx, id)
	if err != nil {
		return err
	}

	_ = copier.Copy(model, body)
	model.ID = id

	return svc.UpdateFromModel(ctx, model)
}

// UpdateFromModel
func (svc *{{ .Model.Name }}Service) UpdateFromModel(ctx context.Context, model *models.{{ .Model.Name }}) error {
	t, q := query.{{ .Model.Name }}, query.{{ .Model.Name }}.WithContext(ctx)
	_, err := q.Where(t.ID.Eq(model.ID)).Updates(model)
	return err
}

// Delete
func (svc *{{ .Model.Name }}Service) Delete(ctx context.Context, id uint64) error {
	t, q := query.{{ .Model.Name }}, query.{{ .Model.Name }}.WithContext(ctx)
	_, err := q.Where(t.ID.Eq(id)).Delete()
	return err
}

// DeletePermanently
func (dao *{{ .Model.Name }}Service) DeletePermanently(ctx context.Context, id uint64) error {
	t, q := query.{{ .Model.Name }}, query.{{ .Model.Name }}.WithContext(ctx)
	_, err := q.Unscoped().Where(t.ID.Eq(id)).Delete()
	return err
}


{{- range .Model.Fields }}
{{ if eq .Name "DeletedAt" }}
func (dao *{{ $.Model.Name }}Service) Restore(ctx context.Context, id uint64) error {
	t, q := query.{{ $.Model.Name }}, query.{{ $.Model.Name }}.WithContext(ctx)
	_, err := q.Unscoped().Where(t.ID.Eq(id)).UpdateSimple(t.DeletedAt.Null())
	return err
}
{{- end }}
{{- end }}

func (svc *{{ .Model.Name }}Service) FindByFilter(
	ctx context.Context,
	queryFilter *dto.{{ .Model.Name }}ListQueryFilter,
	sortFilter *common.SortQueryFilter,
) ([]*models.{{ .Model.Name }}, error) {
	_, q := query.{{ .Model.Name }}, query.{{ .Model.Name }}.WithContext(ctx)

	q = svc.decorateQueryFilter(q, queryFilter)
	q = svc.decorateSortQueryFilter(q, sortFilter)
	return q.Find()
}

func (svc *{{ .Model.Name }}Service) PageByFilter(
	ctx context.Context,
	queryFilter *dto.{{ .Model.Name }}ListQueryFilter,
	pageFilter *common.PageQueryFilter,
	sortFilter *common.SortQueryFilter,
) ([]*models.{{ .Model.Name }}, int64, error) {
	_, q := query.{{ .Model.Name }}, query.{{ .Model.Name }}.WithContext(ctx)

	q = svc.decorateQueryFilter(q, queryFilter)
	q = svc.decorateSortQueryFilter(q, sortFilter)
	return q.FindByPage(pageFilter.Offset(), pageFilter.Limit)
}

func (svc *{{ .Model.Name }}Service) decorateSortQueryFilter(q query.I{{ .Model.Name }}Do, sortFilter *common.SortQueryFilter) query.I{{ .Model.Name }}Do {
	if sortFilter == nil {
		return q
	}
	t := query.{{ .Model.Name }}

	orderExprs := []field.Expr{}
	for _, v := range sortFilter.AscFields() {
		if expr, ok := t.GetFieldByName(v); ok {
			orderExprs = append(orderExprs, expr)
		}
	}
	for _, v := range sortFilter.DescFields() {
		if expr, ok := t.GetFieldByName(v); ok {
			orderExprs = append(orderExprs, expr.Desc())
		}
	}
	return q.Order(orderExprs...)
}

func (dao *{{ .Model.Name }}Service) decorateQueryFilter(q query.I{{ .Model.Name }}Do, queryFilter *dto.{{ .Model.Name }}ListQueryFilter) query.I{{ .Model.Name }}Do {
	if queryFilter == nil {
		return q
	}

	t := query.{{ .Model.Name }}

	{{- range $index, $item := .Model.Fields }}
	{{- if or (eq .Name "ID") (eq .Name "CreatedAt") (eq .Name "DeletedAt") (eq .Name "UpdatedAt")}}
		{{- else }}
	if queryFilter.{{ $item.Name }} != nil {
		{{- if eq $item.Type "*bool" }}
		q = q.Where(t.Is(*queryFilter.{{ $item.Name }}))
		{{- else }}
		q = q.Where(t.{{ $item.Name }}.Eq(*queryFilter.{{ $item.Name }}))
		{{- end }}
	}
	{{- end }}
	{{- end }}

	return q
}
