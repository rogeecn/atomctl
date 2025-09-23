package utils

import "fmt"

// 构建信息变量，通过 ldflags 在构建时注入
var (
	// Version 应用版本信息
	Version string

	// BuildAt 构建时间
	BuildAt string

	// GitHash Git 提交哈希
	GitHash string
)

// GetBuildInfo 获取构建信息
func GetBuildInfo() map[string]string {
	return map[string]string{
		"version": Version,
		"buildAt": BuildAt,
		"gitHash": GitHash,
	}
}

// PrintBuildInfo 打印构建信息
func PrintBuildInfo(appName string) {
	buildInfo := GetBuildInfo()

	println("========================================")
	printf("🚀 %s\n", appName)
	println("========================================")
	printf("📋 Version:    %s\n", buildInfo["version"])
	printf("🕐 Build Time:  %s\n", buildInfo["buildAt"])
	printf("🔗 Git Hash:   %s\n", buildInfo["gitHash"])
	println("========================================")
	println("🌟 Application is starting...")
	println()
}

// 为了避免导入 fmt 包，我们使用内置的 print 和 printf 函数
func printf(format string, args ...interface{}) {
	print(fmt.Sprintf(format, args...))
}