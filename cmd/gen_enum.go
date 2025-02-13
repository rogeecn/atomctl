package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"go.ipao.vip/atomctl/pkg/utils"
	"go.ipao.vip/atomctl/pkg/utils/generator"
)

func CommandGenEnum(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:      "enum",
		Aliases:  []string{"e"},
		Short:    "Generate enums",
		RunE:     commandGenEnumE,
		PostRunE: commandGenProviderE,
	}

	cmd.Flags().BoolP("flag", "f", true, "Flag enum values")
	cmd.Flags().BoolP("marshal", "m", false, "Marshal enum values")
	cmd.Flags().BoolP("sql", "s", true, "SQL driver enum values")

	root.AddCommand(cmd)
}

func commandGenEnumE(cmd *cobra.Command, args []string) error {
	var filenames []string

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	err = filepath.Walk(pwd, func(path string, info fs.FileInfo, err error) error {
		if utils.IsDir(path) {
			return nil
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		var content []byte
		content, err = os.ReadFile(path)
		if err != nil {
			return err
		}

		if strings.Contains(string(content), "ENUM(") && strings.Contains(string(content), "swagger:enum") {
			filenames = append(filenames, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if len(filenames) == 0 {
		return fmt.Errorf("no enum files found in %s", pwd)
	}

	g := generator.NewGenerator()

	if marshal, _ := cmd.Flags().GetBool("marshal"); marshal {
		g.WithMarshal()
	}

	if flag, _ := cmd.Flags().GetBool("flag"); flag {
		g.WithFlag()
	}

	if sql, _ := cmd.Flags().GetBool("sql"); sql {
		g.WithSQLDriver()
		g.WithSQLInt()
		g.WithSQLNullInt()
		g.WithSQLNullStr()
	}

	g.WithNames()
	g.WithValues()

	for _, fileName := range filenames {
		log.Printf("Generating enums for %s", fileName)

		fileName, _ = filepath.Abs(fileName)
		outFilePath := fmt.Sprintf("%s.gen.go", strings.TrimSuffix(fileName, filepath.Ext(fileName)))

		// Parse the file given in arguments
		raw, err := g.GenerateFromFile(fileName)
		if err != nil {
			return fmt.Errorf("failed generating enums\nInputFile=%s\nError=%s", fileName, err)
		}

		// Nothing was generated, ignore the output and don't create a file.
		if len(raw) < 1 {
			continue
		}

		mode := int(0o644)
		err = os.WriteFile(outFilePath, raw, os.FileMode(mode))
		if err != nil {
			return fmt.Errorf("failed writing to file %s: %s", outFilePath, err)
		}
	}

	return nil
}
