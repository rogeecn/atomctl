package {{.PkgName}}

import (
{{- range $pkg, $alias := .Imports }}
	{{- if eq $alias "" }}
	"{{$pkg}}"
	{{- else }}
	{{$alias}} "{{$pkg}}"
	{{- end }}
{{- end }}
)

func Provide(opts ...opt.Option) error {
{{- range .Providers }}
	if err := container.Container.Provide(func(
	{{- range $key, $param := .InjectParams }}
		{{$key}} {{$param.Star}}{{if eq $param.Package ""}}{{$param.Type}}{{else}}{{$param.PackageAlias}}.{{$param.Type}}{{end}},
	{{- end }}
	) ({{.ReturnType}}, error) {
		obj := &{{.StructName}}{
		{{- range $key, $param := .InjectParams }}
			{{$key}}: {{$key}},
		{{- end }}
		}
		{{- if .NeedPrepareFunc }}
		if err := obj.Prepare(); err != nil {
			return nil, err
		}
		{{- end }}
		return obj, nil
	}{{if .ProviderGroup}}, {{.ProviderGroup}}{{end}}); err != nil {
		return err
	}
{{- end }}
	return nil
}