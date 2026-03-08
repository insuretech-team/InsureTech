package s3

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: These tests require actual S3/Spaces credentials to be set in environment variables
// SPACES_ACCESS_KEY_ID, SPACES_SECRET_ACCESS_KEY, SPACES_BUCKET_NAME, SPACES_REGION, SPACES_ENDPOINT

func getTestS3Client(t *testing.T) *Client {
	// Check if required environment variables are set
	accessKeyID := os.Getenv("SPACES_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("SPACES_SECRET_ACCESS_KEY")
	bucket := os.Getenv("SPACES_BUCKET_NAME")
	region := os.Getenv("SPACES_REGION")
	endpoint := os.Getenv("SPACES_ENDPOINT")
	cdnEndpoint := os.Getenv("SPACES_CDN_ENDPOINT")
	rootFolder := os.Getenv("SPACES_ROOT_FOLDER")

	if accessKeyID == "" || secretAccessKey == "" || bucket == "" {
		t.Skip("Skipping S3 tests: S3 credentials not configured in environment")
	}

	// Default values for optional fields
	if region == "" {
		region = "sgp1"
	}
	if endpoint == "" {
		endpoint = "https://sgp1.digitaloceanspaces.com"
	}
	if rootFolder == "" {
		rootFolder = "inscore"
	}

	cfg := Config{
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		Bucket:          bucket,
		Region:          region,
		Endpoint:        endpoint,
		CDNEndpoint:     cdnEndpoint,
		RootFolder:      rootFolder,
	}

	client, err := NewClient(cfg)
	require.NoError(t, err, "Failed to create S3 client")

	return client
}

func TestS3Client_GenerateKey(t *testing.T) {
	client := getTestS3Client(t)

	tests := []struct {
		name          string
		tenantID      string
		referenceType string
		referenceID   string
		fileID        string
		filename      string
		expectPrefix  string
	}{
		{
			name:          "KYC profile",
			tenantID:      "test-tenant",
			referenceType: "USER_KYC_PROFILE",
			referenceID:   "kyc-123",
			fileID:        "file-123",
			filename:      "profile.jpg",
			expectPrefix:  "inscore/test-tenant/kyc/profiles/kyc-123/file-123/",
		},
		{
			name:          "Identity document",
			tenantID:      "test-tenant",
			referenceType: "USER_IDENTITY_DOC",
			referenceID:   "doc-456",
			fileID:        "file-456",
			filename:      "nid.pdf",
			expectPrefix:  "inscore/test-tenant/kyc/identity-docs/doc-456/file-456/",
		},
		{
			name:          "Policy purchase document",
			tenantID:      "test-tenant",
			referenceType: "POLICY_PURCHASE_DOC",
			referenceID:   "policy-789",
			fileID:        "file-789",
			filename:      "policy.pdf",
			expectPrefix:  "inscore/test-tenant/policies/purchase-docs/policy-789/file-789/",
		},
		{
			name:          "Claim document",
			tenantID:      "test-tenant",
			referenceType: "CLAIM_DOC",
			referenceID:   "claim-111",
			fileID:        "file-111",
			filename:      "claim.jpg",
			expectPrefix:  "inscore/test-tenant/claims/documents/claim-111/file-111/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := client.GenerateInsuranceKey(tt.tenantID, tt.fileID, tt.referenceType, tt.referenceID, tt.filename)
			assert.Contains(t, key, tt.expectPrefix, "Generated key should contain expected prefix")
			assert.Contains(t, key, "inscore/", "Generated key should contain root folder")
		})
	}
}

func TestS3Client_GenerateKeyWithReference(t *testing.T) {
	client := getTestS3Client(t)

	key := client.GenerateInsuranceKey("test-tenant", "file-999", "POLICY_PURCHASE_DOC", "policy-123", "invoice.pdf")

	assert.Contains(t, key, "inscore/", "Key should contain root folder")
	assert.Contains(t, key, "policies/purchase-docs/", "Key should contain policy purchase folder")
	assert.Contains(t, key, "policy-123", "Key should contain reference ID")
	assert.Contains(t, key, "file-999", "Key should contain file ID")
}

func TestS3Client_UploadAndDelete(t *testing.T) {
	client := getTestS3Client(t)
	ctx := context.Background()

	// Generate test key
	key := client.GenerateKey("test-tenant", "FILE_TYPE_IMAGE", "test-upload.txt")
	content := []byte("This is a test file for S3 upload")
	contentType := "text/plain"

	// Test Upload
	url, cdnURL, err := client.UploadFile(ctx, key, content, contentType, true)
	require.NoError(t, err, "Upload should succeed")
	assert.NotEmpty(t, url, "URL should not be empty")
	assert.Contains(t, url, key, "URL should contain the key")

	if client.cdnEndpoint != "" {
		assert.NotEmpty(t, cdnURL, "CDN URL should not be empty if CDN endpoint is configured")
		assert.Contains(t, cdnURL, key, "CDN URL should contain the key")
	}

	// Clean up: Delete the uploaded file
	defer func() {
		err := client.DeleteFile(ctx, key)
		assert.NoError(t, err, "Delete should succeed")
	}()
}

func TestS3Client_PresignedURLs(t *testing.T) {
	client := getTestS3Client(t)
	ctx := context.Background()

	// Generate test key
	key := client.GenerateKey("test-tenant", "FILE_TYPE_IMAGE", "test-presigned.jpg")

	// Test presigned upload URL
	uploadURL, err := client.GetPresignedUploadURL(ctx, key, 15*time.Minute)
	require.NoError(t, err, "Presigned upload URL generation should succeed")
	assert.NotEmpty(t, uploadURL, "Presigned upload URL should not be empty")
	assert.Contains(t, uploadURL, client.bucket, "Presigned URL should contain bucket name")

	// Upload a test file first for download URL test
	content := []byte("Test file for presigned download")
	_, _, err = client.UploadFile(ctx, key, content, "text/plain", false)
	require.NoError(t, err, "Upload should succeed")

	defer func() {
		_ = client.DeleteFile(ctx, key)
	}()

	// Test presigned download URL
	downloadURL, err := client.GetPresignedDownloadURL(ctx, key, time.Hour)
	require.NoError(t, err, "Presigned download URL generation should succeed")
	assert.NotEmpty(t, downloadURL, "Presigned download URL should not be empty")
	assert.Contains(t, downloadURL, client.bucket, "Presigned download URL should contain bucket name")
}

func TestS3Client_GetBucket(t *testing.T) {
	client := getTestS3Client(t)

	bucket := client.GetBucket()
	assert.NotEmpty(t, bucket, "Bucket name should not be empty")
	assert.Equal(t, os.Getenv("SPACES_BUCKET_NAME"), bucket, "Bucket should match environment variable")
}
