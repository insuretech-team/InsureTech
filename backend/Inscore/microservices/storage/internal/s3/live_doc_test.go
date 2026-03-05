package s3

import (
	"context"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/newage-saint/insuretech/ops/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

var loadLiveEnvOnce sync.Once

func TestLiveDocumentUploadAndFolderStructure(t *testing.T) {
	loadRootEnvForLiveS3(t)
	aliasStorageEnvFromRoot()
	keepFiles := envTrue("LIVE_DOC_TEST_KEEP_FILES")
	publicUpload := true

	cfg := Config{
		AccessKeyID:     os.Getenv("SPACES_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("SPACES_SECRET_ACCESS_KEY"),
		Bucket:          os.Getenv("SPACES_BUCKET_NAME"),
		Region:          os.Getenv("SPACES_REGION"),
		Endpoint:        os.Getenv("SPACES_ENDPOINT"),
		CDNEndpoint:     os.Getenv("SPACES_CDN_ENDPOINT"),
		RootFolder:      os.Getenv("SPACES_ROOT_FOLDER"),
	}
	if strings.TrimSpace(cfg.AccessKeyID) == "" || strings.TrimSpace(cfg.SecretAccessKey) == "" || strings.TrimSpace(cfg.Bucket) == "" {
		t.Skip("live S3 test skipped: required credentials/bucket not available in root .env")
	}
	if strings.TrimSpace(cfg.Region) == "" {
		cfg.Region = "sgp1"
	}
	if strings.TrimSpace(cfg.Endpoint) == "" {
		cfg.Endpoint = "https://" + cfg.Region + ".digitaloceanspaces.com"
	}

	client, err := NewClient(cfg)
	require.NoError(t, err)

	applyStorageLayout(t, client)

	type tc struct {
		name          string
		referenceType string
		referenceID   string
		filename      string
		expectInKey   string
	}

	cases := []tc{
		{
			name:          "kyc profile",
			referenceType: "USER_KYC_PROFILE",
			referenceID:   "live-kyc-001",
			filename:      "profile.jpg",
			expectInKey:   "/kyc/profiles/",
		},
		{
			name:          "identity doc",
			referenceType: "USER_IDENTITY_DOC",
			referenceID:   "live-doc-001",
			filename:      "nid.pdf",
			expectInKey:   "/kyc/identity-docs/",
		},
		{
			name:          "policy purchase doc",
			referenceType: "POLICY_PURCHASE_DOC",
			referenceID:   "live-policy-001",
			filename:      "policy.pdf",
			expectInKey:   "/policies/",
		},
		{
			name:          "claim doc",
			referenceType: "CLAIM_DOC",
			referenceID:   "live-claim-001",
			filename:      "claim.jpg",
			expectInKey:   "/claims/",
		},
	}

	tenantID := "live-doc-test-tenant"
	ctx := context.Background()

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			fileID := "file-" + strings.ReplaceAll(time.Now().UTC().Format("150405.000"), ".", "")
			key := client.GenerateInsuranceKey(tenantID, fileID, c.referenceType, c.referenceID, c.filename)

			assert.Contains(t, key, "/"+tenantID+"/")
			assert.Contains(t, key, c.expectInKey)
			assert.Contains(t, key, "/"+fileID+".")

			content := []byte("live s3 document path test: " + c.name + " at " + time.Now().UTC().Format(time.RFC3339))
			urlStr, cdnURL, err := client.UploadFile(ctx, key, content, "text/plain", publicUpload)
			require.NoError(t, err)
			assert.NotEmpty(t, urlStr)
			t.Logf("uploaded key=%s", key)
			t.Logf("origin_url=%s", urlStr)
			if strings.TrimSpace(cfg.CDNEndpoint) != "" {
				assert.NotEmpty(t, cdnURL)
				t.Logf("cdn_url=%s", cdnURL)
			}

			size, contentType, err := client.HeadObject(ctx, key)
			require.NoError(t, err)
			assert.Greater(t, size, int64(0))
			assert.NotEmpty(t, contentType)

			if keepFiles {
				t.Logf("keeping uploaded file for manual verification: key=%s", key)
				return
			}

			err = client.DeleteFile(ctx, key)
			require.NoError(t, err)
		})
	}
}

func loadRootEnvForLiveS3(t *testing.T) {
	t.Helper()
	loadLiveEnvOnce.Do(func() {
		envPath, err := config.ResolvePath(".env")
		if err != nil {
			return
		}
		_ = godotenv.Overload(envPath)
	})
}

func setIfEmpty(key, value string) {
	if strings.TrimSpace(os.Getenv(key)) == "" && strings.TrimSpace(value) != "" {
		_ = os.Setenv(key, value)
	}
}

func envTrue(key string) bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	return v == "1" || v == "true" || v == "yes" || v == "y"
}

func aliasStorageEnvFromRoot() {
	setIfEmpty("SPACES_ACCESS_KEY_ID", os.Getenv("ACCESS_KEY_ID"))
	setIfEmpty("SPACES_SECRET_ACCESS_KEY", os.Getenv("SECRET_KEY"))

	mainCDN := strings.TrimSpace(os.Getenv("MAIN_CDN"))
	insuretechCDN := strings.TrimSpace(os.Getenv("INSURETECH_CDN_URL"))

	if mainCDN == "" && insuretechCDN != "" {
		mainCDN = extractOrigin(insuretechCDN)
	}
	setIfEmpty("SPACES_CDN_ENDPOINT", mainCDN)

	if strings.TrimSpace(os.Getenv("SPACES_BUCKET_NAME")) == "" {
		if host := extractHost(mainCDN); host != "" {
			parts := strings.Split(host, ".")
			if len(parts) > 0 {
				setIfEmpty("SPACES_BUCKET_NAME", parts[0])
			}
		}
	}

	if strings.TrimSpace(os.Getenv("SPACES_REGION")) == "" {
		if host := extractHost(mainCDN); host != "" {
			parts := strings.Split(host, ".")
			if len(parts) > 1 {
				setIfEmpty("SPACES_REGION", parts[1])
			}
		}
	}

	if strings.TrimSpace(os.Getenv("SPACES_ENDPOINT")) == "" {
		region := strings.TrimSpace(os.Getenv("SPACES_REGION"))
		if region != "" {
			setIfEmpty("SPACES_ENDPOINT", "https://"+region+".digitaloceanspaces.com")
		}
	}

	if strings.TrimSpace(os.Getenv("SPACES_ROOT_FOLDER")) == "" {
		if segment := firstPathSegment(insuretechCDN); segment != "" {
			setIfEmpty("SPACES_ROOT_FOLDER", segment)
		}
	}
}

func extractOrigin(raw string) string {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Scheme == "" || u.Host == "" {
		return ""
	}
	return u.Scheme + "://" + u.Host
}

func extractHost(raw string) string {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return ""
	}
	return u.Host
}

func firstPathSegment(raw string) string {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return ""
	}
	trimmed := strings.Trim(u.Path, "/")
	if trimmed == "" {
		return ""
	}
	return strings.Split(trimmed, "/")[0]
}

func applyStorageLayout(t *testing.T, client *Client) {
	t.Helper()

	layoutPath, err := config.ResolveConfigPath("storage_layout.yaml")
	require.NoError(t, err)

	f, err := os.Open(layoutPath)
	require.NoError(t, err)
	defer f.Close()

	var layoutCfg struct {
		StorageLayout struct {
			BasePrefix string `yaml:"base_prefix"`
			Categories map[string]struct {
				ReferenceType    string `yaml:"reference_type"`
				FolderTemplate   string `yaml:"folder_template"`
				FilenameTemplate string `yaml:"filename_template"`
			} `yaml:"categories"`
		} `yaml:"storage_layout"`
	}
	require.NoError(t, yaml.NewDecoder(f).Decode(&layoutCfg))

	templates := make(map[string]LayoutTemplate, len(layoutCfg.StorageLayout.Categories))
	for _, cat := range layoutCfg.StorageLayout.Categories {
		if strings.TrimSpace(cat.ReferenceType) == "" {
			continue
		}
		templates[cat.ReferenceType] = LayoutTemplate{
			FolderTemplate:   cat.FolderTemplate,
			FilenameTemplate: cat.FilenameTemplate,
		}
	}
	client.SetLayout(strings.TrimSpace(layoutCfg.StorageLayout.BasePrefix), templates)
}
