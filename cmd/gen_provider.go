package cmd

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.ipao.vip/atomctl/v2/pkg/ast/provider"
	"go.ipao.vip/atomctl/v2/pkg/utils/gomod"
)

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

func commandGenProviderE(cmd *cobra.Command, args []string) error {
	var err error
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	path, _ = filepath.Abs(path)

	err = gomod.Parse(filepath.Join(path, "go.mod"))
	if err != nil {
		return err
	}

	providers := []provider.Provider{}

	// if path is file, then get the dir
	log.Infof("generate providers for dir: %s", path)
	// travel controller to find all controller objects
	_ = filepath.WalkDir(path, func(filepath string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(filepath, ".go") {
			return nil
		}

		if strings.HasSuffix(filepath, "_test.go") {
			return nil
		}

		providers = append(providers, provider.Parse(filepath)...)
		return nil
	})

	// generate files
	groups := lo.GroupBy(providers, func(item provider.Provider) string {
		return item.ProviderFile
	})

	for file, conf := range groups {
		if err := provider.Render(file, conf); err != nil {
			return err
		}
	}
	return nil
}
