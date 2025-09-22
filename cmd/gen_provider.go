package cmd

import (
	"os"
	"path/filepath"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.ipao.vip/atomctl/v2/pkg/ast/provider"
	"go.ipao.vip/atomctl/v2/pkg/utils/gomod"
)

// CommandGenProvider 创建并注册 provider 生成命令到根命令
//
// 命令功能：
//   - 扫描源码中带有 @provider 注释的结构体
//   - 生成 provider.gen.go 文件实现依赖注入与分组注册
//   - 支持多种 Provider 模式（grpc、event、job、cronjob、model）
//   - 自动处理导入依赖和注入配置
//
// 命令设计：
// ┌─────────────────────────────────────────────────────────────┐
// │                  CommandGenProvider()                        │
// ├─────────────────────────────────────────────────────────────┤
// │  1. 命令创建：&cobra.Command{...}                           │
// │  2. 配置设置：Use、Aliases、Short、Long                     │
// │  3. 执行函数：commandGenProviderE                           │
// │  4. 命令注册：root.AddCommand(cmd)                         │
// └─────────────────────────────────────────────────────────────┘
//
// 命令配置：
//   - Use: "provider" - 命令名称
//   - Aliases: ["p"] - 简写别名
//   - Short: "Generate providers" - 简短描述
//   - Long: 详细的帮助文档，包含语法说明和模式特性
//   - RunE: commandGenProviderE - 命令执行函数
//
// 注释语法说明：
//
//	@provider(<mode>):[except|only] [returnType] [group]
//	- mode: grpc|event|job|cronjob|model（可选）
//	- :only: 仅注入字段 tag 为 inject:"true" 的依赖
//	- :except: 注入除标注 inject:"false" 之外的非标量依赖
//	- returnType: Provide 返回类型（如 contracts.Initial）
//	- group: 分组（如 atom.GroupInitial）
//
// 参数：
//   - root: 根命令对象，用于注册子命令
//
// 注意事项：
//   - 命令会在当前目录或指定目录中扫描 Provider
//   - 生成的文件按包分组，每个包生成一个 provider.gen.go
//   - 可与 gen route 命令联动使用
func CommandGenProvider(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "provider",
		Aliases: []string{"p"},
		Short:   "Generate providers",
		Long: `扫描源码中带有 @provider 注释的结构体，生成 provider.gen.go 实现依赖注入与分组注册。

注释语法：
  @provider(<mode>):[except|only] [returnType] [group]
  - mode：grpc|event|job|cronjob|model（可为空）
  - :only：仅注入字段 tag 为 inject:"true" 的依赖
  - :except：注入除标注 inject:"false" 之外的非标量依赖
  - returnType：Provide 返回类型（如 contracts.Initial）
  - group：分组（如 atom.GroupInitial）

模式特性：
  - grpc：注入 providers/grpc.Grpc，设置 GrpcRegisterFunc
  - event：注入 providers/event.PubSub
  - job|cronjob：注入 providers/job.Job，并引入 github.com/riverqueue/river
  - model：需要 Prepare 钩子

行为说明：
  - 忽略标量类型字段，自动补全 imports
  - 以包为单位生成 provider.gen.go
  - 可与 gen route 联动（其 PostRun 会触发本命令）`,
		RunE: commandGenProviderE,
	}

	root.AddCommand(cmd)
}

// commandGenProviderE 实现 provider 生成命令的核心执行逻辑
//
// 执行流程：
// ┌─────────────────────────────────────────────────────────────┐
// │              commandGenProviderE(cmd, args)                │
// ├─────────────────────────────────────────────────────────────┤
// │  1. 路径解析：处理命令行参数，确定扫描目录                 │
// │  2. 路径标准化：转换为绝对路径                             │
// │  3. 模块解析：解析 go.mod 文件获取模块信息                 │
// │  4. Provider 解析：使用重构后的解析器扫描目录              │
// │  5. 分组处理：按输出文件分组 Provider                       │
// │  6. 文件生成：调用 provider.Render() 生成代码               │
// │  7. 错误处理：统一的错误处理和日志记录                     │
// └─────────────────────────────────────────────────────────────┘
//
// 详细步骤说明：
// 步骤 1：路径解析
//   - 检查命令行参数，确定扫描目标目录
//   - 如果没有提供参数，使用当前工作目录
//   - 处理路径获取失败的情况
//
// 步骤 2：路径标准化
//   - 将相对路径转换为绝对路径
//   - 确保路径的一致性和可处理性
//   - 避免相对路径带来的处理问题
//
// 步骤 3：模块解析
//   - 解析目标目录中的 go.mod 文件
//   - 获取模块名称和依赖信息
//   - 为后续的 Provider 解析提供上下文
//
// 步骤 4：Provider 解析（重构后的核心）
//   - 创建 GoParser 实例：provider.NewGoParser()
//   - 解析目录：parser.ParseDir(path)
//   - 错误处理：解析失败时记录详细错误信息
//   - 性能优化：使用重构后的高效解析器
//
// 步骤 5：分组处理
//   - 按输出文件分组：使用 ProviderFile 字段
//   - 使用 samber/lo 库的 GroupBy 函数
//   - 为每个输出文件准备对应的 Provider 配置
//
// 步骤 6：文件生成
//   - 遍历分组后的 Provider 配置
//   - 调用 provider.Render() 生成代码
//   - 处理文件生成过程中的错误
//
// 重构优化说明：
//   - 简化了代码结构：从原来的手动文件遍历简化为 6 行核心代码
//   - 提高了性能：使用重构后的高效解析器
//   - 增强了可维护性：解析逻辑封装在独立的模块中
//   - 统一了错误处理：所有解析错误都有统一的处理方式
//
// 参数：
//   - cmd: 命令对象，包含命令配置和上下文
//   - args: 命令行参数，第一个参数为目录路径
//
// 返回值：
//   - error: 执行过程中的错误，成功时返回 nil
//
// 错误处理策略：
//   - 路径错误：返回 os.Getwd() 的错误
//   - 模块解析错误：返回 gomod.Parse() 的错误
//   - Provider 解析错误：返回解析器的错误，包含详细上下文
//   - 文件生成错误：返回 provider.Render() 的错误
//
// 使用示例：
//
//	# 在当前目录生成 Provider
//	atomctl gen provider
//
//	# 在指定目录生成 Provider
//	atomctl gen provider ./internal/services
//
// 注意事项：
//   - 目标目录必须包含 go.mod 文件
//   - 生成的文件会覆盖现有的 provider.gen.go 文件
//   - 建议在版本控制中提交生成的文件
func commandGenProviderE(cmd *cobra.Command, args []string) error {
	var err error
	var path string

	// === 步骤 1：路径解析 ===
	// 处理命令行参数，确定要扫描的目录
	if len(args) > 0 {
		// 使用用户提供的路径参数
		path = args[0]
	} else {
		// 没有提供参数时，使用当前工作目录
		path, err = os.Getwd()
		if err != nil {
			log.Errorf("Failed to get working directory: %v", err)
			return err
		}
	}

	// === 步骤 2：路径标准化 ===
	// 将路径转换为绝对路径，确保路径的一致性
	path, _ = filepath.Abs(path)

	// === 步骤 3：模块解析 ===
	// 解析 go.mod 文件，获取模块信息和依赖关系
	err = gomod.Parse(filepath.Join(path, "go.mod"))
	if err != nil {
		log.Errorf("Failed to parse go.mod file: %v", err)
		return err
	}

	// === 步骤 4：Provider 解析（重构后的核心） ===
	// 记录开始生成的日志信息
	log.Infof("generate providers for dir: %s", path)

	// 使用重构后的 GoParser 解析目录中的 Provider
	parser := provider.NewGoParser()
	providers, err := parser.ParseDir(path)
	if err != nil {
		log.Errorf("Failed to parse providers from directory %s: %v", path, err)
		return err
	}

	// === 步骤 5：分组处理 ===
	// 按输出文件分组 Provider，每个包生成一个 provider.gen.go
	groups := lo.GroupBy(providers, func(item provider.Provider) string {
		return item.ProviderFile
	})

	// === 步骤 6：文件生成 ===
	// 遍历分组后的 Provider 配置，生成对应的代码文件
	for file, conf := range groups {
		if err := provider.Render(file, conf); err != nil {
			log.Errorf("Failed to render provider file %s: %v", file, err)
			return err
		}
	}

	log.Infof("Successfully generated %d provider files", len(groups))
	return nil
}
