package module

import (
	"embed"
)

//go:embed tpl/*
var Templates embed.FS

var Files = map[string]string{
	"container/container.tpl": "container/container.go",
	"routes/routes.tpl":       "routes/routes.go",
	"controller/keep":         "",
	"dao/keep":                "",
	"dto/keep":                "",
	"service/keep":            "",
}
