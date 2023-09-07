package dto

import (
	"{{ .PkgName }}/common"
	{{- range .Model.Imports }}
	{{ . }}
	{{- end }}
)

type {{ .Model.Name }}Form struct {
	{{- range .Model.Fields }}
	{{- if or (eq .Name "ID") (eq .Name "CreatedAt") (eq .Name "UpdatedAt") (eq .Name "DeletedAt") }}
	{{- else }}
	{{- if .PackageAlias }}
	{{ .Name }} {{ .PackageAlias }}.{{ .Type }} `form:"{{.Tag}}" json:"{{ .Tag }},omitempty"` // {{ .Comment }}
	{{- else }}
	{{ .Name }} {{ .Type }} `form:"{{ .Tag }}" json:"{{ .Tag }},omitempty"` // {{ .Comment }}
	{{- end }}
	{{- end }}
	{{- end }}
}

type {{ .Model.Name }}ListQueryFilter struct {
	{{- range .Model.Fields }}
	{{- if or (eq .Name "ID") (eq .Name "CreatedAt") (eq .Name "UpdatedAt") (eq .Name "DeletedAt") }}
	{{- else }}
	{{- if .PackageAlias }}
	{{ .Name }} *{{ .PackageAlias }}.{{ .Type }} `query:"{{.Tag}}" json:"{{ .Tag }},omitempty"` // {{ .Comment }}
	{{- else }}
	{{ .Name }} *{{ .Type }} `query:"{{.Tag}}" json:"{{ .Tag }},omitempty"` // {{ .Comment }}
	{{- end }}
	{{- end }}
	{{- end }}
}


type {{ .Model.Name }}Item struct {
	{{- range .Model.Fields }}
	{{- if eq .Name "DeletedAt" }}
	{{- else }}
	{{- if .PackageAlias }}
	{{ .Name }} {{ .PackageAlias }}.{{ .Type }} `json:"{{ .Tag }},omitempty"` // {{ .Comment }}
	{{- else }}
	{{ .Name }} {{ .Type }} `json:"{{ .Tag }},omitempty"` // {{ .Comment }}
	{{- end }}
	{{- end }}
	{{- end }}
}
