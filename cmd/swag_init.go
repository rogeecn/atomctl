package cmd

import (
	"os"
	"path/filepath"

	"git.ipao.vip/rogeecn/atomctl/pkg/swag"
	"git.ipao.vip/rogeecn/atomctl/pkg/swag/gen"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CommandSwagInit(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "init",
		Short:   "swag init",
		Aliases: []string{"i"},
		RunE:    commandSwagInitE,
	}

	root.AddCommand(cmd)
}

func commandSwagInitE(cmd *cobra.Command, args []string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if len(args) > 0 {
		pwd = args[0]
	}

	leftDelim, rightDelim := "{{", "}}"

	return gen.New().Build(&gen.Config{
		SearchDir:           pwd,
		Excludes:            "",
		ParseExtension:      "",
		MainAPIFile:         "main.go",
		PropNamingStrategy:  swag.CamelCase,
		OutputDir:           filepath.Join(pwd, "docs"),
		OutputTypes:         []string{"go", "json", "yaml"},
		ParseVendor:         false,
		ParseDependency:     0,
		MarkdownFilesDir:    "",
		ParseInternal:       false,
		Strict:              false,
		GeneratedTime:       false,
		RequiredByDefault:   false,
		CodeExampleFilesDir: "",
		ParseDepth:          100,
		InstanceName:        "",
		OverridesFile:       ".swaggo",
		ParseGoList:         true,
		Tags:                "",
		LeftTemplateDelim:   leftDelim,
		RightTemplateDelim:  rightDelim,
		PackageName:         "",
		Debugger:            log.WithField("module", "swag.init"),
		CollectionFormat:    "csv",
		PackagePrefix:       "",
		State:               "",
		ParseFuncBody:       false,
	})
}
