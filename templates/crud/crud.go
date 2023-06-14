package crud

import (
	"embed"
)

//go:embed tpl/*
var Templates embed.FS

var Files = map[string]string{
	"dto/model.go.tpl":        "dto/{filename}.go",
	"dao/model.go.tpl":        "dao/{filename}.go",
	"controller/model.go.tpl": "controller/{filename}.go",
	"service/model.go.tpl":    "service/{filename}.go",
}
