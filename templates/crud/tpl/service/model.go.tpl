package service

import (
	"context"

	"{{ .PkgName }}/common"
	"{{ .PkgName }}/database/models"
	"{{ .PkgName }}/{{ .Module }}/dao"
	"{{ .PkgName }}/{{ .Module }}/dto"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
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

func (svc *{{ .Model.Name }}Service) GetByID(ctx context.Context, id {{ .Model.IntType }}) (*dto.{{ .Model.Name }}Item, error) {
	model, err := svc.{{ .Model.CamelName }}Dao.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "get by id failed")
	}

	resp := &dto.{{ .Model.Name }}Item{}
	_ = copier.Copy(resp, model)

	return resp, nil
}

func (svc *{{ .Model.Name }}Service) PageByQueryFilter(
	ctx context.Context, 
	queryFilter *dto.{{ .Model.Name }}ListQueryFilter,
	pageFilter *common.PageQueryFilter, 
	sortFilter *common.SortQueryFilter,
) ([]*dto.{{ .Model.Name }}Item, int64, error) {
	models, total, err := svc.{{ .Model.CamelName }}Dao.PageByQueryFilter(ctx, queryFilter, pageFilter.Format(), sortFilter)
	if err != nil {
		return nil, 0, err
	}

	resp := []*dto.{{ .Model.Name }}Item{}
	for _, u := range models {
		item := &dto.{{ .Model.Name }}Item{}
		_ = copier.Copy(item, u)
		resp = append(resp, item)
	}

	return resp, total, nil
}

// Create
func (svc *{{ .Model.Name }}Service) Create(ctx context.Context, body *dto.{{ .Model.Name }}Form) error {
	model := &models.{{ .Model.Name }}{}
	_ = copier.Copy(model, body)
	return svc.{{ .Model.CamelName }}Dao.Create(ctx, model)
}

// Update
func (svc *{{ .Model.Name }}Service) Update(ctx context.Context, id {{ .Model.IntType }}, body *dto.{{ .Model.Name }}Form) error {
	model := &models.{{ .Model.Name }}{}
	_ = copier.Copy(model, body)
	model.ID = id
	return svc.{{ .Model.CamelName }}Dao.Update(ctx, id, model)
}

// Delete
func (svc *{{ .Model.Name }}Service) Delete(ctx context.Context, id {{ .Model.IntType }}) error {
	return svc.{{ .Model.CamelName }}Dao.Delete(ctx, id)
}
