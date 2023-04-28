package module

import (
	"embed"
)

//go:embed tpl/*
var Templates embed.FS

var Files = map[string]string{
	"routes/routes.tpl": "routes/routes.go",
	"controller/keep":   "",
	"dto/keep":          "",
	"service/keep":      "",
	"provider.go.tpl":   "provider.go",
}
