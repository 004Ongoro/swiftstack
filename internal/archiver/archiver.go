/*
Package archiver provides primitives for decompressing and extracting
project slices stored in .tar.zst format.
*/
package archiver

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zstd"
)

// Extract takes a source .tar.zst file and extracts it to the destination path.
func Extract(src io.Reader, dest string) error {
	// 1. Initialize Zstd decoder
	zr, err := zstd.NewReader(src)
	if err != nil {
		return fmt.Errorf("failed to create zstd reader: %w", err)
	}
	defer zr.Close()

	// 2. Initialize Tar reader on top of Zstd
	tr := tar.NewReader(zr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return fmt.Errorf("error reading tar header: %w", err)
		}

		// Clean the path to prevent zip slip vulnerabilities
		target := filepath.Join(dest, filepath.Clean(header.Name))

		switch header.Typeflag {
		case tar.TypeDir:
			// Create directory if it doesn't exist
			if err := os.MkdirAll(target, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", target, err)
			}
		case tar.TypeReg:
			// Ensure parent directory exists
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory for %s: %w", target, err)
			}

			// Create the file
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", target, err)
			}

			// Copy contents from tar to the new file
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return fmt.Errorf("failed to write file content for %s: %w", target, err)
			}
			f.Close()
		}
	}

	return nil
}