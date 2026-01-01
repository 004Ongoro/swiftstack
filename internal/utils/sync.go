/*
Package utils handles synchronization of remote resources.
*/
package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// FetchRemoteManifest downloads the latest manifest from the provided URL 
// and saves it to the local destination.
func FetchRemoteManifest(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("sync: failed to fetch manifest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("sync: registry server returned %d", resp.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("sync: failed to create local manifest: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}