package seeder

import (
	"embed"
)

//go:embed tpl/*
var Templates embed.FS

var Files = map[string]string{
	"seeder.tpl": "./database/seeders/{{.SnakeSeederName}}.go",
}
