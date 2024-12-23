// Code generated by the atomctl ; DO NOT EDIT.

package {{.PackageName}}

import (
{{- range .Imports }}
	{{.}}
{{- end }}
	. "{{.ProjectPackage}}/pkg/f"

	_ "git.ipao.vip/rogeecn/atom"
	_ "git.ipao.vip/rogeecn/atom/contracts"
	"github.com/gofiber/fiber/v3"
	log "github.com/sirupsen/logrus"
)

// @provider contracts.HttpRoute atom.GroupRoutes
type Routes struct {
	log *log.Entry `inject:"false"`
{{- range .Controllers }}
	{{.}}
{{- end }}
}

func (r *Routes) Prepare() error {
	r.log = log.WithField("module", "routes.{{.PackageName}}")
	return nil
}

func (r *Routes) Name() string {
	return "{{.PackageName}}"
}

func (r *Routes) Register(router fiber.Router) {
{{- range $key, $value := .Routes }}
	// 注册路由组: {{$key}}
	{{- range $value }}
	router.{{.Method}}("{{.Route}}", {{.Func}}(
		r.{{.Controller}}.{{.Action}},
		{{- range .Params}}
		{{.}},
		{{- end }}
	))
	{{ end }}
{{- end }}
}
