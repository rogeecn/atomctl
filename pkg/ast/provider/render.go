package provider

import (
	_ "embed"
	"html/template"
	"os"
	"strings"

	"go.ipao.vip/atomctl/pkg/utils/gomod"
	"github.com/samber/lo"
	"golang.org/x/tools/imports"
)

//go:embed provider.go.tpl
var providerTpl string

func Render(filename string, conf []Provider) error {
	defer func() {
		result, err := imports.Process(filename, nil, nil)
		if err == nil {
			os.WriteFile(filename, result, os.ModePerm)
		}
	}()

	imports := map[string]string{
		atomPackage("container"): "",
		atomPackage("opt"):       "",
	}
	lo.ForEach(conf, func(item Provider, _ int) {
		for k, v := range item.Imports {
			// 如果是当前包的引用，直接使用包名
			if strings.HasSuffix(k, "/"+v) {
				imports[k] = ""
				continue
			}

			if gomod.GetPackageModuleName(k) == v {
				imports[k] = ""
				continue
			}

			imports[k] = v
		}
	})

	t := template.Must(template.New("provider").Parse(providerTpl))

	data := struct {
		PkgName   string
		Imports   map[string]string
		Providers []Provider
	}{
		PkgName:   conf[0].PkgName,
		Imports:   imports,
		Providers: conf,
	}

	fd, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer fd.Close()

	return t.Execute(fd, data)
}
