package s3

import (
	"context"
	"testing"
	"time"
	
	"github.com/stretchr/testify/require"
)

// TestUploadTextFile is a manual test to verify S3 upload works end-to-end
func TestUploadTextFile(t *testing.T) {
	client := getTestS3Client(t)
	ctx := context.Background()

	// Create test content
	testContent := []byte(`This is a test file uploaded from the storage service.
Created at: ` + time.Now().Format(time.RFC3339) + `

This file demonstrates successful integration with DigitalOcean Spaces.
- Bucket: merchbd
- Root folder: inscore
- Schema: storage
`)

	// Generate key with inscore prefix
	key := client.GenerateKey("test-tenant", "FILE_TYPE_DOCUMENT", "test-upload.txt")
	
	t.Logf("Uploading test file to key: %s", key)
	
	// Upload file
	url, cdnURL, err := client.UploadFile(ctx, key, testContent, "text/plain", true)
	require.NoError(t, err, "Upload should succeed")
	
	t.Logf("✅ Upload successful!")
	t.Logf("URL: %s", url)
	if cdnURL != "" {
		t.Logf("CDN URL: %s", cdnURL)
	}
	
	// Verify the file exists (optional - uncomment to keep file)
	// t.Logf("File uploaded successfully. You can access it at: %s", url)
	// t.Skip("Skipping cleanup to keep the file for inspection")
	
	// Clean up: Delete the uploaded file
	defer func() {
		t.Log("Cleaning up uploaded file...")
		err := client.DeleteFile(ctx, key)
		if err != nil {
			t.Logf("Warning: Failed to delete file: %v", err)
		} else {
			t.Log("✅ Cleanup successful")
		}
	}()
}
