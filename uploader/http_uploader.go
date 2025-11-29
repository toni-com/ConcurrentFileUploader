package uploader

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

type HttpUploader struct {
	UploadURL string
	Client    *http.Client
}

func (h *HttpUploader) Upload(ctx context.Context, filePath string) error {
	// open file and make sure to close
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// create request and bind to ctx
	req, err := http.NewRequestWithContext(ctx, "POST", h.UploadURL, file)
	if err != nil {
		return err
	}

	// send file
	resp, err := h.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// handle status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned bad status: %d", resp.StatusCode)
	}

	return nil
}
