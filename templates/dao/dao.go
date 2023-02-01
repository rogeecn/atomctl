package dao

import (
	"embed"
)

//go:embed tpl/*
var Templates embed.FS

var Files = map[string]string{
	"dao.tpl": "{{.Name}}.go",
}
