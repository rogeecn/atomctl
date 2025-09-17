package templates

import "embed"

//go:embed project/* project/.*
var Project embed.FS

//go:embed module
var Module embed.FS

//go:embed events
var Events embed.FS

//go:embed jobs
var Jobs embed.FS

//go:embed services
var Services embed.FS

//go:embed providers
var Providers embed.FS
