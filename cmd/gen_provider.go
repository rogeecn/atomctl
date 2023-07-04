package cmd

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"git-open.qianxin-inc.cn/free/csmp/atomctl/utils"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

// pro represents the routes command
var genProviderCmd = &cobra.Command{
	Use:   "provider",
	Short: "generate auto provider",
	RunE:  genProvider,
}

func init() {
	genCmd.AddCommand(genProviderCmd)
}

func genProvider(cmd *cobra.Command, args []string) error {
	var err error
	var path string
	if len(args) == 0 {
		path, err = os.Getwd()
		if err != nil {
			return err
		}
	} else {
		path = args[0]
	}

	path, _ = filepath.Abs(path)

	projectPkgName := getPackage()

	providers := []Provider{}
	// if path is file, then get the dir
	if utils.IsFile(path) && strings.HasSuffix(path, ".go") {
		log.Println("generate providers for file")
		providers = append(providers, astParseProviders(projectPkgName, path)...)
	} else {
		log.Println("generate providers for dir")
		// travel controller to find all controller objects
		_ = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}

			if !strings.HasSuffix(path, ".go") {
				return nil
			}

			providers = append(providers, astParseProviders(projectPkgName, path)...)
			return nil
		})
	}

	// generate files
	groups := lo.GroupBy(providers, func(item Provider) string {
		return item.ProviderFile
	})

	for file, conf := range groups {
		if err := renderFile(file, conf); err != nil {
			return err
		}
	}
	return nil
}

type Provider struct {
	StructName   string
	InjectParams map[string]string
	Imports      []string
	PkgName      string
	ProviderFile string
}

func astParseProviders(projectPkg, source string) []Provider {
	if strings.HasSuffix(source, "_test.go") {
		return []Provider{}
	}

	if strings.HasSuffix(source, "/provider.go") {
		return []Provider{}
	}

	providers := []Provider{}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, source, nil, parser.ParseComments)
	if err != nil {
		log.Println("ERR: ", err)
		return nil
	}

	imports := make(map[string]string)
	for _, imp := range node.Imports {

		name := ""
		if imp.Name == nil {
			paths := strings.Split(strings.Trim(imp.Path.Value, "\""), "/")
			name = paths[len(paths)-1]
		} else {
			name = imp.Name.Name
		}

		var pkg string
		if imp.Name != nil {
			pkg = strings.Trim(imp.Path.Value, `"`)
			if !strings.HasPrefix(pkg, projectPkg) {
				continue
			}
			pkg = fmt.Sprintf("%s %q", name, pkg)
		} else {
			pkg = strings.Trim(imp.Path.Value, "\"")
			if !strings.HasPrefix(pkg, projectPkg) {
				continue
			}
			pkg = fmt.Sprintf("%q", pkg)
		}

		imports[name] = pkg
	}

	// 再去遍历 struct 的方法去
	for _, decl := range node.Decls {
		provider := Provider{}

		decl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		if len(decl.Specs) == 0 {
			continue
		}

		declType, ok := decl.Specs[0].(*ast.TypeSpec)
		if !ok {
			continue
		}

		// 必须包含注释 // @provider:only/except
		if decl.Doc == nil {
			continue
		}

		if len(decl.Doc.List) == 0 {
			continue
		}

		structType, ok := declType.Type.(*ast.StructType)
		if !ok {
			continue
		}
		provider.StructName = declType.Name.Name

		docMark := strings.TrimLeft(decl.Doc.List[len(decl.Doc.List)-1].Text, "/ ")
		if !strings.HasPrefix(docMark, "@provider") {
			continue
		}

		onlyMode := strings.HasSuffix(docMark, ":only")
		exceptMode := strings.HasSuffix(docMark, ":except")
		log.Printf("%s => ONLY: %+v, EXCEPT: %+v", declType.Name.Name, onlyMode, exceptMode)

		for _, field := range structType.Fields.List {
			if provider.InjectParams == nil {
				provider.InjectParams = make(map[string]string)
				provider.Imports = []string{}
			}

			if onlyMode {
				if field.Tag == nil || !strings.Contains(field.Tag.Value, `inject:true`) {
					continue
				}
			}

			if exceptMode {
				if field.Tag != nil && strings.Contains(field.Tag.Value, `inject:false`) {
					continue
				}
			}

			var pkg string
			var typ string
			switch field.Type.(type) {
			case *ast.Ident:
				typ = field.Type.(*ast.Ident).Name
			case *ast.StarExpr:
				paramsType := field.Type.(*ast.StarExpr)
				switch paramsType.X.(type) {
				case *ast.SelectorExpr:
					X := paramsType.X.(*ast.SelectorExpr)

					pkg = X.X.(*ast.Ident).Name
					if _, ok := imports[pkg]; !ok {
						continue
					}

					typ = fmt.Sprintf("*%s.%s", X.X.(*ast.Ident).Name, X.Sel.Name)
				default:
					typ = fmt.Sprintf("*%s", paramsType.X.(*ast.Ident).Name)
				}
			case *ast.SelectorExpr:
				pkg = field.Type.(*ast.SelectorExpr).X.(*ast.Ident).Name
				if _, ok := imports[pkg]; !ok {
					continue
				}
				typ = fmt.Sprintf("%s.%s", field.Type.(*ast.SelectorExpr).X.(*ast.Ident).Name, field.Type.(*ast.SelectorExpr).Sel.Name)
			}

			if strings.HasSuffix(typ, "Context") || strings.HasSuffix(typ, "Ctx") {
				continue
			}
			// todo: ignore int int64 in32 ...

			for _, name := range field.Names {
				provider.InjectParams[name.Name] = typ
			}

			if importPkg, ok := imports[pkg]; ok {
				provider.Imports = append(provider.Imports, importPkg)
			}
		}

		provider.PkgName = node.Name.Name
		provider.ProviderFile = filepath.Join(filepath.Dir(source), "provider.go")

		providers = append(providers, provider)

	}

	return providers
}

func renderFile(filename string, conf []Provider) error {
	imports := []string{
		`"github.com/rogeecn/atom/utils/opt"`,
		`"github.com/rogeecn/atom/container"`,
	}
	lo.ForEach(conf, func(item Provider, _ int) {
		imports = lo.Union(imports, item.Imports)
	})

	// render file
	fd, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, _ = fd.WriteString(fmt.Sprintf("package %s\n\n", conf[0].PkgName))
	_, _ = fd.WriteString("import (\n")
	for _, imp := range imports {
		_, _ = fd.WriteString(fmt.Sprintf("\t%s\n", imp))
	}
	_, _ = fd.WriteString(")\n\n")

	_, _ = fd.WriteString("func Provide(opts ...opt.Option) error {\n")

	lo.ForEach(conf, func(item Provider, _ int) {
		// inject params
		params := []string{}
		structParams := []string{}
		for name, typ := range item.InjectParams {
			params = append(params, fmt.Sprintf("%s %s", name, typ))
			structParams = append(structParams, fmt.Sprintf("\t\t\t%s:%s,", name, name))
		}

		_, _ = fd.WriteString("\tif err:=container.Container.Provide(func(")
		_, _ = fd.WriteString(strings.Join(params, ", "))
		_, _ = fd.WriteString(fmt.Sprintf(") (*%s, error) {\n", item.StructName))
		_, _ = fd.WriteString(fmt.Sprintf("\t\treturn &%s{\n", item.StructName))
		_, _ = fd.WriteString(strings.Join(structParams, "\n") + "\n")
		_, _ = fd.WriteString("\t\t}, nil\n")
		_, _ = fd.WriteString("\t}); err!=nil{\n")
		_, _ = fd.WriteString("\t\treturn nil\n")
		_, _ = fd.WriteString("\t}\n\n")
	})

	_, _ = fd.WriteString("return nil\n")
	_, _ = fd.WriteString("}\n\n")
	return nil
}
