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

	tpl, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil {
		// Publish generation.failed event (non-blocking)
		go func() {
			if s.kafkaPublisher != nil {
				if err := s.kafkaPublisher.PublishGenerationFailed(
					context.Background(),
					uuid.New().String(), // placeholder ID since we couldn't start generation
					tenantID,
					fmt.Sprintf("failed to fetch template: %v", err),
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
	applyBusinessDefaults(tpl.Name, payload)

	generationID := uuid.New().String()
	if includeQRCode {
		qrData, qrErr := buildQRCodeDataURI(fmt.Sprintf("doc:%s|entity:%s|id:%s", generationID, entityType, entityID))
		if qrErr == nil {
			payload["qr_code_data_uri"] = qrData
		}
	}

	renderedHTML, err := renderTemplate(tpl.TemplateContent, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to render template: %w", err)
	}

	fileContent, contentType, fileExt, err := buildOutput(renderedHTML, tpl.OutputFormat)
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
			ReferenceType: strings.ToUpper(entityType),
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
		EntityType:         entityType,
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
				created.DocumentTemplateId,
				tenantID,
				created.EntityId,
				created.EntityType,
				created.QrCodeData,
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

func (s *DocumentService) bootstrapDefaultTemplates(ctx context.Context) error {
	types := map[string]documentv1.DocumentType{
		"invoice":         documentv1.DocumentType_DOCUMENT_TYPE_INVOICE,
		"purchase_order":  documentv1.DocumentType_DOCUMENT_TYPE_RECEIPT,
		"policy_document": documentv1.DocumentType_DOCUMENT_TYPE_POLICY_CERTIFICATE,
	}
	descriptions := map[string]string{
		"invoice":         "Rich invoice template",
		"purchase_order":  "Rich purchase order template",
		"policy_document": "Rich policy document template",
	}

	for name, docType := range types {
		path := filepath.Join(s.templateDirPath, name+".html")
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", path, err)
		}
		vars := extractTemplateVariables(string(content))
		varsJSON, _ := json.Marshal(vars)

		_, err = s.templateRepo.UpsertByName(ctx, &documentv1.DocumentTemplate{
			Id:              uuid.New().String(),
			Name:            name,
			Type:            docType,
			Description:     descriptions[name],
			TemplateContent: string(content),
			OutputFormat:    documentv1.OutputFormat_OUTPUT_FORMAT_HTML,
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

func buildOutput(renderedHTML string, format documentv1.OutputFormat) ([]byte, string, string, error) {
	switch format {
	case documentv1.OutputFormat_OUTPUT_FORMAT_HTML, documentv1.OutputFormat_OUTPUT_FORMAT_UNSPECIFIED:
		return []byte(renderedHTML), "text/html; charset=utf-8", ".html", nil
	case documentv1.OutputFormat_OUTPUT_FORMAT_PDF:
		pdfBytes, err := renderPDF(renderedHTML)
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

func renderPDF(renderedHTML string) ([]byte, error) {
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
	setIfMissing(payload, "company_name", "InsureTech")
	setIfMissing(payload, "company_address", "Dhaka, Bangladesh")
	setIfMissing(payload, "issue_date", now)
	setIfMissing(payload, "terms", "Standard terms and conditions apply.")

	switch strings.ToLower(strings.TrimSpace(templateName)) {
	case "invoice":
		setIfMissing(payload, "invoice_number", "INV-"+time.Now().UTC().Format("20060102150405"))
		setIfMissing(payload, "due_date", time.Now().UTC().Add(7*24*time.Hour).Format("2006-01-02"))
		setIfMissing(payload, "customer_name", "Policy Holder")
		setIfMissing(payload, "subtotal", "0.00")
		setIfMissing(payload, "tax", "0.00")
		setIfMissing(payload, "total", "0.00")
		ensureDefaultItems(payload)
	case "purchase_order":
		setIfMissing(payload, "purchase_order_number", "PO-"+time.Now().UTC().Format("20060102150405"))
		setIfMissing(payload, "vendor_name", "Preferred Vendor")
		setIfMissing(payload, "ship_to_name", "InsureTech Warehouse")
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

func mapTemplateToStorageType(docType documentv1.DocumentType) storageentityv1.FileType {
	if docType == documentv1.DocumentType_DOCUMENT_TYPE_INVOICE {
		return storageentityv1.FileType_FILE_TYPE_INVOICE
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
