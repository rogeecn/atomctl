package cmd

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"git.ipao.vip/rogeecn/atomctl/pkg/ast/provider"
	"git.ipao.vip/rogeecn/atomctl/pkg/utils/gomod"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CommandGenProvider(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "provider",
		Aliases: []string{"p"},
		Short:   "Generate providers",
		Long: `//  @provider
//  @provider:[except|only] [returnType] [group]
//  when except add tag: inject:"false"
//  when only add tag: inject:"true"`,
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
