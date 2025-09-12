# atomctl v2 命令行工具

Atom 应用脚手架与代码生成工具。基于 cobra 构建，集成新项目/模块脚手架、HTTP 路由与 Provider 生成、数据库迁移与模型生成、Swagger 文档、Proto 生成与代码格式化等能力。

## 安装与使用

- 安装已发布版本：`go install go.ipao.vip/atomctl/v2@latest`
- 本地构建：`go build -o atomctl .`
- 开发运行：`go run . --help`
- 运行测试：`go test ./...`

运行 `atomctl --help` 查看总览。

## 顶级命令

- `new`：项目脚手架与代码模板初始化
- `gen`：代码生成（route/provider/model/enum/service）
- `migrate`：数据库迁移（pressly/goose）
- `swag`：Swagger 文档生成与格式化
- `fmt`：Go 代码格式化（自动安装/调用 gofumpt）
- `buf`：Proto 代码生成（自动安装 buf 并执行 generate）

---

## new 脚手架

创建项目/组件初始文件。该命令组带有以下持久化参数（所有子命令通用）：

- `--force, -f`：覆盖已存在文件/目录
- `--dry-run`：仅预览生成动作，不写入文件
- `--dir`：输出基目录（默认 `.`）

子命令：

### new project (别名：p)

- 功能：基于内置模板创建新项目，或在已有 go.mod 的项目内进行就地初始化
- 参数：
  - `[module]`：可选。在空目录中创建项目时必填（例如 `github.com/acme/myapp`）；在已有项目内可省略
- 选项：
  - `--force`：当目标目录存在时强制覆盖
- 行为说明：
  - 识别模板中的隐藏文件占位（模板以 `-` 开头的文件名会渲染为 `.` 前缀，例如 `-.gitignore` -> `.gitignore`）
  - 支持 `.tpl` 渲染与 `.raw` 直写；也支持文件内通过 `atomctl:mode=tpl|raw` 指示
  - 在当前模块内执行时会跳过覆盖已有的 `go.mod`
  - 生成后提示后续步骤（`cd` 与 `go mod tidy`）

示例：

```
atomctl new project github.com/acme/demo
atomctl new -f --dir ./playground project github.com/acme/demo
atomctl new project        # 在已有 go.mod 的项目中就地初始化
```

### new provider

- 功能：创建 Provider 脚手架到 `providers/<name>`
- 参数：
  - `<name>`：必填，provider 名称（例如 `cache`、`email`）
- 选项：继承 `--dry-run`、`--dir`

示例：

```
atomctl new provider email
atomctl new --dry-run --dir ./demo provider cache
```

### new event (别名：e)

- 功能：生成事件发布者与订阅者模板
- 参数：
  - `<name>`：必填，事件名（驼峰自动转换为 snake 作为 topic）
- 选项：
  - `--only`：仅生成某一侧，`publisher` 或 `subscriber`
  - 继承 `--dry-run`、`--dir`
- 行为说明：
  - 生成 publisher 到 `app/events/publishers/<snake>.go`
  - 生成 subscriber 到 `app/events/subscribers/<snake>.go`
  - 追加常量到 `app/events/topics.go`：`const Topic{Name} = "<snake>"`（避免重复）

示例：

```
atomctl new event UserCreated
atomctl new event UserCreated --only=publisher
```

### new job

- 功能：生成任务模板到 `app/jobs/<snake>.go`
- 参数：
  - `<name>`：必填，任务名
- 选项：继承 `--dry-run`、`--dir`

示例：

```
atomctl new job SendDailyReport
```

> 说明：代码中已存在 `new module` 实现，但当前未注册到 `new` 子命令中（处于弃用状态）。

---

## gen 代码生成

命令组持久化参数：

- `-c, --config`：数据库配置文件路径（默认 `config.toml`），用于 `gen model`
- 执行完成后自动调用 `atomctl fmt` 格式化生成代码

子命令：

### gen model (别名：m)

- 功能：连接 PostgreSQL，基于配置生成 GORM 模型与相关文件（委托 `go.ipao.vip/gen`）到 `./database`，并读取 `./database/.transform.yaml` 做类型/命名转换
- 配置文件 `[Database]` 字段：
  - `Username`、`Password`、`Database`、`Host`、`Port`、`Schema`、`SslMode`、`TimeZone`
- 示例 `config.toml`：

```toml
[Database]
Host = "127.0.0.1"
Port = 5432
Username = "postgres"
Password = "secret"
Database = "app"
Schema = "public"
SslMode = "disable"
TimeZone = "Asia/Shanghai"
```

示例：

```
atomctl gen -c config.toml model
```

### gen route

- 功能：扫描项目 `app/http` 目录的控制器，解析注释生成 `routes.gen.go`
- 参数：
  - `--path`：扫描根目录（默认当前工作目录）；也可作为位置参数传入项目根路径
- 注释规则：
  - `@Router <path> [<method>]` 指定路由与方法，如 `@Router /users/:id [get]`
  - `@Bind <name> (<position>) key(<key>) [model(<field[:field_type]>)]`
    - position 枚举：`path|query|body|header|cookie|local|file`
    - model() 定义形式：
      1) `model()` 默认字段 `id`，类型 `int`
      2) `model(id)` 指定字段 `id`，类型默认 `int`
      3) `model(id:int)` 指定字段 `id`，类型 `int`
    - 示例：
      - `@Bind id path key(id)`（PathParam 绑定）
      - `@Bind page query key(page)`（QueryParam 绑定）
      - `@Bind auth header key(Authorization)`
      - `@Bind file file key(upload)`（文件）
      - `@Bind user path key(id) model(id:int)`（从路径 id 加载模型，按 id:int 查询）

示例：

```go
// UserController
// @Router /users/:id [get]
// @Bind id path key(id) model(id:int)
func (c *UserController) Show(ctx context.Context, id int) (*User, error) { /* ... */ }
```

执行：

```
atomctl gen route --path .
```

> 注意：生成完成后会自动触发 `gen provider` 以补全依赖注入。

#### Bind 参数类型与位置支持

- 标量类型（scalar）：`string`、`bool`、`int|int8|int16|int32|int64`、`uint|uint8|uint16|uint32|uint64`、`float32|float64`
- 位置与类型支持：
  - `path`：支持标量；也支持通过 `model()` 从路径值加载模型记录（推荐标量或 model 查找）。
    - 示例：`@Bind id path key(id)`、`@Bind user path key(id) model(id:int)`
  - `query`：支持标量或结构体（结构体将按字段名/标签聚合解析查询参数）。
    - 示例：`@Bind page query key(page)`、`@Bind filter query key(q)` 或 `@Bind opts query key(_)`
  - `header`：支持标量或结构体（按字段名/标签解析请求头）。
    - 示例：`@Bind auth header key(Authorization)`
  - `cookie`：支持标量；`string` 有快捷写法（`CookieParam`）。
    - 示例：`@Bind sid cookie key(session_id)`
  - `body`：支持结构体或标量（通常用于 JSON 体，推荐结构体）。
    - 示例：`@Bind input body key(_)`
  - `file`：固定类型为 `multipart.FileHeader`，用于文件上传。
    - 示例：`@Bind file file key(upload)`
  - `local`：任意类型（从上下文本地存取，而非请求来源）。
    - 示例：`@Bind user local key(currentUser)`


### gen provider (别名：p)

- 功能：扫描 Go 源码查找带 `@provider` 注释的结构体，生成 `provider.gen.go` 实现依赖注入与分组注册
- 注释语法：`@provider(<mode>):[except|only] [returnType] [group]`
  - `mode`：`grpc|event|job|cronjob|model`（留空为默认）
  - `:only`：仅注入标注了结构体字段 tag `inject:"true"` 的依赖
  - `:except`：注入全部非标量字段，除非字段 tag `inject:"false"`
  - `returnType`：Provide 返回类型（如 `contracts.Initial`）
  - `group`：归属分组（如 `atom.GroupInitial`）
- 行为特性：
  - 自动分析字段类型与 import，忽略标量类型
  - `grpc`：注入 `providers/grpc.Grpc`，设置 `GrpcRegisterFunc`
  - `event`：注入 `providers/event.PubSub`
  - `job|cronjob`：注入 `providers/job.Job` 并引入 `github.com/riverqueue/river`
  - `model`：标记需要 `Prepare` 钩子

示例：

```go
// @provider(job):only contracts.Initial atom.GroupInitial
type UserJob struct {
    // 仅在 only 模式下注入
    Repo *repo.UserRepo `inject:"true"`
    Log  *log.Logger    `inject:"true"`
}
```

执行：`atomctl gen provider`（支持在任意子目录执行，自动写入对应包下的 `provider.gen.go`）

### gen enum (别名：e)

- 功能：在工程内扫描包含 `ENUM(...)` 且带 `swagger:enum` 的类型定义文件，为每个文件生成 `*.gen.go`
- 选项：
  - `-f, --flag`：生成 `flag.Value` 支持（默认 true）
  - `-m, --marshal`：生成 JSON 编解码支持
  - `-s, --sql`：生成 `database/sql` 相关驱动与 Null 支持（默认 true）

示例：`atomctl gen enum -m -s`

### gen service

- 功能：根据 `app/services` 下的服务文件聚合生成 `services.gen.go`
- 选项：
  - `--path`：服务目录（默认 `./app/services`）

示例：`atomctl gen service --path ./app/services`

---

## migrate 数据库迁移（goose）

- 用法：`atomctl migrate [action] [args...]`
- 支持 action：`up|up-by-one|up-to|create|down|down-to|fix|redo|reset|status|version`
- 选项：
  - `-c, --config`：数据库配置文件（默认 `config.toml`，读取 `[Database]`）
  - `--dir`：迁移目录（默认 `database/migrations`）
  - `--table`：迁移版本表名（默认 `migrations`）
- 说明：
  - 执行 `create` 时若未指定脚本类型，会自动追加 `sql` 类型
  - 连接配置同上 `gen model` 的 `[Database]` 字段

示例：

```
atomctl migrate status
atomctl migrate up
atomctl migrate create add_users_table
```

---

## swag 文档

### swag init (别名：i)

- 功能：生成 Swagger 文档（go/json/yaml）
- 选项：
  - `--dir`：项目根目录（默认 `.`）
  - `--out`：输出目录（默认 `docs`）
  - `--main`：主入口文件（默认 `main.go`）

示例：`atomctl swag init --dir . --out docs --main cmd/server/main.go`

### swag fmt (别名：f)

- 功能：格式化接口注释与路由
- 选项：
  - `--dir`：扫描目录（默认 `./app/http`）
  - `--main`：主入口文件（默认 `main.go`）

示例：`atomctl swag fmt --dir ./app/http --main main.go`

---

## fmt 代码格式化

- 功能：调用 `gofumpt -extra` 格式化全项目；若未安装将自动 `go install mvdan.cc/gofumpt@latest`
- 选项：
  - `--check`：仅检查不写入，输出未格式化文件列表
  - `--path`：格式化范围（默认 `.`）

示例：

```
atomctl fmt                 # 写入格式化
atomctl fmt --check --path ./pkg
```

---

## buf Proto 生成

- 功能：检查/安装 `buf` 并执行 `buf generate`（在指定目录）
- 选项：
  - `--dir`：执行目录（默认 `.`）
  - `--dry-run`：仅预览不执行
- 说明：若找不到 `buf.yaml` 会给出提示但仍尝试生成

示例：`atomctl buf --dir ./proto`

---

## 配置与约定

- Go 版本：1.23+
- 数据库配置段落：`[Database]`（被 `gen model` 与 `migrate` 读取）
- 代码生成默认位置：
  - 路由：`app/http/**/routes.gen.go`
  - Provider：每个包内的 `provider.gen.go`
  - 模型：`./database`（含 `.transform.yaml`）
  - Enum：原文件同目录生成 `*.gen.go`
  - Services：`app/services/services.gen.go`

---

## 常见问题

- 路由未生成？确保控制器方法注释包含 `@Router`，且形参通过 `@Bind` 标注位置/键名。
- Provider 未注入期望字段？检查 `@provider:only|:except` 与字段 tag `inject:"true|false"` 的组合是否正确。
- 执行 `fmt` 报 gofumpt 不存在？工具会自动安装，如仍失败请检查网络与 `GOBIN`/`PATH`。
- `buf generate` 找不到命令？工具会自动安装 buf，若仍失败请手动安装或检查网络代理。
