package postgres

import (
	"fmt"
	"strconv"
	"time"

	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/opt"
	"gorm.io/gorm/logger"
)

const DefaultPrefix = "Database"

func DefaultProvider() container.ProviderContainer {
	return container.ProviderContainer{
		Provider: Provide,
		Options: []opt.Option{
			opt.Prefix(DefaultPrefix),
		},
	}
}

type Config struct {
	Username     string
	Password     string
	Database     string
	Schema       string
	Host         string
	Port         uint
	SslMode      string
	TimeZone     string
	Prefix       string // 表前缀
	Singular     bool   // 是否开启全局禁用复数，true表示开启
	MaxIdleConns int    // 空闲中的最大连接数
	MaxOpenConns int    // 打开到数据库的最大连接数
	// 可选：连接生命周期配置（0 表示不设置）
	ConnMaxLifetimeSeconds uint
	ConnMaxIdleTimeSeconds uint

	// 可选：GORM 日志与行为配置
	LogLevel               string // silent|error|warn|info（默认info）
	SlowThresholdMs        uint   // 慢查询阈值（毫秒）默认200
	ParameterizedQueries   bool   // 占位符输出，便于日志安全与查询归并
	PrepareStmt            bool   // 预编译语句缓存
	SkipDefaultTransaction bool   // 跳过默认事务

	// 可选：DSN 增强
	UseSearchPath   bool   // 在 DSN 中附带 search_path
	ApplicationName string // application_name
}

func (m Config) GormSlowThreshold() time.Duration {
	if m.SlowThresholdMs == 0 {
		return 200 * time.Millisecond // 默认200ms
	}
	return time.Duration(m.SlowThresholdMs) * time.Millisecond
}

func (m Config) GormLogLevel() logger.LogLevel {
	switch m.LogLevel {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info", "":
		return logger.Info
	default:
		return logger.Info
	}
}

func (m *Config) checkDefault() {
	if m.MaxIdleConns == 0 {
		m.MaxIdleConns = 10
	}

	if m.MaxOpenConns == 0 {
		m.MaxOpenConns = 100
	}

	if m.Username == "" {
		m.Username = "postgres"
	}

	if m.SslMode == "" {
		m.SslMode = "disable"
	}

	if m.TimeZone == "" {
		m.TimeZone = "Asia/Shanghai"
	}

	if m.Port == 0 {
		m.Port = 5432
	}

	if m.Schema == "" {
		m.Schema = "public"
	}
}

func (m *Config) EmptyDsn() string {
	// 基本 DSN
	dsnTpl := "host=%s user=%s password=%s port=%d dbname=%s sslmode=%s TimeZone=%s"
	m.checkDefault()
	base := fmt.Sprintf(dsnTpl, m.Host, m.Username, m.Password, m.Port, m.Database, m.SslMode, m.TimeZone)
	// 附加可选参数
	extras := ""
	if m.UseSearchPath && m.Schema != "" {
		extras += " search_path=" + m.Schema
	}
	if m.ApplicationName != "" {
		extras += " application_name=" + strconv.Quote(m.ApplicationName)
	}
	return base + extras
}

// DSN connection dsn
func (m *Config) DSN() string {
	// 基本 DSN
	dsnTpl := "host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s"
	m.checkDefault()
	base := fmt.Sprintf(dsnTpl, m.Host, m.Username, m.Password, m.Database, m.Port, m.SslMode, m.TimeZone)
	// 附加可选参数
	extras := ""
	if m.UseSearchPath && m.Schema != "" {
		extras += " search_path=" + m.Schema
	}
	if m.ApplicationName != "" {
		extras += " application_name=" + strconv.Quote(m.ApplicationName)
	}
	return base + extras
}
