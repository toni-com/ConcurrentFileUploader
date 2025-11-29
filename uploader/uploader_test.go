package uploader

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestUploaderContract verifies that SimpleUploader satisfies the Uploader interface.
func TestUploaderContract(t *testing.T) {
	var _ Uploader = (*SimpleUploader)(nil)
}

// TestSimpleUploader_Upload verifies the basic behavior.
func TestSimpleUploader_Upload(t *testing.T) {
	u := &SimpleUploader{}
	ctx := context.Background()

	err := u.Upload(ctx, "dummy_file.txt")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// MockUploader is a helper for testing that records which files it processed
type MockUploader struct {
	ShouldFail bool
}

func (m *MockUploader) Upload(ctx context.Context, filePath string) error {
	// Simulate a bit of work
	if m.ShouldFail {
		return fmt.Errorf("simulated error for %s", filePath)
	}
	return nil
}

func TestUploadFiles(t *testing.T) {
	files := []string{"file1.txt", "file2.txt", "file3.txt"}
	mock := &MockUploader{ShouldFail: false}

	results := UploadFiles(context.Background(), files, mock, 5)

	count := 0
	for res := range results {
		if res.Err != nil {
			t.Errorf("expected success for %s, got error: %v", res.FilePath, res.Err)
		}
		count++
	}

	if count != len(files) {
		t.Errorf("expected %d results, got %d", len(files), count)
	}
}

type ConcurrencyMock struct {
	mu          sync.Mutex
	activeCount int
	maxObserved int
}

func (m *ConcurrencyMock) Upload(ctx context.Context, filePath string) error {
	m.mu.Lock()
	m.activeCount++
	// Update the high-water mark of concurrent uploads
	if m.activeCount > m.maxObserved {
		m.maxObserved = m.activeCount
	}
	m.mu.Unlock()

	// Simulate work to force overlap
	time.Sleep(10 * time.Millisecond)

	m.mu.Lock()
	m.activeCount--
	m.mu.Unlock()

	return nil
}

func TestUploadFiles_ConcurrencyLimit(t *testing.T) {
	// want to upload 20 files, only allow 5 at a time
	files := make([]string, 20)
	for i := 0; i < 20; i++ {
		files[i] = fmt.Sprintf("file_%d.txt", i)
	}

	limit := 5
	mock := &ConcurrencyMock{}

	// Update the call to include the limit argument
	results := UploadFiles(context.Background(), files, mock, limit)

	// Drain the channel (must read to let work continue)
	for range results {
	}

	mock.mu.Lock()
	defer mock.mu.Unlock()

	if mock.maxObserved > limit {
		t.Errorf("Concurrency limit exceeded! Expected max %d, but saw %d running at once.", limit, mock.maxObserved)
	}
	if mock.maxObserved <= 1 {
		t.Log("Warning: Tasks ran sequentially. Did you actually spawn workers?")
	}
}
