package provider

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math/rand"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"go.ipao.vip/atomctl/v2/pkg/utils/gomod"
)

// ================================================================================================
// 原有的 Parse 函数 - 保持向后兼容性
// ================================================================================================
//
// 注意：这是重构前的原始解析函数，现在保持向后兼容性。
// 新代码应该使用：
//   - ParseRefactored() - 简化的单文件解析
//   - NewParser().ParseFile() - 完整的解析功能
//   - NewGoParser().ParseDir() - 目录解析功能
//
// 原始函数的缺点：
//   - 单体设计：所有逻辑集中在一个函数中
//   - 难以测试：无法单独测试各个功能模块
//   - 扩展困难：添加新功能需要修改核心函数
//   - 错误处理简单：缺乏详细的错误信息和上下文
// ================================================================================================

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

func atomPackage(suffix string) string {
	root := "go.ipao.vip/atom"
	if suffix != "" {
		return fmt.Sprintf("%s/%s", root, suffix)
	}
	return root
}

// Parse 原始的 Provider 解析函数 - 保持向后兼容性
//
// ⚠️  警告：这是一个遗留函数，建议使用重构后的版本：
//   - 简单使用：ParseRefactored(source)
//   - 完整功能：NewParser().ParseFile(source)
//   - 目录解析：NewGoParser().ParseDir(dir)
//
// 执行流程（原始版本）：
// ┌─────────────────────────────────────────────────────────────┐
// │                      Parse(source)                          │
// ├─────────────────────────────────────────────────────────────┤
// │  1. 文件过滤：跳过测试文件和生成文件                      │
// │  2. AST解析：使用标准库解析 Go 文件                       │
// │  3. 导入处理：构建导入映射表                              │
// │  4. 遍历声明：查找带有 @provider 注解的结构体               │
// │  5. 注解解析：解析 @provider 语法                         │
// │  6. 字段处理：处理结构体字段和注入参数                     │
// │  7. 模式应用：根据模式应用特定逻辑                         │
// │  8. 结果收集：收集所有有效的 Provider                     │
// └─────────────────────────────────────────────────────────────┘
//
// 参数：
//   - source: Go 源文件路径
//
// 返回值：
//   - []Provider: 解析到的 Provider 列表（解析失败时为 nil）
//
// 缺点：
//   - 错误处理不完善：解析失败时返回 nil，丢失错误信息
//   - 单体设计：所有逻辑集中在一个函数中，难以维护
//   - 缺乏扩展性：添加新功能需要修改核心函数
//   - 性能问题：没有缓存和优化机制
//
// 兼容性说明：
//   - 保持原有接口不变
//   - 现有调用代码可以继续工作
//   - 建议逐步迁移到新版本
func Parse(source string) []Provider {
	// === 步骤 1：文件过滤 ===
	// 跳过测试文件（_test.go 后缀）
	if strings.HasSuffix(source, "_test.go") {
		return []Provider{}
	}

	// 跳过生成的 provider 文件（避免循环解析）
	if strings.HasSuffix(source, "/provider.gen.go") {
		return []Provider{}
	}

	// 初始化结果列表
	providers := []Provider{}

	// === 步骤 2：AST 解析 ===
	// 使用 Go 标准库将源文件解析为抽象语法树
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, source, nil, parser.ParseComments)
	if err != nil {
		log.Error("ERR: ", err)
		return nil  // 原始版本在错误时返回 nil
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
			Mode:         ProviderModeBasic, // Default mode
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
			provider.Mode = ProviderModeGrpc

			modePkg := gomod.GetModuleName() + "/providers/grpc"

			provider.Imports[atomPackage("")] = ""
			provider.Imports[atomPackage("contracts")] = ""
			provider.Imports[modePkg] = ""

			provider.ProviderGroup = "atom.GroupInitial"

			// Handle gRPC register function correctly
			if providerDoc.ReturnType != "" && strings.Contains(providerDoc.ReturnType, ".") {
				// User specified a complete register function name, like userv1.RegisterUserServiceServer
				provider.GrpcRegisterFunc = providerDoc.ReturnType
				// Extract package information and add import
				if pkgAlias := getTypePkgName(providerDoc.ReturnType); pkgAlias != "" {
					if importPkg, ok := imports[pkgAlias]; ok {
						provider.Imports[importPkg] = pkgAlias
					}
				}
			} else {
				// Generate default gRPC register function name
				// Example: UserService -> RegisterUserServiceServer
				provider.GrpcRegisterFunc = "Register" + strings.TrimPrefix(provider.ReturnType, "*") + "Server"
			}

			provider.ReturnType = "contracts.Initial"

			provider.InjectParams["__grpc"] = InjectParam{
				Star:         "*",
				Type:         "Grpc",
				Package:      modePkg,
				PackageAlias: "grpc",
			}
		}

		if providerDoc.Mode == "event" {
			provider.Mode = ProviderModeEvent

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
			provider.Mode = ProviderModeJob

			modePkg := gomod.GetModuleName() + "/providers/job"

			provider.Imports["github.com/riverqueue/river"] = ""
			provider.Imports[modePkg] = ""

			provider.InjectParams["__job"] = InjectParam{
				Star:         "*",
				Type:         "Job",
				Package:      modePkg,
				PackageAlias: "job",
			}
		} else if providerDoc.Mode == "cronjob" {
			provider.Mode = ProviderModeCronJob

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

		if providerDoc.Mode == "model" {
			provider.Mode = ProviderModeModel

			provider.ProviderGroup = "atom.GroupInitial"
			provider.ReturnType = "contracts.Initial"
			provider.NeedPrepareFunc = true
		}

		providers = append(providers, provider)

	}

	return providers
}

func (p ProviderDescribe) String() {
	// log.Infof("[%s] %s => ONLY: %+v, EXCEPT: %+v, Type: %s, Group: %s", source, declType.Name.Name, onlyMode, exceptMode, provider.ReturnType, provider.ProviderGroup)
}

// parseProvider 解析 @provider 注解的语法
//
// 支持的语法格式：
//   @provider                    - 基本格式
//   @provider(job)              - 指定模式
//   @provider(job):except       - 排除模式
//   @provider:except           - 排除模式（无模式）
//   @provider:only             - 仅包含模式
//   @provider returnType         - 指定返回类型
//   @provider returnType group   - 指定返回类型和分组
//   @provider(job) returnType group - 完整格式
//
// 解析规则：
//   1. 移除 "@provider" 前缀
//   2. 处理模式（括号内的内容）
//   3. 处理注入模式（:except 或 :only）
//   4. 解析返回类型和分组（剩余部分）
//
// 参数：
//   - doc: @provider 注解字符串
//
// 返回值：
//   - ProviderDescribe: 解析后的注解信息
//
// 示例：
//   输入: "@provider(job) contracts.Initial atom.GroupInitial"
//   输出: ProviderDescribe{
//       Mode: "job",
//       ReturnType: "contracts.Initial",
//       Group: "atom.GroupInitial",
//       IsOnly: false
//   }
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
