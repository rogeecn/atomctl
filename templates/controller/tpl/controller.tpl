package controller

import (
	"atom/providers/config"

	"github.com/gin-gonic/gin"
)

type {{.PascalName}}Controller interface {
	GetName(*gin.Context) (string, error)
}

type {{.CamelName}}ControllerImpl struct {
	conf *config.Config
}

func New{{.PascalName}}Controller(conf *config.Config) {{.PascalName}}Controller {
	return &{{.CamelName}}ControllerImpl{conf: conf}
}

func (c *{{.CamelName}}ControllerImpl) GetName(ctx *gin.Context) (string, error) {
	return "{{.PascalName}}",nil
}
