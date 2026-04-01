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
	"github.com/newage-saint/insuretech/backend/inscore/microservices/docgen/internal/renderer"
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
	// sidecarClient is the optional Python docrender sidecar for DOCX and
	// WeasyPrint PDF generation. May be nil when the sidecar is not configured.
	sidecarClient *renderer.SidecarClient
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

// SetDocRenderer configures the Python docrender sidecar used for DOCX and
// high-quality WeasyPrint PDF generation. sidecarURL example: "http://localhost:8500".
func (s *DocumentService) SetDocRenderer(sidecarURL string, timeout time.Duration) {
	url := strings.TrimSpace(sidecarURL)
	if url == "" {
		return
	}
	s.sidecarClient = renderer.NewSidecarClient(url, timeout)
}

func (s *DocumentService) GenerateDocument(
	ctx context.Context,
	templateID, entityType, entityID string,
	data *structpb.Struct,
	includeQRCode bool,
	tenantID, generatedBy string,
	outputFormatHint string, // per-request override; empty = use template default
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

	// Resolve effective output format: per-request hint overrides template default.
	effectiveFormat := tpl.OutputFormat
	if hint := strings.TrimSpace(strings.ToLower(outputFormatHint)); hint != "" {
		effectiveFormat = parseOutputFormatHint(hint)
	}
	fileContent, contentType, fileExt, err := buildOutput(ctx, renderedHTML, tpl.TemplateContent, payload, tpl.Name, effectiveFormat, s.gotenbergURL, s.pdfTimeout, s.sidecarClient)
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

// GenerateOptions controls per-request generation behaviour.
type GenerateOptions struct {
	// IncludeQRCode embeds a QR code data URI in the payload.
	IncludeQRCode bool
	// OutputFormatHint overrides the template's configured output_format.
	// Accepts: "pdf", "html", "docx", "xlsx" (case-insensitive).
	// Empty string means use the template's default.
	OutputFormatHint string
	// AltMedia signals that the caller wants raw file bytes returned
	// directly in the response (Google API ?$alt=media style).
	AltMedia bool
}

// GenerateResult is returned by GenerateDocumentEx.
type GenerateResult struct {
	DocumentID  string
	FileURL     string
	// FileBytes is populated when AltMedia=true.
	FileBytes   []byte
	ContentType string
	Filename    string
}

// GenerateDocumentEx is the primary generation entry point.
// It honours per-request format overrides and the alt=media flag.
// The legacy GenerateDocument method delegates here.
func (s *DocumentService) GenerateDocumentEx(
	ctx context.Context,
	templateID, entityType, entityID string,
	data *structpb.Struct,
	tenantID, generatedBy string,
	opts GenerateOptions,
) (*GenerateResult, error) {
	doc, err := s.GenerateDocument(ctx, templateID, entityType, entityID, data, opts.IncludeQRCode, tenantID, generatedBy, opts.OutputFormatHint)
	if err != nil {
		return nil, err
	}
	result := &GenerateResult{
		DocumentID: doc.Id,
		FileURL:    doc.FileUrl,
	}
	// When alt=media requested, decode the raw bytes from the stored data.
	if opts.AltMedia && doc.Data != "" {
		var rawData map[string]any
		if jsonErr := json.Unmarshal([]byte(doc.Data), &rawData); jsonErr == nil {
			if b64, ok := rawData["_rendered_content_b64"].(string); ok && b64 != "" {
				if decoded, decErr := base64.StdEncoding.DecodeString(b64); decErr == nil {
					result.FileBytes   = decoded
					result.ContentType = asString(rawData["_content_type"])
					result.Filename    = asString(rawData["_filename"])
				}
			}
		}
	}
	return result, nil
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
		// ── Insurance proposal & claim templates ──────────────────────────────
		{
			Name:        "motor_proposal",
			Relative:    filepath.Join("insurance", "motor_proposal.html"),
			Type:        documentv1.DocumentType_DOCUMENT_TYPE_POLICY_CERTIFICATE,
			Description: "Motor insurance proposal form (Class 1/2/3+)",
		},
		{
			Name:        "fire_proposal",
			Relative:    filepath.Join("insurance", "fire_proposal.html"),
			Type:        documentv1.DocumentType_DOCUMENT_TYPE_POLICY_CERTIFICATE,
			Description: "Fire & allied perils insurance proposal form",
		},
		{
			Name:        "omp_proposal",
			Relative:    filepath.Join("insurance", "omp_proposal.html"),
			Type:        documentv1.DocumentType_DOCUMENT_TYPE_POLICY_CERTIFICATE,
			Description: "Office / commercial package (OMP) proposal form",
		},
		{
			Name:        "motor_claim",
			Relative:    filepath.Join("insurance", "motor_claim.html"),
			Type:        documentv1.DocumentType_DOCUMENT_TYPE_CLAIM_FORM,
			Description: "Motor insurance claim form",
		},
		{
			Name:        "general_claim",
			Relative:    filepath.Join("insurance", "general_claim.html"),
			Type:        documentv1.DocumentType_DOCUMENT_TYPE_CLAIM_FORM,
			Description: "General insurance claim form (fire, OMP, miscellaneous)",
		},
		{
			Name:        "overseas_mediclaim_proposal",
			Relative:    filepath.Join("insurance", "overseas_mediclaim_proposal.json"),
			Type:        documentv1.DocumentType_DOCUMENT_TYPE_POLICY_CERTIFICATE,
			Description: "Overseas Mediclaim Proposal Form (Travel Insurance) — Business & Holidays",
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

// parseOutputFormatHint converts a string hint ("pdf", "docx", "xlsx", "html")
// to the corresponding OutputFormat enum value.
func parseOutputFormatHint(hint string) documentv1.OutputFormat {
	switch strings.ToLower(strings.TrimSpace(hint)) {
	case "pdf":
		return documentv1.OutputFormat_OUTPUT_FORMAT_PDF
	case "docx":
		return documentv1.OutputFormat_OUTPUT_FORMAT_DOCX
	case "xlsx":
		return documentv1.OutputFormat_OUTPUT_FORMAT_XLSX
	case "html":
		return documentv1.OutputFormat_OUTPUT_FORMAT_HTML
	default:
		return documentv1.OutputFormat_OUTPUT_FORMAT_UNSPECIFIED
	}
}

func buildOutput(
	ctx context.Context,
	renderedHTML string,
	templateContent string,
	payload map[string]any,
	templateName string,
	format documentv1.OutputFormat,
	gotenbergURL string,
	timeout time.Duration,
	sidecar *renderer.SidecarClient,
) ([]byte, string, string, error) {
	switch format {
	case documentv1.OutputFormat_OUTPUT_FORMAT_HTML, documentv1.OutputFormat_OUTPUT_FORMAT_UNSPECIFIED:
		return []byte(renderedHTML), "text/html; charset=utf-8", ".html", nil

	case documentv1.OutputFormat_OUTPUT_FORMAT_PDF:
		// Try Gotenberg (Chromium-based, best quality for complex HTML/CSS).
		// Fall back to WeasyPrint sidecar, then to the basic gofpdf fallback.
		if strings.TrimSpace(gotenbergURL) != "" {
			pdf, err := renderPDFWithGotenberg(renderedHTML, gotenbergURL, timeout)
			if err == nil {
				return pdf, "application/pdf", ".pdf", nil
			}
			logger.Warnf("gotenberg pdf conversion failed, trying weasyprint sidecar: %v", err)
		}
		if sidecar != nil {
			pdf, err := sidecar.RenderPDF(ctx, renderedHTML)
			if err == nil {
				return pdf, "application/pdf", ".pdf", nil
			}
			logger.Warnf("weasyprint sidecar pdf conversion failed, falling back to basic renderer: %v", err)
		}
		pdf, err := renderPDFFallback(renderedHTML)
		if err != nil {
			return nil, "", "", err
		}
		return pdf, "application/pdf", ".pdf", nil

	case documentv1.OutputFormat_OUTPUT_FORMAT_DOCX:
		if sidecar == nil {
			return nil, "", "", fmt.Errorf("%w: DOCX requires the docrender sidecar (DOC_RENDERER_URL). Please configure it", ErrUnsupportedOutput)
		}
		docxBytes, err := sidecar.RenderDOCX(ctx, renderer.DocxRequest{
			TemplateContent: templateContent,
			Data:            payload,
			Title:           asString(payload["policy_title"]),
			Author:          "InsureTech",
			Subject:         templateName,
		})
		if err != nil {
			return nil, "", "", fmt.Errorf("docx render failed: %w", err)
		}
		return docxBytes,
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			".docx", nil

	case documentv1.OutputFormat_OUTPUT_FORMAT_XLSX:
		xlsxBytes, err := buildXLSXOutput(payload, templateName)
		if err != nil {
			return nil, "", "", fmt.Errorf("xlsx render failed: %w", err)
		}
		return xlsxBytes,
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			".xlsx", nil

	default:
		return nil, "", "", fmt.Errorf("%w: unknown output format", ErrUnsupportedOutput)
	}
}

// buildXLSXOutput converts a document payload into a rich styled Excel workbook.
// It auto-detects line items, summary fields, and totals from the standard payload keys.
func buildXLSXOutput(payload map[string]any, templateName string) ([]byte, error) {
	// ── Line items ────────────────────────────────────────────────────────────
	var items []map[string]any
	if rawItems, ok := payload["items"]; ok {
		if list, ok := rawItems.([]any); ok {
			for _, it := range list {
				if m, ok := it.(map[string]any); ok {
					items = append(items, m)
				}
			}
		}
	}

	itemCols := []renderer.XLSXColumn{
		{Key: "description", Header: "Description",  Width: 38},
		{Key: "quantity",    Header: "Qty",           Width: 8},
		{Key: "unit_price",  Header: "Unit Price",    Width: 14, IsMoney: true},
		{Key: "amount",      Header: "Amount",        Width: 14, IsMoney: true},
	}

	// ── Totals ────────────────────────────────────────────────────────────────
	var totals []renderer.XLSXTotalRow
	if v := asString(payload["subtotal"]); v != "" && v != "0.00" {
		totals = append(totals, renderer.XLSXTotalRow{Label: "Subtotal", Value: v})
	}
	if v := asString(payload["tax"]); v != "" && v != "0.00" {
		totals = append(totals, renderer.XLSXTotalRow{Label: "Tax", Value: v})
	}
	if v := asString(payload["service_fee"]); v != "" && v != "0.00" {
		totals = append(totals, renderer.XLSXTotalRow{Label: "Service Fee", Value: v})
	}
	if v := asString(payload["shipping_cost"]); v != "" && v != "0.00" {
		totals = append(totals, renderer.XLSXTotalRow{Label: "Shipping", Value: v})
	}
	if v := asString(payload["total"]); v != "" {
		totals = append(totals, renderer.XLSXTotalRow{Label: "TOTAL", Value: v, IsBold: true})
	}

	// ── Summary fields (all non-item, non-private keys) ───────────────────────
	summary := make(map[string]any)
	skipKeys := map[string]bool{
		"items": true, "subtotal": true, "tax": true, "total": true,
		"service_fee": true, "shipping_cost": true, "processing_fee": true,
		"qr_code_data_uri": true,
	}
	for k, v := range payload {
		if strings.HasPrefix(k, "_") || skipKeys[k] {
			continue
		}
		switch v.(type) {
		case string, float64, float32, int, int32, int64, bool:
			summary[k] = v
		}
	}

	title := asString(payload["policy_title"])
	if title == "" {
		title = asString(payload["invoice_number"])
	}
	if title == "" {
		title = templateName
	}

	opts := renderer.XLSXOptions{
		Title:       title,
		Author:      "InsureTech",
		Subject:     templateName,
		Description: asString(payload["description"]),
		Items:       items,
		ItemColumns: itemCols,
		Summary:     summary,
		Totals:      totals,
	}

	return renderer.RenderXLSX(opts)
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

	case "overseas_mediclaim":
		setIfMissing(payload, "company_name", "Pragati Insurance Limited")
		setIfMissing(payload, "proposal_id", "OMP-"+time.Now().UTC().Format("20060102150405"))
		setIfMissing(payload, "generated_at", time.Now().UTC().Format("02 Jan 2006 15:04 UTC"))
		setIfMissing(payload, "plan_type", "")
		setIfMissing(payload, "trip_purpose", "")
		setIfMissing(payload, "departure_date", "")
		setIfMissing(payload, "days_abroad", "")
		setIfMissing(payload, "itinerary", "")
		setIfMissing(payload, "q1_good_health", "")
		setIfMissing(payload, "q2a_nervous", "")
		setIfMissing(payload, "q2b_heart", "")
		setIfMissing(payload, "q2c_hernia", "")
		setIfMissing(payload, "q2d_respiratory", "")
		setIfMissing(payload, "q2e_specialist", "")
		setIfMissing(payload, "q2f_future", "")
		setIfMissing(payload, "q3_additional_facts", "")
		setIfMissing(payload, "q4_winter_sports", "")
		setIfMissing(payload, "known_ailment_1", "")
		setIfMissing(payload, "known_ailment_2", "")
		setIfMissing(payload, "known_ailment_3", "")
		setIfMissing(payload, "known_ailment_4", "")
		setIfMissing(payload, "signature_place", "")
		setIfMissing(payload, "signature_date", "")
		if _, ok := payload["illness_history"]; !ok {
			payload["illness_history"] = []any{
				map[string]any{"nature_of_illness": "", "date_first_treated": "", "practitioner_details": ""},
				map[string]any{"nature_of_illness": "", "date_first_treated": "", "practitioner_details": ""},
				map[string]any{"nature_of_illness": "", "date_first_treated": "", "practitioner_details": ""},
			}
		}
		if _, ok := payload["product_benefits"]; !ok {
			payload["product_benefits"] = []any{
				map[string]any{"number": "01.", "benefit": "Medical Expenses & Hospitalization abroad (Worldwide excl. USA/Canada)", "limit": "US$ 50,000 — Excess USD 100"},
				map[string]any{"number": "02.", "benefit": "Medical Expenses & Hospitalization abroad (Worldwide incl. USA/Canada)", "limit": "US$ 100,000 — Excess USD 100"},
				map[string]any{"number": "03.", "benefit": "Medical Expenses & Hospitalization for Schengen Countries", "limit": "Euro 30,000 — Nil deductible"},
				map[string]any{"number": "04.", "benefit": "Transport or Repatriation in case of Illness or Accident", "limit": "Actual Expenses"},
				map[string]any{"number": "05.", "benefit": "Emergency Dental Care", "limit": "US$ 500 — Excess US$ 50"},
				map[string]any{"number": "06.", "benefit": "Repatriation of Family Member Travelling with the Insured", "limit": "Actual Expenses"},
				map[string]any{"number": "07.", "benefit": "Repatriation of Mortal Remains", "limit": "Actual Expenses"},
				map[string]any{"number": "08.", "benefit": "Travel of one immediate family member", "limit": "US$ 100/day — Max US$ 1,000"},
				map[string]any{"number": "09.", "benefit": "Emergency return home following death of a close family member", "limit": "Actual Expenses"},
			}
		}

	case "proposal":
		setIfMissing(payload, "proposal_id", "PROP-"+time.Now().UTC().Format("20060102150405"))
		setIfMissing(payload, "proposal_number", "PROP-"+time.Now().UTC().Format("20060102150405"))
		setIfMissing(payload, "proposal_date", now)
		setIfMissing(payload, "generated_at", time.Now().UTC().Format("02 Jan 2006 15:04 UTC"))
		setIfMissing(payload, "coverage_start_date", now)
		setIfMissing(payload, "coverage_end_date", time.Now().UTC().AddDate(1, 0, 0).Format("2006-01-02"))
		setIfMissing(payload, "total_premium", "0.00")
		setIfMissing(payload, "stamp_duty", "0.00")
		setIfMissing(payload, "vat_amount", "0.00")
		setIfMissing(payload, "basic_premium", "0.00")

	case "claim":
		setIfMissing(payload, "claim_id", "CLM-"+time.Now().UTC().Format("20060102150405"))
		setIfMissing(payload, "claim_number", "CLM-"+time.Now().UTC().Format("20060102150405"))
		setIfMissing(payload, "claim_date", now)
		setIfMissing(payload, "generated_at", time.Now().UTC().Format("02 Jan 2006 15:04 UTC"))
		setIfMissing(payload, "net_claim_amount", "0.00")
		setIfMissing(payload, "deductible_amount", "0.00")
		setIfMissing(payload, "total_loss_amount", "0.00")
		setIfMissing(payload, "claim_status", "New")
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
	case strings.Contains(name, "mediclaim"):
		return "overseas_mediclaim"
	case strings.Contains(name, "proposal"):
		return "proposal"
	case strings.Contains(name, "claim"):
		return "claim"
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
