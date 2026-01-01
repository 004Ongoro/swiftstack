/*
wizard.go defines the 'wizard' subcommand for an interactive setup.
*/
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/004Ongoro/swiftstack/internal/ui"
)

var wizardCmd = &cobra.Command{
	Use:   "ui",
	Short: "Start the interactive SwiftStack wizard",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(ui.InitialModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(wizardCmd)
}