package v1

import (
	"mime/multipart"

	"{{.ModuleName}}/app/errorx"
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

// Foo
//
//	@Summary		Test
//	@Description	Test
//	@Tags			Test
//	@Accept			json
//	@Produce		json
//
//	@Param			id			path		int											true	"ID"
//	@Param			queryFilter	query		dto.Filter							true	"Filter"
//	@Param			pageFilter	query		request.PageQueryFilter						true	"Pager"
//	@Param			sortFilter	query		request.SortQueryFilter						true	"Sorter"
//	@Success		200			{object}	request.PageDataResponse{list=DataModel}	"成功"
//
//	@Router			/v1/medias/:id [post]
//	@Bind			query query
//	@Bind			header header
//	@Bind			id path
//	@Bind			req body
//	@Bind			file file
//	@Bind			claim local
func (d *demo) Foo(
	ctx fiber.Ctx,
	id int,
	query *FooQuery,
	header *FooHeader,
	claim *jwt.Claims,
	file *multipart.FileHeader,
	req *FooUploadReq,
) error {
	_, err := services.Class.First(ctx)
	if err != nil {
		// 示例：在控制器层自定义错误消息/附加数据
		appErr := errorx.Wrap(err).
			WithMsg("获取班级失败").
			WithData(fiber.Map{"route": "/v1/test"}).
			WithParams("handler", "Test.Hello")
		return appErr
	}
	return nil
}
