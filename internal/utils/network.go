/*
Package utils provides network and file system helpers.
network.go handles high-speed, multi-part downloads with progress tracking.
*/
package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

// DownloadResult represents the status of a chunk download.
type DownloadResult struct {
	Index int
	Error error
}

// DownloadFileConcurrent downloads a file using multiple parallel connections.
// It divides the file into 'chunks' to maximize bandwidth on slow/high-latency links.
func DownloadFileConcurrent(url string, destPath string, chunks int) error {
	// 1. Get the total file size first
	resp, err := http.Head(url)
	if err != nil {
		return fmt.Errorf("network: failed to reach registry: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("network: registry returned status %d", resp.StatusCode)
	}
	size := resp.ContentLength

	// 2. Create the destination file
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	var wg sync.WaitGroup
	chunkSize := size / int64(chunks)
	
	fmt.Printf("Downloading %s in %d parallel chunks...\n", url, chunks)

	for i := 0; i < chunks; i++ {
		wg.Add(1)
		
		start := int64(i) * chunkSize
		end := start + chunkSize - 1
		if i == chunks-1 {
			end = size - 1
		}

		go func(index int, start, end int64) {
			defer wg.Done()
			err := downloadChunk(url, out, start, end)
			if err != nil {
				fmt.Printf("Error downloading chunk %d: %v\n", index, err)
			}
		}(i, start, end)
	}

	wg.Wait()
	return nil
}

// downloadChunk fetches a specific byte range and writes it to the file at the correct offset.
func downloadChunk(url string, out *os.File, start, end int64) error {
	req, _ := http.NewRequest("GET", url, nil)
	rangeHeader := fmt.Sprintf("bytes=%d-%d", start, end)
	req.Header.Add("Range", rangeHeader)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write at the specific offset using WriteAt
	// This allows multiple goroutines to write to the same file concurrently without overlapping.
	_, err = io.Copy(newWriterAt(out, start), resp.Body)
	return err
}

// writerAt adapter to make os.File satisfy io.Writer for a specific offset
type writerAtAdapter struct {
	file   *os.File
	offset int64
}

func (w *writerAtAdapter) Write(p []byte) (n int, err error) {
	n, err = w.file.WriteAt(p, w.offset)
	w.offset += int64(n)
	return
}

func newWriterAt(f *os.File, off int64) io.Writer {
	return &writerAtAdapter{f, off}
}