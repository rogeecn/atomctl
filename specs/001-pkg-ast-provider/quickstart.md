# Quick Start Guide

## 概述

本指南展示如何使用重构后的 pkg/ast/provider 包来解析 Go 源码中的 `@provider` 注释并生成依赖注入代码。

## 前置条件

- Go 1.24.0+
- 理解 Go AST 解析
- 了解依赖注入概念

## 基本用法

### 1. 解析单个文件

```go
package main

import (
    "fmt"
    "log"

    "go.ipao.vip/atomctl/v2/pkg/ast/provider"
)

func main() {
    // 解析单个文件
    providers, err := provider.ParseFile("path/to/your/file.go")
    if err != nil {
        log.Fatal(err)
    }

    // 打印解析结果
    for _, p := range providers {
        fmt.Printf("Provider: %s\n", p.StructName)
        fmt.Printf("  Mode: %s\n", p.Mode)
        fmt.Printf("  Return Type: %s\n", p.ReturnType)
        fmt.Printf("  Inject Params: %d\n", len(p.InjectParams))
    }
}
```

### 2. 批量解析目录

```go
package main

import (
    "fmt"
    "log"

    "go.ipao.vip/atomctl/v2/pkg/ast/provider"
)

func main() {
    // 创建解析器配置
    config := provider.ParserConfig{
        StrictMode:    true,
        AllowTestFile: false,
        IgnorePattern: "*.gen.go",
    }

    // 创建解析器
    parser := provider.NewParser(config)

    // 解析目录
    providers, err := parser.ParseDir("./app")
    if err != nil {
        log.Fatal(err)
    }

    // 生成代码
    for _, p := range providers {
        err := provider.Render(p.ProviderFile, []Provider{p})
        if err != nil {
            log.Printf("Failed to render %s: %v", p.StructName, err)
        }
    }
}
```

## 支持的 Provider 注释格式

### 基本格式
```go
// @provider
type UserService struct {
    // ...
}
```

### 带模式
```go
// @provider(grpc)
type UserService struct {
    // ...
}
```

### 带注入模式
```go
// @provider:only
type UserService struct {
    Repo *UserRepo `inject:"true"`
    Log  *Logger   `inject:"true"`
}
```

### 完整格式
```go
// @provider(grpc):only contracts.Initial atom.GroupInitial
type UserService struct {
    Repo *UserRepo `inject:"true"`
    Log  *Logger   `inject:"true"`
}
```

## 测试指南

### 运行测试
```bash
# 运行所有测试
go test ./pkg/ast/provider/...

# 运行测试并显示覆盖率
go test -cover ./pkg/ast/provider/...

# 运行基准测试
go test -bench=. ./pkg/ast/provider/...
```

### 编写测试
```go
package provider_test

import (
    "testing"

    "go.ipao.vip/atomctl/v2/pkg/ast/provider"
    "github.com/stretchr/testify/assert"
)

func TestParseProvider(t *testing.T) {
    // 准备测试代码
    source := `
package main

// @provider:only contracts.Initial
type TestService struct {
    Repo *TestRepo `inject:"true"`
}
`

    // 创建临时文件
    tmpFile := createTempFile(t, source)
    defer os.Remove(tmpFile)

    // 解析
    providers, err := provider.ParseFile(tmpFile)

    // 验证
    assert.NoError(t, err)
    assert.Len(t, providers, 1)

    p := providers[0]
    assert.Equal(t, "TestService", p.StructName)
    assert.Equal(t, "contracts.Initial", p.ReturnType)
    assert.True(t, p.InjectMode.IsOnly())
}
```

## 重构指南

### 从旧版本迁移

1. **更新导入路径**
```go
// 旧版本
import "go.ipao.vip/atomctl/v2/pkg/ast/provider"

// 新版本（相同的导入路径）
import "go.ipao.vip/atomctl/v2/pkg/ast/provider"
```

2. **使用新的 API**
```go
// 旧版本
providers := provider.Parse("file.go")

// 新版本（向后兼容）
providers := provider.Parse("file.go") // 仍然支持

// 推荐的新方式
parser := provider.NewParser(provider.DefaultConfig())
providers, err := parser.ParseFile("file.go")
```

### 自定义扩展

#### 1. 自定义 Provider 模式
```go
// 实现自定义模式处理器
type CustomModeHandler struct{}

func (h *CustomModeHandler) Handle(ctx *provider.ParserContext, comment *provider.ProviderComment) (*provider.Provider, error) {
    // 自定义处理逻辑
    return &provider.Provider{
        Mode: provider.ProviderMode("custom"),
        // ...
    }, nil
}

// 注册自定义模式
provider.RegisterProviderMode("custom", &CustomModeHandler{})
```

#### 2. 自定义验证器
```go
// 实现自定义验证器
type CustomValidator struct{}

func (v *CustomValidator) Validate(p *provider.Provider) []error {
    var errors []error
    // 自定义验证逻辑
    if p.StructName == "" {
        errors = append(errors, fmt.Errorf("struct name cannot be empty"))
    }
    return errors
}

// 添加到验证链
parser.AddValidator(&CustomValidator{})
```

## 性能优化

### 1. 并行解析
```go
// 使用并行解析提高性能
func ParseProjectParallel(root string) ([]*provider.Provider, error) {
    files, err := findGoFiles(root)
    if err != nil {
        return nil, err
    }

    var wg sync.WaitGroup
    providers := make([]*provider.Provider, 0, len(files))
    errChan := make(chan error, len(files))

    for _, file := range files {
        wg.Add(1)
        go func(f string) {
            defer wg.Done()
            ps, err := provider.ParseFile(f)
            if err != nil {
                errChan <- err
                return
            }
            providers = append(providers, ps...)
        }(file)
    }

    wg.Wait()
    close(errChan)

    // 检查错误
    for err := range errChan {
        if err != nil {
            return nil, err
        }
    }

    return providers, nil
}
```

### 2. 缓存机制
```go
// 使用缓存避免重复解析
type CachedParser struct {
    cache map[string][]*provider.Provider
    parser *provider.Parser
}

func NewCachedParser() *CachedParser {
    return &CachedParser{
        cache: make(map[string][]*provider.Provider),
        parser: provider.NewParser(provider.DefaultConfig()),
    }
}

func (cp *CachedParser) ParseFile(file string) ([]*provider.Provider, error) {
    if providers, ok := cp.cache[file]; ok {
        return providers, nil
    }

    providers, err := cp.parser.ParseFile(file)
    if err != nil {
        return nil, err
    }

    cp.cache[file] = providers
    return providers, nil
}
```

## 故障排除

### 常见错误

1. **解析错误**
```
error: failed to parse provider comment: invalid mode format
```
   解决方案：检查 @provider 注释格式是否正确

2. **导入错误**
```
error: cannot resolve import path "github.com/unknown/pkg"
```
   解决方案：确保所有导入的包都存在

3. **验证错误**
```
error: provider struct has invalid return type
```
   解决方案：确保返回类型是有效的 Go 类型

### 调试技巧

1. **启用详细日志**
```go
provider.SetLogLevel(logrus.DebugLevel)
```

2. **使用解析器上下文**
```go
ctx := provider.NewParserContext(provider.ParserConfig{
    StrictMode:    true,
    ErrorHandler: func(err error) {
        log.Printf("Parse error: %v", err)
    },
})
```

3. **验证生成的代码**
```go
if err := provider.ValidateGeneratedCode(code); err != nil {
    log.Printf("Generated code validation failed: %v", err)
}
```

## 最佳实践

1. **保持注释简洁**
```go
// 推荐
// @provider:only contracts.Initial

// 不推荐
// @provider:only contracts.Initial atom.GroupInitial // 这是一个复杂的服务
```

2. **使用明确的类型**
```go
// 推荐
type UserService struct {
    Repo *UserRepository `inject:"true"`
}

// 不推荐
type UserService struct {
    Repo interface{} `inject:"true"`
}
```

3. **合理组织代码**
```go
// 将相关的 provider 放在同一个文件中
// 使用明确的包名和结构名
// 避免循环依赖
```

## 下一步

- 查看 [data-model.md](data-model.md) 了解详细的数据模型
- 阅读 [research.md](research.md) 了解重构决策过程
- 查看 [contracts/](contracts/) 目录了解 API 契约