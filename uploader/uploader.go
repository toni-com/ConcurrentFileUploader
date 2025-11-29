package uploader

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type UploadResult struct {
	FilePath string
	Err      error
}

type Uploader interface {
	Upload(ctx context.Context, filePath string) error
}

type SimpleUploader struct {
}

func (s *SimpleUploader) Upload(ctx context.Context, filePath string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if !isFilePathValid(filePath) {
		return errors.New("file path is invalid")
	}
	// TODO: Implement Upload
	fmt.Println(filePath)
	return nil
}

func UploadFiles(ctx context.Context, filePaths []string, uploader Uploader, concurrencyLimit int) <-chan UploadResult {
	out := make(chan UploadResult)
	jobs := make(chan string, len(filePaths))
	var wg sync.WaitGroup

	for i := 0; i < concurrencyLimit; i++ {
		wg.Add(1)
		go func(j chan string) {
			defer wg.Done()
			for file := range j {
				if ctx.Err() != nil {
					return
				}
				res := UploadResult{FilePath: file, Err: uploader.Upload(ctx, file)}
				select {
				case out <- res:
				case <-ctx.Done():
					return
				}
			}
		}(jobs)
	}
	for _, filePath := range filePaths {
		jobs <- filePath
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func isFilePathValid(filePath string) bool {
	return true // TODO: Implement filePath validation
}
