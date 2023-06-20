package controller

import (
	. "github.com/rogeecn/fen"
)


type {{.PascalName}}Controller struct {
	conf *config.Config
}

func New{{.PascalName}}Controller() *{{.PascalName}}Controller {
	return &{{.PascalName}}Controller{}
}

func (c *{{.PascalName}}Controller) GetName(ctx *Ctx) (string, error) {
	return "{{.PascalName}}",nil
}
