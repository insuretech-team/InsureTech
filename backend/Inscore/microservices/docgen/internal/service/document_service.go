package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	htmltmpl "html/template"
	"image"
	"image/png"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/docgen/internal/kafka"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/docgen/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	documentv1 "github.com/newage-saint/insuretech/gen/go/insuretech/document/entity/v1"
	storageentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/entity/v1"
	storageservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/service/v1"
	"github.com/newage-saint/insuretech/ops/config"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrUnsupportedOutput = errors.New("unsupported output format")
	ErrTemplateNotFound  = repository.ErrTemplateNotFound
	ErrDocumentNotFound  = repository.ErrDocumentNotFound
	varPattern           = regexp.MustCompile(`\{\{\s*\.([a-zA-Z0-9_]+)\s*\}\}`)
	htmlTagPattern       = regexp.MustCompile(`<[^>]*>`)
	imgSrcAttrPattern    = regexp.MustCompile(`(?i)(<img\b[^>]*?\bsrc\s*=\s*["'])([^"']+)(["'])`)
)

// StorageClient is the subset of storage RPCs used by docgen.
type StorageClient interface {
	UploadFile(ctx context.Context, in *storageservicev1.UploadFileRequest, opts ...grpc.CallOption) (*storageservicev1.UploadFileResponse, error)
	DeleteFile(ctx context.Context, in *storageservicev1.DeleteFileRequest, opts ...grpc.CallOption) (*storageservicev1.DeleteFileResponse, error)
}

// DocumentService provides document generation and template management logic.
type DocumentService struct {
	templateRepo    *repository.DocumentTemplateRepository
	generationRepo  *repository.DocumentGenerationRepository
	storageClient   StorageClient
	kafkaPublisher  *kafka.Publisher
	templateDirPath string
	gotenbergURL    string
	pdfTimeout      time.Duration
}

func NewDocumentService(
	templateRepo *repository.DocumentTemplateRepository,
	generationRepo *repository.DocumentGenerationRepository,
	storageClient StorageClient,
) (*DocumentService, error) {
	templateDir, err := config.ResolvePath(filepath.Join("backend", "inscore", "templates"))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve template directory: %w", err)
	}

	svc := &DocumentService{
		templateRepo:    templateRepo,
		generationRepo:  generationRepo,
		storageClient:   storageClient,
		kafkaPublisher:  nil,
		templateDirPath: templateDir,
		gotenbergURL:    "",
		pdfTimeout:      8 * time.Second,
	}
	if err := svc.bootstrapDefaultTemplates(context.Background()); err != nil {
		return nil, err
	}
	return svc, nil
}

// SetKafkaPublisher injects a Kafka publisher into the service
func (s *DocumentService) SetKafkaPublisher(publisher *kafka.Publisher) {
	s.kafkaPublisher = publisher
}

// SetPDFRenderer configures the external PDF renderer endpoint and timeout.
func (s *DocumentService) SetPDFRenderer(gotenbergURL string, timeout time.Duration) {
	s.gotenbergURL = strings.TrimSpace(gotenbergURL)
	if timeout > 0 {
		s.pdfTimeout = timeout
	}
}

func (s *DocumentService) GenerateDocument(
	ctx context.Context,
	templateID, entityType, entityID string,
	data *structpb.Struct,
	includeQRCode bool,
	tenantID, generatedBy string,
) (*documentv1.DocumentGeneration, error) {
	if strings.TrimSpace(templateID) == "" {
		return nil, fmt.Errorf("%w: template_id is required", ErrInvalidInput)
	}
	if strings.TrimSpace(entityType) == "" || strings.TrimSpace(entityID) == "" {
		return nil, fmt.Errorf("%w: entity_type and entity_id are required", ErrInvalidInput)
	}

	tpl, err := s.resolveTemplateForGeneration(ctx, templateID)
	if err != nil {
		// Publish generation.failed event (non-blocking)
		failedGenerationID := uuid.New().String()
		go func() {
			if s.kafkaPublisher != nil {
				if err := s.kafkaPublisher.PublishGenerationFailed(
					context.Background(),
					failedGenerationID,
					tenantID,
					fmt.Sprintf("failed to fetch template: %v", err),
					failedGenerationID,
				); err != nil {
					logger.Warnf("failed to publish generation.failed event: %v", err)
				}
			}
		}()
		return nil, err
	}
	if !tpl.IsActive {
		return nil, fmt.Errorf("%w: template is inactive", ErrInvalidInput)
	}

	payload := map[string]any{}
	if data != nil {
		payload = data.AsMap()
	}
	payload = ensureMap(payload)
	if err := enrichTemplatePayload(tpl.Name, payload); err != nil {
		return nil, err
	}
	applyBusinessDefaults(tpl.Name, payload)
	normalizeTemplateTotals(tpl.Name, payload)

	generationID := strings.TrimSpace(asString(payload["_generation_id"]))
	if generationID == "" {
		generationID = uuid.New().String()
	}
	correlationID := asString(payload["correlation_id"])
	if correlationID == "" {
		correlationID = generationID
	}
	entityTypeNormalized := strings.ToUpper(strings.TrimSpace(entityType))

	// Publish document generation requested event (non-blocking)
	go func() {
		if s.kafkaPublisher != nil {
			if err := s.kafkaPublisher.PublishGenerationRequested(
				context.Background(),
				generationID,
				tpl.Id,
				tenantID,
				entityTypeNormalized,
				entityID,
				correlationID,
			); err != nil {
				logger.Warnf("failed to publish generation.requested event: %v", err)
			}
		}
	}()

	if includeQRCode {
		qrData, qrErr := buildQRCodeDataURI(fmt.Sprintf("doc:%s|entity:%s|id:%s", generationID, entityTypeNormalized, entityID))
		if qrErr == nil {
			payload["qr_code_data_uri"] = qrData
		}
	}

	renderedHTML, err := renderTemplate(tpl.TemplateContent, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to render template: %w", err)
	}
	renderedHTML = s.inlineTemplateLocalAssets(renderedHTML)

	fileContent, contentType, fileExt, err := buildOutput(renderedHTML, tpl.OutputFormat, s.gotenbergURL, s.pdfTimeout)
	if err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("%s_%s%s", slugify(tpl.Name), generationID, fileExt)
	fileURL := ""
	storageFileID := ""
	if s.storageClient != nil && strings.TrimSpace(tenantID) != "" {
		uploaded, upErr := s.storageClient.UploadFile(ctx, &storageservicev1.UploadFileRequest{
			TenantId:      tenantID,
			Content:       fileContent,
			Filename:      filename,
			ContentType:   contentType,
			FileType:      mapTemplateToStorageType(tpl.Type),
			ReferenceId:   entityID,
			ReferenceType: entityTypeNormalized,
			IsPublic:      false,
			ExpiresAt:     nil,
		})
		if upErr == nil && uploaded.GetFile() != nil {
			storageFileID = uploaded.File.FileId
			if uploaded.File.CdnUrl != "" {
				fileURL = uploaded.File.CdnUrl
			} else {
				fileURL = uploaded.File.Url
			}
		}
	}
	if fileURL == "" {
		fileURL = fmt.Sprintf("inline://documents/%s", generationID)
	}

	rawData := map[string]any{}
	for k, v := range payload {
		rawData[k] = v
	}
	rawData["_rendered_content_b64"] = base64.StdEncoding.EncodeToString(fileContent)
	rawData["_content_type"] = contentType
	rawData["_filename"] = filename
	rawData["_storage_file_id"] = storageFileID
	rawJSON, _ := json.Marshal(rawData)

	doc := &documentv1.DocumentGeneration{
		Id:                 generationID,
		DocumentTemplateId: tpl.Id,
		EntityType:         entityTypeNormalized,
		EntityId:           entityID,
		Data:               string(rawJSON),
		Status:             documentv1.GenerationStatus_GENERATION_STATUS_COMPLETED,
		FileUrl:            fileURL,
		FileSizeBytes:      int64(len(fileContent)),
		QrCodeData:         asString(payload["qr_code_data_uri"]),
		GeneratedBy:        generatedBy,
		GeneratedAt:        timestamppb.Now(),
	}

	created, err := s.generationRepo.Create(ctx, doc)
	if err != nil {
		// Publish generation.failed event (non-blocking)
		go func() {
			if s.kafkaPublisher != nil {
				if err := s.kafkaPublisher.PublishGenerationFailed(
					context.Background(),
					doc.Id,
					tenantID,
					fmt.Sprintf("failed to save generated document: %v", err),
					correlationID,
				); err != nil {
					logger.Warnf("failed to publish generation.failed event: %v", err)
				}
			}
		}()
		return nil, err
	}

	// Publish document.generated event to Kafka (non-blocking)
	go func() {
		if s.kafkaPublisher != nil {
			if err := s.kafkaPublisher.PublishDocumentGenerated(
				context.Background(),
				created.Id,
				tenantID,
				created.EntityId,
				created.EntityType,
				created.FileUrl,
				correlationID,
			); err != nil {
				logger.Warnf("failed to publish document.generated event: %v", err)
			}
		}
	}()

	return created, nil
}

func (s *DocumentService) GetDocument(ctx context.Context, documentID string) (*documentv1.DocumentGeneration, error) {
	if strings.TrimSpace(documentID) == "" {
		return nil, fmt.Errorf("%w: document_id is required", ErrInvalidInput)
	}
	return s.generationRepo.GetByID(ctx, documentID)
}

func (s *DocumentService) ListDocuments(
	ctx context.Context,
	entityType, entityID, status string,
	page, pageSize int,
) ([]*documentv1.DocumentGeneration, int, error) {
	if strings.TrimSpace(entityType) == "" || strings.TrimSpace(entityID) == "" {
		return nil, 0, fmt.Errorf("%w: entity_type and entity_id are required", ErrInvalidInput)
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	var statusFilter *documentv1.GenerationStatus
	if strings.TrimSpace(status) != "" {
		parsed, err := parseGenerationStatus(status)
		if err != nil {
			return nil, 0, err
		}
		statusFilter = &parsed
	}
	offset := (page - 1) * pageSize
	return s.generationRepo.ListByEntity(ctx, entityType, entityID, statusFilter, pageSize, offset)
}

func (s *DocumentService) DownloadDocument(ctx context.Context, documentID string) ([]byte, string, string, error) {
	doc, err := s.GetDocument(ctx, documentID)
	if err != nil {
		return nil, "", "", err
	}
	if strings.TrimSpace(doc.Data) == "" {
		return nil, "", "", fmt.Errorf("%w: document content not available", ErrInvalidInput)
	}
	var data map[string]any
	if err := json.Unmarshal([]byte(doc.Data), &data); err != nil {
		return nil, "", "", fmt.Errorf("failed to parse document data: %w", err)
	}
	encoded := asString(data["_rendered_content_b64"])
	if encoded == "" {
		return nil, "", "", fmt.Errorf("%w: rendered content missing", ErrInvalidInput)
	}
	content, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to decode content: %w", err)
	}
	contentType := asString(data["_content_type"])
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	filename := asString(data["_filename"])
	if filename == "" {
		filename = fmt.Sprintf("document_%s", documentID)
	}
	return content, contentType, filename, nil
}

func (s *DocumentService) DeleteDocument(ctx context.Context, documentID, tenantID string) error {
	doc, err := s.GetDocument(ctx, documentID)
	if err != nil {
		return err
	}
	if s.storageClient != nil && strings.TrimSpace(tenantID) != "" && strings.TrimSpace(doc.Data) != "" {
		var payload map[string]any
		if json.Unmarshal([]byte(doc.Data), &payload) == nil {
			if fileID := asString(payload["_storage_file_id"]); fileID != "" {
				_, _ = s.storageClient.DeleteFile(ctx, &storageservicev1.DeleteFileRequest{
					TenantId: tenantID,
					FileId:   fileID,
				})
			}
		}
	}
	return s.generationRepo.Delete(ctx, documentID)
}

func (s *DocumentService) CreateTemplate(
	ctx context.Context,
	name, typeStr, description, templateContent, outputFormat string,
	variables []string,
	createdBy string,
) (string, error) {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(templateContent) == "" {
		return "", fmt.Errorf("%w: name and template_content are required", ErrInvalidInput)
	}
	docType, err := parseDocumentType(typeStr)
	if err != nil {
		return "", err
	}
	format, err := parseOutputFormat(outputFormat)
	if err != nil {
		return "", err
	}
	if len(variables) == 0 {
		variables = extractTemplateVariables(templateContent)
	}
	variablesJSON, _ := json.Marshal(variables)

	tpl := &documentv1.DocumentTemplate{
		Id:              uuid.New().String(),
		Name:            name,
		Type:            docType,
		Description:     description,
		TemplateContent: templateContent,
		OutputFormat:    format,
		Variables:       string(variablesJSON),
		Version:         1,
		IsActive:        true,
		AuditInfo:       nil,
	}
	if strings.TrimSpace(createdBy) != "" {
		tpl.AuditInfo = nil
	}
	created, err := s.templateRepo.Create(ctx, tpl)
	if err != nil {
		return "", err
	}

	// Publish document.template.created event to Kafka (non-blocking)
	go func() {
		if s.kafkaPublisher != nil {
			if err := s.kafkaPublisher.PublishTemplateCreated(
				context.Background(),
				created.Id,
				"", // tenantID not available in CreateTemplate - can be enhanced later
				created.Name,
			); err != nil {
				logger.Warnf("failed to publish template.created event: %v", err)
			}
		}
	}()

	return created.Id, nil
}

func (s *DocumentService) GetTemplate(ctx context.Context, templateID string) (*documentv1.DocumentTemplate, error) {
	if strings.TrimSpace(templateID) == "" {
		return nil, fmt.Errorf("%w: template_id is required", ErrInvalidInput)
	}
	return s.templateRepo.GetByID(ctx, templateID)
}

func (s *DocumentService) ListTemplates(
	ctx context.Context,
	typeStr string,
	activeOnly bool,
	pageSize int,
	pageToken string,
) ([]*documentv1.DocumentTemplate, string, int, error) {
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	offset := 0
	if strings.TrimSpace(pageToken) != "" {
		n, err := strconv.Atoi(pageToken)
		if err == nil && n >= 0 {
			offset = n
		}
	}

	var docType *documentv1.DocumentType
	if strings.TrimSpace(typeStr) != "" {
		parsed, err := parseDocumentType(typeStr)
		if err != nil {
			return nil, "", 0, err
		}
		docType = &parsed
	}
	items, total, err := s.templateRepo.List(ctx, docType, activeOnly, pageSize, offset)
	if err != nil {
		return nil, "", 0, err
	}
	next := ""
	if offset+len(items) < total {
		next = strconv.Itoa(offset + len(items))
	}
	return items, next, total, nil
}

func (s *DocumentService) UpdateTemplate(ctx context.Context, templateID string, tpl *documentv1.DocumentTemplate) error {
	if strings.TrimSpace(templateID) == "" {
		return fmt.Errorf("%w: template_id is required", ErrInvalidInput)
	}
	if tpl == nil {
		return fmt.Errorf("%w: template payload is required", ErrInvalidInput)
	}
	if strings.TrimSpace(tpl.Name) == "" || strings.TrimSpace(tpl.TemplateContent) == "" {
		return fmt.Errorf("%w: name and template_content are required", ErrInvalidInput)
	}
	if tpl.Type == documentv1.DocumentType_DOCUMENT_TYPE_UNSPECIFIED {
		return fmt.Errorf("%w: valid template type is required", ErrInvalidInput)
	}
	if tpl.OutputFormat == documentv1.OutputFormat_OUTPUT_FORMAT_UNSPECIFIED {
		return fmt.Errorf("%w: valid output_format is required", ErrInvalidInput)
	}
	if strings.TrimSpace(tpl.Variables) == "" {
		vars := extractTemplateVariables(tpl.TemplateContent)
		vb, _ := json.Marshal(vars)
		tpl.Variables = string(vb)
	}
	if tpl.Version == 0 {
		tpl.Version = 1
	}
	err := s.templateRepo.Update(ctx, templateID, tpl)
	if err != nil {
		return err
	}

	// Publish document.template.updated event to Kafka (non-blocking)
	go func() {
		if s.kafkaPublisher != nil {
			if err := s.kafkaPublisher.PublishTemplateUpdated(
				context.Background(),
				templateID,
				"", // tenantID not available in UpdateTemplate - can be enhanced later
			); err != nil {
				logger.Warnf("failed to publish template.updated event: %v", err)
			}
		}
	}()

	return nil
}

func (s *DocumentService) DeactivateTemplate(ctx context.Context, templateID string) error {
	if strings.TrimSpace(templateID) == "" {
		return fmt.Errorf("%w: template_id is required", ErrInvalidInput)
	}
	return s.templateRepo.Deactivate(ctx, templateID)
}

func (s *DocumentService) DeleteTemplate(ctx context.Context, templateID string) error {
	if strings.TrimSpace(templateID) == "" {
		return fmt.Errorf("%w: template_id is required", ErrInvalidInput)
	}
	return s.templateRepo.Delete(ctx, templateID)
}

func (s *DocumentService) resolveTemplateForGeneration(ctx context.Context, templateRef string) (*documentv1.DocumentTemplate, error) {
	ref := strings.TrimSpace(templateRef)
	if ref == "" {
		return nil, repository.ErrTemplateNotFound
	}

	// Allow callers to pass either UUID template IDs or stable template names
	// such as "b2b_pi"/"b2c_po".
	if _, parseErr := uuid.Parse(ref); parseErr == nil {
		tpl, err := s.templateRepo.GetByID(ctx, ref)
		if err == nil {
			return tpl, nil
		}
		if !errors.Is(err, repository.ErrTemplateNotFound) {
			return nil, err
		}
	}

	byName, byNameErr := s.templateRepo.GetByName(ctx, ref)
	if byNameErr == nil {
		return byName, nil
	}
	return nil, repository.ErrTemplateNotFound
}

func (s *DocumentService) bootstrapDefaultTemplates(ctx context.Context) error {
	type templateSpec struct {
		Name        string
		Relative    string
		Type        documentv1.DocumentType
		Description string
	}

	specs := []templateSpec{
		{
			Name:        "b2b_pi",
			Relative:    filepath.Join("b2b", "pi.html"),
			Type:        documentv1.DocumentType_DOCUMENT_TYPE_INVOICE,
			Description: "B2B premium invoice template",
		},
		{
			Name:        "b2c_pi",
			Relative:    filepath.Join("b2c", "pi.html"),
			Type:        documentv1.DocumentType_DOCUMENT_TYPE_INVOICE,
			Description: "B2C premium invoice template",
		},
		{
			Name:        "b2b_po",
			Relative:    filepath.Join("b2b", "po.html"),
			Type:        documentv1.DocumentType_DOCUMENT_TYPE_RECEIPT,
			Description: "B2B policy enrollment order template",
		},
		{
			Name:        "b2c_po",
			Relative:    filepath.Join("b2c", "po.html"),
			Type:        documentv1.DocumentType_DOCUMENT_TYPE_RECEIPT,
			Description: "B2C insurance service order template",
		},
		{
			Name:        "policy_document",
			Relative:    filepath.Join("b2c", "policy_document.html"),
			Type:        documentv1.DocumentType_DOCUMENT_TYPE_POLICY_CERTIFICATE,
			Description: "Rich policy document template",
		},
		// Backward-compatible aliases for existing callers.
		{
			Name:        "invoice",
			Relative:    filepath.Join("b2c", "pi.html"),
			Type:        documentv1.DocumentType_DOCUMENT_TYPE_INVOICE,
			Description: "Legacy invoice alias to B2C PI template",
		},
		{
			Name:        "purchase_order",
			Relative:    filepath.Join("b2b", "po.html"),
			Type:        documentv1.DocumentType_DOCUMENT_TYPE_RECEIPT,
			Description: "Legacy purchase order alias to B2B PO template",
		},
	}

	for _, spec := range specs {
		path := filepath.Join(s.templateDirPath, spec.Relative)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", path, err)
		}
		vars := extractTemplateVariables(string(content))
		varsJSON, _ := json.Marshal(vars)

		_, err = s.templateRepo.UpsertByName(ctx, &documentv1.DocumentTemplate{
			Id:              uuid.New().String(),
			Name:            spec.Name,
			Type:            spec.Type,
			Description:     spec.Description,
			TemplateContent: string(content),
			OutputFormat:    documentv1.OutputFormat_OUTPUT_FORMAT_PDF,
			Variables:       string(varsJSON),
			Version:         1,
			IsActive:        true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func renderTemplate(content string, data map[string]any) (string, error) {
	funcMap := htmltmpl.FuncMap{
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
	}
	tpl, err := htmltmpl.New("doc").Funcs(funcMap).Parse(content)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (s *DocumentService) inlineTemplateLocalAssets(content string) string {
	matches := imgSrcAttrPattern.FindAllStringSubmatchIndex(content, -1)
	if len(matches) == 0 {
		return content
	}

	var out strings.Builder
	out.Grow(len(content) + 256)
	last := 0
	for _, m := range matches {
		if len(m) < 6 {
			continue
		}
		srcStart := m[4]
		srcEnd := m[5]
		if srcStart < last || srcEnd > len(content) || srcStart >= srcEnd {
			continue
		}

		out.WriteString(content[last:srcStart])
		src := content[srcStart:srcEnd]
		if dataURI, ok := s.localAssetToDataURI(src); ok {
			out.WriteString(dataURI)
		} else {
			out.WriteString(src)
		}
		last = srcEnd
	}
	out.WriteString(content[last:])
	return out.String()
}

func (s *DocumentService) localAssetToDataURI(src string) (string, bool) {
	raw := strings.TrimSpace(src)
	if raw == "" {
		return "", false
	}
	lower := strings.ToLower(raw)
	if strings.HasPrefix(lower, "data:") ||
		strings.HasPrefix(lower, "http://") ||
		strings.HasPrefix(lower, "https://") ||
		strings.HasPrefix(lower, "cid:") ||
		strings.HasPrefix(lower, "file:") ||
		strings.HasPrefix(lower, "//") {
		return "", false
	}

	normalized := strings.ReplaceAll(raw, "\\", "/")
	if i := strings.IndexAny(normalized, "?#"); i >= 0 {
		normalized = normalized[:i]
	}
	normalized = strings.TrimSpace(normalized)
	if normalized == "" {
		return "", false
	}
	for strings.HasPrefix(normalized, "./") {
		normalized = strings.TrimPrefix(normalized, "./")
	}
	for strings.HasPrefix(normalized, "../") {
		normalized = strings.TrimPrefix(normalized, "../")
	}
	normalized = strings.TrimPrefix(normalized, "/")
	if normalized == "" {
		return "", false
	}

	candidates := make([]string, 0, 3)
	candidates = append(candidates, filepath.Join(s.templateDirPath, filepath.FromSlash(normalized)))
	if idx := strings.Index(strings.ToLower(normalized), "logos/"); idx >= 0 {
		candidates = append(candidates, filepath.Join(s.templateDirPath, filepath.FromSlash(normalized[idx:])))
	}
	base := filepath.Base(normalized)
	if base != "" && base != "." && base != string(filepath.Separator) {
		candidates = append(candidates, filepath.Join(s.templateDirPath, "logos", base))
	}

	seen := make(map[string]struct{}, len(candidates))
	for _, candidate := range candidates {
		clean := filepath.Clean(candidate)
		if _, ok := seen[clean]; ok {
			continue
		}
		seen[clean] = struct{}{}
		if !isWithinDir(clean, s.templateDirPath) {
			continue
		}

		data, err := os.ReadFile(clean)
		if err != nil || len(data) == 0 {
			continue
		}
		mimeType := mime.TypeByExtension(strings.ToLower(filepath.Ext(clean)))
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}
		return "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(data), true
	}

	return "", false
}

func isWithinDir(path, root string) bool {
	absPath, err1 := filepath.Abs(path)
	absRoot, err2 := filepath.Abs(root)
	if err1 != nil || err2 != nil {
		return false
	}

	absPath = filepath.Clean(absPath)
	absRoot = filepath.Clean(absRoot)
	pathLower := strings.ToLower(absPath)
	rootLower := strings.ToLower(absRoot)

	if pathLower == rootLower {
		return true
	}
	return strings.HasPrefix(pathLower, rootLower+string(filepath.Separator))
}

func buildOutput(renderedHTML string, format documentv1.OutputFormat, gotenbergURL string, timeout time.Duration) ([]byte, string, string, error) {
	switch format {
	case documentv1.OutputFormat_OUTPUT_FORMAT_HTML, documentv1.OutputFormat_OUTPUT_FORMAT_UNSPECIFIED:
		return []byte(renderedHTML), "text/html; charset=utf-8", ".html", nil
	case documentv1.OutputFormat_OUTPUT_FORMAT_PDF:
		pdfBytes, err := renderPDF(renderedHTML, gotenbergURL, timeout)
		if err != nil {
			return nil, "", "", err
		}
		return pdfBytes, "application/pdf", ".pdf", nil
	case documentv1.OutputFormat_OUTPUT_FORMAT_DOCX:
		return nil, "", "", fmt.Errorf("%w: DOCX", ErrUnsupportedOutput)
	default:
		return nil, "", "", fmt.Errorf("%w: unknown output format", ErrUnsupportedOutput)
	}
}

func renderPDF(renderedHTML, gotenbergURL string, timeout time.Duration) ([]byte, error) {
	if strings.TrimSpace(gotenbergURL) != "" {
		pdf, err := renderPDFWithGotenberg(renderedHTML, gotenbergURL, timeout)
		if err == nil {
			return pdf, nil
		}
		logger.Warnf("gotenberg pdf conversion failed, falling back to basic renderer: %v", err)
	}

	return renderPDFFallback(renderedHTML)
}

func renderPDFWithGotenberg(renderedHTML, gotenbergURL string, timeout time.Duration) ([]byte, error) {
	if timeout <= 0 {
		timeout = 8 * time.Second
	}

	var reqBody bytes.Buffer
	writer := multipart.NewWriter(&reqBody)
	filePart, err := writer.CreateFormFile("files", "index.html")
	if err != nil {
		return nil, fmt.Errorf("failed to create gotenberg multipart file: %w", err)
	}
	if _, err := io.WriteString(filePart, renderedHTML); err != nil {
		return nil, fmt.Errorf("failed to write html to gotenberg request: %w", err)
	}
	_ = writer.WriteField("paperWidth", "8.27")
	_ = writer.WriteField("paperHeight", "11.69")
	_ = writer.WriteField("marginTop", "0.25")
	_ = writer.WriteField("marginBottom", "0.25")
	_ = writer.WriteField("marginLeft", "0.25")
	_ = writer.WriteField("marginRight", "0.25")
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to finalize gotenberg request body: %w", err)
	}

	endpoint := strings.TrimRight(strings.TrimSpace(gotenbergURL), "/") + "/forms/chromium/convert/html"
	req, err := http.NewRequest(http.MethodPost, endpoint, &reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create gotenberg request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gotenberg request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("gotenberg returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	pdf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read gotenberg response: %w", err)
	}
	if len(pdf) == 0 {
		return nil, fmt.Errorf("gotenberg returned empty pdf")
	}
	return pdf, nil
}

func renderPDFFallback(renderedHTML string) ([]byte, error) {
	plain := strings.TrimSpace(html.UnescapeString(htmlTagPattern.ReplaceAllString(renderedHTML, " ")))
	plain = regexp.MustCompile(`\s+`).ReplaceAllString(plain, " ")
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(12, 12, 12)
	pdf.AddPage()
	pdf.SetFont("Arial", "", 11)
	pdf.MultiCell(0, 6, plain, "", "L", false)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to render pdf: %w", err)
	}
	return buf.Bytes(), nil
}

func extractTemplateVariables(content string) []string {
	matches := varPattern.FindAllStringSubmatch(content, -1)
	seen := map[string]struct{}{}
	vars := make([]string, 0, len(matches))
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		v := strings.TrimSpace(m[1])
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		vars = append(vars, v)
	}
	return vars
}

func buildQRCodeDataURI(value string) (string, error) {
	code, err := qr.Encode(value, qr.M, qr.Auto)
	if err != nil {
		return "", err
	}
	scaled, err := barcode.Scale(code, 256, 256)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := pngEncode(&buf, scaled); err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func pngEncode(buf *bytes.Buffer, img image.Image) error {
	return png.Encode(buf, img)
}

func applyBusinessDefaults(templateName string, payload map[string]any) {
	now := time.Now().UTC().Format("2006-01-02")
	setIfMissing(payload, "issue_date", now)
	setIfMissing(payload, "terms", "")

	switch templateKind(templateName) {
	case "invoice":
		setIfMissing(payload, "invoice_number", "INV-"+time.Now().UTC().Format("20060102150405"))
		setIfMissing(payload, "due_date", time.Now().UTC().Add(7*24*time.Hour).Format("2006-01-02"))
		setIfMissing(payload, "subtotal", "0.00")
		setIfMissing(payload, "tax", "0.00")
		setIfMissing(payload, "total", "0.00")
		ensureDefaultItems(payload)
	case "purchase_order":
		setIfMissing(payload, "purchase_order_number", "PO-"+time.Now().UTC().Format("20060102150405"))
		setIfMissing(payload, "subtotal", "0.00")
		setIfMissing(payload, "shipping_cost", "0.00")
		setIfMissing(payload, "tax", "0.00")
		setIfMissing(payload, "total", "0.00")
		ensureDefaultItems(payload)
	case "policy_document":
		setIfMissing(payload, "policy_title", "Insurance Policy Certificate")
		setIfMissing(payload, "policy_number", "POL-"+time.Now().UTC().Format("20060102150405"))
		setIfMissing(payload, "policy_holder_name", "Policy Holder")
		setIfMissing(payload, "start_date", now)
		setIfMissing(payload, "end_date", time.Now().UTC().AddDate(1, 0, 0).Format("2006-01-02"))
		setIfMissing(payload, "coverage_amount", "0.00")
		setIfMissing(payload, "premium_amount", "0.00")
		if _, ok := payload["benefits"]; !ok {
			payload["benefits"] = []any{"Coverage details will be provided by insurer"}
		}
	}
}

func templateKind(templateName string) string {
	name := strings.ToLower(strings.TrimSpace(templateName))
	switch {
	case name == "invoice" || strings.HasSuffix(name, "_pi") || strings.Contains(name, "invoice"):
		return "invoice"
	case name == "purchase_order" || strings.HasSuffix(name, "_po") || strings.Contains(name, "purchase_order"):
		return "purchase_order"
	case name == "policy_document":
		return "policy_document"
	default:
		return name
	}
}

func mapTemplateToStorageType(docType documentv1.DocumentType) storageentityv1.FileType {
	if docType == documentv1.DocumentType_DOCUMENT_TYPE_INVOICE {
		return storageentityv1.FileType_FILE_TYPE_INVOICE
	}
	if docType == documentv1.DocumentType_DOCUMENT_TYPE_RECEIPT {
		return storageentityv1.FileType_FILE_TYPE_RECEIPT
	}
	return storageentityv1.FileType_FILE_TYPE_DOCUMENT
}

func parseDocumentType(v string) (documentv1.DocumentType, error) {
	s := strings.TrimSpace(strings.ToUpper(v))
	if s == "" {
		return documentv1.DocumentType_DOCUMENT_TYPE_UNSPECIFIED, fmt.Errorf("%w: template type is required", ErrInvalidInput)
	}
	if !strings.HasPrefix(s, "DOCUMENT_TYPE_") {
		s = "DOCUMENT_TYPE_" + s
	}
	n, ok := documentv1.DocumentType_value[s]
	if !ok || n == int32(documentv1.DocumentType_DOCUMENT_TYPE_UNSPECIFIED) {
		return documentv1.DocumentType_DOCUMENT_TYPE_UNSPECIFIED, fmt.Errorf("%w: invalid template type", ErrInvalidInput)
	}
	return documentv1.DocumentType(n), nil
}

func parseOutputFormat(v string) (documentv1.OutputFormat, error) {
	s := strings.TrimSpace(strings.ToUpper(v))
	if s == "" {
		s = "OUTPUT_FORMAT_HTML"
	}
	if !strings.HasPrefix(s, "OUTPUT_FORMAT_") {
		s = "OUTPUT_FORMAT_" + s
	}
	n, ok := documentv1.OutputFormat_value[s]
	if !ok || n == int32(documentv1.OutputFormat_OUTPUT_FORMAT_UNSPECIFIED) {
		return documentv1.OutputFormat_OUTPUT_FORMAT_UNSPECIFIED, fmt.Errorf("%w: invalid output format", ErrInvalidInput)
	}
	return documentv1.OutputFormat(n), nil
}

func parseGenerationStatus(v string) (documentv1.GenerationStatus, error) {
	s := strings.TrimSpace(strings.ToUpper(v))
	if s == "" {
		return documentv1.GenerationStatus_GENERATION_STATUS_UNSPECIFIED, nil
	}
	if !strings.HasPrefix(s, "GENERATION_STATUS_") {
		s = "GENERATION_STATUS_" + s
	}
	n, ok := documentv1.GenerationStatus_value[s]
	if !ok {
		return documentv1.GenerationStatus_GENERATION_STATUS_UNSPECIFIED, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}
	return documentv1.GenerationStatus(n), nil
}

func slugify(v string) string {
	s := strings.ToLower(strings.TrimSpace(v))
	s = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(s, "_")
	s = strings.Trim(s, "_")
	if s == "" {
		return "document"
	}
	return s
}

func setIfMissing(m map[string]any, key string, value any) {
	if _, ok := m[key]; !ok {
		m[key] = value
	}
}

func ensureDefaultItems(m map[string]any) {
	if _, ok := m["items"]; ok {
		return
	}
	m["items"] = []any{map[string]any{
		"description": "Document item",
		"quantity":    1,
		"unit_price":  "0.00",
		"amount":      "0.00",
	}}
}

func ensureMap(m map[string]any) map[string]any {
	if m == nil {
		return map[string]any{}
	}
	return m
}

func asString(v any) string {
	if v == nil {
		return ""
	}
	s, ok := v.(string)
	if ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}
