/*
Package utils provides cryptographic and security helpers.
*/
package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// VerifyFileHash compares the SHA-256 hash of a file against an expected string.
func VerifyFileHash(filePath, expectedHash string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("crypto: could not open file for hashing: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return fmt.Errorf("crypto: failed to calculate hash: %w", err)
	}

	actualHash := hex.EncodeToString(h.Sum(nil))

	if actualHash != expectedHash {
		return fmt.Errorf("integrity error: hash mismatch!\nExpected: %s\nActual:   %s", expectedHash, actualHash)
	}

	return nil
}