package tests

import (
	"context"
)

type RouteTest struct {
	ctx  context.Context
	Name string
	Age  *int
}

// Create a new item
//
//	@Summary		create new item
//	@Description	create new item
//	@Tags			TODO_ADD_TAGNAME
//	@Accept			json
//	@Produce		json
//	@Param			a	path		string	true	"A"
//	@Param			accountId	path		int	true	"AccountId"
//	@Param			body	body		dto.KubeImageForm	true	"KubeImageForm"
//	@Success		200		{string}	KubeImageID
//	@Router			/clusters/{a}/b/{accountId} [post]
func (s *RouteTest) Hello(ctx context.Context, a, b int, c *string) (*RouteTest, error) {
	return &RouteTest{
		ctx:  ctx,
		Name: "hello",
		Age:  &a,
	}, nil
}
