package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/rogeecn/atomctl/templates/model"
	"github.com/rogeecn/atomctl/utils"
	"github.com/spf13/cobra"
)

type ModelGenerator struct {
	Driver  string
	AppName string
}

var modelGenerator = &ModelGenerator{}

func init() {
	genCmd.AddCommand(genModelCmd)
	genCrudCmd.Flags().StringVarP(&modelGenerator.Driver, "driver", "D", "postgres", "define driver: postgres/mysql/sqlite")
}

var genModelCmd = &cobra.Command{
	Use:     "model",
	Short:   "generate model",
	Example: "atomctl gen model [--driver postgres]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !utils.IsFile("go.mod") {
			return errors.New("run in project root directory")
		}

		pkg := getPackage()
		pkgs := strings.Split(pkg, "/")
		modelGenerator.AppName = pkgs[len(pkgs)-1]

		tpl, err := template.New("model").Parse(model.Templates)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		if err := tpl.Execute(&buf, modelGenerator); err != nil {
			return err
		}

		file := "database/main.go"
		// defer os.Remove(file)
		if err := os.WriteFile(file, buf.Bytes(), os.ModePerm); err != nil {
			return err
		}

		cmder := exec.Command("go", "run", "database/main.go")
		stdout, _ := cmder.StdoutPipe()
		stderr, _ := cmder.StderrPipe()
		if err := cmder.Start(); err != nil {
			return err
		}

		go func() {
			scanner := bufio.NewScanner(stderr)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
		}()

		go func() {
			scanner := bufio.NewScanner(stdout)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
		}()
		return cmder.Wait()
	},
}
