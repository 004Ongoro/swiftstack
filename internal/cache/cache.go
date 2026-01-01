/*
Package cache handles the local storage and retrieval of project slices.
It ensures that slices are stored in the standard OS cache directories.
*/
package cache

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetCacheDir returns the OS-standard path for SwiftStack data.
func GetCacheDir() (string, error) {
	parent, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("cache: could not determine user cache dir: %w", err)
	}

	path := filepath.Join(parent, "swiftstack")
	
	// Ensure the directory exists
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", fmt.Errorf("cache: failed to create cache directory: %w", err)
	}

	return path, nil
}

// GetSlicePath returns the full local path for a specific slice version.
func GetSlicePath(sliceID, version string) (string, error) {
	dir, err := GetCacheDir()
	if err != nil {
		return "", err
	}
	
	// Slices are stored as: id@version.tar.zst
	filename := fmt.Sprintf("%s@%s.tar.zst", sliceID, version)
	return filepath.Join(dir, filename), nil
}