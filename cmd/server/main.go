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
	files := []string{"doc1.txt", "doc2.txt", "doc3.txt", "doc4.txt", "doc5.txt"}
	for i, name := range files {
		err := os.WriteFile(name, []byte(fmt.Sprintf("Hello World: [%d]", i)), 0644)
		if err != nil {
			fmt.Printf("File creation failed: %v\n", err)
			return
		}
		// Defer cleanup
		defer func(name string) {
			err := os.Remove(name)
			if err != nil {
				fmt.Printf("File deletion failed: %v\n", err)
			}
		}(name)
	}

	// Setup Uploader
	up := &uploader.HttpUploader{
		UploadURL: "http://localhost:8080/upload",
		Client:    http.DefaultClient,
	}
	bg := context.Background()

	// Run Concurrent Uploads
	results := uploader.UploadFiles(bg, files, up, 2)

	// Process Results
	for res := range results {
		if res.Err != nil {
			fmt.Printf("%s failed: %v\n", res.FilePath, res.Err)
		} else {
			fmt.Printf("%s uploaded successfully\n", res.FilePath)
		}
	}

	fmt.Println("All uploads finished!")
}
