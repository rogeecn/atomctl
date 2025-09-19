# Data Model Design

## Core Entities

### 1. Provider
代表一个依赖注入提供者的核心实体

```go
type Provider struct {
    StructName       string                 // 结构体名称
    ReturnType       string                 // 返回类型
    Mode             ProviderMode           // 提供者模式
    ProviderGroup    string                 // 提供者分组
    GrpcRegisterFunc string                // gRPC 注册函数
    NeedPrepareFunc  bool                   // 是否需要 Prepare 函数
    InjectParams     map[string]InjectParam // 注入参数
    Imports          map[string]string       // 导入包
    PkgName          string                 // 包名
    ProviderFile     string                 // 生成文件路径
    SourceLocation   SourceLocation         // 源码位置
}
```

### 2. ProviderMode
提供者模式枚举

```go
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

### 3. InjectParam
注入参数描述

```go
type InjectParam struct {
    Star         string // 指针标记 (*)
    Type         string // 类型名称
    Package      string // 包路径
    PackageAlias string // 包别名
    FieldName    string // 字段名称
    InjectTag    string // inject 标签值
}
```

### 4. ProviderComment
Provider 注释描述

```go
type ProviderComment struct {
    RawText     string       // 原始注释文本
    Mode        ProviderMode // 模式
    Injection   InjectionMode // 注入模式
    ReturnType  string       // 返回类型
    Group       string       // 分组
    IsValid     bool         // 是否有效
    Errors      []string     // 解析错误
}
```

### 5. InjectionMode
注入模式枚举

```go
type InjectionMode string

const (
    InjectionDefault InjectionMode = ""
    InjectionOnly    InjectionMode = "only"
    InjectionExcept  InjectionMode = "except"
)
```

### 6. SourceLocation
源码位置信息

```go
type SourceLocation struct {
    File      string // 文件路径
    Line      int    // 行号
    Column    int    // 列号
    StartPos  int    // 开始位置
    EndPos    int    // 结束位置
}
```

### 7. ParserContext
解析器上下文

```go
type ParserContext struct {
    FileSet      *token.FileSet           // 文件集合
    Imports      map[string]string       // 导入映射
    PkgName      string                  // 包名
    ScalarTypes  []string                // 标量类型列表
    ErrorHandler func(error)             // 错误处理器
    Config       ParserConfig            // 解析配置
}
```

### 8. ParserConfig
解析器配置

```go
type ParserConfig struct {
    StrictMode    bool   // 严格模式
    AllowTestFile bool   // 是否允许解析测试文件
    IgnorePattern string // 忽略文件模式
}
```

## Relationships

### Entity Relationships
```
Provider (1) -> (0..*) InjectParam
Provider (1) -> (1) ProviderComment
Provider (1) -> (1) SourceLocation
ProviderComment (1) -> (1) ProviderMode
ProviderComment (1) -> (1) InjectionMode
```

### Data Flow
```
SourceFile -> ParserContext -> ProviderComment -> Provider -> GeneratedCode
```

## Validation Rules

### Provider Validation
- 结构体名称必须有效（符合 Go 标识符规则）
- 返回类型不能为空
- 模式必须为预定义值之一
- 注入参数不能包含标量类型
- 包名必须能正确解析

### Comment Validation
- 注释必须以 @provider 开头
- 模式格式必须正确
- 注入模式只能是 only 或 except
- 返回类型和分组格式必须正确

### Import Validation
- 包路径必须有效
- 包别名不能重复
- 匿名导入必须正确处理

## State Transitions

### Parser States
```
Idle -> Parsing -> Validating -> Generating -> Complete
                     ↓
                   Error
```

### Provider States
```
Discovered -> Parsing -> Validated -> Ready -> Generated
                    ↓
                  Invalid
```

## Performance Considerations

### Memory Usage
- 每个文件创建一个 ParserContext
- Provider 对象在解析完成后可以释放
- 导入映射应该在文件级别共享

### Processing Speed
- 并行解析独立文件
- 缓存常用的标量类型列表
- 延迟验证直到所有信息收集完成

## Error Handling

### Error Types
- ParseError: 解析错误
- ValidationError: 验证错误
- GenerationError: 生成错误
- ConfigurationError: 配置错误

### Error Recovery
- 单个 Provider 错误不影响其他 Provider
- 文件级别错误应该跳过该文件
- 提供详细的错误位置和建议

## Extension Points

### Custom Provider Modes
- 通过 ProviderMode 接口支持自定义模式
- 使用注册机制添加新模式处理器

### Custom Validation Rules
- 通过 Validator 接口支持自定义验证
- 支持链式验证器组合

### Custom Renderers
- 通过 Renderer 接口支持自定义渲染器
- 支持多种输出格式