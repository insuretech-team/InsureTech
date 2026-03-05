# DocGen Microservice - Complete Analysis

## 1. EXACT PACKAGE IMPORT PATHS

### Top-level server.go
```go
"database/sql"
"github.com/jmoiron/sqlx"
"github.com/newage-saint/insuretech/backend/inscore/microservices/docgen/internal/grpc"
"github.com/newage-saint/insuretech/backend/inscore/microservices/docgen/internal/repository"
"github.com/newage-saint/insuretech/backend/inscore/microservices/docgen/internal/service"
"github.com/newage-saint/insuretech/gen/go/insuretech/document/services/v1"
"github.com/newage-saint/insuretech/gen/go/insuretech/storage/service/v1"
"google.golang.org/grpc"
```

### document_service.go
```go
"bytes"
"context"
"encoding/base64"
"encoding/json"
"errors"
"fmt"
"html"
"html/template"
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
"github.com/newage-saint/insuretech/backend/inscore/microservices/docgen/internal/repository"
"github.com/newage-saint/insuretech/gen/go/insuretech/document/entity/v1"
"github.com/newage-saint/insuretech/gen/go/insuretech/storage/entity/v1"
"github.com/newage-saint/insuretech/gen/go/insuretech/storage/service/v1"
"github.com/newage-saint/insuretech/ops/config"
"google.golang.org/grpc"
"google.golang.org/protobuf/types/known/structpb"
"google.golang.org/protobuf/types/known/timestamppb"
```

### generation_repository.go
```go
"context"
"database/sql"
"errors"
"fmt"
"strconv"
"strings"
"time"

"github.com/jmoiron/sqlx"
"github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
"github.com/newage-saint/insuretech/gen/go/insuretech/document/entity/v1"
"google.golang.org/protobuf/types/known/timestamppb"
```

### template_repository.go
```go
"context"
"database/sql"
"encoding/json"
"errors"
"fmt"
"strconv"
"strings"

"github.com/jmoiron/sqlx"
"github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
"github.com/newage-saint/insuretech/gen/go/insuretech/document/entity/v1"
"google.golang.org/protobuf/types/known/timestamppb"
```

### document_handler.go (grpc)
```go
"context"
"errors"
"strings"

"github.com/newage-saint/insuretech/backend/inscore/microservices/docgen/internal/service"
"github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
"github.com/newage-saint/insuretech/gen/go/insuretech/document/services/v1"
"google.golang.org/grpc/metadata"
```

### cmd/docgen/main.go
```go
"context"
"fmt"
"net"
"os"
"os/signal"
"strings"
"syscall"
"time"

"github.com/newage-saint/insuretech/backend/inscore/db"
"github.com/newage-saint/insuretech/backend/inscore/microservices/docgen"
"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
"github.com/newage-saint/insuretech/gen/go/insuretech/document/services/v1"
"github.com/newage-saint/insuretech/ops/config"
"github.com/newage-saint/insuretech/ops/env"
"go.uber.org/zap"
"google.golang.org/grpc"
"google.golang.org/grpc/credentials/insecure"
"google.golang.org/grpc/health"
"google.golang.org/grpc/health/grpc_health_v1"
```

---

## 2. ALL EXISTING RPC IMPLEMENTATIONS AND THEIR SIGNATURES

### File: internal/grpc/document_handler.go
Package: `server` (named as `package server`)

#### RPC Methods (8 total):

1. **GenerateDocument**
   ```go
   func (h *DocumentHandler) GenerateDocument(
       ctx context.Context, 
       req *documentservicev1.GenerateDocumentRequest
   ) (*documentservicev1.GenerateDocumentResponse, error)
   ```
   - Accepts: templateID, entityType, entityID, data (Struct), includeQRCode
   - Returns: documentID, fileURL, message
   - Calls: `h.docService.GenerateDocument()`

2. **GetDocument**
   ```go
   func (h *DocumentHandler) GetDocument(
       ctx context.Context, 
       req *documentservicev1.GetDocumentRequest
   ) (*documentservicev1.GetDocumentResponse, error)
   ```
   - Accepts: documentID
   - Returns: Document proto

3. **ListDocuments**
   ```go
   func (h *DocumentHandler) ListDocuments(
       ctx context.Context, 
       req *documentservicev1.ListDocumentsRequest
   ) (*documentservicev1.ListDocumentsResponse, error)
   ```
   - Accepts: entityType, entityID, status, page, pageSize
   - Returns: documents array, totalCount

4. **DownloadDocument**
   ```go
   func (h *DocumentHandler) DownloadDocument(
       ctx context.Context, 
       req *documentservicev1.DownloadDocumentRequest
   ) (*documentservicev1.DownloadDocumentResponse, error)
   ```
   - Accepts: documentID
   - Returns: content (bytes), contentType, filename

5. **DeleteDocument**
   ```go
   func (h *DocumentHandler) DeleteDocument(
       ctx context.Context, 
       req *documentservicev1.DeleteDocumentRequest
   ) (*documentservicev1.DeleteDocumentResponse, error)
   ```
   - Accepts: documentID
   - Returns: message

6. **CreateDocumentTemplate**
   ```go
   func (h *DocumentHandler) CreateDocumentTemplate(
       ctx context.Context, 
       req *documentservicev1.CreateDocumentTemplateRequest
   ) (*documentservicev1.CreateDocumentTemplateResponse, error)
   ```
   - Accepts: name, type, description, templateContent, outputFormat, variables
   - Returns: templateID, message

7. **GetDocumentTemplate**
   ```go
   func (h *DocumentHandler) GetDocumentTemplate(
       ctx context.Context, 
       req *documentservicev1.GetDocumentTemplateRequest
   ) (*documentservicev1.GetDocumentTemplateResponse, error)
   ```
   - Accepts: templateID
   - Returns: template proto

8. **ListDocumentTemplates**
   ```go
   func (h *DocumentHandler) ListDocumentTemplates(
       ctx context.Context, 
       req *documentservicev1.ListDocumentTemplatesRequest
   ) (*documentservicev1.ListDocumentTemplatesResponse, error)
   ```
   - Accepts: type, activeOnly, pageSize, pageToken
   - Returns: templates array, nextPageToken, totalCount

9. **UpdateDocumentTemplate**
   ```go
   func (h *DocumentHandler) UpdateDocumentTemplate(
       ctx context.Context, 
       req *documentservicev1.UpdateDocumentTemplateRequest
   ) (*documentservicev1.UpdateDocumentTemplateResponse, error)
   ```
   - Accepts: templateID, template proto
   - Returns: message

10. **DeactivateDocumentTemplate**
    ```go
    func (h *DocumentHandler) DeactivateDocumentTemplate(
        ctx context.Context, 
        req *documentservicev1.DeactivateDocumentTemplateRequest
    ) (*documentservicev1.DeactivateDocumentTemplateResponse, error)
    ```
    - Accepts: templateID
    - Returns: message

11. **DeleteDocumentTemplate**
    ```go
    func (h *DocumentHandler) DeleteDocumentTemplate(
        ctx context.Context, 
        req *documentservicev1.DeleteDocumentTemplateRequest
    ) (*documentservicev1.DeleteDocumentTemplateResponse, error)
    ```
    - Accepts: templateID
    - Returns: message

#### Helper Functions in document_handler.go:
```go
func actorFromContext(ctx context.Context, fallback string) string
func tenantFromContext(ctx context.Context) string
func mapErr(err error, notFoundMessage string) *commonv1.Error
```

---

## 3. PDF GENERATION CODE

### Library: **gofpdf**
Import: `"github.com/jung-kurt/gofpdf"`

### Function: renderPDF
**Location**: `internal/service/document_service.go` lines 469-483

```go
func renderPDF(renderedHTML string) ([]byte, error) {
	// Strip HTML tags and unescape HTML entities
	plain := strings.TrimSpace(html.UnescapeString(htmlTagPattern.ReplaceAllString(renderedHTML, " ")))
	plain = regexp.MustCompile(`\s+`).ReplaceAllString(plain, " ")
	
	// Create PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(12, 12, 12)
	pdf.AddPage()
	pdf.SetFont("Arial", "", 11)
	pdf.MultiCell(0, 6, plain, "", "L", false)

	// Render to bytes buffer
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to render pdf: %w", err)
	}
	return buf.Bytes(), nil
}
```

### PDF Generation Integration:
- Called by: `buildOutput()` function (line 452-467)
- Output format check: `documentv1.OutputFormat_OUTPUT_FORMAT_PDF`
- Returns: PDF bytes, content-type `application/pdf`, extension `.pdf`

### Related Functions:
```go
func buildOutput(renderedHTML string, format documentv1.OutputFormat) ([]byte, string, string, error)
// Supports: HTML, PDF, DOCX (unsupported), and unknown formats
// Returns: fileContent ([]byte), contentType (string), fileExtension (string), error
```

---

## 4. KAFKA PUBLISHING CODE

### **NO KAFKA PUBLISHING FOUND**

After thorough examination of all files in the docgen microservice:
- **document_service.go**: No Kafka imports or publishing
- **generation_repository.go**: No Kafka imports or publishing
- **template_repository.go**: No Kafka imports or publishing
- **document_handler.go**: No Kafka imports or publishing
- **server.go**: No Kafka imports or publishing
- **cmd/docgen/main.go**: No Kafka imports or publishing

**Conclusion**: The DocGen microservice does NOT currently implement Kafka publishing. It is purely:
1. A gRPC service
2. Integrates with Storage service via gRPC
3. Uses PostgreSQL database
4. Does NOT publish events to Kafka

---

## 5. SERVICE CONFIG STRUCT

### DocumentService Struct
**Location**: `internal/service/document_service.go` lines 49-55

```go
type DocumentService struct {
	templateRepo    *repository.DocumentTemplateRepository
	generationRepo  *repository.DocumentGenerationRepository
	storageClient   StorageClient
	templateDirPath string
}
```

### Constructor: NewDocumentService
```go
func NewDocumentService(
	templateRepo *repository.DocumentTemplateRepository,
	generationRepo *repository.DocumentGenerationRepository,
	storageClient StorageClient,
) (*DocumentService, error)
```

**Initialization Logic**:
- Resolves template directory path using `config.ResolvePath(filepath.Join("backend", "inscore", "templates"))`
- Bootstraps default templates (invoice, purchase_order, policy_document)
- Returns error if template bootstrap fails

### StorageClient Interface
**Location**: `internal/service/document_service.go` lines 43-47

```go
type StorageClient interface {
	UploadFile(ctx context.Context, in *storageservicev1.UploadFileRequest, opts ...grpc.CallOption) (*storageservicev1.UploadFileResponse, error)
	DeleteFile(ctx context.Context, in *storageservicev1.DeleteFileRequest, opts ...grpc.CallOption) (*storageservicev1.DeleteFileResponse, error)
}
```

### DocumentHandler Struct
**Location**: `internal/grpc/document_handler.go` lines 14-17

```go
type DocumentHandler struct {
	documentservicev1.UnimplementedDocumentServiceServer
	docService *service.DocumentService
}

func NewDocumentHandler(docService *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{docService: docService}
}
```

---

## 6. STORAGE SERVICE CLIENT CALLS

### Storage Service Integration

The docgen service makes 2 types of storage service calls via gRPC:

#### 1. UploadFile
**Location**: `internal/service/document_service.go` lines 129-148

```go
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
```

**Called when**:
- Document is generated AND
- storageClient is not nil AND
- tenantID is not empty

**Maps document type to storage file type**:
```go
func mapTemplateToStorageType(docType documentv1.DocumentType) storageentityv1.FileType {
	if docType == documentv1.DocumentType_DOCUMENT_TYPE_INVOICE {
		return storageentityv1.FileType_FILE_TYPE_INVOICE
	}
	return storageentityv1.FileType_FILE_TYPE_DOCUMENT
}
```

**Response handling**:
```go
if upErr == nil && uploaded.GetFile() != nil {
	storageFileID = uploaded.File.FileId
	if uploaded.File.CdnUrl != "" {
		fileURL = uploaded.File.CdnUrl
	} else {
		fileURL = uploaded.File.Url
	}
}
```

#### 2. DeleteFile
**Location**: `internal/service/document_service.go` lines 254-262

```go
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
```

**Called when**:
- Document is deleted AND
- storageClient is not nil AND
- tenantID is not empty AND
- Document has stored fileID metadata

---

## 7. DATABASE SCHEMA DETAILS

### Tables Used

#### storage_schema.document_generations
**Columns**: generation_id, document_template_id, entity_type, entity_id, data, status, file_url, file_size_bytes, qr_code_data, generated_by, generated_at, audit_info, created_at, updated_at

#### storage_schema.document_templates
**Columns**: template_id, name, type, description, template_content, output_format, variables, version, is_active, audit_info, created_at, updated_at

### Repository Interfaces

#### DocumentGenerationRepository
```go
type DocumentGenerationRepository struct {
	db *sqlx.DB
}

func (r *DocumentGenerationRepository) Create(ctx context.Context, doc *documentv1.DocumentGeneration) (*documentv1.DocumentGeneration, error)
func (r *DocumentGenerationRepository) GetByID(ctx context.Context, documentID string) (*documentv1.DocumentGeneration, error)
func (r *DocumentGenerationRepository) ListByEntity(ctx context.Context, entityType, entityID string, status *documentv1.GenerationStatus, limit, offset int) ([]*documentv1.DocumentGeneration, int, error)
func (r *DocumentGenerationRepository) Delete(ctx context.Context, documentID string) error
```

#### DocumentTemplateRepository
```go
type DocumentTemplateRepository struct {
	db *sqlx.DB
}

func (r *DocumentTemplateRepository) Create(ctx context.Context, tpl *documentv1.DocumentTemplate) (*documentv1.DocumentTemplate, error)
func (r *DocumentTemplateRepository) UpsertByName(ctx context.Context, tpl *documentv1.DocumentTemplate) (*documentv1.DocumentTemplate, error)
func (r *DocumentTemplateRepository) GetByID(ctx context.Context, templateID string) (*documentv1.DocumentTemplate, error)
func (r *DocumentTemplateRepository) GetByName(ctx context.Context, name string) (*documentv1.DocumentTemplate, error)
func (r *DocumentTemplateRepository) List(ctx context.Context, docType *documentv1.DocumentType, activeOnly bool, limit, offset int) ([]*documentv1.DocumentTemplate, int, error)
func (r *DocumentTemplateRepository) Update(ctx context.Context, templateID string, tpl *documentv1.DocumentTemplate) error
func (r *DocumentTemplateRepository) Deactivate(ctx context.Context, templateID string) error
func (r *DocumentTemplateRepository) Delete(ctx context.Context, templateID string) error
```

---

## 8. QR CODE GENERATION

### Library: **barcode** and **barcode/qr**
Imports:
```go
"github.com/boombuler/barcode"
"github.com/boombuler/barcode/qr"
"image"
"image/png"
```

### Function: buildQRCodeDataURI
**Location**: `internal/service/document_service.go` lines 506-520

```go
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
```

### Integration:
- Called in `GenerateDocument()` at line 110
- QR data format: `"doc:{generationID}|entity:{entityType}|id:{entityID}"`
- Resolution: 256x256 pixels
- Error-tolerant: if QR generation fails, document still generates without QR
- Returns: base64-encoded data URI for embedding in HTML/PDF

---

## 9. TEMPLATE RENDERING

### HTML Template Engine
- Library: `"html/template"` (Go standard library)
- Regex patterns for variable extraction: `\{\{\s*\.([a-zA-Z0-9_]+)\s*\}\}`

### Function: renderTemplate
```go
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
```

### Function: extractTemplateVariables
```go
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
```

### Business Defaults Applied
Function: `applyBusinessDefaults()` applies context-specific defaults:

**Global defaults**:
- company_name: "InsureTech"
- company_address: "Dhaka, Bangladesh"
- issue_date: current date (UTC)
- terms: "Standard terms and conditions apply."

**Invoice template defaults**:
- invoice_number: "INV-{timestamp}"
- due_date: 7 days from now
- customer_name: "Policy Holder"
- subtotal, tax, total: "0.00"
- items: default array with one item

**Purchase Order defaults**:
- purchase_order_number: "PO-{timestamp}"
- vendor_name: "Preferred Vendor"
- ship_to_name: "InsureTech Warehouse"
- subtotal, shipping_cost, tax, total: "0.00"
- items: default array

**Policy Document defaults**:
- policy_title: "Insurance Policy Certificate"
- policy_number: "POL-{timestamp}"
- policy_holder_name: "Policy Holder"
- start_date: current date
- end_date: 1 year from now
- coverage_amount, premium_amount: "0.00"
- benefits: ["Coverage details will be provided by insurer"]

---

## 10. ERROR HANDLING

### Custom Error Variables
```go
var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrUnsupportedOutput = errors.New("unsupported output format")
	ErrTemplateNotFound  = repository.ErrTemplateNotFound
	ErrDocumentNotFound  = repository.ErrDocumentNotFound
)
```

### Error Mapping in gRPC Handler
```go
func mapErr(err error, notFoundMessage string) *commonv1.Error {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		return &commonv1.Error{Code: "INVALID_ARGUMENT", Message: err.Error()}
	case errors.Is(err, service.ErrTemplateNotFound), errors.Is(err, service.ErrDocumentNotFound):
		return &commonv1.Error{Code: "NOT_FOUND", Message: notFoundMessage}
	case errors.Is(err, service.ErrUnsupportedOutput):
		return &commonv1.Error{Code: "UNIMPLEMENTED", Message: err.Error()}
	default:
		return &commonv1.Error{Code: "INTERNAL", Message: err.Error()}
	}
}
```

---

## 11. ENTRY POINT CONFIGURATION (cmd/docgen/main.go)

### Database Initialization
- Uses `github.com/newage-saint/insuretech/backend/inscore/db` manager
- Loads `database.yaml` config via `config.ResolveConfigPath()`
- Supports failover with primary and backup databases

### Storage Service Connection
- Environment variable: `STORAGE_GRPC_ADDR`
- Optional: if not set, service runs in "inline-only mode"
- Connection timeout: 5 seconds
- Transport: insecure credentials (no TLS)
- Graceful degradation if storage fails

### gRPC Server Configuration
- Port environment variable: `DOCGEN_GRPC_PORT` (default: `50280`)
- Registers `DocumentServiceServer`
- Includes gRPC health check service
- Graceful shutdown on SIGTERM/SIGINT

---

## 12. FILE STRUCTURE SUMMARY

```
E:\Projects\InsureTech\backend\inscore\microservices\docgen\
├── doc.go                                    (package doc)
├── server.go                                 (NewDocumentServer factory)
├── internal/
│   ├── grpc/
│   │   └── document_handler.go              (11 RPC implementations)
│   ├── repository/
│   │   ├── generation_repository.go         (CRUD for document_generations table)
│   │   └── template_repository.go           (CRUD for document_templates table)
│   └── service/
│       └── document_service.go              (business logic, PDF/QR/template rendering)
└── [no cmd/server/main.go - located at ../cmd/docgen/main.go]
```

---

## 13. KEY OBSERVATIONS

1. **No Kafka Integration**: DocGen is purely gRPC-based, no event publishing
2. **Optional Storage Service**: Works with or without storage integration
3. **PDF Generation**: Uses gofpdf with simple text extraction (no complex HTML to PDF conversion)
4. **QR Code Support**: Embedded as base64 data URIs in templates
5. **Template Variables**: Automatically extracted via regex and validated
6. **Multi-tenant**: TenantID required for storage integration
7. **Database**: PostgreSQL via sqlx with audit info tracking
8. **Bootstrap Templates**: Three default templates loaded on service startup
9. **Output Formats**: HTML (default) and PDF supported; DOCX marked as unsupported

