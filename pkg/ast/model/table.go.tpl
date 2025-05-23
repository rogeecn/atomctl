package model

var {{.CamelTable}}UpdateExcludeColumns = []Column{
	{{- if .HasCreatedAt}}
	table.{{.PascalTable}}.CreatedAt,
	{{- end}}

	{{- if .SoftDelete}}
	table.{{.PascalTable}}.DeletedAt,
	{{- end}}
}
