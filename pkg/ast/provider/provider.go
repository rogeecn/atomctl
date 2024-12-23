package provider

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math/rand"
	"path/filepath"
	"strings"

	"git.ipao.vip/rogeecn/atomctl/pkg/utils/gomod"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
)

func getTypePkgName(typ string) string {
	if strings.Contains(typ, ".") {
		return strings.Split(typ, ".")[0]
	}
	return ""
}

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

func Parse(source string) []Provider {
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
		log.Error("ERR: ", err)
		return nil
	}
	imports := make(map[string]string)
	for _, imp := range node.Imports {
		name := ""
		pkgPath := strings.Trim(imp.Path.Value, "\"")

		if imp.Name != nil {
			// 如果有显式指定包名，直接使用
			name = imp.Name.Name
		} else {
			// 尝试从go.mod中获取真实包名
			name = gomod.GetPackageModuleName(pkgPath)
		}

		// 处理匿名导入的情况
		if name == "_" {
			name = gomod.GetPackageModuleName(pkgPath)

			// 处理重名
			if _, ok := imports[name]; ok {
				name = fmt.Sprintf("%s%d", name, rand.Intn(100))
			}
		}

		imports[name] = pkgPath
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
		// log.Infof("mode: %s, returnType: %s, group: %s", mode, returnType, group)

		if returnType == "#" {
			provider.ReturnType = "*" + provider.StructName
		} else {
			provider.ReturnType = returnType
		}
		onlyMode := mode == "only"
		exceptMode := mode == "except"
		log.Infof("[%s] %s => ONLY: %+v, EXCEPT: %+v, Type: %s, Group: %s", source, declType.Name.Name, onlyMode, exceptMode, provider.ReturnType, provider.ProviderGroup)

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
		provider.ProviderFile = filepath.Join(filepath.Dir(source), "provider.gen.go")

		providers = append(providers, provider)

	}

	return providers
}

func parseDoc(doc string) (string, string, string) {
	// @provider:[except|only] [returnType] [group]
	doc = strings.TrimLeft(doc[len("@provider"):], ":")
	if !strings.HasPrefix(doc, "except") && !strings.HasPrefix(doc, "only") {
		doc = "except " + doc
	}

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
