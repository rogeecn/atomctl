package services

import (
	"gorm.io/gorm"
)

var _db *gorm.DB

// exported CamelCase Services
var (
{{- range . }}
	{{ .CamelName }} *{{ .ServiceName }}
{{- end }}
)

// @provider(model)
type services struct {
	db *gorm.DB
	// define Services
	{{- range . }}
	{{ .ServiceName }} *{{ .ServiceName }}
	{{- end }}
}

func (svc *services) Prepare() error {
	_db = svc.db

	// set exported Services here
	{{- range . }}
	{{ .CamelName }} = svc.{{ .ServiceName }}
	{{- end }}

	return nil
}
