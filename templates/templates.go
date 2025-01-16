package templates

import "embed"

//go:embed project
var Project embed.FS

//go:embed module
var Module embed.FS

//go:embed provider
var Provider embed.FS

//go:embed events
var Events embed.FS

//go:embed jobs
var Jobs embed.FS
