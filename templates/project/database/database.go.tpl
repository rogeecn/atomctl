package database

import (
	"embed"
)

//go:embed migrations/*
var MigrationFS embed.FS
