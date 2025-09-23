# =========================
# 应用基础配置
# =========================
[App]
# 应用运行模式：development | production | testing
Mode = "development"
# 应用基础URI，用于生成完整URL
BaseURI = "http://localhost:8080"

# =========================
# HTTP 服务器配置
# =========================
[Http]
# HTTP服务监听端口
Port = 8080
# 监听地址（可选，默认 0.0.0.0）
# Host = "0.0.0.0"
# 全局路由前缀（可选）
# BaseURI = "/api/v1"

# =========================
# 数据库配置
# =========================
[Database]
# 数据库主机地址
Host = "localhost"
# 数据库端口
Port = 5432
# 数据库名称
Database = "{{.ProjectName}}"
# 数据库用户名
Username = "postgres"
# 数据库密码
Password = "password"
# SSL模式：disable | require | verify-ca | verify-full
SslMode = "disable"
# 时区
TimeZone = "Asia/Shanghai"
# 连接池配置（可选）
MaxIdleConns = 10
MaxOpenConns = 100
ConnMaxLifetime = "1800s"
ConnMaxIdleTime = "300s"

# =========================
# JWT 认证配置
# =========================
[JWT]
# JWT签名密钥（生产环境请使用强密钥）
SigningKey = "your-secret-key-change-in-production"
# Token过期时间，如：72h, 168h, 720h
ExpiresTime = "168h"
# 签发者（可选）
Issuer = "{{.ProjectName}}"

# =========================
# HashIDs 配置
# =========================
[HashIDs]
# 盐值（用于ID加密，请使用随机字符串）
Salt = "your-random-salt-here"
# 最小长度（可选）
MinLength = 8
# 自定义字符集（可选）
# Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

# =========================
# Redis 缓存配置
# =========================
[Redis]
# Redis主机地址
Host = "localhost"
# Redis端口
Port = 6379
# Redis密码（可选）
Password = ""
# 数据库编号
DB = 0
# 连接池配置（可选）
PoolSize = 50
MinIdleConns = 10
MaxRetries = 3
# 超时配置（可选）
DialTimeout = "5s"
ReadTimeout = "3s"
WriteTimeout = "3s"

# =========================
# 日志配置
# =========================
[Log]
# 日志级别：debug | info | warn | error
Level = "info"
# 日志格式：json | text
Format = "json"
# 输出文件（可选，未配置则输出到控制台）
# Output = "./logs/app.log"
# 是否启用调用者信息（文件名:行号）
EnableCaller = true