/*
Package models defines the data structures used across SwiftStack.
This file specifically handles the structure of a Node.js package.json.
*/
package models

// PackageJSON represents the structure of a standard package.json file.
// We use map[string]string for dependencies to allow for easy merging.
type PackageJSON struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Description     string            `json:"description,omitempty"`
	Scripts         map[string]string `json:"scripts,omitempty"`
	Dependencies    map[string]string `json:"dependencies,omitempty"`
	DevDependencies map[string]string `json:"devDependencies,omitempty"`
}
