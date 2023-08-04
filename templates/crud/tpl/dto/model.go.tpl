package dto

import (
	{{- range .Model.Imports }}
	{{ . }}
	{{- end }}
)

type {{ .Model.Name }}Form struct {
	{{ .Model.Name }}Item `json:",inline"`
}

type {{ .Model.Name }}ListQueryFilter struct {
	{{ .Model.Name }}Item `json:",inline"`
}

type {{ .Model.Name }}Item struct {
	{{- range .Model.Fields }}
	{{- if .PackageAlias }}
	{{ .Name }} *{{ .PackageAlias }}.{{ .Type }} {{ .Tag }} // {{ .Comment }}
	{{- else }}
	{{ .Name }} *{{ .Type }} {{ .Tag }} // {{ .Comment }}
	{{- end }}
	{{- end }}
}
