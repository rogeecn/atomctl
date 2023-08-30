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
	"github.com/rogeecn/atomctl/templates/run_migrate"
	"github.com/rogeecn/atomctl/utils"
	"github.com/spf13/cobra"
)

type MigrationRunner struct {
	Pkg         string
	AppName     string
	Driver      string
	MigrateToId string
}

var migrationRunner = &MigrationRunner{}

func init() {
	runCmd.AddCommand(runMigrateCmd)
	runMigrateCmd.AddCommand(runMigrateUpCmd)
	runMigrateCmd.AddCommand(runMigrateDownCmd)

	runMigrateCmd.PersistentFlags().StringVarP(&migrationRunner.Driver, "driver", "D", "mysql", "define driver: postgres/mysql/sqlite")
	runMigrateCmd.PersistentFlags().StringVarP(&migrationRunner.AppName, "app", "A", "", "app name")
	runMigrateCmd.PersistentFlags().StringVar(&migrationRunner.MigrateToId, "to", "", "migrate to id")
}

var runMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate database tables",
	Long:  `migrate database tables`,
}

var runMigrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "migrate up database tables",
	Long:  `migrate up database tables`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !utils.IsFile("go.mod") {
			return errors.New("run in project root directory")
		}

		migrationRunner.Pkg = getPackage()
		if migrationRunner.AppName == "" {
			pkgs := strings.Split(migrationRunner.Pkg, "/")
			migrationRunner.AppName = pkgs[len(pkgs)-1]
		}

		tpl, err := template.New("migrate").Parse(run_migrate.Up)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		if err := tpl.Execute(&buf, migrationRunner); err != nil {
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

var runMigrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "migrate down database tables",
	Long:  `migrate down database tables`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !utils.IsFile("go.mod") {
			return errors.New("run in project root directory")
		}

		migrationRunner.Pkg = getPackage()
		if migrationRunner.AppName == "" {
			pkgs := strings.Split(migrationRunner.Pkg, "/")
			migrationRunner.AppName = pkgs[len(pkgs)-1]
		}

		tpl, err := template.New("migrate").Parse(run_migrate.Down)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		if err := tpl.Execute(&buf, migrationRunner); err != nil {
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
