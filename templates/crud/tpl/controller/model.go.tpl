package controller

import (
	"{{ .PkgName }}/common"
	"{{ .PkgName }}/{{ .Module }}/dto"
	"{{ .PkgName }}/{{ .Module }}/service"

	. "github.com/rogeecn/fen"
)

type {{ .Model.Name }}Controller struct {
	{{ .Model.CamelName }}Svc *service.{{ .Model.Name }}Service
}

func New{{ .Model.Name }}Controller(
	{{ .Model.CamelName }}Svc *service.{{ .Model.Name }}Service,
) *{{ .Model.Name }}Controller {
	return &{{ .Model.Name }}Controller{
		{{ .Model.CamelName }}Svc: {{ .Model.CamelName }}Svc,
	}
}

// Show get single item info
//
//	@Summary		get by id
//	@Description	get info by id
//	@Tags			TODO_ADD_TAGNAME
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"{{ .Model.Name }}ID"
//	@Success		200	{object}	dto.{{ .Model.Name }}Item
//	@Router			/{{ .Model.RouteName }}/{id} [get]
func (c *{{ .Model.Name }}Controller) Show(ctx *Ctx, id {{ .Model.IntType }}) (*dto.{{ .Model.Name }}Item, error) {
	return c.{{ .Model.CamelName }}Svc.GetByID(ctx, id)
}

// List list by query filter
//
//	@Summary		list by query filter
//	@Description	list by query filter
//	@Tags			TODO_ADD_TAGNAME
//	@Accept			json
//	@Produce		json
//	@Param			queryFilter	query		dto.{{ .Model.Name }}ListQueryFilter	true	"{{ .Model.Name }}ListQueryFilter"
//	@Param			pageFilter	query		common.PageQueryFilter	true	"PageQueryFilter"
//	@Param			sortFilter	query		common.SortQueryFilter	true	"SortQueryFilter"
//	@Success		200			{object}	common.PageDataResponse
//	@Router			/{{ .Model.RouteName }} [get]
func (c *{{ .Model.Name }}Controller) List(
	ctx *Ctx, 
	queryFilter *dto.{{ .Model.Name }}ListQueryFilter,
	pageFilter *common.PageQueryFilter, 
	sortFilter *common.SortQueryFilter,
) (*common.PageDataResponse, error) {
	items, total, err := c.{{ .Model.CamelName }}Svc.PageByQueryFilter(ctx, queryFilter, pageFilter, sortFilter)
	if err != nil {
		return nil, err
	}

	return &common.PageDataResponse{
		PageQueryFilter: *pageFilter,
		Total:           total,
		Items:           items,
	}, nil
}

// Create a new item
//
//	@Summary		create new item
//	@Description	create new item
//	@Tags			TODO_ADD_TAGNAME
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.{{ .Model.Name }}Form	true	"{{ .Model.Name }}Form"
//	@Success		200		{string}	{{ .Model.Name }}ID
//	@Router			/{{ .Model.RouteName }} [post]
func (c *{{ .Model.Name }}Controller) Create(ctx *Ctx, body *dto.{{ .Model.Name }}Form) error {
	return c.{{ .Model.CamelName }}Svc.Create(ctx, body)
}

// Update update by id
//
//	@Summary		update by id
//	@Description	update by id
//	@Tags			TODO_ADD_TAGNAME
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int				true	"{{ .Model.Name }}ID"
//	@Param			body	body		dto.{{ .Model.Name }}Form	true	"{{ .Model.Name }}Form"
//	@Success		200		{string}	{{ .Model.Name }}ID
//	@Failure		500		{string}	{{ .Model.Name }}ID
//	@Router			/{{ .Model.RouteName }}/{id} [put]
func (c *{{ .Model.Name }}Controller) Update(ctx *Ctx, id {{ .Model.IntType }}, body *dto.{{ .Model.Name }}Form) error {
	return c.{{ .Model.CamelName }}Svc.Update(ctx, id, body)
}

// Delete delete by id
//
//	@Summary		delete by id
//	@Description	delete by id
//	@Tags			TODO_ADD_TAGNAME
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"{{ .Model.Name }}ID"
//	@Success		200	{string}	{{ .Model.Name }}ID
//	@Failure		500	{string}	{{ .Model.Name }}ID
//	@Router			/{{ .Model.RouteName }}/{id} [delete]
func (c *{{ .Model.Name }}Controller) Delete(ctx *Ctx, id {{ .Model.IntType }}) error {
	return c.{{ .Model.CamelName }}Svc.Delete(ctx, id)
}
