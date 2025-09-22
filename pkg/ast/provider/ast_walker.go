package provider

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// ASTWalker 处理 Go 抽象语法树（AST）的遍历器，专门用于发现 Provider 相关的结构体
//
// 架构设计：
// ┌─────────────────────────────────────────────────────────────┐
// │                    ASTWalker                                │
// ├─────────────────────────────────────────────────────────────┤
// │  fileSet:       *token.FileSet   - 文件集，记录源码位置信息    │
// │  commentParser: *CommentParser   - 注释解析器，处理注解      │
// │  config:        *WalkerConfig    - 遍历器配置               │
// │  visitors:      []NodeVisitor    - 访问者列表，支持扩展     │
// └─────────────────────────────────────────────────────────────┘
//
// 核心功能：
//   - 文件过滤：根据配置跳过测试文件、生成文件等
//   - AST 解析：使用 Go 标准库解析源文件为抽象语法树
//   - 结构化遍历：按照 AST 层次结构进行遍历
//   - 访问者模式：支持多个访问者同时处理不同的分析任务
//   - 错误处理：提供完善的错误处理机制
//
// 执行流程：
//  1. 文件过滤 → 2. AST 解析 → 3. 遍历声明 → 4. 访问者通知 → 5. 结果收集
//
// 设计原则：
//   - 单一职责：专注于 AST 遍历，不涉及具体的业务逻辑
//   - 开闭原则：通过访问者模式支持功能扩展，无需修改核心代码
//   - 接口隔离：定义清晰的访问者接口，避免过度依赖
//   - 配置驱动：通过配置控制遍历行为，提高灵活性
//
// 使用场景：
//   - Provider 注解发现：查找带有 @provider 注解的结构体
//   - 代码分析：分析代码结构、依赖关系等
//   - 代码生成：基于 AST 分析结果生成代码
//   - 重构工具：基于 AST 进行代码重构操作
type ASTWalker struct {
	fileSet       *token.FileSet
	commentParser *CommentParser
	config        *WalkerConfig
	visitors      []NodeVisitor
}

// WalkerConfig 配置 AST 遍历器的行为参数
//
// 配置选项：
//   - IncludeTestFiles: 是否包含测试文件（_test.go 后缀）
//   - IncludeGeneratedFiles: 是否包含生成的文件（.gen.go 后缀）
//   - MaxFileSize: 最大文件大小限制（字节），防止大文件影响性能
//   - StrictMode: 严格模式，启用更严格的验证和错误检查
//
// 使用场景：
//   - 开发环境：包含测试文件，禁用严格模式，便于调试
//   - 生产环境：跳过测试文件，启用严格模式，提高解析质量
//   - 性能优化：设置合理的文件大小限制，跳过不必要的文件
type WalkerConfig struct {
	IncludeTestFiles      bool
	IncludeGeneratedFiles bool
	MaxFileSize           int64
	StrictMode            bool
}

// NodeVisitor 定义 AST 节点访问者接口，支持在遍历过程中处理不同类型的节点
//
// 接口设计原则：
// ┌─────────────────────────────────────────────────────────────┐
// │                  NodeVisitor Interface                       │
// ├─────────────────────────────────────────────────────────────┤
// │  VisitFile:      文件级别处理，在文件开始时调用              │
// │  VisitGenDecl:   通用声明处理，类型、变量、常量声明         │
// │  VisitTypeSpec:  类型规范处理，具体的类型定义                │
// │  VisitStructType: 结构体类型处理，结构体类型定义            │
// │  VisitStructField: 结构体字段处理，结构体中的字段            │
// │  Complete:        完成处理，在文件处理完成时调用              │
// └─────────────────────────────────────────────────────────────┘
//
// 调用时机：
//  1. VisitFile: 开始处理文件时调用，用于初始化文件级上下文
//  2. VisitGenDecl: 遇到通用声明时调用（type、var、const）
//  3. VisitTypeSpec: 遇到类型规范时调用（具体的类型定义）
//  4. VisitStructType: 遇到结构体类型时调用
//  5. VisitStructField: 遍历结构体字段时调用
//  6. Complete: 文件处理完成时调用，用于清理和结果收集
//
// 错误处理：
//   - 如果任何方法返回错误，遍历过程会立即停止
//   - 访问者应该根据错误类型决定是否继续处理
//   - 建议使用 fmt.Errorf 包装错误信息，包含足够的上下文
//
// 扩展性：
//   - 可以实现多个不同的访问者来处理不同的分析任务
//   - 访问者之间相互独立，不会相互影响
//   - 可以通过组合多个访问者来实现复杂的功能
type NodeVisitor interface {
	// VisitFile 在开始处理新文件时调用，用于文件级别的初始化
	//
	// 参数：
	//   - filePath: 当前处理的文件路径
	//   - node: 解析后的 AST 文件节点
	//
	// 返回值：
	//   - error: 如果返回错误，遍历过程会立即停止
	//
	// 使用场景：
	//   - 初始化文件级上下文信息
	//   - 记录文件级别的统计信息
	//   - 设置文件特定的配置
	VisitFile(filePath string, node *ast.File) error

	// VisitGenDecl 在遇到通用声明时调用（type、var、const）
	//
	// 参数：
	//   - filePath: 当前文件路径
	//   - decl: 通用声明节点
	//
	// 返回值：
	//   - error: 如果返回错误，遍历过程会立即停止
	//
	// 使用场景：
	//   - 处理包级别的变量和常量声明
	//   - 分析声明的文档注释
	//   - 收集声明级别的元数据
	VisitGenDecl(filePath string, decl *ast.GenDecl) error

	// VisitTypeSpec 在遇到类型规范时调用（具体的类型定义）
	//
	// 参数：
	//   - filePath: 当前文件路径
	//   - typeSpec: 类型规范节点
	//   - decl: 包含该类型规范的通用声明
	//
	// 返回值：
	//   - error: 如果返回错误，遍历过程会立即停止
	//
	// 使用场景：
	//   - 发现结构体、接口、函数类型等
	//   - 分析类型定义的注释和属性
	//   - 建立类型索引和关系图
	VisitTypeSpec(filePath string, typeSpec *ast.TypeSpec, decl *ast.GenDecl) error

	// VisitStructType 在遇到结构体类型时调用
	//
	// 参数：
	//   - filePath: 当前文件路径
	//   - structType: 结构体类型节点
	//   - typeSpec: 对应的类型规范
	//   - decl: 包含该结构体的通用声明
	//
	// 返回值：
	//   - error: 如果返回错误，遍历过程会立即停止
	//
	// 使用场景：
	//   - 分析结构体定义
	//   - 处理结构体的注解和属性
	//   - 收集结构体级别的元数据
	VisitStructType(filePath string, structType *ast.StructType, typeSpec *ast.TypeSpec, decl *ast.GenDecl) error

	// VisitStructField 在遍历结构体字段时调用
	//
	// 参数：
	//   - filePath: 当前文件路径
	//   - field: 结构体字段节点
	//   - structType: 所属的结构体类型
	//
	// 返回值：
	//   - error: 如果返回错误，遍历过程会立即停止
	//
	// 使用场景：
	//   - 分析字段类型和标签
	//   - 处理注入相关的标签（inject 标签）
	//   - 收集字段级别的依赖信息
	VisitStructField(filePath string, field *ast.Field, structType *ast.StructType) error

	// Complete 在文件处理完成时调用，用于清理和结果收集
	//
	// 参数：
	//   - filePath: 已完成处理的文件路径
	//
	// 返回值：
	//   - error: 如果返回错误，会影响整体处理结果
	//
	// 使用场景：
	//   - 清理文件级的临时数据
	//   - 收集和汇总处理结果
	//   - 执行文件级别的验证和检查
	Complete(filePath string) error
}

// NewASTWalker 使用默认配置创建新的 ASTWalker 实例
//
// 初始化流程：
// ┌─────────────────────────────────────────────────────────────┐
// │                  NewASTWalker()                            │
// ├─────────────────────────────────────────────────────────────┤
// │  1. 创建文件集：token.NewFileSet()                         │
// │  2. 创建注释解析器：NewCommentParser()                     │
// │  3. 设置默认配置：WalkerConfig{...}                       │
// │  4. 初始化访问者列表：make([]NodeVisitor, 0)              │
// │  5. 构建遍历器实例：&ASTWalker{...}                       │
// └─────────────────────────────────────────────────────────────┘
//
// 默认配置：
//   - IncludeTestFiles: false（跳过测试文件）
//   - IncludeGeneratedFiles: false（跳过生成的文件）
//   - MaxFileSize: 10MB（限制文件大小）
//   - StrictMode: false（禁用严格模式）
//
// 返回值：
//   - *ASTWalker: 配置好的遍历器实例，可以直接使用
//
// 使用示例：
//
//	walker := NewASTWalker()
//	visitor := NewProviderDiscoveryVisitor()
//	walker.AddVisitor(visitor)
//	err := walker.WalkFile("user_service.go")
//
// 注意事项：
//   - 必须添加至少一个访问者才能进行有效的分析
//   - 默认配置适用于大多数生产环境
//   - 可以通过 WithConfig 方法修改配置
func NewASTWalker() *ASTWalker {
	return &ASTWalker{
		fileSet:       token.NewFileSet(),
		commentParser: NewCommentParser(),
		config: &WalkerConfig{
			IncludeTestFiles:      false,
			IncludeGeneratedFiles: false,
			MaxFileSize:           10 * 1024 * 1024, // 10MB
			StrictMode:            false,
		},
		visitors: make([]NodeVisitor, 0),
	}
}

// NewASTWalkerWithConfig 使用自定义配置创建新的 ASTWalker 实例
//
// 设计目标：
//   - 提供灵活的配置选项
//   - 支持不同的使用场景（调试、生产、测试）
//   - 保持向后兼容性
//
// 参数处理：
//   - 如果 config 为 nil，自动使用默认配置
//   - 根据配置的 StrictMode 创建对应的注释解析器
//
// 自定义配置场景：
//   - 调试模式：包含测试文件，启用严格模式
//   - 性能优化：设置较小的文件大小限制
//   - 开发模式：包含生成的文件进行分析
//
// 参数：
//   - config: 自定义的遍历器配置（可为 nil）
//
// 返回值：
//   - *ASTWalker: 使用指定配置的遍历器实例
//
// 使用示例：
//
//	config := &WalkerConfig{
//	    IncludeTestFiles: true,
//	    IncludeGeneratedFiles: true,
//	    StrictMode: true,
//	    MaxFileSize: 5 * 1024 * 1024, // 5MB
//	}
//	walker := NewASTWalkerWithConfig(config)
//
// 注意事项：
//   - 自定义配置会覆盖所有默认设置
//   - 建议在特殊场景下使用自定义配置
//   - 确保 MaxFileSize 设置合理，避免性能问题
func NewASTWalkerWithConfig(config *WalkerConfig) *ASTWalker {
	if config == nil {
		return NewASTWalker()
	}

	return &ASTWalker{
		fileSet:       token.NewFileSet(),
		commentParser: NewCommentParserWithStrictMode(config.StrictMode),
		config:        config,
		visitors:      make([]NodeVisitor, 0),
	}
}

// AddVisitor 添加节点访问者到遍历器
//
// 功能说明：
//   - 将指定的访问者添加到访问者列表
//   - 访问者会按照添加的顺序被调用
//   - 支持添加多个访问者来处理不同的分析任务
//
// 使用场景：
//   - 添加 Provider 发现访问者
//   - 添加依赖分析访问者
//   - 添加代码质量检查访问者
//
// 参数：
//   - visitor: 要添加的节点访问者
//
// 注意事项：
//   - 同一个访问者实例只能添加一次
//   - 添加的访问者必须实现 NodeVisitor 接口
//   - 建议在遍历开始前添加所有需要的访问者
func (aw *ASTWalker) AddVisitor(visitor NodeVisitor) {
	aw.visitors = append(aw.visitors, visitor)
}

// RemoveVisitor 从遍历器中移除指定的节点访问者
//
// 功能说明：
//   - 从访问者列表中移除指定的访问者实例
//   - 使用对象引用进行比较，移除第一个匹配的访问者
//   - 如果访问者不存在，不做任何操作
//
// 使用场景：
//   - 动态调整分析任务
//   - 移除不再需要的访问者
//   - 优化性能，移除不必要的处理
//
// 参数：
//   - visitor: 要移除的节点访问者
//
// 注意事项：
//   - 如果有多个相同的访问者实例，只移除第一个
//   - 移除操作会改变访问者列表的顺序
//   - 建议在遍历开始前完成访问者的调整
func (aw *ASTWalker) RemoveVisitor(visitor NodeVisitor) {
	for i, v := range aw.visitors {
		if v == visitor {
			aw.visitors = append(aw.visitors[:i], aw.visitors[i+1:]...)
			break
		}
	}
}

// WalkFile 遍历单个 Go 文件的 AST，执行完整的分析和处理流程
//
// 执行流程：
// ┌─────────────────────────────────────────────────────────────┐
// │                 WalkFile(filePath)                         │
// ├─────────────────────────────────────────────────────────────┤
// │  1. 文件过滤：shouldProcessFile()                          │
// │  2. AST 解析：parser.ParseFile()                          │
// │  3. 文件开始：VisitFile() 通知所有访问者                   │
// │  4. AST 遍历：traverseFile() 递归遍历 AST                 │
// │  5. 文件完成：Complete() 通知所有访问者                   │
// │  6. 返回结果：成功返回 nil，失败返回错误                   │
// └─────────────────────────────────────────────────────────────┘
//
// 详细步骤说明：
// 步骤 1：文件过滤
//   - 检查文件扩展名是否为 .go
//   - 根据配置跳过测试文件和生成的文件
//   - 检查文件大小限制
//
// 步骤 2：AST 解析
//   - 使用 Go 标准库解析源文件
//   - 启用注释解析，保留所有注释信息
//   - 记录源码位置信息到文件集
//
// 步骤 3：文件开始通知
//   - 按顺序通知所有访问者文件开始处理
//   - 传递文件路径和 AST 节点
//   - 如果任何访问者返回错误，立即停止处理
//
// 步骤 4：AST 遍历
//   - 递归遍历 AST 中的所有声明
//   - 按照访问者接口定义的顺序通知访问者
//   - 处理类型声明、结构体、字段等
//
// 步骤 5：文件完成通知
//   - 按顺序通知所有访问者文件处理完成
//   - 访问者可以执行清理和结果收集
//   - 如果任何访问者返回错误，立即停止处理
//
// 参数：
//   - filePath: 要遍历的 Go 源文件路径
//
// 返回值：
//   - error: 成功时返回 nil，失败时返回包装后的错误信息
//
// 错误处理：
//   - 文件不存在或无法读取：返回文件系统错误
//   - 语法错误：返回解析错误，包含行号和列号信息
//   - 访问者错误：返回访问者产生的错误，停止后续处理
//
// 使用示例：
//
//	walker := NewASTWalker()
//	walker.AddVisitor(NewProviderDiscoveryVisitor())
//	err := walker.WalkFile("user_service.go")
//	if err != nil {
//	    log.Fatal("Failed to walk file: ", err)
//	}
//
// 注意事项：
//   - 必须至少添加一个访问者才能进行有效分析
//   - 文件路径必须是有效的 Go 源文件
//   - 错误会立即终止处理过程
func (aw *ASTWalker) WalkFile(filePath string) error {
	// === 步骤 1：文件过滤 ===
	// 检查文件是否符合处理条件，跳过不符合要求的文件
	if !aw.shouldProcessFile(filePath) {
		return nil
	}

	// === 步骤 2：AST 解析 ===
	// 使用 Go 标准库将源文件解析为抽象语法树，保留注释信息
	node, err := parser.ParseFile(aw.fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}

	// === 步骤 3：文件开始通知 ===
	// 按顺序通知所有访问者开始处理文件
	for _, visitor := range aw.visitors {
		if err := visitor.VisitFile(filePath, node); err != nil {
			return err
		}
	}

	// === 步骤 4：AST 遍历 ===
	// 递归遍历 AST 的所有节点，通知访问者处理各种声明
	if err := aw.traverseFile(filePath, node); err != nil {
		return err
	}

	// === 步骤 5：文件完成通知 ===
	// 按顺序通知所有访问者文件处理完成
	for _, visitor := range aw.visitors {
		if err := visitor.Complete(filePath); err != nil {
			return err
		}
	}

	return nil
}

// WalkDir 遍历目录中的所有 Go 文件，批量执行 AST 分析
//
// 执行流程：
// ┌─────────────────────────────────────────────────────────────┐
// │                 WalkDir(dirPath)                          │
// ├─────────────────────────────────────────────────────────────┤
// │  1. 目录遍历：filepath.Walk() 递归遍历                     │
// │  2. 目录过滤：跳过隐藏目录和依赖目录                       │
// │  3. 文件过滤：检查 Go 文件和配置条件                      │
// │  4. 文件处理：调用 WalkFile() 处理每个文件                 │
// │  5. 错误处理：记录错误但继续处理其他文件                   │
// └─────────────────────────────────────────────────────────────┘
//
// 目录过滤规则：
//   - 隐藏目录：以 . 开头的目录（.git、.idea 等）
//   - 依赖目录：node_modules、vendor、testdata
//   - 构建目录：通常被 gitignore 的目录
//
// 文件处理策略：
//   - 只处理 .go 后缀的文件
//   - 根据配置跳过测试文件和生成文件
//   - 单个文件处理失败不影响其他文件
//
// 错误处理策略：
//   - 目录访问错误：立即停止遍历
//   - 文件处理错误：记录警告但继续处理
//   - 使用标准输出打印警告信息
//
// 参数：
//   - dirPath: 要遍历的目录路径
//
// 返回值：
//   - error: 目录遍历错误，文件处理错误不会返回
//
// 使用场景：
//   - 项目分析：分析整个项目的 Provider 定义
//   - 批量处理：对项目中的所有文件进行统一处理
//   - 代码生成：基于项目分析结果生成代码
//
// 注意事项：
//   - 大型项目可能需要较长时间
//   - 建议先配置合适的过滤条件
//   - 错误处理采用宽松策略，不会因为单个文件失败而停止
func (aw *ASTWalker) WalkDir(dirPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		// 处理文件系统错误
		if err != nil {
			return err
		}

		// === 处理目录 ===
		if info.IsDir() {
			// 跳过隐藏目录和常见的依赖/构建目录
			if strings.HasPrefix(info.Name(), ".") ||
				info.Name() == "node_modules" ||
				info.Name() == "vendor" ||
				info.Name() == "testdata" {
				return filepath.SkipDir
			}
			return nil
		}

		// === 处理文件 ===
		// 只处理 Go 文件，且符合配置条件的文件
		if filepath.Ext(path) == ".go" && aw.shouldProcessFile(path) {
			if err := aw.WalkFile(path); err != nil {
				// 单个文件处理失败，记录警告但继续处理其他文件
				fmt.Printf("Warning: failed to process file %s: %v\n", path, err)
			}
		}

		return nil
	})
}

// traverseFile 遍历已解析文件的 AST，处理所有顶级声明
//
// 执行流程：
// ┌─────────────────────────────────────────────────────────────┐
// │              traverseFile(filePath, node)                  │
// ├─────────────────────────────────────────────────────────────┤
// │  1. 遍历声明：node.Decls 列表                              │
// │  2. 声明处理：traverseDeclaration()                       │
// │  3. 错误处理：任何声明处理失败都会停止整个遍历             │
// └─────────────────────────────────────────────────────────────┘
//
// 处理的声明类型：
//   - GenDecl: 通用声明（type、var、const）
//   - FuncDecl: 函数声明（当前版本跳过）
//   - 其他声明类型：根据需要进行扩展
//
// 递归策略：
//   - 深度优先遍历：先处理声明，再处理内部规范
//   - 顺序处理：按照源码中的顺序处理每个声明
//   - 错误传播：任何错误都会立即停止遍历过程
//
// 参数：
//   - filePath: 当前文件路径，用于错误报告
//   - node: 已解析的 AST 文件节点
//
// 返回值：
//   - error: 遍历过程中的错误
//
// 注意事项：
//   - 当前版本主要处理类型声明
//   - 函数声明会被跳过，可根据需要扩展
//   - 错误处理采用快速失败策略
func (aw *ASTWalker) traverseFile(filePath string, node *ast.File) error {
	// === 遍历所有顶级声明 ===
	// ast.File.Decls 包含文件中的所有顶级声明
	for _, decl := range node.Decls {
		// 递归处理每个声明，如果处理失败则立即停止
		if err := aw.traverseDeclaration(filePath, decl); err != nil {
			return err
		}
	}

	return nil
}

// traverseDeclaration 遍历单个声明，处理通用声明及其规范
//
// 执行流程：
// ┌─────────────────────────────────────────────────────────────┐
// │            traverseDeclaration(filePath, decl)             │
// ├─────────────────────────────────────────────────────────────┤
// │  1. 类型检查：decl.(*ast.GenDecl)                          │
// │  2. 访问者通知：VisitGenDecl() 通知访问者                  │
// │  3. 规范遍历：genDecl.Specs 遍历所有规范                   │
// │  4. 错误处理：任何步骤失败都会停止处理                     │
// └─────────────────────────────────────────────────────────────┘
//
// 处理的声明类型：
//   - ast.GenDecl: 通用声明，包含类型、变量、常量声明
//   - ast.FuncDecl: 函数声明（当前跳过）
//   - 其他声明类型：根据需要扩展
//
// 通用声明规范类型：
//   - TypeSpec: 类型规范（struct、interface、func type等）
//   - ValueSpec: 值规范（var、const声明）
//   - ImportSpec: 导入规范（import声明）
//
// 访问者通知策略：
//   - 声明级通知：先通知所有访问者处理声明
//   - 规范级通知：再递归处理声明中的规范
//   - 错误传播：任何访问者错误都会立即停止处理
//
// 参数：
//   - filePath: 当前文件路径，用于错误报告
//   - decl: AST 声明节点
//
// 返回值：
//   - error: 处理过程中的错误
//
// 注意事项：
//   - 当前版本主要处理 GenDecl 类型
//   - FuncDecl 类型会被跳过
//   - 错误处理采用快速失败策略
func (aw *ASTWalker) traverseDeclaration(filePath string, decl ast.Decl) error {
	// === 类型检查 ===
	// 只处理通用声明（type、var、const），跳过函数声明等
	genDecl, ok := decl.(*ast.GenDecl)
	if !ok {
		// 跳过函数声明和其他非通用声明
		return nil
	}

	// === 访问者通知 ===
	// 按顺序通知所有访问者处理通用声明
	for _, visitor := range aw.visitors {
		if err := visitor.VisitGenDecl(filePath, genDecl); err != nil {
			return err
		}
	}

	// === 规范遍历 ===
	// 递归处理声明中的所有规范（TypeSpec、ValueSpec、ImportSpec等）
	for _, spec := range genDecl.Specs {
		if err := aw.traverseSpec(filePath, spec, genDecl); err != nil {
			return err
		}
	}

	return nil
}

// traverseSpec 遍历声明中的规范，处理类型规范和结构体字段
//
// 执行流程：
// ┌─────────────────────────────────────────────────────────────┐
// │              traverseSpec(filePath, spec, decl)            │
// ├─────────────────────────────────────────────────────────────┤
// │  1. 类型检查：spec.(*ast.TypeSpec)                         │
// │  2. 访问者通知：VisitTypeSpec() 通知访问者                  │
// │  3. 结构体检查：typeSpec.Type.(*ast.StructType)           │
// │  4. 结构体通知：VisitStructType() 通知访问者                │
// │  5. 字段遍历：traverseStructFields() 遍历字段              │
// └─────────────────────────────────────────────────────────────┘
//
// 处理的规范类型：
//   - TypeSpec: 类型规范（当前主要处理）
//   - ValueSpec: 值规范（var、const声明，当前跳过）
//   - ImportSpec: 导入规范（import声明，当前跳过）
//
// 类型规范处理：
//   - 结构体类型：特别处理，遍历所有字段
//   - 接口类型：可根据需要扩展处理
//   - 函数类型：可根据需要扩展处理
//   - 其他类型：基本类型、数组、映射等
//
// 结构体处理策略：
//   - 类型识别：检查是否为 ast.StructType
//   - 访问者通知：通知所有访问者处理结构体
//   - 字段遍历：递归处理结构体的所有字段
//
// 参数：
//   - filePath: 当前文件路径，用于错误报告
//   - spec: AST 规范节点
//   - decl: 包含该规范的通用声明
//
// 返回值：
//   - error: 处理过程中的错误
//
// 注意事项：
//   - 当前版本主要处理 TypeSpec 和 StructType
//   - 其他规范类型会被跳过
//   - 结构体字段会进行深度遍历
func (aw *ASTWalker) traverseSpec(filePath string, spec ast.Spec, decl *ast.GenDecl) error {
	// === 类型检查 ===
	// 只处理类型规范，跳过值规范和导入规范
	typeSpec, ok := spec.(*ast.TypeSpec)
	if !ok {
		// 跳过非类型规范
		return nil
	}

	// === 访问者通知 ===
	// 按顺序通知所有访问者处理类型规范
	for _, visitor := range aw.visitors {
		if err := visitor.VisitTypeSpec(filePath, typeSpec, decl); err != nil {
			return err
		}
	}

	// === 结构体检查和处理 ===
	// 检查是否为结构体类型，如果是则进行特殊处理
	structType, ok := typeSpec.Type.(*ast.StructType)
	if ok {
		// === 结构体通知 ===
		// 按顺序通知所有访问者处理结构体类型
		for _, visitor := range aw.visitors {
			if err := visitor.VisitStructType(filePath, structType, typeSpec, decl); err != nil {
				return err
			}
		}

		// === 字段遍历 ===
		// 递归遍历结构体的所有字段
		if err := aw.traverseStructFields(filePath, structType); err != nil {
			return err
		}
	}

	return nil
}

// traverseStructFields 遍历结构体类型中的所有字段
//
// 执行流程：
// ┌─────────────────────────────────────────────────────────────┐
// │          traverseStructFields(filePath, structType)        │
// ├─────────────────────────────────────────────────────────────┤
// │  1. 空值检查：structType.Fields == nil                     │
// │  2. 字段遍历：structType.Fields.List 遍历所有字段          │
// │  3. 访问者通知：VisitStructField() 通知访问者              │
// │  4. 错误处理：任何字段处理失败都会停止遍历                 │
// └─────────────────────────────────────────────────────────────┘
//
// 字段遍历策略：
//   - 顺序遍历：按照源码中的顺序处理每个字段
//   - 批量处理：一个字段列表项可能包含多个同名字段
//   - 深度处理：递归处理字段的类型信息
//
// 字段类型分析：
//   - 基本类型：int、string、bool等
//   - 复合类型：数组、切片、映射、通道
//   - 自定义类型：用户定义的结构体、接口
//   - 指针类型：各种类型的指针
//
// 访问者处理内容：
//   - 字段名称和类型信息
//   - 结构体标签（tag）解析
//   - 注释和文档信息
//   - 注入相关的标签处理
//
// 参数：
//   - filePath: 当前文件路径，用于错误报告
//   - structType: 结构体类型节点
//
// 返回值：
//   - error: 遍历过程中的错误
//
// 注意事项：
//   - 空结构体会被跳过
//   - 字段处理失败会立即停止整个遍历
//   - 匿名字段会被正常处理
func (aw *ASTWalker) traverseStructFields(filePath string, structType *ast.StructType) error {
	// === 空值检查 ===
	// 如果结构体没有字段列表，直接返回
	if structType.Fields == nil {
		return nil
	}

	// === 字段遍历 ===
	// 遍历结构体的所有字段，一个 field.List 可能包含多个同名字段
	for _, field := range structType.Fields.List {
		// === 访问者通知 ===
		// 按顺序通知所有访问者处理结构体字段
		for _, visitor := range aw.visitors {
			if err := visitor.VisitStructField(filePath, field, structType); err != nil {
				return err
			}
		}
	}

	return nil
}

// shouldProcessFile 判断文件是否符合处理条件
//
// 过滤规则：
// ┌─────────────────────────────────────────────────────────────┐
// │                  shouldProcessFile()                        │
// ├─────────────────────────────────────────────────────────────┤
// │  1. 扩展名检查：必须是 .go 文件                             │
// │  2. 测试文件检查：根据配置决定是否处理 _test.go 文件        │
// │  3. 生成文件检查：根据配置决定是否处理 .gen.go 文件        │
// │  4. 文件大小检查：TODO 待实现（需要 os.Stat）             │
// └─────────────────────────────────────────────────────────────┘
//
// 过滤策略：
//   - 严格过滤：只处理符合所有条件的文件
//   - 配置驱动：根据 WalkerConfig 配置决定过滤规则
//   - 性能优化：尽早返回，避免不必要的检查
//
// 文件类型识别：
//   - 源文件：.go 后缀，非测试文件，非生成文件
//   - 测试文件：_test.go 后缀
//   - 生成文件：.gen.go 后缀
//   - 其他文件：非 .go 后缀的文件
//
// 配置影响：
//   - IncludeTestFiles: true 时处理测试文件
//   - IncludeGeneratedFiles: true 时处理生成文件
//   - MaxFileSize: 超过限制的文件会被跳过
//
// 参数：
//   - filePath: 要检查的文件路径
//
// 返回值：
//   - bool: true 表示应该处理，false 表示应该跳过
//
// 使用场景：
//   - WalkFile: 单文件处理前的检查
//   - WalkDir: 目录遍历中的文件过滤
//   - 批量处理：大规模文件处理时的优化
//
// 注意事项：
//   - 文件大小检查尚未实现
//   - 过滤规则是累积的，必须通过所有检查
//   - 配置修改会影响过滤行为
func (aw *ASTWalker) shouldProcessFile(filePath string) bool {
	// === 扩展名检查 ===
	// 只处理 Go 源文件，其他文件一律跳过
	if filepath.Ext(filePath) != ".go" {
		return false
	}

	// === 测试文件检查 ===
	// 如果配置不允许处理测试文件，跳过 _test.go 文件
	if !aw.config.IncludeTestFiles && strings.HasSuffix(filePath, "_test.go") {
		return false
	}

	// === 生成文件检查 ===
	// 如果配置不允许处理生成文件，跳过 .gen.go 文件
	if !aw.config.IncludeGeneratedFiles && strings.HasSuffix(filePath, ".gen.go") {
		return false
	}

	// === 文件大小检查 ===
	// TODO: 待实现文件大小检查功能（需要调用 os.Stat）
	// 需要考虑性能影响，避免过多的文件系统调用

	// 所有检查都通过，文件符合处理条件
	return true
}

// GetFileSet 返回遍历器使用的文件集
//
// 功能说明：
//   - 返回用于记录源码位置信息的文件集
//   - 文件集包含所有已解析文件的位置信息
//   - 可用于将位置信息转换为行号和列号
//
// 返回值：
//   - *token.FileSet: 文件集实例
//
// 使用场景：
//   - 错误报告：将 AST 位置转换为人类可读的位置
//   - 调试信息：显示源码位置的详细信息
//   - 工具集成：与其他需要位置信息的工具集成
//
// 注意事项：
//   - 文件集会在遍历过程中自动更新
//   - 返回的是实例引用，不是副本
//   - 多个文件的位置信息会累积在同一文件集中
func (aw *ASTWalker) GetFileSet() *token.FileSet {
	return aw.fileSet
}

// GetCommentParser 返回遍历器使用的注释解析器
//
// 功能说明：
//   - 返回用于解析注释的注释解析器实例
//   - 注释解析器负责处理 @provider 等注解
//   - 可用于独立的注释解析任务
//
// 返回值：
//   - *CommentParser: 注释解析器实例
//
// 使用场景：
//   - 注解分析：解析特定的注释格式
//   - 工具扩展：基于注释解析器实现自定义功能
//   - 调试和测试：测试注释解析功能
//
// 注意事项：
//   - 注释解析器的配置与遍历器保持一致
//   - 返回的是实例引用，不是副本
//   - 可以用于遍历器外部的注释解析任务
func (aw *ASTWalker) GetCommentParser() *CommentParser {
	return aw.commentParser
}

// GetConfig 返回遍历器的配置信息
//
// 功能说明：
//   - 返回遍历器的当前配置
//   - 配置包含所有影响遍历行为的参数
//   - 可用于检查和修改遍历器配置
//
// 返回值：
//   - *WalkerConfig: 遍历器配置指针
//
// 使用场景：
//   - 配置检查：查看当前的过滤和处理规则
//   - 动态调整：在运行时修改配置
//   - 配置复制：基于当前配置创建新的遍历器
//
// 注意事项：
//   - 返回的是配置的指针，修改会影响遍历器行为
//   - 建议在了解配置含义的情况下进行修改
//   - 配置修改不会影响已处理的文件

// ProviderDiscoveryVisitor 实现 NodeVisitor 接口，专门用于发现 Provider 注解
//
// 架构设计：
// ┌─────────────────────────────────────────────────────────────┐
// │              ProviderDiscoveryVisitor                       │
// ├─────────────────────────────────────────────────────────────┤
// │  commentParser: *CommentParser - 注释解析器                │
// │  providers:     []Provider     - 发现的 Provider 列表      │
// │  currentFile:   string        - 当前处理的文件路径        │
// └─────────────────────────────────────────────────────────────┘
//
// 核心功能：
//   - 注解识别：识别 @provider 注解
//   - Provider 构建：构建完整的 Provider 对象
//   - 结果收集：收集所有发现的 Provider
//   - 文件跟踪：跟踪当前处理的文件
//
// 处理策略：
//   - 结构体优先：只处理结构体类型的注解
//   - 注释解析：使用 CommentParser 解析注释块
//   - 默认值处理：为空值设置合理的默认值
//
// 使用场景：
//   - Provider 发现：在 AST 遍历过程中发现 Provider 注解
//   - 代码生成：为代码生成工具提供输入数据
//   - 项目分析：分析项目中的 Provider 定义
type ProviderDiscoveryVisitor struct {
	commentParser *CommentParser // 注释解析器，用于解析 @provider 注解
	providers     []Provider     // 发现的 Provider 列表
	currentFile   string         // 当前处理的文件路径
}

// NewProviderDiscoveryVisitor 创建新的 ProviderDiscoveryVisitor 实例
//
// 初始化流程：
// ┌─────────────────────────────────────────────────────────────┐
// │            NewProviderDiscoveryVisitor()                    │
// ├─────────────────────────────────────────────────────────────┤
// │  1. 保存注释解析器：commentParser 参数                     │
// │  2. 初始化 Provider 列表：make([]Provider, 0)               │
// │  3. 创建访问者实例：&ProviderDiscoveryVisitor{...}         │
// └─────────────────────────────────────────────────────────────┘
//
// 参数处理：
//   - commentParser: 必须传入有效的注释解析器实例
//   - 不允许为 nil，否则会出现运行时错误
//
// 初始化状态：
//   - providers: 空列表，准备收集发现的 Provider
//   - currentFile: 空字符串，等待文件处理时设置
//
// 参数：
//   - commentParser: 注释解析器实例
//
// 返回值：
//   - *ProviderDiscoveryVisitor: 配置好的访问者实例
//
// 使用示例：
//
//	commentParser := NewCommentParser()
//	visitor := NewProviderDiscoveryVisitor(commentParser)
//	walker := NewASTWalker()
//	walker.AddVisitor(visitor)
//
// 注意事项：
//   - 注释解析器必须提前初始化
//   - 访问者可以重复使用，但需要调用 Reset 清理状态
//   - 建议在遍历开始前添加到 ASTWalker
func NewProviderDiscoveryVisitor(commentParser *CommentParser) *ProviderDiscoveryVisitor {
	return &ProviderDiscoveryVisitor{
		commentParser: commentParser,
		providers:     make([]Provider, 0),
	}
}

// VisitFile 实现 NodeVisitor.VisitFile 接口，处理文件级别的初始化
//
// 功能说明：
//   - 记录当前处理的文件路径
//   - 为文件级别的处理提供上下文
//   - 不执行具体的发现逻辑
//
// 执行时机：
//   - 在开始处理新文件时调用
//   - 在任何 AST 遍历之前调用
//   - 为后续的处理步骤建立上下文
//
// 处理策略：
//   - 简单记录：只保存文件路径，不进行复杂处理
//   - 快速返回：立即返回 nil，不产生错误
//   - 上下文建立：为后续的 VisitStructType 提供文件上下文
//
// 参数：
//   - filePath: 当前处理的文件路径
//   - node: 解析后的 AST 文件节点
//
// 返回值：
//   - error: 总是返回 nil，表示成功
//
// 注意事项：
//   - 此方法主要建立上下文，不执行具体业务逻辑
//   - 文件路径会被后续的方法使用
//   - 错误处理在其他方法中实现
func (pdv *ProviderDiscoveryVisitor) VisitFile(filePath string, node *ast.File) error {
	// 记录当前处理的文件路径，为后续处理提供上下文
	pdv.currentFile = filePath
	return nil
}

// VisitGenDecl 实现 NodeVisitor.VisitGenDecl 接口，处理通用声明
//
// 功能说明：
//   - 当前实现不处理通用声明级别的 Provider 发现
//   - Provider 发现主要在结构体级别进行
//   - 保留接口完整性，为未来扩展做准备
//
// 处理策略：
//   - 空实现：不执行任何处理逻辑
//   - 快速返回：立即返回 nil，继续后续处理
//   - 接口一致性：保持与 NodeVisitor 接口的一致性
//
// 扩展可能性：
//   - 包级别 Provider 注解的处理
//   - 常量和变量级别的注解处理
//   - 声明级别的元数据收集
//
// 参数：
//   - filePath: 当前文件路径
//   - decl: 通用声明节点
//
// 返回值：
//   - error: 总是返回 nil，表示成功
//
// 注意事项：
//   - 当前版本不处理声明级别的注解
//   - 未来版本可能添加相关功能
//   - 保持接口完整性以便扩展
func (pdv *ProviderDiscoveryVisitor) VisitGenDecl(filePath string, decl *ast.GenDecl) error {
	// 当前实现不处理通用声明级别的注解
	// Provider 发现主要在结构体级别进行
	return nil
}

// VisitTypeSpec 实现 NodeVisitor.VisitTypeSpec 接口，处理类型规范
//
// 功能说明：
//   - 当前实现不在类型规范级别处理 Provider 发现
//   - 主要在结构体类型级别进行详细处理
//   - 保留接口完整性，为未来扩展做准备
//
// 处理策略：
//   - 空实现：不执行任何处理逻辑
//   - 快速返回：立即返回 nil，继续后续处理
//   - 结构体优先：将详细处理留给 VisitStructType
//
// 设计考虑：
//   - 类型规范级别的处理可能包含接口、函数类型等
//   - Provider 注解通常与结构体类型关联
//   - 避免重复处理，让专门的 VisitStructType 处理
//
// 扩展可能性：
//   - 接口类型的 Provider 注解处理
//   - 函数类型的注解处理
//   - 类型级别的元数据收集
//
// 参数：
//   - filePath: 当前文件路径
//   - typeSpec: 类型规范节点
//   - decl: 包含该类型的通用声明
//
// 返回值：
//   - error: 总是返回 nil，表示成功
//
// 注意事项：
//   - 当前版本不处理类型规范级别的注解
//   - 结构体类型会在 VisitStructType 中处理
//   - 保持接口完整性以便扩展
func (pdv *ProviderDiscoveryVisitor) VisitTypeSpec(filePath string, typeSpec *ast.TypeSpec, decl *ast.GenDecl) error {
	// 当前实现不在类型规范级别处理 Provider 发现
	// 结构体类型的处理在 VisitStructType 中进行
	return nil
}

// VisitStructType 实现 NodeVisitor.VisitStructType 接口，处理结构体类型的 Provider 注解
//
// 执行流程：
// ┌─────────────────────────────────────────────────────────────┐
// │          VisitStructType(filePath, structType, ...)       │
// ├─────────────────────────────────────────────────────────────┤
// │  1. 注释检查：decl.Doc != nil && len(decl.Doc.List) > 0    │
// │  2. 注释提取：收集所有注释行到 commentLines                │
// │  3. 注解解析：commentParser.ParseCommentBlock()            │
// │  4. Provider 构建：创建 Provider 对象                      │
// │  5. 默认值设置：为空字段设置合理的默认值                  │
// │  6. 结果收集：添加到 providers 列表                       │
// └─────────────────────────────────────────────────────────────┘
//
// 注解发现策略：
//   - 文档检查：检查结构体声明前的文档注释
//   - 注释收集：收集所有相关的注释行
//   - 专业解析：使用 CommentParser 进行专业解析
//   - 结果验证：检查解析结果的有效性
//
// Provider 构建逻辑：
//   - 基本信息：结构体名称、模式、分组、返回类型
//   - 初始化：创建空的注入参数和导入映射
//   - 默认值：为返回类型设置默认值（*结构体名）
//   - 扩展性：为后续的字段处理预留接口
//
// 错误处理策略：
//   - 解析错误：如果解析失败，跳过该结构体
//   - 空结果：如果没有发现注解，不创建 Provider
//   - 静默处理：不产生错误，继续处理其他结构体
//
// 参数：
//   - filePath: 当前文件路径
//   - structType: 结构体类型节点
//   - typeSpec: 对应的类型规范
//   - decl: 包含该结构体的通用声明
//
// 返回值：
//   - error: 总是返回 nil，表示成功
//
// 注意事项：
//   - 只处理带有文档注释的结构体
//   - 注解解析失败不会产生错误
//   - 发现的 Provider 会被收集到内部列表
func (pdv *ProviderDiscoveryVisitor) VisitStructType(filePath string, structType *ast.StructType, typeSpec *ast.TypeSpec, decl *ast.GenDecl) error {
	// === 注释检查 ===
	// 检查结构体是否有文档注释，没有注释则跳过处理
	if decl.Doc != nil && len(decl.Doc.List) > 0 {
		// === 注释提取 ===
		// 提取所有注释行，为后续的注解解析做准备
		commentLines := make([]string, len(decl.Doc.List))
		for i, comment := range decl.Doc.List {
			commentLines[i] = comment.Text
		}

		// === 注解解析 ===
		// 使用注释解析器解析注释块，查找 @provider 注解
		providerComment, err := pdv.commentParser.ParseCommentBlock(commentLines)
		if err == nil && providerComment != nil {
			// === Provider 构建 ===
			// 创建 Provider 对象，设置基本信息
			provider := Provider{
				StructName:    typeSpec.Name.Name,           // 结构体名称
				Mode:          providerComment.Mode,         // Provider 模式
				ProviderGroup: providerComment.Group,        // Provider 分组
				ReturnType:    providerComment.ReturnType,   // 返回类型
				InjectParams:  make(map[string]InjectParam), // 注入参数映射
				Imports:       make(map[string]string),      // 导入包映射
			}

			// === 默认值设置 ===
			// 如果没有指定返回类型，设置默认值为结构体指针
			if provider.ReturnType == "" {
				provider.ReturnType = "*" + provider.StructName
			}

			// === 结果收集 ===
			// 将构建的 Provider 添加到发现列表
			pdv.providers = append(pdv.providers, provider)
		}
	}

	return nil
}

// VisitStructField 实现 NodeVisitor.VisitStructField 接口，处理结构体字段
//
// 功能说明：
//   - 当前实现不处理字段级别的详细信息
//   - 字段处理在后续的 Provider 构建阶段进行
//   - 保留接口完整性，为未来扩展做准备
//
// 当前职责：
//   - 接口实现：保持与 NodeVisitor 接口的一致性
//   - 扩展预留：为字段级别的处理预留接口
//   - 简单返回：不执行复杂处理逻辑
//
// 未来可能的扩展：
//   - 字段类型解析：解析字段的类型信息
//   - 标签处理：处理 struct tags，特别是 inject 标签
//   - 依赖分析：分析字段之间的依赖关系
//   - 验证检查：验证字段配置的正确性
//
// 字段处理策略：
//   - 注入标签：识别 inject:"true" 和 inject:"false" 标签
//   - 类型解析：解析复杂类型，包括指针、切片、映射等
//   - 嵌入字段：处理匿名字段和嵌入结构体
//   - 文档注释：处理字段级别的文档注释
//
// 参数：
//   - filePath: 当前文件路径
//   - field: 结构体字段节点
//   - structType: 所属的结构体类型
//
// 返回值：
//   - error: 总是返回 nil，表示成功
//
// 注意事项：
//   - 当前版本不进行字段级别的处理
//   - 字段处理逻辑在其他模块中实现
//   - 保持接口完整性以便未来扩展
func (pdv *ProviderDiscoveryVisitor) VisitStructField(filePath string, field *ast.Field, structType *ast.StructType) error {
	// 字段级别处理预留接口
	// 当前版本不在此处进行字段处理，相关逻辑在 Provider 构建阶段实现
	// 未来可能的扩展：
	// - 字段类型解析
	// - inject 标签处理
	// - 依赖关系分析
	return nil
}

// Complete 实现 NodeVisitor.Complete 接口，处理文件完成事件
//
// 功能说明：
//   - 当前实现不处理文件完成事件
//   - 主要用于清理和结果汇总工作
//   - 保留接口完整性，为未来扩展做准备
//
// 当前职责：
//   - 接口实现：保持与 NodeVisitor 接口的一致性
//   - 快速返回：不执行复杂处理逻辑
//   - 扩展预留：为文件完成处理预留接口
//
// 未来可能的扩展：
//   - 结果验证：验证文件级别的处理结果
//   - 统计收集：收集文件处理的统计信息
//   - 错误报告：汇总文件处理过程中的错误
//   - 缓存清理：清理文件级别的临时数据
//
// 文件完成处理策略：
//   - 静默完成：不产生错误，继续后续处理
//   - 结果保持：保持已发现的 Provider 结果
//   - 状态重置：为下一个文件的处理准备状态
//
// 参数：
//   - filePath: 已完成处理的文件路径
//
// 返回值：
//   - error: 总是返回 nil，表示成功
//
// 注意事项：
//   - 当前版本不进行文件完成处理
//   - 文件路径参数可用于扩展功能
//   - 保持接口完整性以便未来扩展
func (pdv *ProviderDiscoveryVisitor) Complete(filePath string) error {
	// 文件完成处理预留接口
	// 当前版本不进行特殊的完成处理
	// 未来可能的扩展：
	// - 文件级别结果验证
	// - 统计信息收集
	// - 缓存清理工作
	return nil
}

// GetProviders 返回发现的 Provider 列表
//
// 功能说明：
//   - 返回在 AST 遍历过程中发现的所有 Provider
//   - 返回的是内部列表的副本，外部修改不影响内部状态
//   - 包含完整的 Provider 信息，包括基本属性和配置
//
// 返回数据特征：
//   - 完整性：包含所有发现的 Provider，没有遗漏
//   - 顺序性：按照发现顺序排列，保持源码中的顺序
//   - 结构化：每个 Provider 包含完整的信息结构
//
// 数据内容：
//   - 基本信息：StructName、Mode、ProviderGroup、ReturnType
//   - 注入信息：InjectParams 映射
//   - 导入信息：Imports 映射
//   - 其他属性：NeedPrepareFunc、ProviderFile 等
//
// 使用场景：
//   - 代码生成：基于发现的 Provider 生成代码
//   - 项目分析：分析项目中的 Provider 定义
//   - 工具集成：与其他工具集成处理 Provider 信息
//   - 调试测试：验证 Provider 发现的正确性
//
// 返回值：
//   - []Provider: 发现的 Provider 列表
//
// 注意事项：
//   - 返回的是实际列表，不是副本
//   - 空列表表示没有发现任何 Provider
//   - 外部修改会影响内部状态
func (pdv *ProviderDiscoveryVisitor) GetProviders() []Provider {
	return pdv.providers
}

// Reset 清理已发现的 Provider，重置访问者状态
//
// 功能说明：
//   - 清空已发现的 Provider 列表
//   - 重置当前文件路径
//   - 重置访问者到初始状态，可以重新使用
//
// 重置策略：
//   - 列表重置：创建新的空 Provider 列表
//   - 路径重置：将当前文件路径设为空字符串
//   - 状态清理：清理所有与文件处理相关的状态
//
// 使用场景：
//   - 重复使用：在多次遍历中重复使用同一个访问者
//   - 状态清理：清理之前的处理结果，避免混淆
//   - 测试环境：在单元测试中重置状态
//   - 批量处理：在批量处理多个目录时重置状态
//
// 重置前后对比：
//   - 重置前：providers 包含之前发现的结果，currentFile 包含上一个文件
//   - 重置后：providers 为空列表，currentFile 为空字符串
//
// 注意事项：
//   - 重置操作不可逆，之前的结果会丢失
//   - 重置后可以安全地用于新的遍历任务
//   - 建议在开始新的遍历前调用 Reset
func (pdv *ProviderDiscoveryVisitor) Reset() {
	// 清空已发现的 Provider 列表
	pdv.providers = make([]Provider, 0)
	// 重置当前文件路径
	pdv.currentFile = ""
}
