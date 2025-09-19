# Validation Rules Contract

## 概述

定义 pkg/ast/provider 包的验证规则契约，确保 provider 配置的正确性和一致性。

## 验证规则分类

### 1. 结构验证规则

#### 命名规则
```go
// StructNameRule 验证结构体名称
type StructNameRule struct{}

func (r *StructNameRule) Name() string {
    return "struct_name"
}

func (r *StructNameRule) Validate(p *Provider) error {
    if p.StructName == "" {
        return &ValidationError{
            Field:   "StructName",
            Value:   "",
            Message: "struct name cannot be empty",
        }
    }

    if !isValidGoIdentifier(p.StructName) {
        return &ValidationError{
            Field:   "StructName",
            Value:   p.StructName,
            Message: "struct name must be a valid Go identifier",
        }
    }

    if !isExported(p.StructName) {
        return &ValidationError{
            Field:   "StructName",
            Value:   p.StructName,
            Message: "struct name must be exported (start with uppercase letter)",
        }
    }

    return nil
}
```

#### 包名规则
```go
// PackageNameRule 验证包名
type PackageNameRule struct{}

func (r *PackageNameRule) Name() string {
    return "package_name"
}

func (r *PackageNameRule) Validate(p *Provider) error {
    if p.PkgName == "" {
        return &ValidationError{
            Field:   "PkgName",
            Value:   "",
            Message: "package name cannot be empty",
        }
    }

    if !isValidGoIdentifier(p.PkgName) {
        return &ValidationError{
            Field:   "PkgName",
            Value:   p.PkgName,
            Message: "package name must be a valid Go identifier",
        }
    }

    return nil
}
```

### 2. 类型验证规则

#### 返回类型规则
```go
// ReturnTypeRule 验证返回类型
type ReturnTypeRule struct{}

func (r *ReturnTypeRule) Name() string {
    return "return_type"
}

func (r *ReturnTypeRule) Validate(p *Provider) error {
    if p.ReturnType == "" {
        return &ValidationError{
            Field:   "ReturnType",
            Value:   "",
            Message: "return type cannot be empty",
        }
    }

    if !isValidGoType(p.ReturnType) {
        return &ValidationError{
            Field:   "ReturnType",
            Value:   p.ReturnType,
            Message: "return type must be a valid Go type",
        }
    }

    return nil
}
```

#### 模式验证规则
```go
// ProviderModeRule 验证 provider 模式
type ProviderModeRule struct {
    validModes map[ProviderMode]bool
}

func NewProviderModeRule() *ProviderModeRule {
    return &ProviderModeRule{
        validModes: map[ProviderMode]bool{
            ModeDefault: true,
            ModeGRPC:    true,
            ModeEvent:   true,
            ModeJob:     true,
            ModeCronJob: true,
            ModeModel:   true,
        },
    }
}

func (r *ProviderModeRule) Name() string {
    return "provider_mode"
}

func (r *ProviderModeRule) Validate(p *Provider) error {
    if !r.validModes[p.Mode] {
        return &ValidationError{
            Field:   "Mode",
            Value:   string(p.Mode),
            Message: fmt.Sprintf("invalid provider mode: %s", p.Mode),
        }
    }

    return nil
}
```

### 3. 依赖注入验证规则

#### 注入参数规则
```go
// InjectParamsRule 验证注入参数
type InjectParamsRule struct{}

func (r *InjectParamsRule) Name() string {
    return "inject_params"
}

func (r *InjectParamsRule) Validate(p *Provider) error {
    // 验证注入参数不为标量类型
    for name, param := range p.InjectParams {
        if isScalarType(param.Type) {
            return &ValidationError{
                Field:   fmt.Sprintf("InjectParams.%s", name),
                Value:   param.Type,
                Message: fmt.Sprintf("scalar type '%s' cannot be injected", param.Type),
            }
        }
    }

    // 验证 only 模式下的 inject:true 标签
    if p.InjectionMode == InjectionOnly {
        for name, param := range p.InjectParams {
            if param.InjectTag != "true" {
                return &ValidationError{
                    Field:   fmt.Sprintf("InjectParams.%s", name),
                    Value:   param.InjectTag,
                    Message: "all fields must have inject:\"true\" tag in 'only' mode",
                }
            }
        }
    }

    // 验证 except 模式下的 inject:false 标签
    if p.InjectionMode == InjectionExcept {
        for name, param := range p.InjectParams {
            if param.InjectTag == "false" {
                return &ValidationError{
                    Field:   fmt.Sprintf("InjectParams.%s", name),
                    Value:   param.InjectTag,
                    Message: "fields with inject:\"false\" tag should not be included in 'except' mode",
                }
            }
        }
    }

    return nil
}
```

#### 包别名规则
```go
// PackageAliasRule 验证包别名
type PackageAliasRule struct{}

func (r *PackageAliasRule) Name() string {
    return "package_alias"
}

func (r *PackageAliasRule) Validate(p *Provider) error {
    // 收集所有包别名
    aliases := make(map[string]string)
    for alias, pkg := range p.Imports {
        if existing, exists := aliases[alias]; exists && existing != pkg {
            return &ValidationError{
                Field:   "Imports",
                Value:   alias,
                Message: fmt.Sprintf("duplicate package alias '%s' for different packages", alias),
            }
        }
        aliases[alias] = pkg
    }

    // 验证注入参数的包别名
    for name, param := range p.InjectParams {
        if param.PackageAlias != "" {
            if _, exists := aliases[param.PackageAlias]; !exists {
                return &ValidationError{
                    Field:   fmt.Sprintf("InjectParams.%s", name),
                    Value:   param.PackageAlias,
                    Message: fmt.Sprintf("undefined package alias '%s'", param.PackageAlias),
                }
            }
        }
    }

    return nil
}
```

### 4. 模式特定验证规则

#### gRPC 模式规则
```go
// GRPCModeRule 验证 gRPC 模式特定规则
type GRPCModeRule struct{}

func (r *GRPCModeRule) Name() string {
    return "grpc_mode"
}

func (r *GRPCModeRule) Validate(p *Provider) error {
    if p.Mode != ModeGRPC {
        return nil
    }

    // 验证 gRPC 注册函数
    if p.GrpcRegisterFunc == "" {
        return &ValidationError{
            Field:   "GrpcRegisterFunc",
            Value:   "",
            Message: "gRPC register function cannot be empty in gRPC mode",
        }
    }

    // 验证返回类型
    if p.ReturnType != "contracts.Initial" {
        return &ValidationError{
            Field:   "ReturnType",
            Value:   p.ReturnType,
            Message: "return type must be 'contracts.Initial' in gRPC mode",
        }
    }

    // 验证 provider 组
    if p.ProviderGroup != "atom.GroupInitial" {
        return &ValidationError{
            Field:   "ProviderGroup",
            Value:   p.ProviderGroup,
            Message: "provider group must be 'atom.GroupInitial' in gRPC mode",
        }
    }

    // 验证必须包含 __grpc 注入参数
    if _, exists := p.InjectParams["__grpc"]; !exists {
        return &ValidationError{
            Field:   "InjectParams",
            Value:   "",
            Message: "gRPC provider must include '__grpc' injection parameter",
        }
    }

    return nil
}
```

#### Event 模式规则
```go
// EventModeRule 验证 Event 模式特定规则
type EventModeRule struct{}

func (r *EventModeRule) Name() string {
    return "event_mode"
}

func (r *EventModeRule) Validate(p *Provider) error {
    if p.Mode != ModeEvent {
        return nil
    }

    // 验证返回类型
    if p.ReturnType != "contracts.Initial" {
        return &ValidationError{
            Field:   "ReturnType",
            Value:   p.ReturnType,
            Message: "return type must be 'contracts.Initial' in event mode",
        }
    }

    // 验证 provider 组
    if p.ProviderGroup != "atom.GroupInitial" {
        return &ValidationError{
            Field:   "ProviderGroup",
            Value:   p.ProviderGroup,
            Message: "provider group must be 'atom.GroupInitial' in event mode",
        }
    }

    // 验证必须包含 __event 注入参数
    if _, exists := p.InjectParams["__event"]; !exists {
        return &ValidationError{
            Field:   "InjectParams",
            Value:   "",
            Message: "event provider must include '__event' injection parameter",
        }
    }

    return nil
}
```

#### Job 模式规则
```go
// JobModeRule 验证 Job 模式特定规则
type JobModeRule struct{}

func (r *JobModeRule) Name() string {
    return "job_mode"
}

func (r *JobModeRule) Validate(p *Provider) error {
    if p.Mode != ModeJob && p.Mode != ModeCronJob {
        return nil
    }

    // 验证返回类型
    if p.ReturnType != "contracts.Initial" {
        return &ValidationError{
            Field:   "ReturnType",
            Value:   p.ReturnType,
            Message: "return type must be 'contracts.Initial' in job mode",
        }
    }

    // 验证 provider 组
    if p.ProviderGroup != "atom.GroupInitial" {
        return &ValidationError{
            Field:   "ProviderGroup",
            Value:   p.ProviderGroup,
            Message: "provider group must be 'atom.GroupInitial' in job mode",
        }
    }

    // 验证必须包含 __job 注入参数
    if _, exists := p.InjectParams["__job"]; !exists {
        return &ValidationError{
            Field:   "InjectParams",
            Value:   "",
            Message: "job provider must include '__job' injection parameter",
        }
    }

    return nil
}
```

#### Model 模式规则
```go
// ModelModeRule 验证 Model 模式特定规则
type ModelModeRule struct{}

func (r *ModelModeRule) Name() string {
    return "model_mode"
}

func (r *ModelModeRule) Validate(p *Provider) error {
    if p.Mode != ModeModel {
        return nil
    }

    // 验证返回类型
    if p.ReturnType != "contracts.Initial" {
        return &ValidationError{
            Field:   "ReturnType",
            Value:   p.ReturnType,
            Message: "return type must be 'contracts.Initial' in model mode",
        }
    }

    // 验证 provider 组
    if p.ProviderGroup != "atom.GroupInitial" {
        return &ValidationError{
            Field:   "ProviderGroup",
            Value:   p.ProviderGroup,
            Message: "provider group must be 'atom.GroupInitial' in model mode",
        }
    }

    // 验证必须设置 NeedPrepareFunc
    if !p.NeedPrepareFunc {
        return &ValidationError{
            Field:   "NeedPrepareFunc",
            Value:   "false",
            Message: "model provider must set NeedPrepareFunc to true",
        }
    }

    return nil
}
```

### 5. 注释验证规则

#### 注释格式规则
```go
// CommentFormatRule 验证注释格式
type CommentFormatRule struct{}

func (r *CommentFormatRule) Name() string {
    return "comment_format"
}

func (r *CommentFormatRule) Validate(p *Provider) error {
    if p.Comment == nil {
        return nil
    }

    comment := p.Comment

    // 验证注释必须以 @provider 开头
    if !strings.HasPrefix(comment.RawText, "@provider") {
        return &ValidationError{
            Field:   "Comment.RawText",
            Value:   comment.RawText,
            Message: "provider comment must start with '@provider'",
        }
    }

    // 验证模式格式
    if comment.Mode != "" && !isValidProviderMode(comment.Mode) {
        return &ValidationError{
            Field:   "Comment.Mode",
            Value:   string(comment.Mode),
            Message: fmt.Sprintf("invalid provider mode: %s", comment.Mode),
        }
    }

    // 验证注入模式
    if comment.Injection != "" && !isValidInjectionMode(comment.Injection) {
        return &ValidationError{
            Field:   "Comment.Injection",
            Value:   string(comment.Injection),
            Message: fmt.Sprintf("invalid injection mode: %s", comment.Injection),
        }
    }

    return nil
}
```

### 6. 文件路径验证规则

#### 文件路径规则
```go
// FilePathRule 验证文件路径
type FilePathRule struct{}

func (r *FilePathRule) Name() string {
    return "file_path"
}

func (r *FilePathRule) Validate(p *Provider) error {
    if p.ProviderFile == "" {
        return &ValidationError{
            Field:   "ProviderFile",
            Value:   "",
            Message: "provider file path cannot be empty",
        }
    }

    // 验证文件扩展名
    if !strings.HasSuffix(p.ProviderFile, ".go") {
        return &ValidationError{
            Field:   "ProviderFile",
            Value:   p.ProviderFile,
            Message: "provider file must have .go extension",
        }
    }

    // 验证文件路径安全性
    if !isSafeFilePath(p.ProviderFile) {
        return &ValidationError{
            Field:   "ProviderFile",
            Value:   p.ProviderFile,
            Message: "provider file path is not safe",
        }
    }

    return nil
}
```

## 验证器注册

### 默认验证器
```go
// DefaultValidator 返回默认验证器
func DefaultValidator() Validator {
    validator := NewCompositeValidator()

    // 注册所有验证规则
    validator.AddRule(&StructNameRule{})
    validator.AddRule(&PackageNameRule{})
    validator.AddRule(&ReturnTypeRule{})
    validator.AddRule(NewProviderModeRule())
    validator.AddRule(&InjectParamsRule{})
    validator.AddRule(&PackageAliasRule{})
    validator.AddRule(&GRPCModeRule{})
    validator.AddRule(&EventModeRule{})
    validator.AddRule(&JobModeRule{})
    validator.AddRule(&ModelModeRule{})
    validator.AddRule(&CommentFormatRule{})
    validator.AddRule(&FilePathRule{})

    return validator
}
```

### 自定义验证器
```go
// CompositeValidator 组合验证器
type CompositeValidator struct {
    rules []ValidationRule
}

func NewCompositeValidator() *CompositeValidator {
    return &CompositeValidator{
        rules: make([]ValidationRule, 0),
    }
}

func (v *CompositeValidator) AddRule(rule ValidationRule) {
    v.rules = append(v.rules, rule)
}

func (v *CompositeValidator) Validate(p *Provider) []error {
    var errors []error

    for _, rule := range v.rules {
        if err := rule.Validate(p); err != nil {
            errors = append(errors, err)
        }
    }

    return errors
}
```

## 验证结果

### 验证报告
```go
// ValidationReport 验证报告
type ValidationReport struct {
    Provider  *Provider  `json:"provider"`
    Errors    []error    `json:"errors"`
    Warnings  []error    `json:"warnings"`
    IsValid   bool       `json:"is_valid"`
    Timestamp time.Time  `json:"timestamp"`
}

// ValidateProvider 验证 provider 并生成报告
func ValidateProvider(p *Provider) *ValidationReport {
    report := &ValidationReport{
        Provider:  p,
        Timestamp: time.Now(),
    }

    validator := DefaultValidator()
    errors := validator.Validate(p)

    if len(errors) == 0 {
        report.IsValid = true
    } else {
        report.Errors = errors
        report.IsValid = false
    }

    return report
}
```

## 性能考虑

### 验证性能优化
- 快速失败：遇到第一个错误立即返回
- 缓存验证结果
- 并行验证独立规则
- 懒加载验证器

### 内存优化
- 重用验证器实例
- 避免重复分配错误对象
- 使用对象池管理验证器

## 扩展性

### 自定义规则注册
```go
// RegisterValidationRule 注册全局验证规则
func RegisterValidationRule(name string, rule ValidationRule) error {
    // 实现规则注册逻辑
}

// GetValidationRule 获取验证规则
func GetValidationRule(name string) (ValidationRule, bool) {
    // 实现规则获取逻辑
}
```

### 规则优先级
```go
// PrioritizedRule 优先级规则
type PrioritizedRule struct {
    Rule     ValidationRule
    Priority int
}

// PriorityValidator 优先级验证器
type PriorityValidator struct {
    rules []PrioritizedRule
}

func (v *PriorityValidator) Validate(p *Provider) []error {
    // 按优先级执行验证
}
```