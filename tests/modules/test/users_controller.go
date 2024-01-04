package test

import (
	"a/common"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

// @provider
type UserController struct {
	userSvc *UserService
}

// Show get single item info
//
//	@Summary		Show
//	@Tags			DEFAULT_TAG_NAME
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"UserID"
//	@Success		200	{object}	UserItem
//	@Router			/users/{id} [get]
func (c *UserController) Show(ctx *fiber.Ctx, id uint64) (*UserItem, error) {
	item, err := c.userSvc.FirstByID(ctx.Context(), id)
	if err != nil {
		return nil, err
	}

	return c.userSvc.DecorateItem(item, 0), nil
}

// List list by query filter
//
//	@Summary		List
//	@Tags			DEFAULT_TAG_NAME
//	@Accept			json
//	@Produce		json
//	@Param			queryFilter	query		UserListQueryFilter	true	"UserListQueryFilter"
//	@Param			pageFilter	query		common.PageQueryFilter	true	"PageQueryFilter"
//	@Param			sortFilter	query		common.SortQueryFilter	true	"SortQueryFilter"
//	@Success		200			{object}	common.PageDataResponse{list=UserItem}
//	@Router			/users [get]
func (c *UserController) List(
	ctx *fiber.Ctx,
	queryFilter *UserListQueryFilter,
	pageFilter *common.PageQueryFilter,
	sortFilter *common.SortQueryFilter,
) (*common.PageDataResponse, error) {
	items, total, err := c.userSvc.PageByFilter(ctx.Context(), queryFilter, pageFilter, sortFilter.DescID())
	if err != nil {
		return nil, err
	}

	return &common.PageDataResponse{
		PageQueryFilter: *pageFilter,
		Total:           total,
		Items:           lo.Map(items, c.userSvc.DecorateItem),
	}, nil
}

// Create a new item
//
//	@Summary		Create
//	@Tags			DEFAULT_TAG_NAME
//	@Accept			json
//	@Produce		json
//	@Param			body	body		UserForm	true	"UserForm"
//	@Success		200		{string}	UserID
//	@Router			/users [post]
func (c *UserController) Create(ctx *fiber.Ctx, body *UserForm) error {
	return c.userSvc.Create(ctx.Context(), body)
}

// Update by id
//
//	@Summary		update by id
//	@Tags			DEFAULT_TAG_NAME
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int				true	"UserID"
//	@Param			body	body		UserForm	true	"UserForm"
//	@Success		200		{string}	UserID
//	@Router			/users/{id} [put]
func (c *UserController) Update(ctx *fiber.Ctx, id uint64, body *UserForm) error {
	return c.userSvc.Update(ctx.Context(), id, body)
}

// Delete by id
//
//	@Summary		Delete
//	@Tags			DEFAULT_TAG_NAME
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"UserID"
//	@Success		200	{string}	UserID
//	@Router			/users/{id} [delete]
func (c *UserController) Delete(ctx *fiber.Ctx, id uint64) error {
	return c.userSvc.Delete(ctx.Context(), id)
}
