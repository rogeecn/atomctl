package migration

import (
	"embed"
)

//go:embed tpl/*
var Templates embed.FS

var Files = map[string]string{
	"migration.tpl": "./database/migrations/{{.ID}}_{{.SnakeMigrationName}}.go",
}
