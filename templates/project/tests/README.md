# 测试指南

本项目的测试使用 **Convey 框架**，分为三个层次：单元测试、集成测试和端到端测试。

## 测试结构

```
tests/
├── setup_test.go          # 测试设置和通用工具
├── unit/                  # 单元测试
│   ├── config_test.go     # 配置测试
│   └── ...                # 其他单元测试
├── integration/           # 集成测试
│   ├── database_test.go   # 数据库集成测试
│   └── ...                # 其他集成测试
└── e2e/                   # 端到端测试
    ├── api_test.go        # API 测试
    └── ...                # 其他 E2E 测试
```

## Convey 框架概述

Convey 是一个 BDD 风格的 Go 测试框架，提供直观的语法和丰富的断言。

### 核心概念

- **Convey**: 定义测试上下文，类似于 `Describe` 或 `Context`
- **So**: 断言函数，验证预期结果
- **Reset**: 清理函数，在每个测试后执行

### 基本语法

```go
Convey("测试场景描述", t, func() {
    Convey("当某个条件发生时", func() {
        // 准备测试数据
        result := SomeFunction()

        Convey("那么应该得到预期结果", func() {
            So(result, ShouldEqual, "expected")
        })
    })

    Reset(func() {
        // 清理测试数据
    })
})
```

## 运行测试

### 运行所有测试
```bash
go test ./tests/... -v
```

### 运行特定类型的测试
```bash
# 单元测试
go test ./tests/unit/... -v

# 集成测试
go test ./tests/integration/... -v

# 端到端测试
go test ./tests/e2e/... -v
```

### 运行带覆盖率报告的测试
```bash
go test ./tests/... -v -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### 运行基准测试
```bash
go test ./tests/... -bench=. -v
```

## 测试环境配置

### 单元测试
- 不需要外部依赖
- 使用内存数据库或模拟对象
- 快速执行

### 集成测试
- 需要数据库连接
- 使用测试数据库 `{{.ProjectName}}_test`
- 需要启动 Redis 等服务

### 端到端测试
- 需要完整的应用环境
- 测试真实的 HTTP 请求
- 可能需要 Docker 环境

## Convey 测试最佳实践

### 1. 测试结构设计
- 使用描述性的中文场景描述
- 遵循 `当...那么...` 的语义结构
- 嵌套 Convey 块来组织复杂测试逻辑

```go
Convey("用户认证测试", t, func() {
    var user *User
    var token string

    Convey("当用户注册时", func() {
        user = &User{Name: "测试用户", Email: "test@example.com"}
        err := user.Register()
        So(err, ShouldBeNil)

        Convey("那么用户应该被创建", func() {
            So(user.ID, ShouldBeGreaterThan, 0)
        })
    })

    Convey("当用户登录时", func() {
        token, err := user.Login("password")
        So(err, ShouldBeNil)
        So(token, ShouldNotBeEmpty)

        Convey("那么应该获得有效的访问令牌", func() {
            So(len(token), ShouldBeGreaterThan, 0)
        })
    })
})
```

### 2. 断言使用
- 使用丰富的 So 断言函数
- 提供有意义的错误消息
- 验证所有重要的方面

### 3. 数据管理
- 使用 `Reset` 函数进行清理
- 每个测试独立准备数据
- 确保测试间不相互影响

### 4. 异步测试
- 使用适当的超时设置
- 处理并发测试
- 使用 channel 进行同步

### 5. 错误处理
- 测试错误情况
- 验证错误消息
- 确保错误处理逻辑正确

## 常用 Convey 断言

### 相等性断言
```go
So(value, ShouldEqual, expected)
So(value, ShouldNotEqual, expected)
So(value, ShouldResemble, expected)  // 深度比较
So(value, ShouldNotResemble, expected)
```

### 类型断言
```go
So(value, ShouldBeNil)
So(value, ShouldNotBeNil)
So(value, ShouldBeTrue)
So(value, ShouldBeFalse)
So(value, ShouldBeZeroValue)
```

### 数值断言
```go
So(value, ShouldBeGreaterThan, expected)
So(value, ShouldBeLessThan, expected)
So(value, ShouldBeBetween, lower, upper)
```

### 集合断言
```go
So(slice, ShouldHaveLength, expected)
So(slice, ShouldContain, expected)
So(slice, ShouldNotContain, expected)
So(map, ShouldContainKey, key)
```

### 字符串断言
```go
So(str, ShouldContainSubstring, substr)
So(str, ShouldStartWith, prefix)
So(str, ShouldEndWith, suffix)
So(str, ShouldMatch, regexp)
```

### 错误断言
```go
So(err, ShouldBeNil)
So(err, ShouldNotBeNil)
So(err, ShouldError, expectedError)
```

## 测试工具

- `goconvey/convey` - BDD 测试框架
- `gomock` - Mock 生成器
- `httptest` - HTTP 测试
- `sqlmock` - 数据库 mock
- `testify` - 辅助测试工具（可选）

## 测试示例

### 配置测试示例
```go
Convey("配置加载测试", t, func() {
    var config *Config

    Convey("当从文件加载配置时", func() {
        config, err := LoadConfig("config.toml")
        So(err, ShouldBeNil)
        So(config, ShouldNotBeNil)

        Convey("那么配置应该正确加载", func() {
            So(config.App.Mode, ShouldEqual, "development")
            So(config.Http.Port, ShouldEqual, 8080)
        })
    })
})
```

### 数据库测试示例
```go
Convey("数据库操作测试", t, func() {
    var db *gorm.DB

    Convey("当连接数据库时", func() {
        db = SetupTestDB()
        So(db, ShouldNotBeNil)

        Convey("那么应该能够创建记录", func() {
            user := User{Name: "测试用户", Email: "test@example.com"}
            result := db.Create(&user)
            So(result.Error, ShouldBeNil)
            So(user.ID, ShouldBeGreaterThan, 0)
        })
    })

    Reset(func() {
        if db != nil {
            CleanupTestDB(db)
        }
    })
})
```

### API 测试示例
```go
Convey("API 端点测试", t, func() {
    var server *httptest.Server

    Convey("当启动测试服务器时", func() {
        server = httptest.NewServer(NewApp())
        So(server, ShouldNotBeNil)

        Convey("那么健康检查端点应该正常工作", func() {
            resp, err := http.Get(server.URL + "/health")
            So(err, ShouldBeNil)
            So(resp.StatusCode, ShouldEqual, http.StatusOK)

            var result map[string]interface{}
            json.NewDecoder(resp.Body).Decode(&result)
            So(result["status"], ShouldEqual, "ok")
        })
    })

    Reset(func() {
        if server != nil {
            server.Close()
        }
    })
})
```

## CI/CD 集成

测试会在以下情况下自动运行：
- 代码提交时
- 创建 Pull Request 时
- 合并到主分支时

测试结果会影响代码合并决策。Convey 的详细输出有助于快速定位问题。