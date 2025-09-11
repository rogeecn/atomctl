# =========================
# gRPC Server (providers/grpc)
# =========================
[Grpc]
# 必填
Port = 9090 # gRPC 监听端口
# 可选
# Host = "0.0.0.0"            # 监听地址（默认 0.0.0.0）
EnableReflection = true # 开启服务反射（开发/调试友好）
EnableHealth = true # 注册 gRPC health 服务
ShutdownTimeoutSeconds = 10 # 优雅关停超时，超时后强制 Stop

# 说明：
# - 统一的拦截器、ServerOption 可通过 providers/grpc/options.go 的
#   UseUnaryInterceptors/UseStreamInterceptors/UseOptions 动态注入。
# =========================
# HTTP Server (providers/http)
# =========================
[Http]
# 必填
Port = 8080 # HTTP 监听端口

# 可选
# BaseURI = "/api"             # 全局前缀
# StaticRoute = "/static"      # 静态路由路径
# StaticPath = "./public"      # 静态文件目录
[Http.Tls]

# Cert = "server.crt"
# Key  = "server.key"
[Http.Cors]

# Mode = "enabled"             # "enabled"|"disabled"（默认按 Whitelist 推断）
# 白名单项示例（按需追加多条）
# [[Http.Cors.Whitelist]]
# AllowOrigin = "https://example.com"
# AllowHeaders = "Authorization,Content-Type"
# AllowMethods = "GET,POST,PUT,DELETE"
# ExposeHeaders = "X-Request-Id"
# AllowCredentials = true
# =========================
# Connection Multiplexer (providers/cmux)
# 用于同端口同时暴露 HTTP + gRPC：cmux -> 分发到 Http/Grpc
# =========================
[Cmux]
# 必填
Port = 8081 # cmux 监听端口

# 可选
# Host = "0.0.0.0"
# =========================
# Events / PubSub (providers/event)
# gochannel 为默认内存通道（始终启用）
# 如需 Kafka / Redis Stream / SQL，请按需开启对应小节
# =========================
[Events]

# Kafka（可选）
[Events.Kafka]
# 必填（启用时）
Brokers = ["127.0.0.1:9092"]
# 可选
ConsumerGroup = "my-group"

# Redis Stream（可选）
[Events.Redis]
# 必填（启用时）
ConsumerGroup = "my-group"
# 可选
Streams = ["mystream"] # 订阅的 streams；可在 Handler 侧具体指定

# SQL（可选，基于 PostgreSQL）
[Events.Sql]
# 必填（启用时）
ConsumerGroup = "my-group"

# =========================
# Job / Queue (providers/job)
# 基于 River（Postgres）队列
# =========================
[Job]

# 可选：每队列并发数（默认 high/default/low 均 10）
#QueueWorkers = { high = 20, default = 10, low = 5 }
# 说明：
# - 需要启用 providers/postgres 以提供数据库连接
# - 通过 Add/AddWithID 入队，AddPeriodicJob 注册定时任务
# =========================
# JWT (providers/jwt)
# =========================
[JWT]
# 必填
SigningKey = "your-signing-key" # 密钥
ExpiresTime = "168h" # 过期时间，形如 "72h", "168h"
# 可选
Issuer = "my-service"

# =========================
# HashIDs (providers/hashids)
# =========================
[HashIDs]
# 必填
Salt = "your-salt"

# 可选
# Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
# MinLength = 8
# =========================
# Redis (providers/redis)
# =========================
[Redis]
# 必填（若不填 Host/Port，将默认 localhost:6379）
Host = "127.0.0.1"
Port = 6379

# 可选
# Username = ""
# Password = ""
# DB = 0
# ClientName = "my-service"
# 连接池/重试（可选）
# PoolSize = 50
# MinIdleConns = 10
# MaxRetries = 2
# 超时（秒，可选）
# DialTimeoutSeconds = 5
# ReadTimeoutSeconds = 3
# WriteTimeoutSeconds = 3
# 连接生命周期（秒，可选）
# ConnMaxIdleTimeSeconds = 300
# ConnMaxLifetimeSeconds = 1800
# 探活（秒，可选，默认 5）
# PingTimeoutSeconds = 5
# =========================
# PostgreSQL / GORM (providers/postgres)
# =========================
[Database]
# 必填
Host = "127.0.0.1"
Port = 5432
Database = "app"

# 可选（未填 Username 默认 postgres；其它有默认值见代码）
# Username = "postgres"
# Password = ""
# SslMode = "disable"          # "disable"|"require"|...
# TimeZone = "Asia/Shanghai"
# Schema = "public"
# Prefix = ""                   # 表前缀
# Singular = false              # 表名是否使用单数
# 连接池（可选）
# MaxIdleConns = 10
# MaxOpenConns = 100
# ConnMaxLifetimeSeconds = 1800
# ConnMaxIdleTimeSeconds = 300
# DSN 增强（可选）
# UseSearchPath = true
# ApplicationName = "my-service"
# =========================
# HTTP Client (providers/req)
# =========================
[HttpClient]

# 可选
# DevMode = true
# CookieJarFile = "./data/cookies.jar"
# RootCa = ["./ca.crt"]
# UserAgent = "my-service/1.0"
# InsecureSkipVerify = false
# BaseURL = "https://api.example.com"
# ContentType = "application/json"
# Timeout = 10                         # 秒
# CommonHeaders = { X-Request-From = "service" }
# CommonQuery = { locale = "zh-CN" }
[HttpClient.AuthBasic]

# Username = ""
# Password = ""
# 其它认证 / 代理 / 跳转策略（可选）
# AuthBearerToken = "Bearer <token>"   # 或仅 <token>，内部会设置 Bearer
# ProxyURL = "http://127.0.0.1:7890"
# RedirectPolicy = ["Max:10","SameHost"]
# =========================
# OpenTracing (Jaeger) (providers/tracing)
# =========================
[Tracing]
# 必填
Name = "my-service"
# 可选（Agent / Collector 至少配一个；未配时走默认本地端口）
Reporter_LocalAgentHostPort = "127.0.0.1:6831"
Reporter_CollectorEndpoint = "http://127.0.0.1:14268/api/traces"

# 行为开关（可选）
# Disabled = false
# Gen128Bit = true
# ZipkinSharedRPCSpan = true
# RPCMetrics = false
# 采样器（可选）
# Sampler_Type = "const"               # const|probabilistic|ratelimiting|remote
# Sampler_Param = 1.0
# Sampler_SamplingServerURL = ""
# Sampler_MaxOperations = 0
# Sampler_RefreshIntervalSec = 0
# Reporter（可选）
# Reporter_LogSpans = true
# Reporter_BufferFlushMs = 100
# Reporter_QueueSize = 1000
# 进程 Tags（可选）
# [Tracing.Tags]
# version = "1.0.0"
# zone = "az1"
# =========================
# OpenTelemetry (providers/otel)
# =========================
[OTEL]
# 必填（建议设置）
ServiceName = "my-service"
# 可选（版本/环境）
Version = "1.0.0"
Env = "dev"
# 导出端点（二选一，若都填优先 HTTP）
# EndpointGRPC = "127.0.0.1:4317"
# EndpointHTTP = "127.0.0.1:4318"
# 认证（可选，支持仅填纯 token，将自动补齐 Bearer）
# Token = "Bearer <token>"
# 安全（可选）
# InsecureGRPC = true
# InsecureHTTP = true
# 采样（可选）
# Sampler = "always"                   # "always"|"ratio"
# SamplerRatio = 0.1                   # 0..1，仅当 Sampler="ratio" 生效
# 批处理（毫秒，可选）
# BatchTimeoutMs = 5000
# ExportTimeoutMs = 10000
# MaxQueueSize = 2048
# MaxExportBatchSize = 512
# 指标（毫秒，可选）
# MetricReaderIntervalMs = 10000       # 指标导出周期
# RuntimeReadMemStatsIntervalMs = 5000 # 运行时指标采集最小周期
