# Parser API Contract

## 概述

定义 pkg/ast/provider 包的解析器 API 契约，确保向后兼容性和一致的接口设计。

## 核心接口

### Parser Interface
```go
// Parser 定义 provider 解析器接口
type Parser interface {
    // ParseFile 解析单个 Go 文件
    ParseFile(filename string) ([]*Provider, error)

    // ParseDir 解析目录中的所有 Go 文件
    ParseDir(dirname string) ([]*Provider, error)

    // ParseString 解析字符串内容
    ParseString(content string, filename string) ([]*Provider, error)

    // SetConfig 设置解析器配置
    SetConfig(config ParserConfig) error

    // GetConfig 获取当前配置
    GetConfig() ParserConfig
}
```

### Validator Interface
```go
// Validator 定义 provider 验证器接口
type Validator interface {
    // Validate 验证 provider 配置
    Validate(p *Provider) []error

    // ValidateComment 验证 provider 注释
    ValidateComment(comment *ProviderComment) []error
}
```

### Renderer Interface
```go
// Renderer 定义 provider 渲染器接口
type Renderer interface {
    // Render 渲染 provider 代码
    Render(filename string, providers []*Provider) error

    // RenderTemplate 使用自定义模板渲染
    RenderTemplate(filename string, providers []*Provider, template string) error
}
```

## 数据结构

### Provider
```go
// Provider 表示一个依赖注入提供者
type Provider struct {
    StructName       string            `json:"struct_name"`
    ReturnType       string            `json:"return_type"`
    Mode             ProviderMode      `json:"mode"`
    ProviderGroup    string            `json:"provider_group"`
    GrpcRegisterFunc string            `json:"grpc_register_func"`
    NeedPrepareFunc  bool              `json:"need_prepare_func"`
    InjectParams     map[string]InjectParam `json:"inject_params"`
    Imports          map[string]string  `json:"imports"`
    PkgName          string            `json:"pkg_name"`
    ProviderFile     string            `json:"provider_file"`
    SourceLocation   SourceLocation    `json:"source_location"`
}
```

### ProviderComment
```go
// ProviderComment 表示解析的 provider 注释
type ProviderComment struct {
    RawText     string       `json:"raw_text"`
    Mode        ProviderMode `json:"mode"`
    Injection   InjectionMode `json:"injection"`
    ReturnType  string       `json:"return_type"`
    Group       string       `json:"group"`
    IsValid     bool         `json:"is_valid"`
    Errors      []string     `json:"errors"`
    Location    SourceLocation `json:"location"`
}
```

### InjectParam
```go
// InjectParam 表示注入参数
type InjectParam struct {
    Star         string `json:"star"`
    Type         string `json:"type"`
    Package      string `json:"package"`
    PackageAlias string `json:"package_alias"`
    FieldName    string `json:"field_name"`
    InjectTag    string `json:"inject_tag"`
}
```

## 枚举类型

### ProviderMode
```go
// ProviderMode 定义 provider 模式
type ProviderMode string

const (
    ModeDefault ProviderMode = ""
    ModeGRPC   ProviderMode = "grpc"
    ModeEvent  ProviderMode = "event"
    ModeJob    ProviderMode = "job"
    ModeCronJob ProviderMode = "cronjob"
    ModeModel  ProviderMode = "model"
)
```

### InjectionMode
```go
// InjectionMode 定义注入模式
type InjectionMode string

const (
    InjectionDefault InjectionMode = ""
    InjectionOnly    InjectionMode = "only"
    InjectionExcept  InjectionMode = "except"
)
```

## 配置结构

### ParserConfig
```go
// ParserConfig 定义解析器配置
type ParserConfig struct {
    StrictMode     bool   `json:"strict_mode"`     // 严格模式
    AllowTestFile  bool   `json:"allow_test_file"` // 是否允许解析测试文件
    IgnorePattern  string `json:"ignore_pattern"`  // 忽略文件模式
    MaxFileSize    int64  `json:"max_file_size"`   // 最大文件大小
    EnableCache    bool   `json:"enable_cache"`    // 启用缓存
    CacheTTL       int    `json:"cache_ttl"`       // 缓存 TTL（秒）
}
```

### DefaultConfig
```go
// DefaultConfig 返回默认配置
func DefaultConfig() ParserConfig {
    return ParserConfig{
        StrictMode:     false,
        AllowTestFile:  false,
        IgnorePattern:  "*.gen.go",
        MaxFileSize:    1024 * 1024, // 1MB
        EnableCache:    false,
        CacheTTL:       300, // 5 分钟
    }
}
```

## 工厂函数

### NewParser
```go
// NewParser 创建新的解析器实例
func NewParser(config ParserConfig) (Parser, error)

// NewParserWithContext 创建带上下文的解析器
func NewParserWithContext(ctx context.Context, config ParserConfig) (Parser, error)

// NewDefaultParser 创建默认配置的解析器
func NewDefaultParser() (Parser, error)
```

### NewValidator
```go
// NewValidator 创建新的验证器实例
func NewValidator() (Validator, error)

// NewCustomValidator 创建自定义验证器
func NewCustomValidator(rules []ValidationRule) (Validator, error)
```

### NewRenderer
```go
// NewRenderer 创建新的渲染器实例
func NewRenderer() (Renderer, error)

// NewRendererWithTemplate 创建带自定义模板的渲染器
func NewRendererWithTemplate(template string) (Renderer, error)
```

## 错误处理

### Error Types
```go
// ParseError 解析错误
type ParseError struct {
    File    string `json:"file"`
    Line    int    `json:"line"`
    Column  int    `json:"column"`
    Message string `json:"message"`
}

// ValidationError 验证错误
type ValidationError struct {
    Field   string `json:"field"`
    Value   string `json:"value"`
    Message string `json:"message"`
}

// GenerationError 生成错误
type GenerationError struct {
    File    string `json:"file"`
    Message string `json:"message"`
}

// ConfigurationError 配置错误
type ConfigurationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}
```

### Error Interface
```go
// ProviderError 定义 provider 相关错误接口
type ProviderError interface {
    error
    Code() string
    Details() map[string]interface{}
    IsRetryable() bool
}
```

## 兼容性保证

### 向后兼容
- 所有公共 API 保持向后兼容
- 结构体字段只能添加，不能删除或重命名
- 接口方法只能添加，不能删除或修改签名

### 废弃策略
- 废弃的 API 会标记为 deprecated
- 废弃的 API 会在至少 2 个主版本后移除
- 提供迁移指南和工具

## 版本控制

### 版本格式
- 遵循语义化版本控制 (SemVer)
- 主版本：不兼容的 API 变更
- 次版本：向下兼容的功能性新增
- 修订号：向下兼容的问题修正

### 版本检查
```go
// Version 返回当前版本
func Version() string

// Compatible 检查版本兼容性
func Compatible(version string) bool
```

## 性能要求

### 解析性能
- 单个文件解析时间 < 100ms
- 目录解析时间 < 1s（100 个文件以内）
- 内存使用 < 50MB（正常情况）

### 生成性能
- 代码生成时间 < 200ms
- 模板渲染时间 < 50ms
- 生成的代码大小合理（< 10KB per provider）

## 安全考虑

### 输入验证
- 所有输入参数必须验证
- 文件路径必须安全检查
- 防止路径遍历攻击

### 资源限制
- 最大文件大小限制
- 最大递归深度限制
- 最大并发处理数限制

## 扩展点

### 自定义模式
```go
// RegisterProviderMode 注册自定义 provider 模式
func RegisterProviderMode(mode string, handler ProviderModeHandler) error

// ProviderModeHandler provider 模式处理器接口
type ProviderModeHandler interface {
    Parse(comment string) (*ProviderComment, error)
    Validate(comment *ProviderComment) error
    Generate(comment *ProviderComment) (*Provider, error)
}
```

### 自定义验证器
```go
// ValidationRule 验证规则接口
type ValidationRule interface {
    Name() string
    Validate(p *Provider) error
}

// RegisterValidationRule 注册验证规则
func RegisterValidationRule(rule ValidationRule) error
```

### 自定义渲染器
```go
// TemplateFunction 模板函数类型
type TemplateFunction func(args ...interface{}) (interface{}, error)

// RegisterTemplateFunction 注册模板函数
func RegisterTemplateFunction(name string, fn TemplateFunction) error
```

## 测试契约

### 单元测试
```go
// TestParserContract 测试解析器契约
func TestParserContract(t *testing.T) {
    parser := NewDefaultParser()

    // 测试基本功能
    providers, err := parser.ParseFile("testdata/simple.go")
    assert.NoError(t, err)
    assert.NotEmpty(t, providers)

    // 测试错误处理
    _, err = parser.ParseFile("testdata/invalid.go")
    assert.Error(t, err)
}
```

### 集成测试
```go
// TestFullWorkflow 测试完整工作流
func TestFullWorkflow(t *testing.T) {
    // 解析 -> 验证 -> 生成 -> 验证结果
}
```

### 基准测试
```go
// BenchmarkParser 基准测试
func BenchmarkParser(b *testing.B) {
    parser := NewDefaultParser()

    for i := 0; i < b.N; i++ {
        _, err := parser.ParseFile("testdata/large.go")
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

## 监控和日志

### 日志接口
```go
// Logger 定义日志接口
type Logger interface {
    Debug(msg string, fields ...interface{})
    Info(msg string, fields ...interface{})
    Warn(msg string, fields ...interface{})
    Error(msg string, fields ...interface{})
}

// SetLogger 设置日志器
func SetLogger(logger Logger)
```

### 指标接口
```go
// Metrics 定义指标接口
type Metrics interface {
    Counter(name string, value int64, tags ...string)
    Timer(name string, value time.Duration, tags ...string)
    Gauge(name string, value float64, tags ...string)
}

// SetMetrics 设置指标收集器
func SetMetrics(metrics Metrics)
```