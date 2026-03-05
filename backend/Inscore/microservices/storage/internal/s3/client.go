package s3

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

// Client handles S3/Spaces operations
type Client struct {
	s3Client    *s3.Client
	bucket      string
	region      string
	endpoint    string
	cdnEndpoint string
	rootFolder  string
	layout      map[string]LayoutTemplate
	layoutBase  string
}

// LayoutTemplate configures folder and filename templates for a reference type.
type LayoutTemplate struct {
	FolderTemplate   string
	FilenameTemplate string
}

// Config for S3 client
type Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	Region          string
	Endpoint        string
	CDNEndpoint     string
	RootFolder      string
}

// NewClient creates a new S3 client
func NewClient(cfg Config) (*Client, error) {
	// Create custom resolver for DigitalOcean Spaces
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           cfg.Endpoint,
			SigningRegion: cfg.Region,
		}, nil
	})

	// Load AWS config with custom credentials
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
		config.WithRegion(cfg.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client with path-style addressing
	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &Client{
		s3Client:    s3Client,
		bucket:      cfg.Bucket,
		region:      cfg.Region,
		endpoint:    cfg.Endpoint,
		cdnEndpoint: cfg.CDNEndpoint,
		rootFolder:  cfg.RootFolder,
		layout:      map[string]LayoutTemplate{},
	}, nil
}

// UploadFile uploads a file to S3/Spaces
func (c *Client) UploadFile(ctx context.Context, key string, content []byte, contentType string, isPublic bool) (string, string, error) {
	// Set ACL based on public flag
	var acl types.ObjectCannedACL
	if isPublic {
		acl = types.ObjectCannedACLPublicRead
	} else {
		acl = types.ObjectCannedACLPrivate
	}

	// Upload to S3
	_, err := c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(content),
		ContentType: aws.String(contentType),
		ACL:         acl,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to upload file: %w", err)
	}

	url, cdnURL := c.BuildObjectURLs(key)

	return url, cdnURL, nil
}

// DeleteFile deletes a file from S3/Spaces
func (c *Client) DeleteFile(ctx context.Context, key string) error {
	_, err := c.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// HeadObject returns object size and content type for a key.
func (c *Client) HeadObject(ctx context.Context, key string) (int64, string, error) {
	out, err := c.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return 0, "", fmt.Errorf("failed to head object: %w", err)
	}

	var size int64
	if out.ContentLength != nil {
		size = *out.ContentLength
	}

	contentType := ""
	if out.ContentType != nil {
		contentType = *out.ContentType
	}

	return size, contentType, nil
}

// GetPresignedUploadURL generates a presigned URL for uploading
func (c *Client) GetPresignedUploadURL(ctx context.Context, key string, expiresIn time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(c.s3Client)

	req, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiresIn
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned upload URL: %w", err)
	}

	return req.URL, nil
}

// GetPresignedDownloadURL generates a presigned URL for downloading
func (c *Client) GetPresignedDownloadURL(ctx context.Context, key string, expiresIn time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(c.s3Client)

	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiresIn
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned download URL: %w", err)
	}

	return req.URL, nil
}

// GenerateKey generates a unique S3 key for a file based on file type
// Uses existing DigitalOcean Spaces folder structure: assets, catalogs, orders, shipments
func (c *Client) GenerateKey(tenantID string, fileType string, filename string) string {
	return c.GenerateInsuranceKey(tenantID, uuid.New().String(), "", "", filename)
}

// GenerateKeyWithType generates a unique S3 key with optional subfolder
func (c *Client) GenerateKeyWithType(tenantID string, fileType string, subfolder string, filename string) string {
	return c.GenerateInsuranceKey(tenantID, uuid.New().String(), subfolder, "", filename)
}

// GenerateKeyWithReference generates a key with reference ID (for better organization)
func (c *Client) GenerateKeyWithReference(tenantID string, fileType string, referenceID string, filename string) string {
	return c.GenerateInsuranceKey(tenantID, uuid.New().String(), fileType, referenceID, filename)
}

// GenerateInsuranceKey generates insurance-domain S3 keys.
// Expected reference_type values:
// USER_KYC_PROFILE, USER_IDENTITY_DOC, POLICY_PURCHASE_DOC, CLAIM_DOC,
// POLICY_ATTACHMENT, CLAIM_ATTACHMENT
// For layout templates, reference_id may optionally include key-value metadata
// (e.g. "user_id=u1;document_id=d1;document_type=passport") to populate
// typed placeholders without overloading all IDs with the same value.
func (c *Client) GenerateInsuranceKey(tenantID string, fileID string, referenceType string, referenceID string, filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".bin"
	}
	if fileID == "" {
		fileID = uuid.New().String()
	}
	if referenceID == "" {
		referenceID = "unassigned"
	}

	now := time.Now().UTC()
	refType := strings.ToUpper(strings.TrimSpace(referenceType))

	if tmpl, ok := c.layout[refType]; ok {
		basePrefix := c.layoutBase
		if basePrefix == "" {
			basePrefix = strings.Trim(c.rootFolder, "/")
		}
		placeholder := c.layoutPlaceholders(basePrefix, tenantID, fileID, refType, referenceID, now, ext)

		folder := strings.Trim(c.renderTemplate(tmpl.FolderTemplate, placeholder), "/")
		filenameTmpl := tmpl.FilenameTemplate
		if strings.TrimSpace(filenameTmpl) == "" {
			filenameTmpl = "{file_id}.{ext}"
		}
		leaf := c.renderTemplate(filenameTmpl, placeholder)
		return strings.Trim(strings.Trim(folder, "/")+"/"+strings.Trim(leaf, "/"), "/")
	}

	domainPath := "misc/files"
	switch refType {
	case "USER_KYC_PROFILE":
		domainPath = "kyc/profiles"
	case "USER_IDENTITY_DOC":
		domainPath = "kyc/identity-docs"
	case "POLICY_PURCHASE_DOC":
		domainPath = "policies/purchase-docs"
	case "CLAIM_DOC":
		domainPath = "claims/documents"
	case "POLICY_ATTACHMENT":
		domainPath = "policies/attachments"
	case "CLAIM_ATTACHMENT":
		domainPath = "claims/attachments"
	}

	key := fmt.Sprintf("%s/%s/%s/%s/%04d/%02d/%02d/%s%s",
		tenantID,
		domainPath,
		referenceID,
		fileID,
		now.Year(),
		now.Month(),
		now.Day(),
		fileID,
		ext,
	)

	if c.rootFolder != "" {
		return fmt.Sprintf("%s/%s", strings.Trim(c.rootFolder, "/"), key)
	}
	return key
}

// SetLayout configures reference_type specific key templates loaded from storage_layout.yaml.
func (c *Client) SetLayout(basePrefix string, templates map[string]LayoutTemplate) {
	c.layoutBase = strings.Trim(basePrefix, "/")
	c.layout = map[string]LayoutTemplate{}
	for k, v := range templates {
		c.layout[strings.ToUpper(strings.TrimSpace(k))] = v
	}
}

func (c *Client) renderTemplate(template string, values map[string]string) string {
	out := template
	for k, v := range values {
		out = strings.ReplaceAll(out, "{"+k+"}", v)
	}
	return out
}

func (c *Client) layoutPlaceholders(
	basePrefix string,
	tenantID string,
	fileID string,
	referenceType string,
	referenceID string,
	now time.Time,
	ext string,
) map[string]string {
	values := map[string]string{
		"base_prefix":         basePrefix,
		"tenant_id":           tenantID,
		"reference_id":        referenceID,
		"user_id":             "unknown-user",
		"kyc_verification_id": "unknown-kyc",
		"document_id":         "unknown-document",
		"policy_id":           "unknown-policy",
		"claim_id":            "unknown-claim",
		"document_type":       "default",
		"yyyy":                fmt.Sprintf("%04d", now.Year()),
		"mm":                  fmt.Sprintf("%02d", now.Month()),
		"dd":                  fmt.Sprintf("%02d", now.Day()),
		"file_id":             fileID,
		"ext":                 strings.TrimPrefix(ext, "."),
	}

	meta := parseReferenceMetadata(referenceID)
	for key, value := range meta {
		values[key] = value
	}

	switch referenceType {
	case "USER_KYC_PROFILE":
		if values["user_id"] == "unknown-user" {
			values["user_id"] = referenceID
		}
		if values["kyc_verification_id"] == "unknown-kyc" {
			values["kyc_verification_id"] = fileID
		}
		values["document_type"] = "kyc_profile"
	case "USER_IDENTITY_DOC":
		if values["document_id"] == "unknown-document" {
			values["document_id"] = referenceID
		}
		if values["document_type"] == "default" {
			values["document_type"] = "identity"
		}
	case "POLICY_PURCHASE_DOC", "POLICY_ATTACHMENT":
		if values["policy_id"] == "unknown-policy" {
			values["policy_id"] = referenceID
		}
		if values["document_type"] == "default" {
			values["document_type"] = "policy"
		}
	case "CLAIM_DOC", "CLAIM_ATTACHMENT":
		if values["claim_id"] == "unknown-claim" {
			values["claim_id"] = referenceID
		}
		if values["document_type"] == "default" {
			values["document_type"] = "claim"
		}
	}

	return values
}

func parseReferenceMetadata(referenceID string) map[string]string {
	meta := map[string]string{}
	raw := strings.TrimSpace(referenceID)
	if raw == "" {
		return meta
	}
	if !strings.ContainsAny(raw, "=:;,|&") {
		return meta
	}

	segments := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ';' || r == ',' || r == '|' || r == '&'
	})
	for _, segment := range segments {
		segment = strings.TrimSpace(segment)
		if segment == "" {
			continue
		}

		kv := strings.SplitN(segment, "=", 2)
		if len(kv) != 2 {
			kv = strings.SplitN(segment, ":", 2)
			if len(kv) != 2 {
				continue
			}
		}
		key := normalizeMetadataKey(kv[0])
		value := strings.TrimSpace(kv[1])
		if key == "" || value == "" {
			continue
		}
		meta[key] = value
	}
	return meta
}

func normalizeMetadataKey(key string) string {
	k := strings.ToLower(strings.TrimSpace(key))
	k = strings.ReplaceAll(k, "-", "_")
	k = strings.ReplaceAll(k, " ", "_")

	switch k {
	case "user", "userid", "user_id":
		return "user_id"
	case "kyc", "kyc_id", "kyc_verification_id":
		return "kyc_verification_id"
	case "doc", "document", "document_id":
		return "document_id"
	case "doctype", "doc_type", "document_type":
		return "document_type"
	case "policy", "policy_id":
		return "policy_id"
	case "claim", "claim_id":
		return "claim_id"
	case "reference", "reference_id", "ref":
		return "reference_id"
	default:
		return ""
	}
}

// GetBucket returns the bucket name
func (c *Client) GetBucket() string {
	return c.bucket
}

// BuildObjectURLs returns canonical object and CDN URLs for a storage key.
func (c *Client) BuildObjectURLs(key string) (string, string) {
	url := fmt.Sprintf("%s/%s", strings.TrimRight(c.endpoint, "/"), strings.TrimLeft(key, "/"))
	cdnURL := ""
	if c.cdnEndpoint != "" {
		cdnURL = fmt.Sprintf("%s/%s", strings.TrimRight(c.cdnEndpoint, "/"), strings.TrimLeft(key, "/"))
	}
	return url, cdnURL
}
