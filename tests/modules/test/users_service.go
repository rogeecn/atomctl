package test

import (
	"context"

	"a/common"
	"a/database/models"
	"a/database/query"

	"github.com/jinzhu/copier"
	"gorm.io/gen/field"
)

// @provider
type UserService struct{}

func (svc *UserService) DecorateItem(model *models.User, id int) *UserItem {
	return &UserItem{
		ID:        model.ID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		Username:  model.Username,
		Password:  model.Password,
		Age:       model.Age,
	}
}

func (svc *UserService) FirstByID(ctx context.Context, id uint64) (*models.User, error) {
	t, q := query.User, query.User.WithContext(ctx)
	return q.Where(t.ID.Eq(id)).First()
}

func (svc *UserService) FindByIDs(ctx context.Context, id []uint64) ([]*models.User, error) {
	t, q := query.User, query.User.WithContext(ctx)
	return q.Where(t.ID.In(id...)).Find()
}

// Create
func (svc *UserService) Create(ctx context.Context, body *UserForm) error {
	model := &models.User{}
	_ = copier.Copy(model, body)

	return svc.CreateFromModel(ctx, model)
}

// CreateFromModel
func (svc *UserService) CreateFromModel(ctx context.Context, model *models.User) error {
	_, q := query.User, query.User.WithContext(ctx)
	return q.Create(model)
}

// Update
func (svc *UserService) Update(ctx context.Context, id uint64, body *UserForm) error {
	model, err := svc.FirstByID(ctx, id)
	if err != nil {
		return err
	}

	_ = copier.Copy(model, body)
	model.ID = id

	return svc.UpdateFromModel(ctx, model)
}

// UpdateFromModel
func (svc *UserService) UpdateFromModel(ctx context.Context, model *models.User) error {
	t, q := query.User, query.User.WithContext(ctx)
	_, err := q.Where(t.ID.Eq(model.ID)).Updates(model)
	return err
}

// Delete
func (svc *UserService) Delete(ctx context.Context, id uint64) error {
	t, q := query.User, query.User.WithContext(ctx)
	_, err := q.Where(t.ID.Eq(id)).Delete()
	return err
}

// DeletePermanently
func (dao *UserService) DeletePermanently(ctx context.Context, id uint64) error {
	t, q := query.User, query.User.WithContext(ctx)
	_, err := q.Unscoped().Where(t.ID.Eq(id)).Delete()
	return err
}

func (dao *UserService) Restore(ctx context.Context, id uint64) error {
	t, q := query.User, query.User.WithContext(ctx).Unscoped()
	_, err := q.Where(t.ID.Eq(id)).UpdateSimple(t.DeletedAt.Null())
	return err
}

func (svc *UserService) FindByFilter(
	ctx context.Context,
	queryFilter *UserListQueryFilter,
	sortFilter *common.SortQueryFilter,
) ([]*models.User, error) {
	_, q := query.User, query.User.WithContext(ctx)

	q = svc.decorateQueryFilter(q, queryFilter)
	q = svc.decorateSortQueryFilter(q, sortFilter)
	return q.Find()
}

func (svc *UserService) PageByFilter(
	ctx context.Context,
	queryFilter *UserListQueryFilter,
	pageFilter *common.PageQueryFilter,
	sortFilter *common.SortQueryFilter,
) ([]*models.User, int64, error) {
	_, q := query.User, query.User.WithContext(ctx)

	q = svc.decorateQueryFilter(q, queryFilter)
	q = svc.decorateSortQueryFilter(q, sortFilter)
	return q.FindByPage(pageFilter.Offset(), pageFilter.Limit)
}

func (svc *UserService) decorateSortQueryFilter(q query.IUserDo, sortFilter *common.SortQueryFilter) query.IUserDo {
	if sortFilter == nil {
		return q
	}
	t := query.User

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

func (dao *UserService) decorateQueryFilter(q query.IUserDo, queryFilter *UserListQueryFilter) query.IUserDo {
	if queryFilter == nil {
		return q
	}

	t := query.User
	if queryFilter.Username != nil {
		q = q.Where(t.Username.Eq(*queryFilter.Username))
	}
	if queryFilter.Password != nil {
		q = q.Where(t.Password.Eq(*queryFilter.Password))
	}
	if queryFilter.Age != nil {
		q = q.Where(t.Age.Eq(*queryFilter.Age))
	}

	return q
}
