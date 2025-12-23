package provider

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"go.ipao.vip/atomctl/v2/pkg/utils/gomod"
)

// MainParser 重构后的主解析器，采用组合模式协调各个子组件完成 Provider 解析
//
// 架构设计：
// ┌─────────────────────────────────────────────────────────────┐
// │                    MainParser (协调器)                        │
// ├─────────────────────────────────────────────────────────────┤
// │  commentParser  │  importResolver  │  astWalker               │
// │  (注释解析器)   │   (导入解析器)   │  (AST遍历器)            │
// ├─────────────────────────────────────────────────────────────┤
// │  builder        │  validator      │  config                 │
// │  (构建器)      │   (验证器)      │  (配置管理)             │
// └─────────────────────────────────────────────────────────────┘
//
// 执行流程：
// 1. 文件过滤 → 2. AST解析 → 3. 导入处理 → 4. 注解发现 → 5. Provider构建 → 6. 验证
type MainParser struct {
	commentParser  *CommentParser   // 负责解析 @provider 注解
	importResolver *ImportResolver  // 负责处理 Go 文件的导入信息
	astWalker      *ASTWalker       // 负责遍历 AST 发现 Provider 注解
	builder        *ProviderBuilder // 负责从 AST 节点构建 Provider 对象
	validator      *GoValidator     // 负责验证 Provider 配置的正确性
	config         *ParserConfig    // 负责解析器配置管理
}

// NewParser 创建一个使用默认配置的 MainParser 实例
//
// 初始化流程：
// ┌─────────────────────────────────────────────────────────────┐
// │                    NewParser()                              │
// ├─────────────────────────────────────────────────────────────┤
// │  创建各个子组件实例：                                        │
// │  - CommentParser: 解析 @provider 注释                     │
// │  - ImportResolver: 处理导入依赖                            │
// │  - ASTWalker: 遍历 AST 发现结构体                          │
// │  - ProviderBuilder: 构建 Provider 对象                    │
// │  - GoValidator: 验证 Provider 配置                        │
// │  - ParserConfig: 默认配置                                  │
// └─────────────────────────────────────────────────────────────┘
//
// 返回值：配置好的 MainParser 实例，可直接调用 ParseFile() 或 ParseDir()
func NewParser() *MainParser {
	return &MainParser{
		commentParser:  NewCommentParser(),   // 初始化注释解析器
		importResolver: NewImportResolver(),  // 初始化导入解析器
		astWalker:      NewASTWalker(),       // 初始化 AST 遍历器
		builder:        NewProviderBuilder(), // 初始化 Provider 构建器
		validator:      NewGoValidator(),     // 初始化验证器
		config:         NewParserConfig(),    // 初始化默认配置
	}
}

// NewParserWithConfig creates a new MainParser with custom configuration
func NewParserWithConfig(config *ParserConfig) *MainParser {
	if config == nil {
		return NewParser()
	}

	return &MainParser{
		commentParser:  NewCommentParser(),
		importResolver: NewImportResolver(),
		astWalker: NewASTWalkerWithConfig(&WalkerConfig{
			StrictMode: config.StrictMode,
		}),
		builder: NewProviderBuilderWithConfig(&BuilderConfig{
			EnableValidation:          config.StrictMode,
			StrictMode:                config.StrictMode,
			DefaultProviderMode:       ProviderModeBasic,
			DefaultInjectionMode:      InjectionModeAuto,
			AutoGenerateReturnTypes:   true,
			ResolveImportDependencies: true,
		}),
		validator: NewGoValidator(),
		config:    config,
	}
}

// ParseRefactored 重构后的单文件解析函数 - 替代原有的 Parse 函数
//
// 执行流程：
// ┌─────────────────────────────────────────────────────────────┐
// │                 ParseRefactored(source)                     │
// ├─────────────────────────────────────────────────────────────┤
// │  1. 创建解析器实例：NewParser()                            │
// │  2. 调用文件解析：parser.ParseFile(source)                   │
// │  3. 错误处理：如果解析失败，记录错误并返回空切片            │
// │  4. 返回结果：返回发现的 Provider 切片                    │
// └─────────────────────────────────────────────────────────────┘
//
// 参数：
//   - source: Go 源文件路径
//
// 返回值：
//   - []Provider: 解析到的 Provider 列表（如果失败则为空）
//
// 设计原则：
//   - 简化接口：提供便捷的单函数调用方式
//   - 错误处理：内部处理错误，避免调用者需要处理复杂错误
//   - 兼容性：保持与原有 Parse 函数相同的签名，方便迁移
func ParseRefactored(source string) []Provider {
	// 创建解析器实例
	parser := NewParser()

	// 调用详细的文件解析方法
	providers, err := parser.ParseFile(source)
	// 错误处理：记录错误并返回空结果
	if err != nil {
		log.Error("Parse error: ", err)
		return []Provider{}
	}

	// 返回解析结果
	return providers
}

// ParseFile 单文件解析的核心方法 - 详细解析 Go 源文件中的 Provider 注解
//
// 执行流程图：
// ┌─────────────────────────────────────────────────────────────┐
// │                 ParseFile(source)                           │
// ├─────────────────────────────────────────────────────────────┤
// │  1. 文件过滤：shouldProcessFile()                          │
// │  2. AST解析：parser.ParseFile()                            │
// │  3. 创建上下文：NewParserContext()                         │
// │  4. 导入解析：ResolveFileImports()                         │
// │  5. 构建上下文：BuilderContext{...}                        │
// │  6. 注解发现：ProviderDiscoveryVisitor                     │
// │  7. AST遍历：astWalker.WalkFile()                          │
// │  8. Provider构建：buildProviderFromDiscovery()             │
// │  9. 验证检查：validator.Validate()                        │
// │ 10. 日志记录：记录警告和错误                               │
// └─────────────────────────────────────────────────────────────┘
//
// 详细步骤说明：
// 步骤 1-2：文件预处理
//   - 检查文件是否应该被处理（跳过测试文件和生成文件）
//   - 使用 Go 标准库解析文件为 AST
//
// 步骤 3-5：上下文准备
//   - 创建解析器上下文，包含工作目录和模块名
//   - 解析文件的所有导入信息，建立包名映射
//   - 创建构建器上下文，包含解析所需的所有信息
//
// 步骤 6-7：注解发现
//   - 创建专门的访问器来发现 @provider 注解
//   - 遍历 AST，收集所有带有 @provider 注解的结构体
//
// 步骤 8-9：Provider 构建
//   - 对每个发现的注解，构建完整的 Provider 对象
//   - 如果启用严格模式，验证 Provider 配置的正确性
//
// 步骤 10：结果处理
//   - 记录解析过程中的所有警告和错误
//   - 返回构建成功的 Provider 列表
func (p *MainParser) ParseFile(source string) ([]Provider, error) {
	// === 步骤 1：文件过滤 ===
	// 检查文件是否应该被处理，跳过测试文件和生成文件
	if !p.shouldProcessFile(source) {
		return []Provider{}, nil
	}

	// === 步骤 2：AST 解析 ===
	// 使用 Go 标准库将源文件解析为抽象语法树
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, source, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", source, err)
	}

	// === 步骤 3：创建解析器上下文 ===
	// 创建包含工作目录和模块信息的上下文
	context := NewParserContext(p.config)
	context.WorkingDir = filepath.Dir(source)
	context.ModuleName = gomod.GetModuleName()

	// === 步骤 4：导入解析 ===
	// 解析文件的所有导入信息，建立包名到路径的映射
	importContext, err := p.importResolver.ResolveFileImports(node, source)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve imports: %w", err)
	}

	// === 步骤 5：创建构建器上下文 ===
	// 创建包含解析所需所有信息的构建器上下文
	builderContext := &BuilderContext{
		FilePath:       source,                // 当前文件路径
		PackageName:    node.Name.Name,        // 包名
		ImportContext:  importContext,         // 导入信息上下文
		ASTFile:        node,                  // AST 节点
		ProcessedTypes: make(map[string]bool), // 已处理的类型，避免重复
		Errors:         make([]error, 0),      // 错误列表
		Warnings:       make([]string, 0),     // 警告列表
	}

	// === 步骤 6：创建注解发现访问器 ===
	// 创建专门的访问器来发现 @provider 注解
	visitor := NewProviderDiscoveryVisitor(p.commentParser)
	p.astWalker.AddVisitor(visitor)

	// === 步骤 7：AST 遍历 ===
	// 遍历 AST，发现所有带有 @provider 注解的结构体
	if err := p.astWalker.WalkFile(source); err != nil {
		return nil, fmt.Errorf("failed to walk AST: %w", err)
	}

	// === 步骤 8：Provider 构建和验证 ===
	// 初始化结果列表
	providers := make([]Provider, 0)
	discoveredProviders := visitor.GetProviders()

	// 对每个发现的注解构建 Provider 对象
	for _, discoveredProvider := range discoveredProviders {
		// 查找对应的 AST 节点
		provider, err := p.buildProviderFromDiscovery(discoveredProvider, node, builderContext)
		if err != nil {
			context.AddError(
				source,
				0,
				0,
				fmt.Sprintf("failed to build provider %s: %v", discoveredProvider.StructName, err),
				"error",
			)
			continue
		}

		// 如果启用严格模式，验证 Provider 配置
		if p.config.StrictMode {
			if err := p.validator.Validate(&provider); err != nil {
				context.AddError(
					source,
					0,
					0,
					fmt.Sprintf("validation failed for provider %s: %v", provider.StructName, err),
					"error",
				)
				continue
			}
		}

		// 添加到结果列表
		providers = append(providers, provider)
	}

	// === 步骤 9：日志记录 ===
	// 记录解析过程中的所有警告和错误
	for _, parseErr := range context.GetErrors("warning") {
		log.Warnf("Warning while parsing %s: %s", source, parseErr.Message)
	}
	for _, parseErr := range context.GetErrors("error") {
		log.Errorf("Error while parsing %s: %s", source, parseErr.Message)
	}

	// 返回构建成功的 Provider 列表
	return providers, nil
}

// ParseDir parses all Go files in a directory and returns discovered providers
func (p *MainParser) ParseDir(dir string) ([]Provider, error) {
	var allProviders []Provider

	// Use AST walker to traverse the directory
	if err := p.astWalker.WalkDir(dir); err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	// Note: This would need to be enhanced to collect providers from all files
	// For now, we'll return an empty slice and log a warning
	log.Warn("ParseDir not fully implemented yet")
	return allProviders, nil
}

// shouldProcessFile determines if a file should be processed
func (p *MainParser) shouldProcessFile(source string) bool {
	// Skip test files
	if strings.HasSuffix(source, "_test.go") {
		return false
	}

	// Skip generated provider files
	if strings.HasSuffix(source, "provider.gen.go") {
		return false
	}

	return true
}

// buildProviderFromDiscovery builds a complete Provider from a discovered provider annotation
func (p *MainParser) buildProviderFromDiscovery(
	discoveredProvider Provider,
	node *ast.File,
	context *BuilderContext,
) (Provider, error) {
	// Find the corresponding type specification in the AST
	var typeSpec *ast.TypeSpec
	var genDecl *ast.GenDecl

	for _, decl := range node.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		if len(gd.Specs) == 0 {
			continue
		}

		ts, ok := gd.Specs[0].(*ast.TypeSpec)
		if !ok {
			continue
		}

		if ts.Name.Name == discoveredProvider.StructName {
			typeSpec = ts
			genDecl = gd
			break
		}
	}

	if typeSpec == nil {
		return Provider{}, fmt.Errorf("type specification not found for %s", discoveredProvider.StructName)
	}

	// Use the builder to construct the complete provider
	provider, err := p.builder.BuildFromTypeSpec(typeSpec, genDecl, context)
	if err != nil {
		return Provider{}, fmt.Errorf("failed to build provider: %w", err)
	}

	// Apply legacy compatibility transformations
	result, err := p.applyLegacyCompatibility(&provider)
	if err != nil {
		return Provider{}, fmt.Errorf("failed to apply legacy compatibility: %w", err)
	}

	return result, nil
}

// applyLegacyCompatibility applies transformations to maintain backward compatibility
func (p *MainParser) applyLegacyCompatibility(provider *Provider) (Provider, error) {
	// Set provider file path
	provider.ProviderFile = filepath.Join(filepath.Dir(provider.ProviderFile), "provider.gen.go")

	// Apply mode-specific transformations based on the original logic
	switch provider.Mode {
	case ProviderModeGrpc:
		p.applyGrpcCompatibility(provider)
	case ProviderModeEvent:
		p.applyEventCompatibility(provider)
	case ProviderModeJob, ProviderModeCronJob:
		p.applyJobCompatibility(provider)
	case ProviderModeModel:
		p.applyModelCompatibility(provider)
	}

	return *provider, nil
}

// applyGrpcCompatibility applies gRPC-specific compatibility transformations
func (p *MainParser) applyGrpcCompatibility(provider *Provider) {
	modePkg := gomod.GetModuleName() + "/providers/grpc"

	// Add required imports
	provider.Imports[createAtomPackage("")] = ""
	provider.Imports[createAtomPackage("contracts")] = ""
	provider.Imports[modePkg] = ""

	// Set provider group
	if provider.ProviderGroup == "" {
		provider.ProviderGroup = "atom.GroupInitial"
	}

	// Set return type and register function
	// Important: Save the original return type before setting it to contracts.Initial
	originalReturnType := provider.ReturnType
	provider.ReturnType = "contracts.Initial"

	// Set gRPC register function name if not already set
	if provider.GrpcRegisterFunc == "" {
		if originalReturnType != "" && strings.Contains(originalReturnType, ".") {
			// User specified a complete register function name, like userv1.RegisterUserServiceServer
			provider.GrpcRegisterFunc = originalReturnType
		} else {
			// Generate default gRPC register function name
			// Example: UserService -> RegisterUserServiceServer
			provider.GrpcRegisterFunc = "Register" + strings.TrimPrefix(originalReturnType, "*") + "Server"
		}
	}

	// Note: Package import handling for gRPC register functions is now done
	// in resolveImportDependencies to ensure access to original file imports

	// Add gRPC injection parameter
	provider.InjectParams["__grpc"] = InjectParam{
		Star:         "*",
		Type:         "Grpc",
		Package:      modePkg,
		PackageAlias: "grpc",
	}
}

// applyEventCompatibility applies event-specific compatibility transformations
func (p *MainParser) applyEventCompatibility(provider *Provider) {
	modePkg := gomod.GetModuleName() + "/providers/event"

	// Add required imports
	provider.Imports[createAtomPackage("")] = ""
	provider.Imports[createAtomPackage("contracts")] = ""
	provider.Imports[modePkg] = ""

	// Set provider group
	if provider.ProviderGroup == "" {
		provider.ProviderGroup = "atom.GroupInitial"
	}

	// Set return type
	provider.ReturnType = "contracts.Initial"

	// Add event injection parameter
	provider.InjectParams["__event"] = InjectParam{
		Star:         "*",
		Type:         "PubSub",
		Package:      modePkg,
		PackageAlias: "event",
	}
}

// applyJobCompatibility applies job-specific compatibility transformations
func (p *MainParser) applyJobCompatibility(provider *Provider) {
	modePkg := gomod.GetModuleName() + "/providers/job"

	// Add required imports
	provider.Imports[createAtomPackage("")] = ""
	provider.Imports[createAtomPackage("contracts")] = ""
	provider.Imports["github.com/riverqueue/river"] = ""
	provider.Imports[modePkg] = ""

	// Set provider group
	if provider.ProviderGroup == "" {
		provider.ProviderGroup = "atom.GroupInitial"
	}

	// Set return type
	provider.ReturnType = "contracts.Initial"

	// Add job injection parameter
	provider.InjectParams["__job"] = InjectParam{
		Star:         "*",
		Type:         "Job",
		Package:      modePkg,
		PackageAlias: "job",
	}
}

// applyModelCompatibility applies model-specific compatibility transformations
func (p *MainParser) applyModelCompatibility(provider *Provider) {
	// Set provider group
	if provider.ProviderGroup == "" {
		provider.ProviderGroup = "atom.GroupInitial"
	}

	// Set return type
	provider.ReturnType = "contracts.Initial"

	// Ensure prepare function is needed
	provider.NeedPrepareFunc = true
}

// GetCommentParser returns the comment parser used by this parser
func (p *MainParser) GetCommentParser() *CommentParser {
	return p.commentParser
}

// GetImportResolver returns the import resolver used by this parser
func (p *MainParser) GetImportResolver() *ImportResolver {
	return p.importResolver
}

// GetASTWalker returns the AST walker used by this parser
func (p *MainParser) GetASTWalker() *ASTWalker {
	return p.astWalker
}

// GetBuilder returns the provider builder used by this parser
func (p *MainParser) GetBuilder() *ProviderBuilder {
	return p.builder
}

// GetValidator returns the validator used by this parser
func (p *MainParser) GetValidator() *GoValidator {
	return p.validator
}

// GetConfig returns the parser configuration
func (p *MainParser) GetConfig() *ParserConfig {
	return p.config
}

// SetConfig updates the parser configuration
func (p *MainParser) SetConfig(config *ParserConfig) {
	p.config = config
}

// Helper function to create atom package paths
func createAtomPackage(suffix string) string {
	root := "go.ipao.vip/atom"
	if suffix != "" {
		return fmt.Sprintf("%s/%s", root, suffix)
	}
	return root
}
