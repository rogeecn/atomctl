/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
)

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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.atomctl.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
