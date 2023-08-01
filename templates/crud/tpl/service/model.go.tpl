package service

import (
	"context"

	"{{ .PkgName }}/common"
	"{{ .PkgName }}/database/models"
	"{{ .PkgName }}/{{ .Module }}/dao"
	"{{ .PkgName }}/{{ .Module }}/dto"

	"github.com/jinzhu/copier"
)

type {{ .Model.Name }}Service struct {
	{{ .Model.CamelName }}Dao *dao.{{ .Model.Name }}Dao
}

func New{{ .Model.Name }}Service(
	{{ .Model.CamelName }}Dao *dao.{{ .Model.Name }}Dao,
) *{{ .Model.Name }}Service {
	return &{{ .Model.Name }}Service{
		{{ .Model.CamelName }}Dao: {{ .Model.CamelName }}Dao,
	}
}

func (svc *{{ .Model.Name }}Service) DecorateItem(model *models.{{ .Model.Name }}, id int) *dto.{{ .Model.Name }}Item {
	var dtoItem *dto.{{ .Model.Name }}Item
	_ = copier.Copy(dtoItem, model)

	return dtoItem
}

func (svc *{{ .Model.Name }}Service) GetByID(ctx context.Context, id {{ .Model.IntType }}) (*models.{{ .Model.Name }}, error) {
	return svc.{{ .Model.CamelName }}Dao.GetByID(ctx, id)
}

func (svc *{{ .Model.Name }}Service) FindByQueryFilter(
	ctx context.Context, 
{{- range $i, $field := .Model.PathFields }} 
	{{ $field.Name}} {{ $field.Type }}, 
{{- end}}
	queryFilter *dto.{{ .Model.Name }}ListQueryFilter,
	sortFilter *common.SortQueryFilter,
) ([]*models.{{ .Model.Name }}, error) {
	return svc.{{ .Model.CamelName }}Dao.FindByQueryFilter(ctx, queryFilter, sortFilter)
}

func (svc *{{ .Model.Name }}Service) PageByQueryFilter(
	ctx context.Context, 
{{- range $i, $field := .Model.PathFields }} 
	{{ $field.Name}} {{ $field.Type }}, 
{{- end}}
	queryFilter *dto.{{ .Model.Name }}ListQueryFilter,
	pageFilter *common.PageQueryFilter, 
	sortFilter *common.SortQueryFilter,
) ([]*models.{{ .Model.Name }}, int64, error) {
	return svc.{{ .Model.CamelName }}Dao.PageByQueryFilter(ctx, queryFilter, pageFilter.Format(), sortFilter)
}

// CreateFromModel
func (svc *{{ .Model.Name }}Service) CreateFromModel(ctx context.Context, model *models.{{ .Model.Name }}) error {
	return svc.{{ .Model.CamelName }}Dao.Create(ctx, model)
}

// Create
func (svc *{{ .Model.Name }}Service) Create(ctx context.Context,{{ range $i, $field := .Model.PathFields }} {{ $field.Name}} {{ $field.Type }}, {{end}} body *dto.{{ .Model.Name }}Form) error {
	model := &models.{{ .Model.Name }}{}
	_ = copier.Copy(model, body)
	return svc.{{ .Model.CamelName }}Dao.Create(ctx, model)
}

// Update
func (svc *{{ .Model.Name }}Service) Update(ctx context.Context,{{ range $i, $field := .Model.PathFields }} {{ $field.Name}} {{ $field.Type }}, {{end}} id {{ .Model.IntType }}, body *dto.{{ .Model.Name }}Form) error {
	model, err := svc.GetByID(ctx, id)
	if err != nil {
		return err
	}

	_ = copier.Copy(model, body)
	model.ID = id
	return svc.{{ .Model.CamelName }}Dao.Update(ctx, model)
}

// UpdateFromModel
func (svc *{{ .Model.Name }}Service) UpdateFromModel(ctx context.Context, model *models.{{ .Model.Name }}) error {
	return svc.{{ .Model.CamelName }}Dao.Update(ctx, model)
}

// Delete
func (svc *{{ .Model.Name }}Service) Delete(ctx context.Context,{{ range $i, $field := .Model.PathFields }} {{ $field.Name}} {{ $field.Type }}, {{end}} id {{ .Model.IntType }}) error {
	return svc.{{ .Model.CamelName }}Dao.Delete(ctx, id)
}
