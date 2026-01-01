/*
Package engine handles the core logic of stitching project slices together.
resolver.go specifically manages semantic version comparisons.
*/
package engine

import (
	"strings"

	"github.com/blang/semver/v4"
)

// ResolveVersion takes two version strings and returns the most suitable one.
// It prioritizes the higher version to ensure modern features are available.
// If versions cannot be parsed (e.g., "latest", "next", or URLs), the slice version is preferred.
func ResolveVersion(baseVer, sliceVer string) string {
	// Clean the strings by removing common prefix symbols (^ or ~) for comparison
	bV := strings.TrimPrefix(strings.TrimPrefix(baseVer, "^"), "~")
	sV := strings.TrimPrefix(strings.TrimPrefix(sliceVer, "^"), "~")

	v1, err1 := semver.Make(bV)
	v2, err2 := semver.Make(sV)

	// If parsing fails for either version, we default to the slice version
	// as it's likely the "newer" addition to the project.
	if err1 != nil || err2 != nil {
		return sliceVer
	}

	// Compare and return the string representation of the higher version
	if v1.Compare(v2) >= 0 {
		return baseVer
	}
	return sliceVer
}
