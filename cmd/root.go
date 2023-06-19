/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "atomctl",
	Short: "atomctl",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getPackage() string {
	content, err := ioutil.ReadFile("go.mod")
	if err != nil {
		log.Fatal(err)
	}

	// Parse the go.mod file
	f, err := modfile.Parse("go.mod", content, nil)
	if err != nil {
		log.Fatal(err)
	}

	return f.Module.Mod.Path
}

func dotToModule(raw string) (string, string) {
	items := append([]string{""}, strings.Split(raw, ".")...)
	path := strings.Join(items, "/modules/")
	name := items[len(items)-1]
	return strings.Trim(path, "/"), name
}
