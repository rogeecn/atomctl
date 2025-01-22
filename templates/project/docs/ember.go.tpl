package docs

import (
	_ "embed"

	_ "github.com/rogeecn/swag"
)

//go:embed swagger.json
var SwaggerSpec string
