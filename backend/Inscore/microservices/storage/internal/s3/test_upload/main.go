package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/storage/internal/s3"
	"github.com/newage-saint/insuretech/ops/config"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	envPath, err := config.ResolvePath(".env")
	if err != nil {
		fmt.Printf("Warning: Could not resolve .env path: %v\n", err)
	} else {
		if err := godotenv.Load(envPath); err != nil {
			fmt.Printf("Warning: Could not load .env file: %v\n", err)
		} else {
			fmt.Println("✅ Loaded .env file")
		}
	}
	// Load S3 config from environment
	cfg := s3.Config{
		AccessKeyID:     os.Getenv("SPACES_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("SPACES_SECRET_ACCESS_KEY"),
		Bucket:          os.Getenv("SPACES_BUCKET_NAME"),
		Region:          os.Getenv("SPACES_REGION"),
		Endpoint:        os.Getenv("SPACES_ENDPOINT"),
		CDNEndpoint:     os.Getenv("SPACES_CDN_ENDPOINT"),
		RootFolder:      os.Getenv("SPACES_ROOT_FOLDER"),
	}

	if cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" || cfg.Bucket == "" {
		fmt.Println("❌ Missing S3 credentials in environment")
		fmt.Println("Make sure SPACES_ACCESS_KEY_ID, SPACES_SECRET_ACCESS_KEY, and SPACES_BUCKET_NAME are set")
		os.Exit(1)
	}

	fmt.Printf("Creating S3 client for bucket: %s\n", cfg.Bucket)
	client, err := s3.NewClient(cfg)
	if err != nil {
		fmt.Printf("❌ Failed to create S3 client: %v\n", err)
		os.Exit(1)
	}

	// Create test content
	content := []byte(fmt.Sprintf("Test file uploaded at %s\n", time.Now().Format(time.RFC3339)))
	
	// Generate key
	key := client.GenerateKey("test-tenant", "FILE_TYPE_DOCUMENT", "test-upload.txt")
	fmt.Printf("Uploading to key: %s\n", key)

	// Upload
	ctx := context.Background()
	url, cdnURL, err := client.UploadFile(ctx, key, content, "text/plain", true)
	if err != nil {
		fmt.Printf("❌ Upload failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Upload successful!\n")
	fmt.Printf("URL: %s\n", url)
	if cdnURL != "" {
		fmt.Printf("CDN URL: %s\n", cdnURL)
	}

	fmt.Println("\nPress Enter to delete the file or Ctrl+C to keep it...")
	fmt.Scanln()

	// Delete
	err = client.DeleteFile(ctx, key)
	if err != nil {
		fmt.Printf("⚠️  Delete failed: %v\n", err)
	} else {
		fmt.Println("✅ File deleted")
	}
}
