/*
Package builder provides tools for creating SwiftStack slices.
It handles the compression of directories into .tar.zst format.
*/
package builder

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zstd"
)

// CreateSlice packs a source directory into a target .tar.zst file.
func CreateSlice(srcDir, targetFile string) error {
	// 1. Create the output file
	f, err := os.Create(targetFile)
	if err != nil {
		return fmt.Errorf("builder: failed to create output file: %w", err)
	}
	defer f.Close()

	// 2. Initialize Zstd writer
	zw, err := zstd.NewWriter(f)
	if err != nil {
		return fmt.Errorf("builder: failed to create zstd writer: %w", err)
	}
	defer zw.Close()

	// 3. Initialize Tar writer
	tw := tar.NewWriter(zw)
	defer tw.Close()

	// 4. Walk the source directory and add files to tar
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create a tar header
		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		// Use relative paths inside the archive
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		header.Name = relPath

		// Write header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// If it's a directory, we're done with this iteration
		if info.IsDir() {
			return nil
		}

		// If it's a file, copy content
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(tw, file)
		return err
	})
}
