package s3

import (
	"strings"
	"testing"
)

func TestParseReferenceMetadata(t *testing.T) {
	t.Parallel()

	meta := parseReferenceMetadata("user_id=user-1;document_type=passport;document_id=doc-1")
	if meta["user_id"] != "user-1" {
		t.Fatalf("expected user_id=user-1, got %q", meta["user_id"])
	}
	if meta["document_type"] != "passport" {
		t.Fatalf("expected document_type=passport, got %q", meta["document_type"])
	}
	if meta["document_id"] != "doc-1" {
		t.Fatalf("expected document_id=doc-1, got %q", meta["document_id"])
	}
}

func TestGenerateInsuranceKey_LayoutPlaceholders(t *testing.T) {
	t.Parallel()

	client := &Client{
		rootFolder: "insuretech",
		layout: map[string]LayoutTemplate{
			"USER_IDENTITY_DOC": {
				FolderTemplate:   "{base_prefix}/{tenant_id}/kyc/identity-docs/{user_id}/{document_type}/{document_id}/{yyyy}/{mm}/{dd}/",
				FilenameTemplate: "{file_id}.{ext}",
			},
		},
		layoutBase: "insuretech",
	}

	key := client.GenerateInsuranceKey("tenant-1", "file-1", "USER_IDENTITY_DOC", "doc-ref-1", "nid.pdf")
	if !strings.Contains(key, "/kyc/identity-docs/unknown-user/identity/doc-ref-1/") {
		t.Fatalf("unexpected fallback key structure: %s", key)
	}
	if strings.Contains(key, "/doc-ref-1/doc-ref-1/doc-ref-1/") {
		t.Fatalf("placeholder collapse detected (all placeholders mapped to reference_id): %s", key)
	}

	key2 := client.GenerateInsuranceKey(
		"tenant-1",
		"file-2",
		"USER_IDENTITY_DOC",
		"user_id=user-22;document_type=passport;document_id=doc-55",
		"passport.jpg",
	)
	if !strings.Contains(key2, "/kyc/identity-docs/user-22/passport/doc-55/") {
		t.Fatalf("structured metadata not applied: %s", key2)
	}
}

func TestGenerateInsuranceKey_KYCDefaults(t *testing.T) {
	t.Parallel()

	client := &Client{
		rootFolder: "insuretech",
		layout: map[string]LayoutTemplate{
			"USER_KYC_PROFILE": {
				FolderTemplate:   "{base_prefix}/{tenant_id}/kyc/profiles/{user_id}/{kyc_verification_id}/{yyyy}/{mm}/{dd}/",
				FilenameTemplate: "{file_id}.{ext}",
			},
		},
		layoutBase: "insuretech",
	}

	key := client.GenerateInsuranceKey("tenant-abc", "file-123", "USER_KYC_PROFILE", "user-900", "selfie.png")
	if !strings.Contains(key, "/kyc/profiles/user-900/file-123/") {
		t.Fatalf("unexpected KYC placeholder defaults: %s", key)
	}
}
