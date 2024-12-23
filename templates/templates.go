package templates

import "embed"

//go:embed project
var Project embed.FS

//go:embed module
var Module embed.FS

//go:embed provider
var Provider embed.FS
