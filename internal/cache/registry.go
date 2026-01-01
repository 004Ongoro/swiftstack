/*
Package cache handles local storage and remote resolution.
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

// ResolveAlias searches the manifest for a specific slice's URL.
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