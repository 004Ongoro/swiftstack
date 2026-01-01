/*
Package models defines the data structures for SwiftStack metadata.
*/
package models

// SliceMetadata defines the properties of a project slice.
type SliceMetadata struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Version     string `json:"version"`
	Hash        string `json:"hash"` // SHA-256 hash for integrity verification
}

// RemoteManifest is the structure of the master list hosted online.
type RemoteManifest struct {
	Bases  []SliceMetadata `json:"bases"`
	Addons []SliceMetadata `json:"addons"`
}
