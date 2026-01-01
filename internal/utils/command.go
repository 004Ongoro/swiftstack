/*
Package utils provides helper functions for system operations.
command.go specifically handles executing external processes like npm.
*/
package utils

import (
	"fmt"
	"os/exec"
)

// RunNpmLockUpdate runs 'npm install --package-lock-only' in the specified directory.
// This ensures that the lockfile is regenerated to match our merged package.json
// without performing a full network download.
func RunNpmLockUpdate(dir string) error {
	// Check if npm is even installed first
	_, err := exec.LookPath("npm")
	if err != nil {
		return fmt.Errorf("npm not found in PATH: please install Node.js to use SwiftStack")
	}

	// Prepare the command
	cmd := exec.Command("npm", "install", "--package-lock-only")
	cmd.Dir = dir

	// Execute and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("npm lock update failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}