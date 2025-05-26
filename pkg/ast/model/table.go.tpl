package model

var tbl{{.PascalTable}}UpdateMutableColumns = tbl{{.PascalTable}}.MutableColumns.Except(
	{{- if .HasCreatedAt}}
	tbl{{.PascalTable}}.CreatedAt,
	{{- end}}

	{{- if .SoftDelete}}
	tbl{{.PascalTable}}.DeletedAt,
	{{- end}}
)
