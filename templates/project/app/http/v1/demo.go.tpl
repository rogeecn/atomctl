package v1

import (
	"mime/multipart"

	"{{.ModuleName}}/app/errorx"
	"{{.ModuleName}}/app/requests"
	"{{.ModuleName}}/app/services"
	"{{.ModuleName}}/providers/jwt"

	"github.com/gofiber/fiber/v3"
)

// @provider
type demo struct{}

type FooUploadReq struct {
	Folder string `json:"folder" form:"folder"` // 上传到指定文件夹
}

type FooQuery struct {
	Search string `query:"search"` // 搜索关键词
}

type FooHeader struct {
	ContentType string `header:"Content-Type"` // 内容类型
}
type Filter struct {
	Name string `query:"name"` // 名称
	Age  int    `query:"age"`  // 年龄
}

type ResponseItem struct{}

// Foo
//
//	@Summary		Test
//	@Description	Test
//	@Tags			Test
//	@Accept			json
//	@Produce		json
//
//	@Param			id		path		int									true	"ID"
//	@Param			query	query		Filter								true	"Filter"
//	@Param			pager	query		requests.Pagination					true	"Pager"
//	@Success		200		{object}	requests.Pager{list=ResponseItem}	"成功"
//
//	@Router			/v1/medias/:id [post]
//	@Bind			query query
//	@Bind			pager query
//	@Bind			header header
//	@Bind			id path
//	@Bind			req body
//	@Bind			file file
//	@Bind			claim local
func (d *demo) Foo(
	ctx fiber.Ctx,
	id int,
	pager *requests.Pagination,
	query *FooQuery,
	header *FooHeader,
	claim *jwt.Claims,
	file *multipart.FileHeader,
	req *FooUploadReq,
) error {
	_, err := services.Test.Test(ctx)
	if err != nil {
		// 示例：在控制器层自定义错误消息/附加数据
		appErr := errorx.Wrap(err).
			WithMsg("获取测试失败").
			WithData(fiber.Map{"route": "/v1/test"}).
			WithParams("handler", "Test.Hello")
		return appErr
	}
	return nil
}
