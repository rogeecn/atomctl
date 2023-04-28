package http

import (
	"embed"
)

//go:embed tpl/*
var Templates embed.FS

var Files = map[string]string{
	"database/migrations/migration.go.tpl": "database/migrations/migrations.go",
	"database/seeders/seeder.go.tpl":       "database/seeders/seeders.go",
	"database/models/model.go.tpl":         "database/models/models.go",
	"database/query/query.go.tpl":          "database/query/query.go",
	"modules/modules.go.tpl":               "modules/modules.go",
	"config.toml":                          "config.toml",
	"main.go.tpl":                          "main.go",
	"Makefile":                             "Makefile",
	"go.mod.tpl":                           "go.mod",
	"gitignore":                            ".gitignore",
	"golangci.yaml":                        ".golangci.yaml",
}
