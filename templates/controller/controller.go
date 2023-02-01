package controller

import (
	"embed"
)

//go:embed tpl/*
var Templates embed.FS

var Files = map[string]string{
	"controller.tpl": "{{.Name}}.go",
}
