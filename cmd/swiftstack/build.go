/*
build.go defines the 'build' subcommand and now outputs the SHA-256 hash.
*/
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/004Ongoro/swiftstack/internal/builder"
)

var buildCmd = &cobra.Command{
	Use:   "build [source_dir] [output_file.tar.zst]",
	Short: "Compress a directory and output its SHA-256 hash",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		src := args[0]
		dest := args[1]

		fmt.Printf("Building slice from %s...\n", src)
		if err := builder.CreateSlice(src, dest); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Calculate hash for the manifest
		hash, err := calculateHash(dest)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error calculating hash: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nSuccessfully created slice!\nLocation: %s\nSHA-256:  %s\n", dest, hash)
		fmt.Println("\nhash to  registry.json to ensure security.")
	},
}

func calculateHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func init() {
	rootCmd.AddCommand(buildCmd)
}