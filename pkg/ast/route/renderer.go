package route

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

var routerTmpl = template.Must(template.New("route").
	Funcs(sprig.FuncMap()).
	Option("missingkey=error").
	Parse(routeTpl),
)

func renderTemplate(data RenderData) ([]byte, error) {
	var buf bytes.Buffer
	if err := routerTmpl.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
