package suite

import (
	"embed"
)

//go:embed tpl/*
var Templates embed.FS

var Files = map[string]string{
	"suite.tpl": "{{.Name}}_test.go",
}
