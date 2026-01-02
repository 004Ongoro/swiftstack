/*
Package cache handles local storage and remote resolution of the project manifest.
*/
package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/004Ongoro/swiftstack/internal/models"
)

// GetManifestPath returns the local path to the synced manifest file.
func GetManifestPath() (string, error) {
	dir, err := GetCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "manifest.json"), nil
}

// LoadManifest reads the manifest from the local cache.
func LoadManifest() (*models.RemoteManifest, error) {
	path, err := GetManifestPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		// If manifest doesn't exist, return empty but not error
		// so the UI can prompt for a sync.
		return &models.RemoteManifest{}, nil
	}

	var m models.RemoteManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("registry: failed to parse manifest: %w", err)
	}
	return &m, nil
}

// GetAvailableBases is a helper for the UI to get the list of base templates.
func GetAvailableBases() []models.SliceMetadata {
	m, _ := LoadManifest()
	return m.Bases
}

// GetAvailableAddons is a helper for the UI to get the list of optional addons.
func GetAvailableAddons() []models.SliceMetadata {
	m, _ := LoadManifest()
	return m.Addons
}

// ResolveAlias searches the manifest for a specific slice's URL and Hash.
func ResolveAlias(alias string) (string, error) {
	m, err := LoadManifest()
	if err != nil {
		return "", err
	}

	// Search in Bases
	for _, b := range m.Bases {
		if b.ID == alias {
			return b.URL, nil
		}
	}

	// Search in Addons
	for _, a := range m.Addons {
		if a.ID == alias {
			return a.URL, nil
		}
	}

	return "", fmt.Errorf("alias '%s' not found in registry. Try running 'swiftstack sync'", alias)
}
