/*
Package engine handles the core logic of stitching project slices together.
merger.go manages the reading and writing of package.json files.
*/
package engine

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/004Ongoro/swiftstack/internal/models"
)

// MergePackageJSON reads two package.json files and merges their dependencies,
// devDependencies, and scripts. It uses ResolveVersion to handle conflicts.
func MergePackageJSON(basePath, slicePath string) error {
	// Load the base package.json (the target)
	baseData, err := readJSON(basePath)
	if err != nil {
		return fmt.Errorf("merger: base file error: %w", err)
	}

	// Load the slice package.json (the source of new features)
	sliceData, err := readJSON(slicePath)
	if err != nil {
		return fmt.Errorf("merger: slice file error: %w", err)
	}

	// 1. Merge standard Dependencies
	if baseData.Dependencies == nil {
		baseData.Dependencies = make(map[string]string)
	}
	for pkg, ver := range sliceData.Dependencies {
		if existingVer, exists := baseData.Dependencies[pkg]; exists {
			baseData.Dependencies[pkg] = ResolveVersion(existingVer, ver)
		} else {
			baseData.Dependencies[pkg] = ver
		}
	}

	// 2. Merge DevDependencies
	if baseData.DevDependencies == nil {
		baseData.DevDependencies = make(map[string]string)
	}
	for pkg, ver := range sliceData.DevDependencies {
		if existingVer, exists := baseData.DevDependencies[pkg]; exists {
			baseData.DevDependencies[pkg] = ResolveVersion(existingVer, ver)
		} else {
			baseData.DevDependencies[pkg] = ver
		}
	}

	// 3. Merge Scripts
	if baseData.Scripts == nil {
		baseData.Scripts = make(map[string]string)
	}
	for name, command := range sliceData.Scripts {
		// Note: We currently overwrite scripts if they have the same name.
		baseData.Scripts[name] = command
	}

	// Write the final merged object back to the base path
	return writeJSON(basePath, baseData)
}

// readJSON is a private helper to read and unmarshal a package.json file.
func readJSON(path string) (*models.PackageJSON, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("readJSON: failed to read %s: %w", path, err)
	}

	var pkg models.PackageJSON
	if err := json.Unmarshal(file, &pkg); err != nil {
		return nil, fmt.Errorf("readJSON: failed to unmarshal %s: %w", path, err)
	}
	return &pkg, nil
}

// writeJSON is a private helper to marshal and save the package.json with 2-space indentation.
func writeJSON(path string, pkg *models.PackageJSON) error {
	data, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return fmt.Errorf("writeJSON: failed to marshal: %w", err)
	}

	// Ensure we end with a newline to follow standard JSON formatting
	data = append(data, '\n')

	return os.WriteFile(path, data, 0644)
}
