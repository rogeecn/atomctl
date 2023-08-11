/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rogeecn/atomctl/utils"
	"github.com/rogeecn/atomctl/utils/generator"
	"github.com/spf13/cobra"
)

var genEnumCmd = &cobra.Command{
	Use:     "enum",
	Short:   "generate enum for consts",
	Example: "atomctl gen enum",
	RunE: func(cmd *cobra.Command, args []string) error {
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

			content, err := os.ReadFile(path)
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

		// g.WithMarshal()
		g.WithFlag()
		g.WithSQLDriver()
		g.WithSQLInt()
		g.WithSQLNullInt()
		g.WithSQLNullStr()
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
	},
}

func init() {
	genCmd.AddCommand(genEnumCmd)
}
