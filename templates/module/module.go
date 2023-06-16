package module

import (
	"embed"
)

//go:embed tpl/*
var Templates embed.FS

var Files = map[string]string{
	"routes/routes.tpl":          "routes/routes.go",
	"controller/provider.go.tpl": "controller/provider.go",
	"dao/provider.go.tpl":        "dao/provider.go",
	"dto/keep":                   "",
	"service/provider.go.tpl":    "service/provider.go",
	"provider.go.tpl":            "provider.go",
}
