package provider

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math/rand"
	"path/filepath"
	"strings"

	"go.ipao.vip/atomctl/pkg/utils/gomod"
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
	StructName       string
	ReturnType       string
	Mode             string
	ProviderGroup    string
	GrpcRegisterFunc string
	NeedPrepareFunc  bool
	InjectParams     map[string]InjectParam
	Imports          map[string]string
	PkgName          string
	ProviderFile     string
}

func atomPackage(suffix string) string {
	root := gomod.GetModuleName() + "/pkg/atom"
	if suffix != "" {
		return fmt.Sprintf("%s/%s", root, suffix)
	}
	return root
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

		// mode, returnType, group := parseDoc(docMark)
		// // log.Infof("mode: %s, returnType: %s, group: %s", mode, returnType, group)

		// if returnType == "#" {
		// 	provider.ReturnType = "*" + provider.StructName
		// } else {
		// 	provider.ReturnType = returnType
		// }

		// onlyMode := mode == "only"
		// exceptMode := mode == "except"
		// log.Infof("[%s] %s => ONLY: %+v, EXCEPT: %+v, Type: %s, Group: %s", source, declType.Name.Name, onlyMode, exceptMode, provider.ReturnType, provider.ProviderGroup)

		providerDoc := parseProvider(docMark)
		log.Infof("[%s] %s %+v", source, declType.Name.Name, providerDoc)
		provider.ProviderGroup = providerDoc.Group
		provider.ReturnType = providerDoc.ReturnType
		if provider.ReturnType == "" {
			provider.ReturnType = "*" + provider.StructName
		}

		for _, field := range structType.Fields.List {
			if field.Names == nil {
				continue
			}

			if field.Tag != nil {
				provider.NeedPrepareFunc = true
			}

			if providerDoc.IsOnly {
				if field.Tag == nil || !strings.Contains(field.Tag.Value, `inject:"true"`) {
					continue
				}
			} else {
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

		if providerDoc.Mode == "grpc" {
			provider.Mode = "grpc"

			modePkg := gomod.GetModuleName() + "/providers/grpc"

			provider.Imports[atomPackage("")] = ""
			provider.Imports[atomPackage("contracts")] = ""
			provider.Imports[modePkg] = ""

			provider.ProviderGroup = "atom.GroupInitial"
			provider.GrpcRegisterFunc = provider.ReturnType
			provider.ReturnType = "contracts.Initial"

			provider.InjectParams["__grpc"] = InjectParam{
				Star:         "*",
				Type:         "Grpc",
				Package:      modePkg,
				PackageAlias: "grpc",
			}
		}

		if providerDoc.Mode == "event" {
			provider.Mode = "event"

			modePkg := gomod.GetModuleName() + "/providers/event"

			provider.Imports[atomPackage("")] = ""
			provider.Imports[atomPackage("contracts")] = ""
			provider.Imports[modePkg] = ""

			provider.ProviderGroup = "atom.GroupInitial"
			provider.ReturnType = "contracts.Initial"

			provider.InjectParams["__event"] = InjectParam{
				Star:         "*",
				Type:         "PubSub",
				Package:      modePkg,
				PackageAlias: "event",
			}
		}

		if providerDoc.Mode == "job" {
			provider.Mode = "job"

			modePkg := gomod.GetModuleName() + "/providers/job"

			provider.Imports[atomPackage("")] = ""
			provider.Imports[atomPackage("contracts")] = ""
			provider.Imports["github.com/riverqueue/river"] = ""
			provider.Imports[modePkg] = ""

			provider.ProviderGroup = "atom.GroupInitial"
			provider.ReturnType = "contracts.Initial"

			provider.InjectParams["__job"] = InjectParam{
				Star:         "*",
				Type:         "Job",
				Package:      modePkg,
				PackageAlias: "job",
			}
		}

		providers = append(providers, provider)

	}

	return providers
}

// @provider(mode):[except|only] [returnType] [group]
type ProviderDescribe struct {
	IsOnly     bool
	Mode       string // job
	ReturnType string
	Group      string
}

func (p ProviderDescribe) String() {
	// log.Infof("[%s] %s => ONLY: %+v, EXCEPT: %+v, Type: %s, Group: %s", source, declType.Name.Name, onlyMode, exceptMode, provider.ReturnType, provider.ProviderGroup)
}

// @provider
// @provider(job)
// @provider(job):except
// @provider:except
// @provider:only
// @provider returnType
// @provider returnType group
// @provider(job) returnType group
func parseProvider(doc string) ProviderDescribe {
	result := ProviderDescribe{IsOnly: false}

	// Remove @provider prefix
	doc = strings.TrimSpace(strings.TrimPrefix(doc, "@provider"))

	// Handle empty case
	if doc == "" {
		return result
	}

	// Handle :except and :only
	if strings.Contains(doc, ":except") {
		result.IsOnly = false
		doc = strings.Replace(doc, ":except", "", 1)
	} else if strings.Contains(doc, ":only") {
		result.IsOnly = true
		doc = strings.Replace(doc, ":only", "", 1)
	}

	// Handle mode in parentheses
	if strings.Contains(doc, "(") && strings.Contains(doc, ")") {
		start := strings.Index(doc, "(")
		end := strings.Index(doc, ")")
		result.Mode = doc[start+1 : end]
		doc = doc[:start] + doc[end+1:]
	}

	// Handle remaining parts (returnType and group)
	parts := strings.Fields(strings.TrimSpace(doc))
	if len(parts) >= 1 {
		result.ReturnType = parts[0]
	}
	if len(parts) >= 2 {
		result.Group = parts[1]
	}

	return result
}
