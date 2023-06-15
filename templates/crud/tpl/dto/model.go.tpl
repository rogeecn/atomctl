package dto

import (
	"time"
	"gorm.io/gorm"
)

type {{ .Model.Name }}Form struct {
	{{ .Model.Name }}Item `json:",inline"`
}

type {{ .Model.Name }}ListQueryFilter struct {
	{{ .Model.Name }}Item `json:",inline"`
}

type {{ .Model.Name }}Item struct {
	{{- range .Model.Fields }}
	{{ .Name }} {{ .Type }} {{ .Tag }} // {{ .Comment }}
	{{- end }}
}
