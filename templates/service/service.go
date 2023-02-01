package service

import (
	"embed"
)

//go:embed tpl/*
var Templates embed.FS

var Files = map[string]string{
	"service.tpl": "{{.Name}}.go",
}
