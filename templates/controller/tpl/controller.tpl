package controller

import (
	"github.com/gin-gonic/gin"
)


type {{.PascalName}}Controller struct {
	conf *config.Config
}

func New{{.PascalName}}Controller() *{{.PascalName}}Controller {
	return &{{.PascalName}}Controller{}
}

func (c *{{.PascalName}}Controller) GetName(ctx *gin.Context) (string, error) {
	return "{{.PascalName}}",nil
}
