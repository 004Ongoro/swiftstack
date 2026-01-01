/*
create.go defines the primary user interface for generating new projects.
*/
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/004Ongoro/swiftstack/internal/engine"
)

var (
	projectName string
	baseAlias   string
	addonsList  []string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project using slice aliases",
	Example: "swiftstack create --name my-app --base next-base --addons tailwind-ui",
	Run: func(cmd *cobra.Command, args []string) {
		if projectName == "" {
			fmt.Println("Error: Project name is required (--name)")
			os.Exit(1)
		}

		if baseAlias == "" {
			fmt.Println("Error: Base slice alias is required (--base)")
			os.Exit(1)
		}

		options := engine.ProjectOptions{
			Name:        projectName,
			OutputPath:  ".",
			BaseSlice:   baseAlias,
			AddonSlices: addonsList,
		}

		fmt.Printf("🚀 Starting SwiftStack assembly for '%s'...\n", projectName)
		
		if err := engine.GenerateProject(options); err != nil {
			fmt.Fprintf(os.Stderr, "\n❌ Assembly Failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\n✨ Successfully assembled '%s' in record time!\n", projectName)
	},
}

func init() {
	createCmd.Flags().StringVarP(&projectName, "name", "n", "", "Name of the project")
	createCmd.Flags().StringVarP(&baseAlias, "base", "b", "", "Alias of the base slice (e.g., next-base)")
	createCmd.Flags().StringSliceVarP(&addonsList, "addons", "a", []string{}, "Comma-separated aliases (e.g., tailwind,auth)")

	rootCmd.AddCommand(createCmd)
}