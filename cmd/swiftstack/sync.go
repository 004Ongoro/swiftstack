/*
sync.go defines the command to update the local registry manifest.
*/
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/004Ongoro/swiftstack/internal/cache"
	"github.com/004Ongoro/swiftstack/internal/utils"
)

// For now, this points to your personal repo or a placeholder
const ManifestURL = "https://raw.githubusercontent.com/004Ongoro/swiftstack/main/registry.json"

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Update the local slice registry",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Syncing with remote registry...")

		dest, err := cache.GetManifestPath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if err := utils.FetchRemoteManifest(ManifestURL, dest); err != nil {
			fmt.Fprintf(os.Stderr, "Sync failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Registry updated successfully!")
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
