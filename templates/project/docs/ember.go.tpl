package docs

import (
	_ "embed"

	_ "git.ipao.vip/rogeecn/atomctl/pkg/swag"
)

//go:embed swagger.json
var SwaggerSpec string
