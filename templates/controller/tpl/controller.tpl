package controller

import (
	"github.com/rogeecn/atom/providers/config"

	"github.com/gin-gonic/gin"
)


type {{.PascalName}}Controller struct {
	conf *config.Config
}

func New{{.PascalName}}Controller(conf *config.Config) *{{.PascalName}}Controller {
	return &{{.PascalName}}ControllerImpl{conf: conf}
}

func (c *{{.PascalName}}Controller) GetName(ctx *gin.Context) (string, error) {
	return "{{.PascalName}}",nil
}
