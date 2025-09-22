package provider

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"go.ipao.vip/atomctl/v2/pkg/utils/gomod"
)

// Parser 定义了解析器接口，为 Provider 注解解析提供统一的抽象层
//
// 接口设计原则：
// ┌─────────────────────────────────────────────────────────────┐
// │                    Parser Interface                        │
// ├─────────────────────────────────────────────────────────────┤
// │  核心功能：                                               │
// │  - ParseFile: 单文件解析                                   │
// │  - ParseDir: 目录解析                                      │
// │  - ParseString: 字符串解析                                 │
// │  配置管理：                                               │
// │  - SetConfig/GetConfig: 配置设置和获取                    │
// │  - GetContext: 获取解析上下文                             │
// └─────────────────────────────────────────────────────────────┘
//
// 设计理念：
//   - 接口隔离：只包含必要的方法，避免过度设计
//   - 扩展性：支持不同的解析器实现
//   - 配置化：支持运行时配置修改
//   - 上下文感知：提供解析过程的上下文信息
type Parser interface {
	// ParseFile 解析单个 Go 文件并返回发现的 Provider 列表
	//
	// 参数：
	//   - filePath: Go 源文件路径
	//
	// 返回值：
	//   - []Provider: 解析到的 Provider 列表
	//   - error: 解析过程中的错误（如文件不存在、语法错误等）
	//
	// 使用场景：
	//   - 需要解析单个文件的 Provider 注解
	//   - 精确控制要解析的文件
	//   - 集成到构建工具中
	ParseFile(filePath string) ([]Provider, error)

	// ParseDir 解析目录中的所有 Go 文件并返回发现的 Provider 列表
	//
	// 执行流程：
	//   1. 遍历目录中的所有文件
	//   2. 过滤出 Go 源文件（.go 后缀）
	//   3. 跳过测试文件（_test.go）和隐藏目录
	//   4. 对每个文件调用 ParseFile
	//   5. 汇总所有文件的 Provider 结果
	//
	// 参数：
	//   - dirPath: 要解析的目录路径
	//
	// 返回值：
	//   - []Provider: 目录中所有文件的 Provider 列表
	//   - error: 目录遍历或解析过程中的错误
	//
	// 使用场景：
	//   - 批量解析整个项目的 Provider
	//   - 代码生成工具的输入
	//   - 项目分析工具
	ParseDir(dirPath string) ([]Provider, error)

	// ParseString 从字符串解析 Go 代码并返回发现的 Provider 列表
	//
	// 参数：
	//   - code: Go 源代码字符串
	//
	// 返回值：
	//   - []Provider: 解析到的 Provider 列表
	//   - error: 解析过程中的错误
	//
	// 使用场景：
	//   - 动态生成的代码解析
	//   - IDE 插件的实时分析
	//   - 单元测试中的模拟数据
	ParseString(code string) ([]Provider, error)

	// SetConfig 设置解析器配置
	//
	// 支持的配置项：
	//   - StrictMode: 严格模式，启用更严格的验证
	//   - CacheEnabled: 启用解析结果缓存
	//   - SourceLocations: 包含源码位置信息
	//   - Mode: 解析模式（带注释、不带注释等）
	//
	// 参数：
	//   - config: 新的解析器配置
	//
	// 使用场景：
	//   - 根据不同环境调整解析行为
	//   - 性能优化时启用缓存
	//   - 调试时启用更详细的信息
	SetConfig(config *ParserConfig)

	// GetConfig 获取当前的解析器配置
	//
	// 返回值：
	//   - *ParserConfig: 当前配置的副本
	//
	// 使用场景：
	//   - 检查当前配置状态
	//   - 保存和恢复配置
	//   - 调试和诊断
	GetConfig() *ParserConfig

	// GetContext 获取当前的解析器上下文
	//
	// 返回值：
	//   - *ParserContext: 包含解析过程中的状态信息
	//     - FilesProcessed: 已处理的文件数量
	//     - FilesSkipped: 跳过的文件数量
	//     - ProvidersFound: 发现的 Provider 数量
	//     - Cache: 缓存的解析结果
	//     - Imports: 导入信息映射
	//
	// 使用场景：
	//   - 监控解析进度
	//   - 性能分析
	//   - 调试解析问题
	GetContext() *ParserContext
}

// GoParser Parser 接口的具体实现，专门用于解析 Go 源文件中的 Provider 注解
//
// 架构设计：
// ┌─────────────────────────────────────────────────────────────┐
// │                     GoParser                                │
// ├─────────────────────────────────────────────────────────────┤
// │  config:    *ParserConfig    - 解析器配置                  │
// │  context:   *ParserContext   - 解析上下文                  │
// │  mu:        sync.RWMutex    - 读写锁，保证并发安全        │
// └─────────────────────────────────────────────────────────────┘
//
// 核心特性：
//   - 线程安全：使用读写锁保护共享状态
//   - 缓存支持：可选的解析结果缓存
//   - 并发处理：支持多文件并发解析
//   - 错误恢复：单个文件解析失败不影响其他文件
//
// 适用场景：
//   - 大型项目的批量 Provider 解析
//   - 需要高性能的代码生成工具
//   - 需要线程安全的解析环境
type GoParser struct {
	config  *ParserConfig  // 解析器配置，控制解析行为
	context *ParserContext // 解析上下文，包含解析状态信息
	mu      sync.RWMutex  // 读写锁，保护 config 和 context 的并发访问
}

// NewGoParser 创建一个使用默认配置的 GoParser 实例
//
// 初始化流程：
// ┌─────────────────────────────────────────────────────────────┐
// │                  NewGoParser()                              │
// ├─────────────────────────────────────────────────────────────┤
// │  1. 创建默认配置：NewParserConfig()                        │
// │  2. 创建解析上下文：NewParserContext()                     │
// │  3. 初始化文件集：token.NewFileSet()                       │
// │  4. 构建解析器实例：&GoParser{...}                         │
// └─────────────────────────────────────────────────────────────┘
//
// 默认配置特点：
//   - 解析模式：带注释解析 (parser.ParseComments)
//   - 缓存：默认关闭
//   - 源码位置：默认关闭
//   - 严格模式：默认关闭
//
// 返回值：
//   - *GoParser: 配置好的解析器实例，可以直接使用
//
// 使用示例：
//   parser := NewGoParser()
//   providers, err := parser.ParseFile("user_service.go")
func NewGoParser() *GoParser {
	// 创建默认配置
	config := NewParserConfig()

	// 创建解析上下文
	context := NewParserContext(config)

	// 初始化文件集（用于记录源码位置信息）
	if config.FileSet == nil {
		config.FileSet = token.NewFileSet()
		context.FileSet = config.FileSet
	}

	// 构建并返回解析器实例
	return &GoParser{
		config:  config,
		context: context,
	}
}

// NewGoParserWithConfig 使用自定义配置创建 GoParser 实例
//
// 设计目标：
//   - 提供灵活的配置选项
//   - 支持不同的使用场景（调试、生产、测试）
//   - 保持向后兼容性
//
// 参数处理：
//   - 如果 config 为 nil，自动使用默认配置
//   - 如果 config.FileSet 为 nil，自动创建新的文件集
//
// 自定义配置场景：
//   - 调试模式：启用详细日志和源码位置
//   - 性能优化：启用缓存和并发解析
//   - 严格验证：启用严格模式进行代码质量检查
//
// 参数：
//   - config: 自定义的解析器配置（可为 nil）
//
// 返回值：
//   - *GoParser: 使用指定配置的解析器实例
//
// 使用示例：
//   config := &ParserConfig{
//       CacheEnabled: true,
//       StrictMode: true,
//       SourceLocations: true,
//   }
//   parser := NewGoParserWithConfig(config)
func NewGoParserWithConfig(config *ParserConfig) *GoParser {
	// 处理 nil 配置，保持向后兼容
	if config == nil {
		return NewGoParser()
	}

	// 使用自定义配置创建解析上下文
	context := NewParserContext(config)

	// 初始化文件集（如果未提供）
	if config.FileSet == nil {
		config.FileSet = token.NewFileSet()
		context.FileSet = config.FileSet
	}

	// 构建并返回解析器实例
	return &GoParser{
		config:  config,
		context: context,
	}
}

// ParseFile implements Parser.ParseFile
func (p *GoParser) ParseFile(filePath string) ([]Provider, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Check if file should be included
	if !p.context.ShouldIncludeFile(filePath) {
		p.context.FilesSkipped++
		return []Provider{}, nil
	}

	// Check cache if enabled
	if p.config.CacheEnabled {
		if cached, found := p.context.Cache[filePath]; found {
			if providers, ok := cached.([]Provider); ok {
				return providers, nil
			}
		}
	}

	// Parse the file
	node, err := parser.ParseFile(p.context.FileSet, filePath, nil, p.config.Mode)
	if err != nil {
		p.context.AddError(filePath, 0, 0, err.Error(), "error")
		return nil, err
	}

	// Parse providers from the file
	providers, err := p.parseFileContent(filePath, node)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if p.config.CacheEnabled {
		p.context.Cache[filePath] = providers
	}

	p.context.FilesProcessed++
	p.context.ProvidersFound += len(providers)

	return providers, nil
}

// ParseDir implements Parser.ParseDir
func (p *GoParser) ParseDir(dirPath string) ([]Provider, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var allProviders []Provider

	// Walk through directory
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			// Skip hidden directories and common build/dependency directories
			if strings.HasPrefix(info.Name(), ".") ||
				info.Name() == "node_modules" ||
				info.Name() == "vendor" ||
				info.Name() == "testdata" {
				return filepath.SkipDir
			}
			return nil
		}

		// Parse Go files
		if filepath.Ext(path) == ".go" && p.context.ShouldIncludeFile(path) {
			providers, err := p.ParseFile(path)
			if err != nil {
				log.Warnf("Failed to parse file %s: %v", path, err)
				// Continue with other files
				return nil
			}
			allProviders = append(allProviders, providers...)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return allProviders, nil
}

// ParseString implements Parser.ParseString
func (p *GoParser) ParseString(code string) ([]Provider, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Parse the code string
	node, err := parser.ParseFile(p.context.FileSet, "", strings.NewReader(code), p.config.Mode)
	if err != nil {
		return nil, err
	}

	// Parse providers from the AST
	return p.parseFileContent("<string>", node)
}

// SetConfig implements Parser.SetConfig
func (p *GoParser) SetConfig(config *ParserConfig) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.config = config
	p.context = NewParserContext(config)
}

// GetConfig implements Parser.GetConfig
func (p *GoParser) GetConfig() *ParserConfig {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.config
}

// GetContext implements Parser.GetContext
func (p *GoParser) GetContext() *ParserContext {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.context
}

// parseFileContent parses providers from a parsed AST node
func (p *GoParser) parseFileContent(filePath string, node *ast.File) ([]Provider, error) {
	// Extract imports
	imports := make(map[string]string)
	for _, imp := range node.Imports {
		name := ""
		pkgPath := strings.Trim(imp.Path.Value, "\"")

		if imp.Name != nil {
			name = imp.Name.Name
		} else {
			name = gomod.GetPackageModuleName(pkgPath)
		}

		// Handle anonymous imports
		if name == "_" {
			name = gomod.GetPackageModuleName(pkgPath)
			// Handle duplicates
			if _, ok := imports[name]; ok {
				name = fmt.Sprintf("%s%d", name, rand.Intn(100))
			}
		}

		imports[name] = pkgPath
		p.context.AddImport(name, pkgPath)
	}

	var providers []Provider

	// Parse providers from declarations
	for _, decl := range node.Decls {
		provider, err := p.parseProviderDecl(filePath, node, decl, imports)
		if err != nil {
			p.context.AddError(filePath, 0, 0, err.Error(), "warning")
			continue
		}

		if provider != nil {
			providers = append(providers, *provider)
		}
	}

	return providers, nil
}

// parseProviderDecl parses a provider from an AST declaration
func (p *GoParser) parseProviderDecl(filePath string, fileNode *ast.File, decl ast.Decl, imports map[string]string) (*Provider, error) {
	genDecl, ok := decl.(*ast.GenDecl)
	if !ok {
		return nil, nil
	}

	if len(genDecl.Specs) == 0 {
		return nil, nil
	}

	typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec)
	if !ok {
		return nil, nil
	}

	// Check if it's a struct type
	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return nil, nil
	}

	// Check for provider annotation
	if genDecl.Doc == nil || len(genDecl.Doc.List) == 0 {
		return nil, nil
	}

	docText := strings.TrimLeft(genDecl.Doc.List[len(genDecl.Doc.List)-1].Text, "/ \t")
	if !strings.HasPrefix(docText, "@provider") {
		return nil, nil
	}

	// Parse provider annotation
	providerDoc := parseProvider(docText)

	// Create provider struct
	provider := &Provider{
		StructName:    typeSpec.Name.Name,
		ReturnType:    providerDoc.ReturnType,
		Mode:          ProviderModeBasic,
		ProviderGroup: providerDoc.Group,
		InjectParams:  make(map[string]InjectParam),
		Imports:       make(map[string]string),
		PkgName:       fileNode.Name.Name,
		ProviderFile:  filepath.Join(filepath.Dir(filePath), "provider.gen.go"),
	}

	// Set default return type if not specified
	if provider.ReturnType == "" {
		provider.ReturnType = "*" + provider.StructName
	}

	// Parse provider mode
	if providerDoc.Mode != "" {
		if IsValidProviderMode(providerDoc.Mode) {
			provider.Mode = ProviderMode(providerDoc.Mode)
		} else {
			return nil, fmt.Errorf("invalid provider mode: %s", providerDoc.Mode)
		}
	}

	// Parse struct fields for injection
	if err := p.parseStructFields(structType, imports, provider, providerDoc.IsOnly); err != nil {
		return nil, err
	}

	// Handle special provider modes
	p.handleProviderModes(provider, providerDoc.Mode, imports)

	// Add source location if enabled
	if p.config.SourceLocations {
		if genDecl.Doc != nil && len(genDecl.Doc.List) > 0 {
			position := p.context.FileSet.Position(genDecl.Doc.List[0].Pos())
			provider.Location = SourceLocation{
				File:   position.Filename,
				Line:   position.Line,
				Column: position.Column,
			}
		}
	}

	return provider, nil
}

// parseStructFields parses struct fields for injection parameters
func (p *GoParser) parseStructFields(structType *ast.StructType, imports map[string]string, provider *Provider, onlyMode bool) error {
	for _, field := range structType.Fields.List {
		if field.Names == nil {
			continue
		}

		// Check for struct tags
		if field.Tag != nil {
			provider.NeedPrepareFunc = true
		}

		// Check injection mode
		shouldInject := true
		if onlyMode {
			shouldInject = field.Tag != nil && strings.Contains(field.Tag.Value, `inject:"true"`)
		} else {
			shouldInject = field.Tag == nil || !strings.Contains(field.Tag.Value, `inject:"false"`)
		}

		if !shouldInject {
			continue
		}

		// Parse field type
		star, pkg, pkgAlias, typ, err := p.parseFieldType(field.Type, imports)
		if err != nil {
			continue
		}

		// Skip scalar types
		if lo.Contains(scalarTypes, typ) {
			continue
		}

		// Add injection parameter
		for _, name := range field.Names {
			provider.InjectParams[name.Name] = InjectParam{
				Star:         star,
				Type:         typ,
				Package:      pkg,
				PackageAlias: pkgAlias,
			}

			// Add to imports
			if pkg != "" && pkgAlias != "" {
				provider.Imports[pkg] = pkgAlias
			}
		}
	}

	return nil
}

// parseFieldType parses a field type and returns its components
func (p *GoParser) parseFieldType(expr ast.Expr, imports map[string]string) (star, pkg, pkgAlias, typ string, err error) {
	switch t := expr.(type) {
	case *ast.Ident:
		typ = t.Name
	case *ast.StarExpr:
		star = "*"
		return p.parseFieldType(t.X, imports)
	case *ast.SelectorExpr:
		if x, ok := t.X.(*ast.Ident); ok {
			pkgAlias = x.Name
			if path, ok := imports[pkgAlias]; ok {
				pkg = path
			}
			typ = t.Sel.Name
		}
	default:
		return "", "", "", "", errors.New("unsupported field type")
	}

	return star, pkg, pkgAlias, typ, nil
}

// handleProviderModes applies special handling for different provider modes
func (p *GoParser) handleProviderModes(provider *Provider, mode string, imports map[string]string) {
	moduleName := gomod.GetModuleName()

	switch provider.Mode {
	case ProviderModeGrpc:
		modePkg := moduleName + "/providers/grpc"
		provider.ProviderGroup = "atom.GroupInitial"
		// Save the original return type before changing it
		originalReturnType := provider.ReturnType
		provider.GrpcRegisterFunc = originalReturnType
		provider.ReturnType = "contracts.Initial"

		provider.Imports[atomPackage("")] = ""
		provider.Imports[atomPackage("contracts")] = ""
		provider.Imports[modePkg] = ""

		// Extract and add gRPC service package import
		if originalReturnType != "" && strings.Contains(originalReturnType, ".") {
			// Extract package alias from gRPC register function name (e.g., "userv1" from "userv1.RegisterUserServiceServer")
			if pkgAlias := getTypePkgName(originalReturnType); pkgAlias != "" {
				// Look for this package in the original file's imports
				if importPath, exists := imports[pkgAlias]; exists {
					// Use the exact import path from the original file
					provider.Imports[importPath] = pkgAlias
				} else {
					// Fallback: try to infer the common pattern
					if moduleName != "" {
						// Common pattern: {module}/pkg/proto/{service}/v1
						servicePkg := moduleName + "/pkg/proto/" + strings.ToLower(pkgAlias)
						provider.Imports[servicePkg] = pkgAlias
					}
				}
			}
		}

		provider.InjectParams["__grpc"] = InjectParam{
			Star:         "*",
			Type:         "Grpc",
			Package:      modePkg,
			PackageAlias: "grpc",
		}

	case ProviderModeEvent:
		modePkg := moduleName + "/providers/event"
		provider.ProviderGroup = "atom.GroupInitial"
		provider.ReturnType = "contracts.Initial"

		provider.Imports[atomPackage("")] = ""
		provider.Imports[atomPackage("contracts")] = ""
		provider.Imports[modePkg] = ""

		provider.InjectParams["__event"] = InjectParam{
			Star:         "*",
			Type:         "PubSub",
			Package:      modePkg,
			PackageAlias: "event",
		}

	case ProviderModeJob, ProviderModeCronJob:
		modePkg := moduleName + "/providers/job"
		provider.ProviderGroup = "atom.GroupInitial"
		provider.ReturnType = "contracts.Initial"

		provider.Imports[atomPackage("")] = ""
		provider.Imports[atomPackage("contracts")] = ""
		provider.Imports["github.com/riverqueue/river"] = ""
		provider.Imports[modePkg] = ""

		provider.InjectParams["__job"] = InjectParam{
			Star:         "*",
			Type:         "Job",
			Package:      modePkg,
			PackageAlias: "job",
		}

	case ProviderModeModel:
		provider.ProviderGroup = "atom.GroupInitial"
		provider.ReturnType = "contracts.Initial"
		provider.NeedPrepareFunc = true
	}

	// Handle return type and group package imports
	if pkgAlias := getTypePkgName(provider.ReturnType); pkgAlias != "" {
		if importPkg, ok := p.context.Imports[pkgAlias]; ok {
			provider.Imports[importPkg] = pkgAlias
		}
	}

	if pkgAlias := getTypePkgName(provider.ProviderGroup); pkgAlias != "" {
		if importPkg, ok := p.context.Imports[pkgAlias]; ok {
			provider.Imports[importPkg] = pkgAlias
		}
	}
}
