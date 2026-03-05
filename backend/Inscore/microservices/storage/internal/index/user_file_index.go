package index

import (
	"sort"
	"strings"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	storageentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/entity/v1"
)

// UserFileIndex is an in-memory secondary index for fast tenant+user document listing.
// It is warmed lazily from DB and updated dynamically by write paths.
type UserFileIndex struct {
	mu sync.RWMutex
	// tenant_id -> file_id -> file
	files map[string]map[string]*storageentityv1.StoredFile
	// tenant_id -> uploaded_by -> ordered files (created_at desc)
	userFiles map[string]map[string][]*storageentityv1.StoredFile
	// tenant_id -> uploaded_by -> fully warmed snapshot flag
	warmed map[string]map[string]bool
}

func NewUserFileIndex() *UserFileIndex {
	return &UserFileIndex{
		files:     map[string]map[string]*storageentityv1.StoredFile{},
		userFiles: map[string]map[string][]*storageentityv1.StoredFile{},
		warmed:    map[string]map[string]bool{},
	}
}

func (i *UserFileIndex) Upsert(file *storageentityv1.StoredFile) {
	if file == nil {
		return
	}
	tenantID := strings.TrimSpace(file.GetTenantId())
	fileID := strings.TrimSpace(file.GetFileId())
	if tenantID == "" || fileID == "" {
		return
	}
	userID := normalizeUserKey(file.GetUploadedBy(), tenantID)

	i.mu.Lock()
	defer i.mu.Unlock()

	i.ensureTenant(tenantID)

	prev, hadPrev := i.files[tenantID][fileID]
	i.files[tenantID][fileID] = cloneFile(file)

	// Keep warm indexes dynamically updated.
	if hadPrev {
		prevUser := normalizeUserKey(prev.GetUploadedBy(), tenantID)
		if i.warmed[tenantID][prevUser] {
			i.userFiles[tenantID][prevUser] = removeByFileID(i.userFiles[tenantID][prevUser], fileID)
		}
	}
	if i.warmed[tenantID][userID] {
		i.userFiles[tenantID][userID] = upsertOrdered(i.userFiles[tenantID][userID], file)
	}
}

func (i *UserFileIndex) Delete(tenantID string, fileID string) {
	tenantID = strings.TrimSpace(tenantID)
	fileID = strings.TrimSpace(fileID)
	if tenantID == "" || fileID == "" {
		return
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	if _, ok := i.files[tenantID]; !ok {
		return
	}
	existing := i.files[tenantID][fileID]
	delete(i.files[tenantID], fileID)
	if existing == nil {
		return
	}

	userID := normalizeUserKey(existing.GetUploadedBy(), tenantID)
	if i.warmed[tenantID][userID] {
		i.userFiles[tenantID][userID] = removeByFileID(i.userFiles[tenantID][userID], fileID)
	}
}

// WarmUser replaces a user index snapshot with DB-truth and marks it as warmed.
func (i *UserFileIndex) WarmUser(tenantID string, userID string, files []*storageentityv1.StoredFile) {
	tenantID = strings.TrimSpace(tenantID)
	userID = normalizeUserKey(userID, tenantID)
	if tenantID == "" || userID == "" {
		return
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	i.ensureTenant(tenantID)

	old := i.userFiles[tenantID][userID]
	newList := make([]*storageentityv1.StoredFile, 0, len(files))
	newIDs := make(map[string]struct{}, len(files))
	for _, f := range files {
		if f == nil {
			continue
		}
		if strings.TrimSpace(f.GetTenantId()) != tenantID {
			continue
		}
		id := strings.TrimSpace(f.GetFileId())
		if id == "" {
			continue
		}
		cloned := cloneFile(f)
		newList = append(newList, cloned)
		i.files[tenantID][id] = cloned
		newIDs[id] = struct{}{}
	}
	sortByCreatedDesc(newList)

	// Remove stale file map entries that previously belonged to this user snapshot.
	for _, oldFile := range old {
		if oldFile == nil {
			continue
		}
		oldID := strings.TrimSpace(oldFile.GetFileId())
		if oldID == "" {
			continue
		}
		if _, ok := newIDs[oldID]; ok {
			continue
		}
		current := i.files[tenantID][oldID]
		if current != nil && normalizeUserKey(current.GetUploadedBy(), tenantID) == userID {
			delete(i.files[tenantID], oldID)
		}
	}

	i.userFiles[tenantID][userID] = newList
	i.warmed[tenantID][userID] = true
}

func (i *UserFileIndex) List(
	tenantID string,
	userID string,
	fileType storageentityv1.FileType,
	referenceID string,
	referenceType string,
	limit int32,
	offset int32,
) ([]*storageentityv1.StoredFile, int, bool) {
	tenantID = strings.TrimSpace(tenantID)
	userID = normalizeUserKey(userID, tenantID)
	referenceID = strings.TrimSpace(referenceID)
	referenceType = strings.TrimSpace(referenceType)
	if tenantID == "" || userID == "" {
		return nil, 0, false
	}

	i.mu.RLock()
	defer i.mu.RUnlock()

	if !i.warmed[tenantID][userID] {
		return nil, 0, false
	}

	src := i.userFiles[tenantID][userID]
	filtered := make([]*storageentityv1.StoredFile, 0, len(src))
	for _, file := range src {
		if file == nil {
			continue
		}
		if fileType != storageentityv1.FileType_FILE_TYPE_UNSPECIFIED && file.GetFileType() != fileType {
			continue
		}
		if referenceID != "" && file.GetReferenceId() != referenceID {
			continue
		}
		if referenceType != "" && file.GetReferenceType() != referenceType {
			continue
		}
		filtered = append(filtered, cloneFile(file))
	}

	total := len(filtered)
	start := int(offset)
	if start < 0 {
		start = 0
	}
	if start > total {
		start = total
	}
	end := total
	if limit > 0 {
		end = start + int(limit)
		if end > total {
			end = total
		}
	}
	return filtered[start:end], total, true
}

func (i *UserFileIndex) ensureTenant(tenantID string) {
	if _, ok := i.files[tenantID]; !ok {
		i.files[tenantID] = map[string]*storageentityv1.StoredFile{}
	}
	if _, ok := i.userFiles[tenantID]; !ok {
		i.userFiles[tenantID] = map[string][]*storageentityv1.StoredFile{}
	}
	if _, ok := i.warmed[tenantID]; !ok {
		i.warmed[tenantID] = map[string]bool{}
	}
}

func normalizeUserKey(userID string, tenantID string) string {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return strings.TrimSpace(tenantID)
	}
	return userID
}

func removeByFileID(files []*storageentityv1.StoredFile, fileID string) []*storageentityv1.StoredFile {
	if len(files) == 0 {
		return files
	}
	out := make([]*storageentityv1.StoredFile, 0, len(files))
	for _, f := range files {
		if f == nil {
			continue
		}
		if strings.TrimSpace(f.GetFileId()) == fileID {
			continue
		}
		out = append(out, f)
	}
	return out
}

func upsertOrdered(files []*storageentityv1.StoredFile, file *storageentityv1.StoredFile) []*storageentityv1.StoredFile {
	if file == nil {
		return files
	}
	fileID := strings.TrimSpace(file.GetFileId())
	if fileID == "" {
		return files
	}
	files = removeByFileID(files, fileID)
	files = append(files, cloneFile(file))
	sortByCreatedDesc(files)
	return files
}

func sortByCreatedDesc(files []*storageentityv1.StoredFile) {
	sort.SliceStable(files, func(a, b int) bool {
		return createdAt(files[a]).After(createdAt(files[b]))
	})
}

func createdAt(file *storageentityv1.StoredFile) time.Time {
	if file == nil || file.GetCreatedAt() == nil {
		return time.Time{}
	}
	return file.GetCreatedAt().AsTime()
}

func cloneFile(file *storageentityv1.StoredFile) *storageentityv1.StoredFile {
	if file == nil {
		return nil
	}
	return proto.Clone(file).(*storageentityv1.StoredFile)
}
