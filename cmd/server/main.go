package main

import (
	"ConcurrentFileUploader/uploader"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	// start Background Server
	go func() {
		http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
			// Simulate processing time
			time.Sleep(500 * time.Millisecond)

			// Read body to simulate receiving
			_, err := io.Copy(io.Discard, r.Body)
			if err != nil {
				fmt.Printf("Receiving failed: %v\n", err)
				return
			}
			err = r.Body.Close()
			if err != nil {
				fmt.Printf("Receiving failed: %v\n", err)
				return
			}

			fmt.Println("Server: Received upload request")
			w.WriteHeader(http.StatusOK)
		})

		fmt.Println("Server: Listening on Port:8080...")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Printf("Server failed: %v\n", err)
		}
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Create Dummy Files
	fileCount := 100
	files := make([]string, fileCount)
	fmt.Printf("Generating %d dummy files...\n", fileCount)

	for i := 0; i < fileCount; i++ {
		name := fmt.Sprintf("job_%03d.dat", i)
		files[i] = name
		_ = os.WriteFile(name, []byte("data"), 0644)
		// Clean up when main exits
		defer func(name string) {
			err := os.Remove(name)
			if err != nil {
				fmt.Printf("Error closing files: %v\n", err)
			}
		}(name)
	}
	start := time.Now()

	// Setup Uploader
	up := &uploader.HttpUploader{
		UploadURL: "http://localhost:8080/upload",
		Client:    http.DefaultClient,
	}

	// Run Concurrent Uploads
	concurrency := 10
	results := uploader.UploadFiles(context.Background(), files, up, concurrency)

	// Process Results
	successCount := 0
	failCount := 0
	for res := range results {
		if res.Err != nil {
			fmt.Printf("%s failed: %v\n", res.FilePath, res.Err)
			failCount++
		} else {
			fmt.Printf("%s uploaded successfully\n", res.FilePath)
			successCount++
		}
	}

	totalTime := time.Since(start)
	fmt.Println("\n" + "----------------------------------------")
	fmt.Println("    UPLOAD SUMMARY")
	fmt.Println("----------------------------------------")
	fmt.Printf("Total Files:    %d\n", fileCount)
	fmt.Printf("Concurrency:    %d workers\n", concurrency)
	fmt.Printf("Successful:     %d\n", successCount)
	fmt.Printf("Failed:         %d\n", failCount)
	fmt.Printf("Total Duration: %s\n", totalTime)
	fmt.Printf("Throughput:     %.2f files/sec\n", float64(fileCount)/totalTime.Seconds())
	fmt.Println("----------------------------------------")
}
