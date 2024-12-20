package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/swaggo/swag"
	"github.com/swaggo/swag/gen"
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
	leftDelim, rightDelim := "{{", "}}"

	return gen.New().Build(&gen.Config{
		SearchDir:           "./",
		Excludes:            "",
		ParseExtension:      "",
		MainAPIFile:         "main.go",
		PropNamingStrategy:  swag.CamelCase,
		OutputDir:           "./docs",
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
