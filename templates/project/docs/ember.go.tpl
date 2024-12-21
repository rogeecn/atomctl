package docs

import (
	_ "embed"

	_ "github.com/swaggo/swag"
)

//go:embed swagger.json
var SwaggerSpec string
