# DocGen & Media Service Enhanced Implementation Plan V3
## Based on Comprehensive Codebase Analysis

**Date:** February 27, 2026  
**Version:** 3.0  
**Status:** SECOND-PASS VERIFIED: IMPLEMENTED BASELINE WITH HARDENING REMAINING

---

## Executive Summary

Based on thorough analysis of the **existing codebase**, this enhanced plan focuses on **completing and productionizing** the DocGen and Media services. Key findings:

### ✅ **ALREADY COMPLETE:**
1. **Storage Service** - Production-ready S3 integration
2. **All Proto Definitions** - Document & Media services fully defined
3. **DocGen Runtime** - Template CRUD, generation flow, async worker wiring, and Gotenberg renderer integration already exist
4. **Media Runtime** - CRUD, processing job tracking, worker pool, image processing hooks, and virus scan path already exist
5. **Database Schemas** - Proto-defined with auto-migration
6. **AuthZ Integration** - JWT middleware & Casbin enforcer ready
7. **Kafka Publishers** - DocGen and Media both already publish baseline events
8. **Gateway Patterns** - Validation middleware established

### 🎯 **ENHANCEMENT FOCUS:**
1. **DocGen Hardening** - Richer templates, retry strategy, and stronger E2E/load validation
2. **Media Completion** - Replace OCR stub with real Tesseract-backed extraction; tighten quality validation
3. **Eventing Completion** - Verify topic coverage, tenant context propagation, and consumer-side handling
4. **Async Job Reliability** - Robust queue management, retries, and failure visibility
5. **Integration Testing** - End-to-end workflow validation
6. **Performance Optimization** - Load testing and scaling

> Second-pass correction: the codebase already contains baseline Gotenberg wiring, Kafka publishers, media processing workers, and ClamAV integration points. The remaining phases below should be read as hardening, rollout, and gap-closure work, not as blank-slate implementation.

---

## 1. Current State Analysis

### 1.1 DocGen Service Status

**✅ COMPLETE:**
- Service structure (`server.go`, `document_service.go`)
- Template CRUD operations
- HTML rendering + Gotenberg PDF generation path with fallback renderer in service code
- QR code generation
- Storage service integration
- Repository layer
- Kafka event publishing for template and generation activity
- Async generation worker wiring in server startup

**⚠️ NEEDS ENHANCEMENT:**
- **Templates**: Limited HTML → Rich templates with stronger branding/layout controls
- **Performance**: baseline async support exists, but queue durability/retry/load behavior needs hardening
- **Events**: publishers exist, but downstream consumption and tenant-context completeness need verification
- **Testing**: expand from current unit/integration coverage to stronger E2E and load tests

### 1.2 Media Service Status

**✅ COMPLETE:**
- Service structure (`server.go`, `media_service.go`)
- Media CRUD operations
- Basic validation (file size)
- Processing job tracking
- Storage service integration
- Repository layer
- Worker pool startup and processing job execution
- Image processing and virus scanning integration points
- Kafka publishing for processing lifecycle events

**⚠️ NEEDS ENHANCEMENT:**
- **OCR**: worker pipeline exists, but `ocr_processor.go` is still a stub until Tesseract/gosseract is wired for real extraction
- **Validation**: Basic size check → DPI, format, quality validation
- **Async Jobs**: pipeline exists, but retry/priority/observability behavior still needs hardening
- **CDN Integration**: Basic URL resolution → Full CDN support

### 1.3 Infrastructure Status

**✅ COMPLETE:**
- **Storage Service**: S3 upload/download, pre-signed URLs
- **Database**: PostgreSQL with proto-based migrations
- **AuthZ**: JWT validation, Casbin RBAC
- **Gateway**: Request validation, rate limiting
- **Proto Definitions**: Complete for both services
- **Code-level Gotenberg wiring**: already present in DocGen service startup
- **Code-level ClamAV wiring**: already present in Media service processor/worker path
- **Code-level Kafka publishing**: already present for both services

**⚠️ NEEDS SETUP:**
- **Gotenberg**: runtime deployment and health monitoring in target environments
- **ClamAV**: runtime deployment and signature updates
- **Tesseract**: real OCR installation/configuration
- **ImageMagick/libvips**: stronger optimization stack if pure-Go processing is insufficient
- **Kafka**: topic provisioning, DLQ policy, and consumer rollout

---

## 2. Enhanced Architecture Design

### 2.1 Production-Ready DocGen Service

**Core Enhancements:**
1. **Gotenberg Integration** - Docker-based PDF generation
2. **Template Management** - Rich HTML templates with variables
3. **Async Processing** - Job queue for batch generation
4. **Event Publishing** - Kafka events for all operations
5. **Performance Caching** - Template and result caching
6. **Multi-language Support** - Bengali/English templates

**Technology Stack:**
- **PDF Engine**: Gotenberg (Docker)
- **Template Engine**: Go templates + HTML/CSS
- **Async Processing**: Kafka + internal job queue
- **Caching**: Redis (optional)
- **Storage**: S3 via Storage service
- **Database**: PostgreSQL (proto-generated)

### 2.2 Production-Ready Media Service

**Core Enhancements:**
1. **Image Processing** - ImageMagick/libvips integration
2. **OCR Processing** - Tesseract with Bengali support
3. **Virus Scanning** - ClamAV integration
4. **Quality Validation** - DPI, blur, format checks
5. **Async Pipeline** - Multi-stage processing queue
6. **CDN Integration** - Automatic CDN URL generation

**Technology Stack:**
- **Image Processing**: libvips (Go bindings)
- **OCR**: Tesseract (Bengali/English)
- **Virus Scan**: ClamAV (Docker)
- **Async Processing**: Kafka + worker pool
- **CDN**: CloudFront/S3 integration
- **Storage**: S3 via Storage service
- **Database**: PostgreSQL (proto-generated)

---

## 3. Enhanced Implementation Phases

### Phase 1: Infrastructure Setup (Day 1-2)
**Goal:** Production dependencies ready

**Tasks:**
1. **Gotenberg Setup** - Docker container, health checks
2. **ClamAV Setup** - Docker container, virus definitions
3. **Tesseract Setup** - Docker with Bengali language pack
4. **libvips Setup** - Go bindings installation
5. **Kafka Topics** - Create required topics
6. **Monitoring** - Basic health monitoring

**Deliverables:**
- All external services running in Docker
- Health check endpoints working
- Kafka topics created
- Monitoring dashboards

### Phase 2: DocGen Gotenberg Integration (Day 3-4)
**Goal:** Production PDF generation

**Tasks:**
1. **Gotenberg Client** - Robust HTTP client with retries
2. **Template Rendering** - Enhanced Go template engine
3. **Error Handling** - Comprehensive error recovery
4. **Performance** - Connection pooling, timeouts
5. **Testing** - Integration tests with Gotenberg

**Deliverables:**
- Gotenberg client library
- 5 production templates
- Performance benchmarks
- Integration test suite

### Phase 3: Media Processing Pipeline (Day 5-7)
**Goal:** Complete media processing

**Tasks:**
1. **Image Processing** - libvips integration for optimization
2. **OCR Processing** - Tesseract integration
3. **Virus Scanning** - ClamAV integration
4. **Quality Validation** - DPI, format, blur detection
5. **Worker Pool** - Concurrent processing workers

**Deliverables:**
- Image optimization working
- OCR extraction working
- Virus scanning working
- Worker pool implementation
- Processing pipeline tests

### Phase 4: Kafka Event Integration (Day 8-9)
**Goal:** Full event-driven architecture

**Tasks:**
1. **Event Publishers** - Implement all defined events
2. **Event Consumers** - Notification service integration
3. **Error Handling** - Dead letter queues, retries
4. **Monitoring** - Event metrics, lag monitoring
5. **Testing** - Event publishing/consuming tests

**Deliverables:**
- All proto events published
- Notification service integration
- Event monitoring dashboard
- Comprehensive event tests

### Phase 5: Async Job Management (Day 10-11)
**Goal:** Robust async processing

**Tasks:**
1. **Job Queue** - Priority-based job scheduling
2. **Retry Logic** - Exponential backoff, max retries
3. **Status Tracking** - Real-time job status updates
4. **Cancellation** - Job cancellation support
5. **Scaling** - Horizontal scaling support

**Deliverables:**
- Job queue implementation
- Retry mechanism
- Status tracking API
- Cancellation support
- Scaling documentation

### Phase 6: Integration & Load Testing (Day 12-13)
**Goal:** Production readiness validation

**Tasks:**
1. **End-to-End Tests** - Complete workflow testing
2. **Load Testing** - 1000 concurrent PDF generations
3. **Load Testing** - 500 concurrent media uploads
4. **Performance Optimization** - Identify bottlenecks
5. **Failure Testing** - Service failure scenarios

**Deliverables:**
- E2E test suite
- Load test results
- Performance benchmarks
- Failure recovery procedures
- Production readiness report

### Phase 7: Production Deployment (Day 14-15)
**Goal:** Services in production

**Tasks:**
1. **Deployment** - ECS/Fargate deployment
2. **Monitoring** - CloudWatch metrics, alerts
3. **Documentation** - Runbooks, API documentation
4. **Rollout** - Gradual traffic rollout
5. **Support** - On-call procedures, escalation

**Deliverables:**
- Services in production
- Monitoring dashboards
- Complete documentation
- Rollback procedures
- Support runbooks

---

## 4. Enhanced Service Specifications

### 4.1 DocGen Service - Enhanced RPCs

**Existing RPCs (✅ Complete):**
- GenerateDocument
- GetDocument
- ListDocuments
- DownloadDocument
- DeleteDocument
- CreateDocumentTemplate
- GetDocumentTemplate
- ListDocumentTemplates
- UpdateDocumentTemplate
- DeactivateDocumentTemplate
- DeleteDocumentTemplate

**Enhanced Features:**
1. **Async Generation** - `GenerateDocumentAsync` with webhook
2. **Batch Generation** - `BatchGenerateDocuments`
3. **Template Preview** - `PreviewTemplate` (HTML preview)
4. **Template Variables** - `GetTemplateVariables` (schema)
5. **Document Signing** - `SignDocument` (digital signatures)
6. **Version Management** - `ListDocumentVersions`

### 4.2 Media Service - Enhanced RPCs

**Existing RPCs (✅ Complete):**
- UploadMedia
- GetMedia
- ListMedia
- DownloadMedia
- DownloadOptimized
- DownloadThumbnail
- DeleteMedia
- ValidateMedia
- RequestProcessing
- GetProcessingJob
- ListProcessingJobs

**Enhanced Features:**
1. **Streaming Upload** - `UploadMediaStream` (large files)
2. **Batch Upload** - `UploadMediaBatch`
3. **Smart Cropping** - `CropMedia` (face detection)
4. **Format Conversion** - `ConvertMedia` (format changes)
5. **Quality Analysis** - `AnalyzeMediaQuality`
6. **OCR Languages** - `SetOCRLanguage` (Bengali/English)

---

## 5. Enhanced Database Schema

### 5.1 DocGen Schema Enhancements

**Existing Tables (✅ Auto-generated):**
- `storage_schema.document_templates`
- `storage_schema.document_generations`

**Enhanced Indexes:**
```sql
-- Performance optimizations
CREATE INDEX CONCURRENTLY idx_document_generations_composite 
ON storage_schema.document_generations(entity_type, entity_id, status, generated_at DESC);

CREATE INDEX CONCURRENTLY idx_document_templates_active 
ON storage_schema.document_templates(type, is_active, version DESC);

-- Partitioning for large tables
CREATE TABLE storage_schema.document_generations_y2026m03 
PARTITION OF storage_schema.document_generations
FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');
```

### 5.2 Media Schema Enhancements

**Existing Tables (✅ Auto-generated):**
- `media_schema.media_files`
- `media_schema.processing_jobs`

**Enhanced Indexes:**
```sql
-- Performance optimizations
CREATE INDEX CONCURRENTLY idx_media_files_composite 
ON media_schema.media_files(entity_type, entity_id, media_type, validation_status);

CREATE INDEX CONCURRENTLY idx_processing_jobs_composite 
ON media_schema.processing_jobs(media_id, status, priority DESC, created_at);

-- GIN indexes for JSON fields
CREATE INDEX CONCURRENTLY idx_media_files_validation_errors_gin 
ON media_schema.media_files USING GIN (validation_errors);

CREATE INDEX CONCURRENTLY idx_processing_jobs_result_data_gin 
ON media_schema.processing_jobs USING GIN (result_data);
```

---

## 6. Enhanced Kafka Integration

### 6.1 Document Events (Enhanced)

**Existing Events (✅ Defined):**
- DocumentGenerationRequestedEvent
- DocumentGeneratedEvent
- DocumentGenerationFailedEvent

**Enhanced Events:**
1. **DocumentTemplateCreatedEvent** - Template management
2. **DocumentTemplateUpdatedEvent** - Template changes
3. **DocumentDownloadedEvent** - Usage tracking
4. **DocumentSignedEvent** - Digital signature events
5. **DocumentVersionCreatedEvent** - Version management

**Topic Strategy:**
- `insuretech.document.v1.generated` (existing)
- `insuretech.document.v1.template.*` (new)
- `insuretech.document.v1.download.*` (new)
- `insuretech.document.v1.signature.*` (new)

### 6.2 Media Events (Enhanced)

**Existing Events (✅ Defined):**
- MediaFileUploadedEvent
- MediaValidationCompletedEvent
- MediaVirusScanCompletedEvent
- MediaProcessingJobCreatedEvent
- MediaProcessingCompletedEvent
- MediaProcessingFailedEvent
- MediaOCRCompletedEvent

**Enhanced Events:**
1. **MediaQualityAnalyzedEvent** - Quality metrics
2. **MediaFormatConvertedEvent** - Format conversion
3. **MediaCDNPublishedEvent** - CDN distribution
4. **MediaBatchProcessedEvent** - Batch completion
5. **MediaDeletedEvent** - Deletion tracking

**Topic Strategy:**
- `insuretech.media.v1.uploaded` (existing)
- `insuretech.media.v1.processed.*` (enhanced)
- `insuretech.media.v1.quality.*` (new)
- `insuretech.media.v1.cdn.*` (new)

---

## 7. Enhanced Security & Compliance

### 7.1 Authentication & Authorization

**Existing (✅ Complete):**
- JWT validation middleware
- Casbin RBAC policies
- Public/private method configuration

**Enhancements:**
1. **Tenant Isolation** - Row-level security (RLS)
2. **File Access Control** - Signed URL expiration
3. **Audit Logging** - Comprehensive audit trails
4. **Compliance** - GDPR, data retention policies
5. **Encryption** - At-rest and in-transit encryption

### 7.2 Data Protection

**Document Service:**
- PDF encryption support
- Digital signature validation
- Watermarking for sensitive documents
- Redaction capabilities
- Audit trail for document access

**Media Service:**
- EXIF data stripping
- Face detection blurring (privacy)
- Virus scanning mandatory
- Content moderation integration
- Retention policy enforcement

---

## 8. Enhanced Monitoring & Observability

### 8.1 Metrics Collection

**DocGen Metrics:**
- `docgen_generation_duration_seconds` (histogram)
- `docgen_template_render_errors_total` (counter)
- `docgen_async_jobs_queue_size` (gauge)
- `docgen_pdf_generation_success_rate` (gauge)
- `docgen_cache_hit_ratio` (gauge)

**Media Metrics:**
- `media_upload_duration_seconds` (histogram)
- `media_processing_duration_seconds{type}` (histogram)
- `media_ocr_accuracy_score` (gauge)
- `media_virus_detection_rate` (gauge)
- `media_cdn_cache_hit_ratio` (gauge)

### 8.2 Alerting Rules

**Critical Alerts (PagerDuty):**
- Service down > 1 minute
- PDF generation failure rate > 5%
- Virus detected in upload
- Storage quota > 90%
- Kafka consumer lag > 1000 messages

**Warning Alerts (Slack):**
- PDF generation latency p95 > 3 seconds
- Media upload latency p95 > 5 seconds
- OCR accuracy < 80%
- CDN cache hit ratio < 70%
- Database connection pool > 80%

---

## 9. Enhanced Testing Strategy

### 9.1 Unit Tests

**DocGen:**
- Template rendering engine
- Gotenberg client integration
- QR code generation
- Variable substitution
- Error handling scenarios

**Media:**
- Image optimization algorithms
- OCR text extraction
- Virus scanning integration
- Quality validation rules
- Format conversion logic

**Target:** 85% code coverage

### 9.2 Integration Tests

**Document Generation Flow:**
1. Template creation → validation
2. Data preparation → template rendering
3. PDF generation → storage upload
4. Event publishing → notification
5. Download → access validation

**Media Processing Flow:**
1. File upload → validation
2. Virus scan → optimization
3. OCR processing → text extraction
4. Thumbnail generation → CDN distribution
5. Event publishing → analytics

### 9.3 Load Tests

**Performance Targets:**
- PDF generation: 1000 concurrent requests
- Media upload: 500 concurrent requests
- OCR processing: 200 documents/minute
- Template rendering: 5000 templates/hour
- Event publishing: 10,000 events/second

**Success Criteria:**
- 95th percentile latency < 3 seconds
- Error rate < 1%
- No memory leaks (24-hour soak)
- Graceful degradation under load

---

## 10. Enhanced Deployment Strategy

### 10.1 Infrastructure Requirements

**DocGen Service:**
- **Compute**: 4 vCPU, 8GB RAM (autoscale 2-8 instances)
- **Gotenberg**: 2 vCPU, 4GB RAM per instance (2 instances)
- **Redis**: Cache layer (optional, for performance)
- **Database**: PostgreSQL (shared, read replicas)
- **Storage**: S3 via storage service

**Media Service:**
- **Compute**: 8 vCPU, 16GB RAM (autoscale 4-12 instances)
- **ClamAV**: 4 vCPU, 8GB RAM (2 instances)
- **Tesseract**: 2 vCPU, 4GB RAM (per worker)
- **libvips**: Bundled with service
- **CDN**: CloudFront distribution

### 10.2 Deployment Phases

**Phase 1: Canary Deployment (10% traffic)**
- Monitor error rates, latency
- Validate integration points
- Gather performance baselines

**Phase 2: Gradual Rollout (50% traffic)**
- Scale based on metrics
- Validate under production load
- Fix any issues discovered

**Phase 3: Full Deployment (100% traffic)**
- Complete traffic shift
- Final performance validation
- Documentation completion

**Phase 4: Post-Deployment (Monitoring)**
- 24/7 monitoring established
- Alerting tuned
- Runbooks verified
- Support handoff

---

## 11. Enhanced Success Metrics

### 11.1 Business Metrics

**Document Generation:**
- Policy issuance time reduction (target: 50%)
- Customer satisfaction (CSAT) improvement
- Operational cost per document (target: < $0.01)
- Template reuse rate (target: > 70%)
- Digital signature adoption rate

**Media Processing:**
- Claims processing time reduction (target: 40%)
- Document validation accuracy (target: > 95%)
- OCR extraction accuracy (target: > 85%)
- Storage cost optimization (target: 30% reduction)
- CDN cache efficiency (target: > 80%)

### 11.2 Technical Metrics

**Performance:**
- PDF generation latency p95: < 3 seconds
- Media upload latency p95: < 5 seconds
- Service availability: > 99.9%
- Error rate: < 0.5%
- API response time p99: < 1 second

**Reliability:**
- Mean time between failures (MTBF): > 30 days
- Mean time to recovery (MTTR): < 5 minutes
- Data consistency: 100%
- Backup success rate: 100%
- Disaster recovery RTO: < 1 hour

---

## 12. Enhanced Risk Mitigation

### 12.1 Technical Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Gotenberg performance | High | Medium | Load testing, horizontal scaling |
| OCR accuracy for Bengali | High | High | Multiple OCR engines, manual review |
| Virus scanning false positives | Medium | Medium | Whitelist, manual review queue |
| S3 cost overrun | Medium | Low | Lifecycle policies, monitoring |
| Kafka event loss | High | Low | Idempotent consumers, DLQ |

### 12.2 Operational Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Team skill gap | Medium | Low | Training, documentation |
| Integration complexity | High | Medium | Incremental rollout, testing |
| Compliance requirements | High | Medium | Early legal review |
| Third-party dependencies | Medium | Medium | Fallback mechanisms |
| Scalability limitations | High | Low | Load testing, auto-scaling |

---

## 13. Enhanced Timeline Summary

| Phase | Days | Focus | Deliverables |
|-------|------|-------|--------------|
| Infrastructure Setup | 1-2 | External services | Gotenberg, ClamAV, Tesseract ready |
| DocGen Enhancement | 3-4 | PDF generation | Gotenberg integration, templates |
| Media Processing | 5-7 | Advanced features | OCR, optimization, virus scan |
| Kafka Integration | 8-9 | Event-driven | Full event publishing/consuming |
| Async Job Management | 10-11 | Job queues | Robust async processing |
| Integration Testing | 12-13 | E2E validation | Load tests, performance optimization |
| Production Deployment | 14-15 | Go-live | Services in production, monitoring |

**Total:** 15 working days (3 weeks)
**Target Completion:** March 15, 2026

---

## 14. Enhanced Next Steps

### Immediate Actions (Week 1):

1. **Infrastructure Provisioning**
   - Setup Gotenberg Docker containers
   - Deploy ClamAV with latest definitions
   - Install Tesseract with Bengali support
   - Configure libvips Go bindings
   - Create Kafka topics

2. **Development Environment**
   - Update docker-compose for local development
   - Setup integration test environment
   - Configure monitoring locally
   - Create development documentation

3. **Team Preparation**
   - Review enhanced plan with team
   - Assign implementation tasks
   - Setup daily standups
   - Establish communication channels

### Week 1 Sprint Planning:

**Sprint Goal:** Infrastructure ready, DocGen Gotenberg integration

**User Stories:**
1. As a developer, I can generate PDFs using Gotenberg
2. As a system, I can process 100 PDFs concurrently
3. As a business, I have 5 production-ready templates
4. As an operator, I can monitor service health

**Acceptance Criteria:**
- PDF generation < 2 seconds
- Error rate < 1%
- Templates render correctly
- Health checks pass

---

## 15. Enhanced Critical Success Factors

1. **Leverage Existing Infrastructure** - Use established patterns
2. **Focus on Production Readiness** - Not just functionality
3. **Implement Comprehensive Testing** - Unit, integration, load
4. **Establish Robust Monitoring** - Metrics, alerts, dashboards
5. **Document Everything** - APIs, deployment, troubleshooting
6. **Plan for Scale** - From day one design for growth
7. **Security First** - Authentication, authorization, encryption
8. **Compliance Ready** - GDPR, data protection, audit trails

---

## Appendix A: Enhanced Template Examples

### A.1 Policy Document Template (Enhanced)

```html
<!DOCTYPE html>
<html lang="{{language}}">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{company_name}} - Policy Document</title>
    <style>
        /* Enhanced styling with responsive design */
        @page { margin: 2cm; }
        body { font-family: 'Noto Sans Bengali', 'Arial', sans-serif; line-height: 1.6; }
        .header { text-align: center; border-bottom: 2px solid #2c3e50; padding-bottom: 20px; }
        .policy-details { background: #f8f9fa; padding: 20px; border-radius: 5px; }
        .signature-area { margin-top: 50px; border-top: 1px dashed #ccc; padding-top: 20px; }
        .qr-code { text-align: center; margin: 20px 0; }
        @media print {
            .no-print { display: none; }
            body { font-size: 12pt; }
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>{{company_name}}</h1>
        <p>{{company_address}}</p>
        <p>Phone: {{company_phone}} | Email: {{company_email}}</p>
    </div>
    
    <div class="policy-details">
        <h2>Policy Certificate</h2>
        <table width="100%">
            <tr><td><strong>Policy Number:</strong></td><td>{{policy_number}}</td></tr>
            <tr><td><strong>Issue Date:</strong></td><td>{{issue_date}}</td></tr>
            <tr><td><strong>Policy Type:</strong></td><td>{{policy_type}}</td></tr>
            <tr><td><strong>Sum Assured:</strong></td><td>{{sum_assured}}</td></tr>
            <tr><td><strong>Premium:</strong></td><td>{{premium}}</td></tr>
            <tr><td><strong>Term:</strong></td><td>{{term}}</td></tr>
        </table>
    </div>
    
    <div class="insured-details">
        <h3>Insured Person Details</h3>
        <p><strong>Name:</strong> {{customer_name}}</p>
        <p><strong>Date of Birth:</strong> {{customer_dob}}</p>
        <p><strong>National ID:</strong> {{customer_nid}}</p>
        <p><strong>Address:</strong> {{customer_address}}</p>
    </div>
    
    <div class="coverage-details">
        <h3>Coverage Details</h3>
        <ul>
            {{#each benefits}}
            <li>{{this}}</li>
            {{/each}}
        </ul>
    </div>
    
    <div class="qr-code">
        <img src="{{qr_code_data_uri}}" alt="Verification QR Code" width="200" height="200">
        <p>Scan to verify authenticity</p>
    </div>
    
    <div class="signature-area">
        <table width="100%">
            <tr>
                <td width="50%">
                    <p><strong>Authorized Signatory</strong></p>
                    <p>_________________________</p>
                    <p>{{authorized_name}}</p>
                    <p>{{authorized_title}}</p>
                </td>
                <td width="50%">
                    <p><strong>Date:</strong> {{signature_date}}</p>
                    <p><strong>Document ID:</strong> {{document_id}}</p>
                    <p><strong>Version:</strong> {{version}}</p>
                </td>
            </tr>
        </table>
    </div>
    
    <div class="footer no-print">
        <p>Generated on: {{generation_timestamp}} | Document ID: {{document_id}}</p>
    </div>
</body>
</html>
```

### A.2 Template Variable Schema

```json
{
  "required_variables": [
    "policy_number",
    "customer_name", 
    "sum_assured",
    "premium",
    "term"
  ],
  "optional_variables": [
    "customer_dob",
    "customer_nid",
    "customer_address",
    "benefits",
    "qr_code_data_uri"
  ],
  "default_values": {
    "company_name": "InsureTech Bangladesh",
    "company_address": "Dhaka, Bangladesh",
    "language": "bn",
    "version": "1.0"
  },
  "validation_rules": {
    "policy_number": "regex:^[A-Z]{3}-\\d{4}-\\d{6}<file name=\"DOCGEN_MEDIA_ENHANCED_PLAN_V3.md\" language=\"markdown\" >\n<content>\n",
    "sum_assured": "numeric:min=1000",
    "premium": "numeric:min=0",
    "term": "regex:^\\d+\\s+(years|months)<file name=\"DOCGEN_MEDIA_ENHANCED_PLAN_V3.md\" language=\"markdown\" >\n<content>\n"
  }
}
```

---

## Appendix B: Enhanced Media Processing Rules

### B.1 Claims Document Validation (Enhanced)

```yaml
validation_rules:
  file_size:
    max: 10485760  # 10 MB
    min: 1024      # 1 KB
    
  dimensions:
    min_dpi: 300
    min_resolution: [800, 600]
    max_resolution: [4000, 4000]
    
  formats:
    allowed: ["application/pdf", "image/jpeg", "image/png", "image/tiff"]
    preferred: "application/pdf"
    
  quality:
    blur_threshold: 0.8
    brightness_range: [50, 200]
    contrast_min: 0.5
    
  content:
    ocr_required: true
    languages: ["ben", "eng"]
    min_text_length: 10
    confidence_threshold: 0.7
    
  security:
    virus_scan: required
    exif_strip: true
    max_files_per_entity: 20
    
processing_pipeline:
  - stage: validation
    timeout: 30s
    
  - stage: virus_scan
    timeout: 60s
    
  - stage: optimization
    actions: ["compress", "format_convert"]
    
  - stage: ocr
    languages: ["ben", "eng"]
    
  - stage: thumbnail
    sizes: ["256x256", "512x512"]
    
  - stage: cdn_distribution
    cache_ttl: 86400
```

### B.2 Image Optimization Settings

```yaml
optimization_profiles:
  high_quality:
    format: "webp"
    quality: 85
    compression: "lossy"
    resize: "max_width:1920"
    strip_metadata: true
    
  web_optimized:
    format: "webp"
    quality: 75
    compression: "lossy"
    resize: "max_width:1200"
    strip_metadata: true
    
  thumbnail:
    format: "jpeg"
    quality: 70
    compression: "lossy"
    resize: "width:256,height:256,crop:center"
    strip_metadata: true
    
  print_ready:
    format: "tiff"
    quality: 100
    compression: "lossless"
    dpi: 300
    color_space: "CMYK"
    
performance_settings:
  concurrent_workers: 4
  memory_limit: "512MB"
  timeout_per_stage: 300
  retry_attempts: 3
  retry_delay: "5s"
```

---

**END OF ENHANCED IMPLEMENTATION PLAN**

**Document Version:** 3.0  
**Last Updated:** 2026-02-27  
**Status:** APPROVED - Ready for Implementation

---

## Appendix C: Enhanced Authentication & Authorization Integration

### C.1 JWT Token Integration

**Current AuthZ Service Analysis:**
Based on the existing `authz/cmd/server/main.go`, the authentication pattern is:

```go
// JWT Interceptor with public methods
publicMethods := []string{
    "/insuretech.authz.services.v1.AuthZService/CheckPermission",
    "/insuretech.authz.services.v1.AuthZService/CheckAccess",
    "/insuretech.authz.services.v1.AuthZService/GetJWKS",
    "/grpc.health.v1.Health/Check",
    "/grpc.health.v1.Health/Watch",
}

jwtInterceptor := middleware.NewJWTInterceptor(publicKey, publicMethods)
```

**Enhanced Integration for DocGen & Media:**

```go
// DocGen Service - JWT protected methods
func (s *DocGenService) GenerateDocument(ctx context.Context, req *GenerateDocumentRequest) (*GenerateDocumentResponse, error) {
    // Extract JWT claims from context
    claims, err := middleware.ExtractClaims(ctx)
    if err != nil {
        return nil, status.Error(codes.Unauthenticated, "authentication required")
    }
    
    // Check permission via AuthZ service
    hasPermission, err := s.authzClient.CheckPermission(ctx, &authzv1.CheckPermissionRequest{
        Subject:   claims.Subject,
        Resource:  fmt.Sprintf("document:template:%s", req.TemplateId),
        Action:    "generate",
        TenantId:  req.TenantId,
    })
    
    if err != nil || !hasPermission.Allowed {
        return nil, status.Error(codes.PermissionDenied, "insufficient permissions")
    }
    
    // Proceed with generation
    // ...
}
```

### C.2 Casbin RBAC Policies

**Policy Definitions for Document Service:**

```csv
# Document Template Policies
p, role:admin, document:template:*, create, allow
p, role:admin, document:template:*, read, allow
p, role:admin, document:template:*, update, allow
p, role:admin, document:template:*, delete, allow

p, role:editor, document:template:*, create, allow
p, role:editor, document:template:*, read, allow
p, role:editor, document:template:*, update, allow

p, role:viewer, document:template:*, read, allow

# Document Generation Policies
p, role:admin, document:generation:*, generate, allow
p, role:admin, document:generation:*, read, allow
p, role:admin, document:generation:*, delete, allow

p, role:agent, document:generation:policy, generate, allow
p, role:agent, document:generation:receipt, generate, allow
p, role:agent, document:generation:*, read, allow

p, role:customer, document:generation:own, read, allow
```

**Policy Definitions for Media Service:**

```csv
# Media Upload Policies
p, role:admin, media:file:*, upload, allow
p, role:admin, media:file:*, read, allow
p, role:admin, media:file:*, delete, allow

p, role:agent, media:file:claim, upload, allow
p, role:agent, media:file:kyc, upload, allow
p, role:agent, media:file:*, read, allow

p, role:customer, media:file:own, upload, allow
p, role:customer, media:file:own, read, allow

# Media Processing Policies
p, role:admin, media:processing:*, request, allow
p, role:admin, media:processing:*, read, allow

p, role:processor, media:processing:optimization, request, allow
p, role:processor, media:processing:ocr, request, allow
```

### C.3 Tenant Isolation with Row-Level Security

**PostgreSQL RLS Policies:**

```sql
-- Media files RLS
ALTER TABLE media_schema.media_files ENABLE ROW LEVEL SECURITY;

CREATE POLICY media_files_tenant_isolation ON media_schema.media_files
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);

CREATE POLICY media_files_admin_access ON media_schema.media_files
    USING (current_setting('app.current_user_role') = 'admin');

-- Document templates RLS
ALTER TABLE storage_schema.document_templates ENABLE ROW LEVEL SECURITY;

CREATE POLICY document_templates_tenant_isolation ON storage_schema.document_templates
    USING (EXISTS (
        SELECT 1 FROM authn_schema.users u 
        WHERE u.user_id = current_setting('app.current_user_id')::uuid
        AND u.tenant_id = current_setting('app.current_tenant_id')::uuid
    ));
```

### C.4 Enhanced Service-to-Service Authentication

**gRPC Interceptor Chain:**

```go
// Enhanced gRPC server setup for DocGen/Media services
func NewServer(authzConn *grpc.ClientConn) *grpc.Server {
    // Create authz client
    authzClient := authzservicev1.NewAuthZServiceClient(authzConn)
    
    // Create JWT validator
    validator := middleware.NewJWTValidator(publicKey)
    
    // Create permission checker
    permissionChecker := middleware.NewPermissionChecker(authzClient)
    
    // Build interceptor chain
    chain := grpc.ChainUnaryInterceptor(
        validator.UnaryInterceptor(),
        permissionChecker.UnaryInterceptor(),
        tenant.EnforceTenantInterceptor(),
        audit.LoggingInterceptor(),
        rate.LimiterInterceptor(100, 200),
        recovery.RecoveryInterceptor(),
    )
    
    return grpc.NewServer(chain)
}
```

### C.5 API Gateway Integration

**Enhanced Gateway Configuration:**

```yaml
# gateway/config/routes.yaml
routes:
  - path: /v1/documents
    methods: [POST]
    service: docgen
    authentication: required
    authorization:
      resource: "document:generation:*"
      action: "generate"
    rate_limit: 100/hour
    validation:
      schema: "document_generation.json"
      
  - path: /v1/media/upload
    methods: [POST]
    service: media
    authentication: required
    authorization:
      resource: "media:file:*"
      action: "upload"
    rate_limit: 50/hour
    validation:
      max_file_size: 10485760
      allowed_mime_types: ["image/*", "application/pdf"]
      
  - path: /v1/documents/{id}/download
    methods: [GET]
    service: docgen
    authentication: required
    authorization:
      resource: "document:generation:{id}"
      action: "read"
    cache:
      ttl: 3600
      public: false
```

### C.6 Audit Trail Integration

**Comprehensive Audit Logging:**

```go
// Audit event structure
type AuditEvent struct {
    EventID      string    `json:"event_id"`
    Timestamp    time.Time `json:"timestamp"`
    Service      string    `json:"service"`      // "docgen" or "media"
    Action       string    `json:"action"`       // "generate", "upload", "delete"
    ResourceType string    `json:"resource_type"` // "document", "template", "media"
    ResourceID   string    `json:"resource_id"`
    TenantID     string    `json:"tenant_id"`
    UserID       string    `json:"user_id"`
    UserRole     string    `json:"user_role"`
    IPAddress    string    `json:"ip_address"`
    UserAgent    string    `json:"user_agent"`
    Details      string    `json:"details"`      // JSON details
    Status       string    `json:"status"`       // "success", "failure"
    Error        string    `json:"error,omitempty"`
}

// Audit logging in services
func (s *DocGenService) GenerateDocument(ctx context.Context, req *GenerateDocumentRequest) (*GenerateDocumentResponse, error) {
    // Start audit event
    auditEvent := &AuditEvent{
        EventID:      uuid.New().String(),
        Timestamp:    time.Now(),
        Service:      "docgen",
        Action:       "generate",
        ResourceType: "document",
        TenantID:     req.TenantId,
        UserID:       claims.Subject,
        UserRole:     claims.Role,
        Details:      fmt.Sprintf(`{"template_id": "%s", "entity_type": "%s"}`, req.TemplateId, req.EntityType),
    }
    
    defer func() {
        auditEvent.Status = "success"
        if err != nil {
            auditEvent.Status = "failure"
            auditEvent.Error = err.Error()
        }
        s.auditLogger.Log(auditEvent)
    }()
    
    // Business logic
    // ...
}
```

### C.7 Enhanced Error Responses

**Structured Error Handling:**

```protobuf
// Enhanced error response
message EnhancedError {
  string code = 1;                    // Error code (e.g., "PERMISSION_DENIED")
  string message = 2;                 // User-friendly message
  string detail = 3;                  // Technical details
  repeated string violations = 4;     // Validation violations
  string trace_id = 5;                // Correlation ID
  string documentation_url = 6;       // Help documentation
  map<string, string> metadata = 7;   // Additional context
}

// Permission denied error example
{
  "error": {
    "code": "PERMISSION_DENIED",
    "message": "You don't have permission to generate this document",
    "detail": "User lacks 'generate' permission on resource 'document:template:policy'",
    "violations": [
      "Missing role:agent or role:admin",
      "Tenant context mismatch"
    ],
    "trace_id": "trace-123456",
    "documentation_url": "https://docs.insuretech.com/errors/permission-denied",
    "metadata": {
      "required_role": "agent",
      "resource": "document:template:policy",
      "action": "generate"
    }
  }
}
```

---

## Appendix D: Enhanced Security Implementation

### D.1 File Upload Security

**Enhanced Media Upload Validation:**

```go
func (s *MediaService) validateUpload(file []byte, filename string, contentType string) error {
    // 1. File type validation
    if !s.isAllowedFileType(contentType, filename) {
        return fmt.Errorf("file type not allowed: %s", contentType)
    }
    
    // 2. File size validation
    if len(file) > 10*1024*1024 { // 10MB
        return fmt.Errorf("file size exceeds 10MB limit")
    }
    
    // 3. Magic number validation (not just extension)
    detectedType, err := filetype.Match(file)
    if err != nil || detectedType.MIME.Value != contentType {
        return fmt.Errorf("file content doesn't match declared type")
    }
    
    // 4. Image-specific validation
    if strings.HasPrefix(contentType, "image/") {
        if err := s.validateImage(file); err != nil {
            return err
        }
    }
    
    // 5. PDF-specific validation
    if contentType == "application/pdf" {
        if err := s.validatePDF(file); err != nil {
            return err
        }
    }
    
    // 6. Virus scanning (async)
    go s.scanForViruses(file, filename)
    
    return nil
}

func (s *MediaService) validateImage(file []byte) error {
    // Decode image
    img, format, err := image.Decode(bytes.NewReader(file))
    if err != nil {
        return fmt.Errorf("invalid image: %v", err)
    }
    
    // Check dimensions
    bounds := img.Bounds()
    width := bounds.Dx()
    height := bounds.Dy()
    
    if width < 100 || height < 100 {
        return fmt.Errorf("image too small: %dx%d", width, height)
    }
    
    if width > 4000 || height > 4000 {
        return fmt.Errorf("image too large: %dx%d", width, height)
    }
    
    // Check for embedded scripts (SVG)
    if format == "svg" {
        if strings.Contains(string(file), "<script>") {
            return fmt.Errorf("SVG contains scripts")
        }
    }
    
    return nil
}
```

### D.2 PDF Security Features

**Enhanced PDF Generation Security:**

```go
func (s *DocGenService) generateSecurePDF(htmlContent string, options PDFOptions) ([]byte, error) {
    // Generate PDF with security features
    pdfOptions := gotenberg.PDFOptions{
        Margins: gotenberg.Margins{
            Top:    "1in",
            Bottom: "1in",
            Left:   "1in",
            Right:  "1in",
        },
        Header:  s.generateHeader(),
        Footer:  s.generateFooter(),
        // Security features
        PDFFormat: "PDF/A-2b", // Archival standard
        PDFProperties: gotenberg.PDFProperties{
            Title:              options.Title,
            Author:             options.Author,
            Subject:            options.Subject,
            Keywords:           options.Keywords,
            Creator:            "InsureTech DocGen Service",
            Producer:           "Gotenberg",
            CreationDate:       time.Now().Format(time.RFC3339),
            ModDate:            time.Now().Format(time.RFC3339),
        },
        // Encryption (if required)
        PDFEncryption: gotenberg.PDFEncryption{
            UserPassword:      options.UserPassword,
            OwnerPassword:     options.OwnerPassword,
            CanPrint:          true,
            CanModify:         false,
            CanCopy:           true,
            CanAnnotate:       false,
            CanFillForms:      false,
            CanExtractContent: true,
            CanAssemble:       false,
        },
        // Watermark for sensitive documents
        Watermark: gotenberg.Watermark{
            Text:     "CONFIDENTIAL",
            FontSize: 48,
            Opacity:  0.1,
            Angle:    45,
        },
    }
    
    // Generate PDF
    pdfBytes, err := s.gotenbergClient.ConvertHTML(htmlContent, pdfOptions)
    if err != nil {
        return nil, fmt.Errorf("PDF generation failed: %w", err)
    }
    
    // Add digital signature if required
    if options.Sign {
        pdfBytes, err = s.signPDF(pdfBytes, options.SigningCertificate)
        if err != nil {
            return nil, fmt.Errorf("PDF signing failed: %w", err)
        }
    }
    
    return pdfBytes, nil
}
```

### D.3 Data Retention & Compliance

**Enhanced Retention Policies:**

```go
type RetentionPolicy struct {
    PolicyType    string        `json:"policy_type"`    // "document", "media"
    RetentionDays int           `json:"retention_days"` // Days to retain
    LegalHold     bool          `json:"legal_hold"`     // Legal hold flag
    AutoDelete    bool          `json:"auto_delete"`    // Auto-delete enabled
    ArchivePath   string        `json:"archive_path"`   // Archive location
    Compliance    []string      `json:"compliance"`     // ["GDPR", "HIPAA", "PCI"]
}

func (s *DocGenService) applyRetentionPolicy(documentID string) error {
    // Get document
    doc, err := s.generationRepo.GetByID(ctx, documentID)
    if err != nil {
        return err
    }
    
    // Determine retention policy based on document type
    var policy RetentionPolicy
    switch doc.EntityType {
    case "POLICY":
        policy = RetentionPolicy{
            PolicyType:    "document",
            RetentionDays: 3650, // 10 years for policies
            LegalHold:     true,
            AutoDelete:    false,
            ArchivePath:   "s3://archive/policies/",
            Compliance:    []string{"GDPR", "InsuranceRegulation"},
        }
    case "CLAIM":
        policy = RetentionPolicy{
            PolicyType:    "document",
            RetentionDays: 1825, // 5 years for claims
            LegalHold:     true,
            AutoDelete:    true,
            ArchivePath:   "s3://archive/claims/",
            Compliance:    []string{"GDPR"},
        }
    case "RECEIPT":
        policy = RetentionPolicy{
            PolicyType:    "document",
            RetentionDays: 730, // 2 years for receipts
            LegalHold:     false,
            AutoDelete:    true,
            ArchivePath:   "s3://archive/receipts/",
            Compliance:    []string{"GDPR"},
        }
    }
    
    // Apply retention policy
    expiryDate := time.Now().AddDate(0, 0, policy.RetentionDays)
    
    // Update document with expiry
    doc.ExpiresAt = timestamppb.New(expiryDate)
    doc.RetentionPolicy = policy.ToJSON()
    
    // Schedule auto-delete if enabled
    if policy.AutoDelete {
        s.scheduleDeletion(documentID, expiryDate)
    }
    
    return s.generationRepo.Update(ctx, doc)
}
```

---

## Appendix E: Enhanced Performance Optimization

### E.1 Caching Strategy

**Multi-level Caching Architecture:**

```go
type CacheManager struct {
    // L1: In-memory cache (fast, limited)
    memoryCache *ristretto.Cache
    
    // L2: Redis cache (distributed)
    redisClient *redis.Client
    
    // L3: CDN cache (edge)
    cdnClient *CDNClient
    
    // Cache policies
    policies map[string]CachePolicy
}

type CachePolicy struct {
    TTL             time.Duration `json:"ttl"`
    MaxSize         int64         `json:"max_size"`
    Strategy        string        `json:"strategy"` // "write-through", "write-behind"
    Invalidation    []string      `json:"invalidation"` // Events that invalidate
    Prefetch        bool          `json:"prefetch"`     // Prefetch on read
}

// Cache policies for different resources
var cachePolicies = map[string]CachePolicy{
    "document:templates": {
        TTL:      24 * time.Hour,
        MaxSize:  1000,
        Strategy: "write-through",
        Invalidation: []string{
            "template.created",
            "template.updated",
            "template.deleted",
        },
        Prefetch: true,
    },
    "document:generated": {
        TTL:      1 * time.Hour,
        MaxSize:  10000,
        Strategy: "write-behind",
        Invalidation: []string{
            "document.regenerated",
            "document.deleted",
        },
        Prefetch: false,
    },
    "media:thumbnails": {
        TTL:      7 * 24 * time.Hour,
        MaxSize:  100000,
        Strategy: "write-through",
        Invalidation: []string{
            "media.updated",
            "media.deleted",
        },
        Prefetch: true,
    },
}
```

### E.2 Connection Pooling

**Enhanced gRPC Connection Management:**

```go
type ConnectionPool struct {
    connections map[string]*grpc.ClientConn
    mu          sync.RWMutex
    config      PoolConfig
}

type PoolConfig struct {
    MaxIdleConns     int           `json:"max_idle_conns"`
    MaxOpenConns     int           `json:"max_open_conns"`
    ConnMaxLifetime  time.Duration `json:"conn_max_lifetime"`
    ConnMaxIdleTime  time.Duration `json:"conn_max_idle_time"`
    HealthCheckInterval time.Duration `json:"health_check_interval"`
}

// Service-specific connection pools
var servicePools = map[string]PoolConfig{
    "storage": {
        MaxIdleConns:     10,
        MaxOpenConns:     100,
        ConnMaxLifetime:  5 * time.Minute,
        ConnMaxIdleTime:  1 * time.Minute,
        HealthCheckInterval: 30 * time.Second,
    },
    "authz": {
        MaxIdleConns:     5,
        MaxOpenConns:     50,
        ConnMaxLifetime:  10 * time.Minute,
        ConnMaxIdleTime:  2 * time.Minute,
        HealthCheckInterval: 30 * time.Second,
    },
    "notification": {
        MaxIdleConns:     3,
        MaxOpenConns:     20,
        ConnMaxLifetime:  5 * time.Minute,
        ConnMaxIdleTime:  1 * time.Minute,
        HealthCheckInterval: 30 * time.Second,
    },
}
```

### E.3 Database Optimization

**Enhanced Query Patterns:**

```go
// Optimized repository methods
type OptimizedDocumentRepository struct {
    db *sqlx.DB
    // Prepared statements
    stmts map[string]*sqlx.Stmt
}

func (r *OptimizedDocumentRepository) initPreparedStatements() error {
    statements := map[string]string{
        "GetByID": `
            SELECT * FROM storage_schema.document_generations 
            WHERE generation_id = $1 
            AND tenant_id = $2
            LIMIT 1
        `,
        "ListByEntity": `
            SELECT * FROM storage_schema.document_generations 
            WHERE entity_type = $1 
            AND entity_id = $2 
            AND tenant_id = $3
            ORDER BY generated_at DESC
            LIMIT $4 OFFSET $5
        `,
        "BatchInsert": `
            INSERT INTO storage_schema.document_generations 
            (generation_id, document_template_id, entity_type, entity_id, data, status)
            VALUES %s
            ON CONFLICT (generation_id) DO UPDATE SET
            updated_at = CURRENT_TIMESTAMP
        `,
    }
    
    for name, query := range statements {
        stmt, err := r.db.Preparex(query)
        if err != nil {
            return fmt.Errorf("failed to prepare statement %s: %w", name, err)
        }
        r.stmts[name] = stmt
    }
    
    return nil
}

// Use prepared statements
func (r *OptimizedDocumentRepository) GetByID(ctx context.Context, id, tenantID string) (*documentv1.DocumentGeneration, error) {
    var doc documentv1.DocumentGeneration
    err := r.stmts["GetByID"].GetContext(ctx, &doc, id, tenantID)
    if err != nil {
        return nil, err
    }
    return &doc, nil
}
```

---

## Appendix F: Enhanced Deployment Configuration

### F.1 Kubernetes Deployment

**DocGen Service Deployment:**

```yaml
# k8s/deployments/docgen.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: docgen-service
  namespace: insuretech
spec:
  replicas: 3
  selector:
    matchLabels:
      app: docgen
      tier: backend
  template:
    metadata:
      labels:
        app: docgen
        tier: backend
    spec:
      containers:
      - name: docgen
        image: insuretech/docgen:latest
        ports:
        - containerPort: 8080
          name: grpc
        - containerPort: 8081
          name: http
        env:
        - name: ENVIRONMENT
          value: "production"
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-secrets
              key: connection-string
        - name: STORAGE_SERVICE_URL
          value: "storage-service.insuretech.svc.cluster.local:8080"
        - name: AUTHZ_SERVICE_URL
          value: "authz-service.insuretech.svc.cluster.local:8080"
        - name: KAFKA_BROKERS
          value: "kafka-0.kafka.insuretech.svc.cluster.local:9092,kafka-1.kafka.insuretech.svc.cluster.local:9092"
        resources:
          requests:
            memory: "512Mi"
            cpu: "250m"
          limits:
            memory: "1Gi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 5
      # Sidecar for Gotenberg
      - name: gotenberg
        image: gotenberg/gotenberg:7.10
        ports:
        - containerPort: 3000
          name: gotenberg
        resources:
          requests:
            memory: "256Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "200m"
```

### F.2 Horizontal Pod Autoscaler

```yaml
# k8s/hpa/docgen-hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: docgen-hpa
  namespace: insuretech
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: docgen-service
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  - type: Pods
    pods:
      metric:
        name: pdf_generation_queue_length
      target:
        type: AverageValue
        averageValue: 100
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 10
        periodSeconds: 60
    scaleUp:
      stabilizationWindowSeconds: 60
      policies:
      - type: Percent
        value: 100
        periodSeconds: 60
```

### F.3 Service Mesh Configuration

```yaml
# istio/virtual-service.yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: docgen-virtual-service
  namespace: insuretech
spec:
  hosts:
  - docgen-service.insuretech.svc.cluster.local
  - docgen.insuretech.com
  http:
  - match:
    - uri:
        prefix: /v1/documents
    route:
    - destination:
        host: docgen-service.insuretech.svc.cluster.local
        port:
          number: 8080
    timeout: 30s
    retries:
      attempts: 3
      perTryTimeout: 10s
    corsPolicy:
      allowOrigin:
      - "*.insuretech.com"
      allowMethods:
      - GET
      - POST
      - PUT
      - DELETE
      allowHeaders:
      - authorization
      - content-type
      - x-tenant-id
    fault:
      delay:
        percentage:
          value: 0.1
        fixedDelay: 5s
  - match:
    - uri:
        prefix: /health
    route:
    - destination:
        host: docgen-service.insuretech.svc.cluster.local
        port:
          number: 8081
```

---

## Appendix G: Enhanced Monitoring & Alerting

### G.1 Prometheus Metrics

**Enhanced Metrics Collection:**

```go
// Metrics initialization
func initMetrics() {
    // Document generation metrics
    docgenGenerationDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "docgen_generation_duration_seconds",
            Help:    "Time taken to generate documents",
            Buckets: prometheus.DefBuckets,
        },
        []string{"template_type", "status"},
    )
    
    docgenTemplateCache = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "docgen_template_cache_size",
            Help: "Number of templates in cache",
        },
        []string{"cache_level"},
    )
    
    docgenAsyncQueue = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "docgen_async_queue_size",
            Help: "Number of documents in async queue",
        },
    )
    
    // Media processing metrics
    mediaProcessingDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "media_processing_duration_seconds",
            Help:    "Time taken to process media files",
            Buckets: prometheus.DefBuckets,
        },
        []string{"processing_type", "status"},
    )
    
    mediaOCRAccuracy = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "media_ocr_accuracy_score",
            Help: "OCR accuracy score (0-1)",
        },
        []string{"language"},
    )
    
    // Register all metrics
    prometheus.MustRegister(
        docgenGenerationDuration,
        docgenTemplateCache,
        docgenAsyncQueue,
        mediaProcessingDuration,
        mediaOCRAccuracy,
    )
}
```

### G.2 Grafana Dashboards

**Enhanced Dashboard Configuration:**

```json
{
  "dashboard": {
    "title": "DocGen & Media Services",
    "panels": [
      {
        "title": "PDF Generation Performance",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(docgen_generation_duration_seconds_bucket[5m]))",
            "legendFormat": "95th percentile"
          },
          {
            "expr": "histogram_quantile(0.50, rate(docgen_generation_duration_seconds_bucket[5m]))",
            "legendFormat": "50th percentile"
          }
        ],
        "thresholds": [
          {
            "value": 3,
            "color": "red",
            "op": "gt"
          },
          {
            "value": 2,
            "color": "yellow",
            "op": "gt"
          }
        ]
      },
      {
        "title": "Media Processing Queue",
        "type": "stat",
        "targets": [
          {
            "expr": "docgen_async_queue_size",
            "legendFormat": "Queue Size"
          }
        ],
        "thresholds": [
          {
            "value": 1000,
            "color": "red"
          },
          {
            "value": 500,
            "color": "yellow"
          }
        ]
      },
      {
        "title": "OCR Accuracy by Language",
        "type": "bar",
        "targets": [
          {
            "expr": "media_ocr_accuracy_score",
            "legendFormat": "{{language}}"
          }
        ],
        "thresholds": [
          {
            "value": 0.8,
            "color": "green"
          },
          {
            "value": 0.6,
            "color": "yellow"
          }
        ]
      }
    ]
  }
}
```

### G.3 AlertManager Configuration

```yaml
# alertmanager/config.yaml
route:
  group_by: ['alertname', 'cluster', 'service']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 4h
  receiver: 'slack-notifications'
  routes:
  - match:
      severity: critical
    receiver: 'pagerduty'
  - match:
      severity: warning
    receiver: 'slack-notifications'

receivers:
- name: 'slack-notifications'
  slack_configs:
  - channel: '#alerts-docgen-media'
    title: '{{ .GroupLabels.alertname }}'
    text: '{{ .CommonAnnotations.summary }}'
    
- name: 'pagerduty'
  pagerduty_configs:
  - service_key: '{{ .PagerDutyKey }}'
    description: '{{ .CommonAnnotations.description }}'
    severity: 'critical'

inhibit_rules:
- source_match:
    severity: 'critical'
  target_match:
    severity: 'warning'
  equal: ['alertname', 'cluster', 'service']
```

---

## Conclusion

This enhanced implementation plan builds upon the **solid foundation** already established in the codebase while addressing the **critical gaps** needed for production readiness. The plan focuses on:

1. **Completing existing implementations** with production-grade features
2. **Integrating with established authentication/authorization** patterns
3. **Adding comprehensive monitoring and observability**
4. **Implementing robust security and compliance** measures
5. **Ensuring scalability and performance** from day one
6. **Providing thorough testing and validation** strategies

The **15-day timeline** is aggressive but achievable given the existing codebase foundation. Each phase builds upon the previous, with clear deliverables and success criteria.

**Key Success Factors:**
- Leverage existing proto definitions and patterns
- Integrate with established AuthZ service
- Implement comprehensive testing early
- Focus on production readiness metrics
- Document all APIs and deployment procedures

This enhanced plan positions the DocGen and Media services for **successful production deployment** with the reliability, security, and performance required for the InsureTech platform.
