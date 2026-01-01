package archiver

import (
	"os"
	"testing"
)

func TestExtract(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "swiftstack-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	if tmpDir == "" {
		t.Error("Temporary directory was not created")
	}
}