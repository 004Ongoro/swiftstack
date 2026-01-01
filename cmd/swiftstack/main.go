/*
Package main is the entry point for the SwiftStack CLI.
SwiftStack is a high-speed project scaffolder designed to bypass
heavy dependency installation by stitching pre-built binary slices.
*/
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "swiftstack",
	Short: "SwiftStack is a blazingly fast project scaffolder",
	Long: `A high-performance tool built in Go that assembles web projects 
using pre-cached binary slices, eliminating the need for slow npm installs 
on every new project setup.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no arguments are provided, show help
		cmd.Help()
	},
}