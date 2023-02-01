package controller

import (
	"atom/providers/config"

	"github.com/gin-gonic/gin"
)

type {{.PascalName}}Controller interface {
	GetName(*gin.Context) (string, error)
}

type {{.PascalName}}ControllerImpl struct {
	conf *config.Config
}

func New{{.PascalName}}Controller(Conf *config.Config) {{.PascalName}}Controller {
	return &{{.PascalName}}ControllerImpl{conf: Conf}
}

func (c *{{.PascalName}}ControllerImpl) GetName(ctx *gin.Context) (string, error) {
	return "{{.PascalName}}",nil
}
