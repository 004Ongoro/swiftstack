/*
Package utils provides helper functions for system operations.
fs.go handles complex file system operations like moving directories and backups.
*/
package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// MoveWithBackup moves files from src to dst. If a file exists in dst, 
// it renames the original to filename.bak before placing the new one.
func MoveWithBackup(srcDir, dstDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate the relative path from the source directory
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if relPath == "." {
			return nil
		}

		targetPath := filepath.Join(dstDir, relPath)

		if info.IsDir() {
			// Ensure the directory exists in the destination
			return os.MkdirAll(targetPath, 0755)
		}

		// If it's a file and it exists, create a backup
		if _, err := os.Stat(targetPath); err == nil {
			backupPath := targetPath + ".bak"
			if err := os.Rename(targetPath, backupPath); err != nil {
				return fmt.Errorf("fs: failed to create backup for %s: %w", targetPath, err)
			}
		}

		// Move the file
		return os.Rename(path, targetPath)
	})
}

// CopyFile is a helper for manual file copies if rename fails across partitions
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}