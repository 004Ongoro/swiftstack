/*
Package engine orchestrates the entire stitching process.
This version includes SHA-256 integrity verification for all slices.
*/
package engine

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/004Ongoro/swiftstack/internal/archiver"
	"github.com/004Ongoro/swiftstack/internal/cache"
	"github.com/004Ongoro/swiftstack/internal/utils"
)

type ProjectOptions struct {
	Name        string
	OutputPath  string
	BaseSlice   string
	AddonSlices []string
}

func GenerateProject(opts ProjectOptions) error {
	fullPath := filepath.Join(opts.OutputPath, opts.Name)
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return fmt.Errorf("engine: failed to create project dir: %w", err)
	}

	var success bool
	defer func() {
		if !success {
			fmt.Printf("Cleaning up failed project directory: %s\n", fullPath)
			os.RemoveAll(fullPath)
		}
	}()

	// 1. Resolve and Verify Slices
	basePath, err := ensureSlice(opts.BaseSlice)
	if err != nil {
		return err
	}

	var addonPaths []string
	for _, addon := range opts.AddonSlices {
		path, err := ensureSlice(addon)
		if err != nil {
			return err
		}
		addonPaths = append(addonPaths, path)
	}

	// 2. Extract Base
	if err := extractSlice(basePath, fullPath); err != nil {
		return err
	}

	// 3. Process Addons
	for i, slicePath := range addonPaths {
		tempAddonDir := filepath.Join(fullPath, ".swiftstack_temp")
		os.MkdirAll(tempAddonDir, 0755)

		if err := extractSlice(slicePath, tempAddonDir); err != nil {
			os.RemoveAll(tempAddonDir)
			return err
		}

		basePkg := filepath.Join(fullPath, "package.json")
		slicePkg := filepath.Join(tempAddonDir, "package.json")
		
		if _, err := os.Stat(slicePkg); err == nil {
			MergePackageJSON(basePkg, slicePkg)
			os.Remove(slicePkg)
		}

		utils.MoveWithBackup(tempAddonDir, fullPath)
		os.RemoveAll(tempAddonDir)
	}

	// 4. Finalize
	fmt.Println("Finalizing project structure...")
	utils.RunNpmLockUpdate(fullPath)

	success = true
	return nil
}

func ensureSlice(alias string) (string, error) {
	// 1. Resolve metadata from manifest
	m, err := cache.LoadManifest()
	if err != nil {
		return "", err
	}

	var targetHash, url string
	found := false

	// Search manifest for the hash and URL
	for _, b := range m.Bases {
		if b.ID == alias {
			targetHash = b.Hash
			url = b.URL
			found = true
			break
		}
	}
	if !found {
		for _, a := range m.Addons {
			if a.ID == alias {
				targetHash = a.Hash
				url = a.URL
				found = true
				break
			}
		}
	}

	if !found {
		return "", fmt.Errorf("engine: slice %s not found in registry", alias)
	}

	cachePath, _ := cache.GetSlicePath(alias, "latest")

	// 2. Download if missing
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		fmt.Printf("Downloading %s...\n", alias)
		if err := utils.DownloadFileConcurrent(url, cachePath, 4); err != nil {
			return "", err
		}
	}

	// 3. VERIFY INTEGRITY (Senior Level Security)
	fmt.Printf("Verifying integrity of %s...\n", alias)
	if err := utils.VerifyFileHash(cachePath, targetHash); err != nil {
		// If hash fails, delete the corrupted file so it can be re-downloaded
		os.Remove(cachePath)
		return "", fmt.Errorf("security alert: %w", err)
	}

	return cachePath, nil
}

func extractSlice(slicePath, dest string) error {
	file, err := os.Open(slicePath)
	if err != nil {
		return err
	}
	defer file.Close()
	return archiver.Extract(file, dest)
}