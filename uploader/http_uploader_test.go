package uploader

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHttpUploader_Upload(t *testing.T) {
	// Create Fake Server that expects a POST
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		// Verify file content was sent
		body, _ := io.ReadAll(r.Body)
		if string(body) != "test content" {
			t.Errorf("expected 'test content', got %s", body)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// create a temporary file to upload
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Write some data to the file
	if _, err := tmpFile.WriteString("test content"); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Initialize real implementation pointing to the fake server
	u := &HttpUploader{
		UploadURL: ts.URL,
		Client:    http.DefaultClient,
	}

	// Run Upload
	err = u.Upload(context.Background(), tmpFile.Name())
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}
