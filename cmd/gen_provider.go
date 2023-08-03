package cmd

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/rogeecn/atomctl/utils"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var scalarTypes = []string{
	"float32",
	"float64",
	"int",
	"int8",
	"int16",
	"int32",
	"int64",
	"uint",
	"uint8",
	"uint16",
	"uint32",
	"uint64",
	"bool",
	"uintptr",
	"complex64",
	"complex128",
}

// pro represents the routes command
var genProviderCmd = &cobra.Command{
	Use:   "provider",
	Short: "generate auto provider, @provider:[except|only] [returnType] [group]",
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

type InjectParam struct {
	Star         string
	Type         string
	Package      string
	PackageAlias string
}
type Provider struct {
	StructName      string
	ReturnType      string
	ProviderGroup   string
	NeedPrepareFunc bool
	InjectParams    map[string]InjectParam
	Imports         map[string]string
	PkgName         string
	ProviderFile    string
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
	vPattern := regexp.MustCompile(`^v\d+$`)
	imports := make(map[string]string)
	for _, imp := range node.Imports {

		name := ""
		if imp.Name == nil {
			paths := strings.Split(strings.Trim(imp.Path.Value, "\""), "/")
			name = paths[len(paths)-1]
			if vPattern.MatchString(name) {
				name = paths[len(paths)-2]
			}
		} else {
			name = imp.Name.Name
		}

		if name == "_" {
			paths := strings.Split(strings.Trim(imp.Path.Value, "\""), "/")
			name = paths[len(paths)-1]
			if vPattern.MatchString(name) {
				name = paths[len(paths)-2]
			}
			if _, ok := imports[name]; ok {
				name = fmt.Sprintf("%s%d", name, rand.Intn(100))
			}
		}
		pkg := strings.Trim(imp.Path.Value, `"`)
		imports[name] = pkg
	}

	// 再去遍历 struct 的方法去
	for _, decl := range node.Decls {
		provider := Provider{
			InjectParams: make(map[string]InjectParam),
			Imports:      make(map[string]string),
		}

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

		docMark := strings.TrimLeft(decl.Doc.List[len(decl.Doc.List)-1].Text, "/ \t")
		if !strings.HasPrefix(docMark, "@provider") {
			continue
		}
		mode, returnType, group := parseDoc(docMark)
		if group != "" {
			provider.ProviderGroup = group
		}
		// fmt.Println(mode, returnType, group, provider.ProviderGroup)

		if returnType == "#" {
			provider.ReturnType = "*" + provider.StructName
		} else {
			provider.ReturnType = returnType
		}
		onlyMode := mode == "only"
		exceptMode := mode == "except"
		log.Printf("[%s] %s => ONLY: %+v, EXCEPT: %+v, Type: %s, Group: %s", source, declType.Name.Name, onlyMode, exceptMode, provider.ReturnType, provider.ProviderGroup)

		for _, field := range structType.Fields.List {
			if field.Names == nil {
				continue
			}

			if field.Tag != nil {
				provider.NeedPrepareFunc = true
			}

			if onlyMode {
				if field.Tag == nil || !strings.Contains(field.Tag.Value, `inject:"true"`) {
					continue
				}
			}

			if exceptMode {
				if field.Tag != nil && strings.Contains(field.Tag.Value, `inject:"false"`) {
					continue
				}
			}

			var star string
			var pkg string
			var pkgAlias string
			var typ string
			switch field.Type.(type) {
			case *ast.Ident:
				typ = field.Type.(*ast.Ident).Name
			case *ast.StarExpr:
				star = "*"
				paramsType := field.Type.(*ast.StarExpr)
				switch paramsType.X.(type) {
				case *ast.SelectorExpr:
					X := paramsType.X.(*ast.SelectorExpr)

					pkgAlias = X.X.(*ast.Ident).Name
					p, ok := imports[pkgAlias]
					if !ok {
						continue
					}
					pkg = p

					typ = X.Sel.Name
				default:
					typ = paramsType.X.(*ast.Ident).Name
				}
			case *ast.SelectorExpr:
				pkgAlias = field.Type.(*ast.SelectorExpr).X.(*ast.Ident).Name
				p, ok := imports[pkgAlias]
				if !ok {
					continue
				}
				pkg = p
				typ = field.Type.(*ast.SelectorExpr).Sel.Name
			}

			if lo.Contains(scalarTypes, typ) {
				continue
			}

			for _, name := range field.Names {
				provider.InjectParams[name.Name] = InjectParam{
					Star:         star,
					Type:         typ,
					Package:      pkg,
					PackageAlias: pkgAlias,
				}
			}

			if importPkg, ok := imports[pkgAlias]; ok {
				provider.Imports[importPkg] = pkgAlias
			}
		}

		if pkgAlias := getTypePkgName(provider.ReturnType); pkgAlias != "" {
			if importPkg, ok := imports[pkgAlias]; ok {
				provider.Imports[importPkg] = pkgAlias
			}
		}

		if pkgAlias := getTypePkgName(provider.ProviderGroup); pkgAlias != "" {
			if importPkg, ok := imports[pkgAlias]; ok {
				provider.Imports[importPkg] = pkgAlias
			}
		}

		provider.PkgName = node.Name.Name
		provider.ProviderFile = filepath.Join(filepath.Dir(source), "provider.go")

		providers = append(providers, provider)

	}

	return providers
}

func renderFile(filename string, conf []Provider) error {
	imports := map[string]string{
		"seccloud/cspm/pkg/container": "",
		"seccloud/cspm/pkg/utils/opt": "",
	}
	lo.ForEach(conf, func(item Provider, _ int) {
		for k, v := range item.Imports {
			imports[k] = v
		}
	})

	// render file
	fd, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, _ = fd.WriteString(fmt.Sprintf("package %s\n\n", conf[0].PkgName))
	_, _ = fd.WriteString("import (\n")
	for pkg, alias := range imports {
		if strings.HasSuffix(pkg, "/"+alias) || pkg == alias {
			fd.WriteString(fmt.Sprintf("\t%q\n", pkg))
			continue
		}
		fd.WriteString(fmt.Sprintf("\t%s %q\n", alias, pkg))
	}
	_, _ = fd.WriteString(")\n\n")

	_, _ = fd.WriteString("func Provide(opts ...opt.Option) error {\n")

	lo.ForEach(conf, func(item Provider, _ int) {
		// sort params
		keys := lo.Keys(item.InjectParams)
		sort.Strings(keys)

		fd.WriteString("\tif err := container.Container.Provide(func(")
		if len(keys) > 0 {
			fd.WriteString("\n")
			for _, key := range keys {
				name, param := key, item.InjectParams[key]
				fd.WriteString(fmt.Sprintf("\t\t%s %s", name, param.Star))

				if param.Package == "" {
					fd.WriteString(fmt.Sprintf("%s,\n", param.Type))
					continue
				}

				if alias, ok := imports[param.Package]; ok {
					fd.WriteString(fmt.Sprintf("%s.%s,\n", alias, param.Type))
				}
			}
			fd.WriteString("\t")
		}
		fd.WriteString(fmt.Sprintf(") (%s, error) {\n", item.ReturnType))

		fd.WriteString(fmt.Sprintf("\t\tobj := &%s{", item.StructName))
		if len(keys) > 0 {
			fd.WriteString("\n")
			for _, name := range keys {
				fd.WriteString(fmt.Sprintf("\t\t\t%s: %s,\n", name, name))
			}
			fd.WriteString("\t\t")
		}
		fd.WriteString("}\n")

		if item.NeedPrepareFunc {
			_, _ = fd.WriteString("\t\tif err := obj.Prepare(); err != nil {\n\t\t\treturn nil, err\n}\n")
		}
		_, _ = fd.WriteString("\t\treturn obj, nil\n")
		_, _ = fd.WriteString("\t}")
		if item.ProviderGroup != "" {
			_, _ = fd.WriteString(fmt.Sprintf(", %s", item.ProviderGroup))
		}
		_, _ = fd.WriteString("); err != nil {\n")
		_, _ = fd.WriteString("\t\treturn err\n")
		_, _ = fd.WriteString("\t}\n\n")
	})

	_, _ = fd.WriteString("\treturn nil\n")
	_, _ = fd.WriteString("}\n\n")
	return nil
}

func parseDoc(doc string) (string, string, string) {
	// @provider:[except|only] [returnType] [group]
	doc = strings.TrimLeft(doc[len("@provider"):], ":")
	doc = strings.ReplaceAll(doc, "\t", " ")
	cmds := strings.Split(doc, " ")
	cmds = lo.Filter(cmds, func(item string, idx int) bool {
		return strings.TrimSpace(item) != ""
	})

	if len(cmds) == 0 {
		return "except", "#", ""
	}

	if len(cmds) == 1 {
		return cmds[0], "#", ""
	}

	if len(cmds) == 2 {
		return cmds[0], cmds[1], ""
	}

	return cmds[0], cmds[1], cmds[2]
}

func getTypePkgName(typ string) string {
	if strings.Contains(typ, ".") {
		return strings.Split(typ, ".")[0]
	}
	return ""
}
