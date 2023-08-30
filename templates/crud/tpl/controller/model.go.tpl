package controller

import (
	"{{ .PkgName }}/common"
	"{{ .PkgName }}/{{ .Module }}/dto"
	"{{ .PkgName }}/{{ .Module }}/service"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

// @provider
type {{ .Model.Name }}Controller struct {
	{{ .Model.CamelName }}Svc *service.{{ .Model.Name }}Service
}

// Filters 
//
//	@Summary		Filters
//	@Tags			{{ .Model.TagName }}
//	@Accept			json
//	@Produce		json
//	@Success		200			{array}	common.Filter
//	@Router			/{{ .Model.RouteName }}/filters [get]
func (c *{{ .Model.Name }}Controller) Filters(ctx *fiber.Ctx) ([]common.Filter, error) {
	return dto.{{ .Model.Name }}ListQueryFilters(), nil
}

// Columns 
//
//	@Summary		columns
//	@Tags			{{ .Model.TagName }}
//	@Accept			json
//	@Produce		json
//	@Success		200			{object}	common.Columns
//	@Router			/{{ .Model.RouteName }}/columns [get]
func (c *{{ .Model.Name }}Controller) Columns(ctx *fiber.Ctx) (common.Columns, error) {
	columns := []common.TableColumnData{
	{{- range .Model.Fields }}
		{Title: "{{ .Comment }}", DataIndex: {{ .Tag }}},
	{{- end }}
		{Title: "操作", DataIndex: "operations", Align: lo.ToPtr("right")},
	}

	return common.NewColumns(columns), nil
}

// LabelShow
//
//	@Summary		LabelShow
//	@Tags			{{ .Model.TagName }}
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"ID"
//	@Success		200	{object}	dto.UserItem
//	@Router			/{{ .Model.RouteName }}/{id}/label [get]
func (c *{{ .Model.Name }}Controller) LabelShow(ctx *fiber.Ctx, id int64) ([]common.LabelItem, error) {
	item, err := c.{{ .Model.CamelName }}Svc.GetByID(ctx.Context(), id)
	if err != nil{
		return nil, err
	}

	return []common.LabelItem{
	{{- range .Model.Fields }}
		{Label: "{{ .Comment }}", Value: item.{{ .Name }}},
	{{- end }}
	}, nil
}

// Show get single item info
//
//	@Summary		get by id
//	@Description	get info by id
//	@Tags			{{ .Model.TagName }}
//	@Accept			json
//	@Produce		json
{{- range $i, $field := .Model.PathFields }}
//	@Param			{{ $field.Name }}	path		{{ $field.Type }}	true	"{{ $field.Comment }}"
{{- end}}
//	@Param			id	path		int	true	"{{ .Model.Name }}ID"
//	@Success		200	{object}	dto.{{ .Model.Name }}Item
//	@Router			/{{ .Model.RouteName }}/{id} [get]
func (c *{{ .Model.Name }}Controller) Show(ctx *fiber.Ctx,{{ range $i, $field := .Model.PathFields }} {{ $field.Name}} {{ $field.Type }}, {{end}} id {{ .Model.IntType }}) (*dto.{{ .Model.Name }}Item, error) {
	item, err := c.{{ .Model.CamelName }}Svc.GetByID(ctx.Context(), id)
	if err != nil{
		return nil, err
	}

	return c.{{ .Model.CamelName }}Svc.DecorateItem(item, 0), nil
}

// List list by query filter
//
//	@Summary		list by query filter
//	@Tags			{{ .Model.TagName }}
//	@Accept			json
//	@Produce		json
{{- range $i, $field := .Model.PathFields }}
//	@Param			{{ $field.Name }}	path		{{ $field.Type }}	true	"{{ $field.Comment }}"
{{- end}}
//	@Param			queryFilter	query		dto.{{ .Model.Name }}ListQueryFilter	true	"{{ .Model.Name }}ListQueryFilter"
//	@Param			pageFilter	query		common.PageQueryFilter	true	"PageQueryFilter"
//	@Param			sortFilter	query		common.SortQueryFilter	true	"SortQueryFilter"
//	@Success		200			{object}	common.PageDataResponse{list=dto.{{ .Model.Name }}Item}
//	@Router			/{{ .Model.RouteName }} [get]
func (c *{{ .Model.Name }}Controller) List(
	ctx *fiber.Ctx, 
{{- range $i, $field := .Model.PathFields }} 
	{{ $field.Name}} {{ $field.Type }}, 
{{- end}}
	queryFilter *dto.{{ .Model.Name }}ListQueryFilter,
	pageFilter *common.PageQueryFilter, 
	sortFilter *common.SortQueryFilter,
) (*common.PageDataResponse, error) {
	items, total, err := c.{{ .Model.CamelName }}Svc.PageByQueryFilter(ctx.Context(), {{ range $i, $field := .Model.PathFields }} {{ $field.Name}},{{ end }}queryFilter, pageFilter, sortFilter)
	if err != nil {
		return nil, err
	}

	return &common.PageDataResponse{
		PageQueryFilter: *pageFilter,
		Total:           total,
		Items:           lo.Map(items, c.{{ .Model.CamelName }}Svc.DecorateItem),
	}, nil
}

// Create a new item
//
//	@Summary		create new item
//	@Tags			{{ .Model.TagName }}
//	@Accept			json
//	@Produce		json
{{- range $i, $field := .Model.PathFields }}
//	@Param			{{ $field.Name }}	path		{{ $field.Type }}	true	"{{ $field.Comment }}"
{{- end}}
//	@Param			body	body		dto.{{ .Model.Name }}Form	true	"{{ .Model.Name }}Form"
//	@Success		200		{string}	{{ .Model.Name }}ID
//	@Router			/{{ .Model.RouteName }} [post]
func (c *{{ .Model.Name }}Controller) Create(ctx *fiber.Ctx,{{ range $i, $field := .Model.PathFields }} {{ $field.Name}} {{ $field.Type }}, {{end}} body *dto.{{ .Model.Name }}Form) error {
	return c.{{ .Model.CamelName }}Svc.Create(ctx.Context(), {{ range $i, $field := .Model.PathFields }} {{ $field.Name}},{{ end }}body)
}

// Update by id
//
//	@Summary		update by id
//	@Tags			{{ .Model.TagName }}
//	@Accept			json
//	@Produce		json
{{- range $i, $field := .Model.PathFields }}
//	@Param			{{ $field.Name }}	path		{{ $field.Type }}	true	"{{ $field.Comment }}"
{{- end}}
//	@Param			id		path		int				true	"{{ .Model.Name }}ID"
//	@Param			body	body		dto.{{ .Model.Name }}Form	true	"{{ .Model.Name }}Form"
//	@Success		200		{string}	{{ .Model.Name }}ID
//	@Router			/{{ .Model.RouteName }}/{id} [put]
func (c *{{ .Model.Name }}Controller) Update(ctx *fiber.Ctx,{{ range $i, $field := .Model.PathFields }} {{ $field.Name}} {{ $field.Type }}, {{end}}id {{ .Model.IntType }}, body *dto.{{ .Model.Name }}Form) error {
	return c.{{ .Model.CamelName }}Svc.Update(ctx.Context(), {{ range $i, $field := .Model.PathFields }} {{ $field.Name}},{{ end }}id, body)
}

// Delete by id
//
//	@Summary		delete by id
//	@Tags			{{ .Model.TagName }}
//	@Accept			json
//	@Produce		json
{{- range $i, $field := .Model.PathFields }}
//	@Param			{{ $field.Name }}	path		{{ $field.Type }}	true	"{{ $field.Comment }}"
{{- end}}
//	@Param			id	path		int	true	"{{ .Model.Name }}ID"
//	@Success		200	{string}	{{ .Model.Name }}ID
//	@Router			/{{ .Model.RouteName }}/{id} [delete]
func (c *{{ .Model.Name }}Controller) Delete(ctx *fiber.Ctx,{{ range $i, $field := .Model.PathFields }} {{ $field.Name}} {{ $field.Type }}, {{end}}id {{ .Model.IntType }}) error {
	return c.{{ .Model.CamelName }}Svc.Delete(ctx.Context(),{{ range $i, $field := .Model.PathFields }} {{ $field.Name}},{{ end }} id)
}
