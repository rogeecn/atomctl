package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/rogeecn/atomctl/templates/run_seeder"
	"github.com/rogeecn/atomctl/utils"
	"github.com/spf13/cobra"
)

func init() {
	runCmd.AddCommand(runSeederCmd)
	runSeederCmd.Flags().StringVarP(&seederRunner.Driver, "driver", "D", "postgres", "define driver: postgres/mysql/sqlite")
	runSeederCmd.Flags().StringVarP(&seederRunner.AppName, "app", "A", "", "app name")
}

type SeederRunner struct {
	Pkg     string
	AppName string
	Driver  string
}

var seederRunner = &SeederRunner{}

var runSeederCmd = &cobra.Command{
	Use:   "seeder",
	Short: "atomctl run seeder [--driver postgres]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !utils.IsFile("go.mod") {
			return errors.New("run in project root directory")
		}

		seederRunner.Pkg = getPackage()
		if seederRunner.AppName == "" {
			pkgs := strings.Split(seederRunner.Pkg, "/")
			seederRunner.AppName = pkgs[len(pkgs)-1]
		}

		tpl, err := template.New("seeder").Parse(run_seeder.Templates)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		if err := tpl.Execute(&buf, seederRunner); err != nil {
			return err
		}

		file := "database/main.go"
		defer os.Remove(file)
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
