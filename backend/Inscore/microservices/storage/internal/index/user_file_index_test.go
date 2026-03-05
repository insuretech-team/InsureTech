package index

import (
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	storageentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/entity/v1"
)

func TestUserFileIndexWarmAndList(t *testing.T) {
	t.Parallel()

	idx := NewUserFileIndex()
	now := time.Now().UTC()
	idx.WarmUser("tenant-1", "user-1", []*storageentityv1.StoredFile{
		{
			FileId:        "f1",
			TenantId:      "tenant-1",
			UploadedBy:    "user-1",
			ReferenceType: "USER_IDENTITY_DOC",
			FileType:      storageentityv1.FileType_FILE_TYPE_DOCUMENT,
			CreatedAt:     timestamppb.New(now.Add(-1 * time.Minute)),
		},
		{
			FileId:        "f2",
			TenantId:      "tenant-1",
			UploadedBy:    "user-1",
			ReferenceType: "USER_KYC_PROFILE",
			FileType:      storageentityv1.FileType_FILE_TYPE_IMAGE,
			CreatedAt:     timestamppb.New(now),
		},
	})

	files, total, ok := idx.List("tenant-1", "user-1", storageentityv1.FileType_FILE_TYPE_UNSPECIFIED, "", "", 50, 0)
	if !ok {
		t.Fatal("expected warmed index hit")
	}
	if total != 2 {
		t.Fatalf("expected total=2 got=%d", total)
	}
	if len(files) != 2 {
		t.Fatalf("expected len=2 got=%d", len(files))
	}
	if files[0].FileId != "f2" {
		t.Fatalf("expected newest file first, got %s", files[0].FileId)
	}
}

func TestUserFileIndexDynamicUpdateAndDelete(t *testing.T) {
	t.Parallel()

	idx := NewUserFileIndex()
	idx.WarmUser("tenant-1", "user-1", nil)

	idx.Upsert(&storageentityv1.StoredFile{
		FileId:        "f1",
		TenantId:      "tenant-1",
		UploadedBy:    "user-1",
		ReferenceType: "USER_IDENTITY_DOC",
		FileType:      storageentityv1.FileType_FILE_TYPE_DOCUMENT,
		CreatedAt:     timestamppb.Now(),
	})

	files, total, ok := idx.List("tenant-1", "user-1", storageentityv1.FileType_FILE_TYPE_UNSPECIFIED, "", "", 10, 0)
	if !ok || total != 1 || len(files) != 1 {
		t.Fatalf("unexpected list after upsert: ok=%v total=%d len=%d", ok, total, len(files))
	}

	idx.Delete("tenant-1", "f1")
	files, total, ok = idx.List("tenant-1", "user-1", storageentityv1.FileType_FILE_TYPE_UNSPECIFIED, "", "", 10, 0)
	if !ok || total != 0 || len(files) != 0 {
		t.Fatalf("unexpected list after delete: ok=%v total=%d len=%d", ok, total, len(files))
	}
}
