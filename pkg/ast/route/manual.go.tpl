package {{.PackageName}}

func (r *Routes) Path() string {
	return "/{{.PackageName}}"
}

func (r *Routes) Middlewares() []any{
	return []any{}
}
