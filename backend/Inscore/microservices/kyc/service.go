package kyc

// service.go — KYCService gRPC implementation for InsureTech.
//
// This package implements the insuretech.kyc.services.v1.KYCService RPC
// contract backed by the authn_schema.kyc_verifications table.
//
// Design:
//   - InsureTech kyc_verifications.id  = canonical InsureTech UUID (never exposed externally as FLVE session)
//   - InsureTech kyc_verifications.provider_reference = FLVE session_id (opaque, internal only)
//   - This service is the gRPC-accessible KYC backend; the authn service is the
//     user-facing orchestrator that calls FLVE directly via FLVEAdapter.
//   - S3 storage credentials come from environment (ACCESS_KEY / SECRET_KEY / DO_BUCKET / DO_REGION / MAIN_CDN).

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/grpcmeta"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	kycv1 "github.com/newage-saint/insuretech/gen/go/insuretech/kyc/entity/v1"
	kycservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/kyc/services/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// KYCService implements kycservicev1.KYCServiceServer.
type KYCService struct {
	kycservicev1.UnimplementedKYCServiceServer
	db          *gorm.DB
	flveAdapter FLVEAdapter
}

// NewKYCService creates a new KYCService backed by the given GORM DB.
func NewKYCService(db *gorm.DB) *KYCService {
	return &KYCService{db: db}
}

// SetFLVEAdapter sets the FLVE adapter for starting eKYC sessions.
func (s *KYCService) SetFLVEAdapter(adapter FLVEAdapter) {
	s.flveAdapter = adapter
}

// ── RPC implementations ──────────────────────────────────────────────────────

// StartKYCVerification creates a new KYC verification record.
// When method is FLVE_EKYC and the FLVE adapter is configured, it synchronously
// starts an eKYC session on the FLVE HuggingFace Space and stores the session_id
// as provider_reference. The KycVerificationId returned is the InsureTech canonical
// UUID that the caller (authn) can use to look up the session_id via GetKYCVerification.
func (s *KYCService) StartKYCVerification(ctx context.Context, req *kycservicev1.StartKYCVerificationRequest) (*kycservicev1.StartKYCVerificationResponse, error) {
	if req.GetEntityId() == "" {
		return nil, status.Error(codes.InvalidArgument, "entity_id is required")
	}
	if req.GetEntityType() == "" {
		return nil, status.Error(codes.InvalidArgument, "entity_type is required")
	}

	kycID := uuid.New().String()
	method := mapVerificationMethod(req.GetMethod())
	vtype := mapVerificationType(req.GetType())

	// For FLVE_EKYC — call the HF Space synchronously to get the session_id.
	// This makes the response useful: callers can immediately get the session_id
	// by calling GetKYCVerification after this returns.
	var flveSessionID string
	isFLVE := method == kycv1.VerificationMethod_VERIFICATION_METHOD_FLVE_EKYC
	if isFLVE && s.flveAdapter != nil {
		tenantID := grpcmeta.TenantID(ctx, "insuretech")
		userType := grpcmeta.FirstOfFromContext(ctx, "x-user-type", "user-type")
		portal := grpcmeta.NormalizePortal(grpcmeta.FirstOfFromContext(ctx, "x-portal", "portal"))
		if portal == "" {
			portal = "b2c"
		}
		flveResp, err := s.flveAdapter.StartEKYC(ctx, &FLVEStartRequest{
			UserID:            req.GetEntityId(),
			TenantID:          tenantID,
			UserType:          userType,
			Portal:            portal,
			KYCVerificationID: kycID,
		})
		if err != nil {
			logger.Errorf("StartKYCVerification: FLVE call failed: %v", err)
			return nil, status.Errorf(codes.Internal, "FLVE eKYC session failed: %v", err)
		}
		flveSessionID = flveResp.SessionID
		logger.Infof("FLVE session started: kyc_id=%s session_id=%s", kycID, flveSessionID)
	}

	k := &kycv1.KYCVerification{
		Id:                kycID,
		Type:              vtype,
		EntityType:        req.GetEntityType(),
		EntityId:          req.GetEntityId(),
		Method:            method,
		Status:            kycv1.VerificationStatus_VERIFICATION_STATUS_IN_PROGRESS,
		Provider:          "FLVE",
		ProviderReference: flveSessionID,
	}

	if err := s.create(ctx, k); err != nil {
		logger.Errorf("StartKYCVerification: %v", err)
		return nil, status.Error(codes.Internal, "failed to create KYC verification")
	}

	logger.Infof("KYC verification started: id=%s entity_id=%s method=%s session_id=%s", kycID, req.GetEntityId(), req.GetMethod(), flveSessionID)
	return &kycservicev1.StartKYCVerificationResponse{
		KycVerificationId: kycID,
		Message:           flveSessionID, // Encode session_id in Message field for caller
	}, nil
}

// GetKYCVerification retrieves a KYC verification record by ID.
func (s *KYCService) GetKYCVerification(ctx context.Context, req *kycservicev1.GetKYCVerificationRequest) (*kycservicev1.GetKYCVerificationResponse, error) {
	if req.GetKycVerificationId() == "" {
		return nil, status.Error(codes.InvalidArgument, "kyc_verification_id is required")
	}

	k, err := s.getByID(ctx, req.GetKycVerificationId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "KYC verification not found")
		}
		logger.Errorf("GetKYCVerification: %v", err)
		return nil, status.Error(codes.Internal, "failed to get KYC verification")
	}

	return &kycservicev1.GetKYCVerificationResponse{
		KycVerification: k,
	}, nil
}

// UploadDocument stores a document reference against a KYC verification.
// In the FLVE flow the authn service proxies frames directly to FLVE; this
// method handles the non-FLVE legacy path where document URLs are submitted.
func (s *KYCService) UploadDocument(ctx context.Context, req *kycservicev1.UploadDocumentRequest) (*kycservicev1.UploadDocumentResponse, error) {
	if req.GetKycVerificationId() == "" {
		return nil, status.Error(codes.InvalidArgument, "kyc_verification_id is required")
	}

	docID := uuid.New().String()

	// Append document reference to the verification_result JSONB column.
	snapshot := `{"doc_id":"` + docID + `","doc_type":"` + req.GetDocumentType() + `","doc_number":"` + req.GetDocumentNumber() + `","doc_url":"[REDACTED]","uploaded_at":"` + time.Now().UTC().Format(time.RFC3339) + `"}`
	if err := s.appendVerificationResult(ctx, req.GetKycVerificationId(), snapshot); err != nil {
		logger.Warnf("UploadDocument append result: %v", err)
		// non-fatal — doc_id is still returned
	}

	// Advance status to IN_PROGRESS if still PENDING.
	_ = s.advanceStatusIfPending(ctx, req.GetKycVerificationId())

	return &kycservicev1.UploadDocumentResponse{
		DocumentVerificationId: docID,
		Message:                "Document uploaded successfully",
	}, nil
}

// VerifyKYC marks a KYC verification as VERIFIED.
func (s *KYCService) VerifyKYC(ctx context.Context, req *kycservicev1.VerifyKYCRequest) (*kycservicev1.VerifyKYCResponse, error) {
	if req.GetKycVerificationId() == "" {
		return nil, status.Error(codes.InvalidArgument, "kyc_verification_id is required")
	}
	if req.GetVerifiedBy() == "" {
		return nil, status.Error(codes.InvalidArgument, "verified_by is required")
	}

	now := time.Now()
	upd := map[string]any{
		"status":      "VERIFIED",
		"verified_by": req.GetVerifiedBy(),
		"verified_at": now,
		"updated_at":  now,
	}
	if req.GetVerificationResult() != "" {
		snapshot := `{"result":"` + req.GetVerificationResult() + `","verified_by":"` + req.GetVerifiedBy() + `","verified_at":"` + now.UTC().Format(time.RFC3339) + `"}`
		_ = s.appendVerificationResult(ctx, req.GetKycVerificationId(), snapshot)
	}

	if err := s.db.WithContext(ctx).Table("authn_schema.kyc_verifications").
		Where("verification_id = ?", req.GetKycVerificationId()).
		Updates(upd).Error; err != nil {
		logger.Errorf("VerifyKYC: %v", err)
		return nil, status.Error(codes.Internal, "failed to verify KYC")
	}

	logger.Infof("KYC verified: id=%s by=%s", req.GetKycVerificationId(), req.GetVerifiedBy())
	return &kycservicev1.VerifyKYCResponse{
		Message: "KYC verification completed successfully",
	}, nil
}

// RejectKYC marks a KYC verification as REJECTED.
func (s *KYCService) RejectKYC(ctx context.Context, req *kycservicev1.RejectKYCRequest) (*kycservicev1.RejectKYCResponse, error) {
	if req.GetKycVerificationId() == "" {
		return nil, status.Error(codes.InvalidArgument, "kyc_verification_id is required")
	}

	now := time.Now()
	upd := map[string]any{
		"status":           "REJECTED",
		"rejection_reason": req.GetReason(),
		"updated_at":       now,
	}
	if err := s.db.WithContext(ctx).Table("authn_schema.kyc_verifications").
		Where("verification_id = ?", req.GetKycVerificationId()).
		Updates(upd).Error; err != nil {
		logger.Errorf("RejectKYC: %v", err)
		return nil, status.Error(codes.Internal, "failed to reject KYC")
	}

	logger.Infof("KYC rejected: id=%s reason=%s", req.GetKycVerificationId(), req.GetReason())
	return &kycservicev1.RejectKYCResponse{
		Message: "KYC verification rejected",
	}, nil
}

// ListPendingVerifications returns KYC records awaiting review.
func (s *KYCService) ListPendingVerifications(ctx context.Context, req *kycservicev1.ListPendingVerificationsRequest) (*kycservicev1.ListPendingVerificationsResponse, error) {
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	// Statuses that require human review — using the canonical DB strings from the proto enum.
	reviewStatuses := []string{
		strings.TrimPrefix(kycv1.VerificationStatus_VERIFICATION_STATUS_PENDING.String(), "VERIFICATION_STATUS_"),
		strings.TrimPrefix(kycv1.VerificationStatus_VERIFICATION_STATUS_IN_PROGRESS.String(), "VERIFICATION_STATUS_"),
		strings.TrimPrefix(kycv1.VerificationStatus_VERIFICATION_STATUS_PENDING_REVIEW.String(), "VERIFICATION_STATUS_"),
	}

	cols := `verification_id, type, entity_type, entity_id, method, provider, provider_reference, status, rejection_reason, verified_by, verified_at, expires_at`
	q := `select ` + cols + ` from authn_schema.kyc_verifications where status = ANY(?) order by verification_id desc limit ?`

	rows, err := s.db.WithContext(ctx).Raw(q, reviewStatuses, pageSize).Rows()
	if err != nil {
		logger.Errorf("ListPendingVerifications: %v", err)
		return nil, status.Error(codes.Internal, "failed to list pending verifications")
	}
	defer rows.Close()

	var verifications []*kycv1.KYCVerification
	for rows.Next() {
		k, err := scanKYC(rows)
		if err != nil {
			logger.Warnf("ListPendingVerifications scan: %v", err)
			continue
		}
		verifications = append(verifications, k)
	}
	if err := rows.Err(); err != nil {
		logger.Errorf("ListPendingVerifications rows: %v", err)
	}

	return &kycservicev1.ListPendingVerificationsResponse{
		Verifications: verifications,
		TotalCount:    int32(len(verifications)),
	}, nil
}

// StartFLVESession initiates an eKYC session via FLVE and stores the session ID.
func (s *KYCService) StartFLVESession(ctx context.Context, kycID, userID, tenantID, userType, portal string) (string, error) {
	if s.flveAdapter == nil {
		return "", errors.New("FLVE adapter not configured")
	}

	req := &FLVEStartRequest{
		UserID:            userID,
		TenantID:          tenantID,
		UserType:          userType,
		Portal:            portal,
		KYCVerificationID: kycID,
	}

	resp, err := s.flveAdapter.StartEKYC(ctx, req)
	if err != nil {
		logger.Errorf("StartFLVESession: FLVE call failed: %v", err)
		return "", err
	}

	// Update the KYC verification record with the FLVE session ID
	now := time.Now()
	if err := s.db.WithContext(ctx).Table("authn_schema.kyc_verifications").
		Where("verification_id = ?", kycID).
		Updates(map[string]any{
			"provider":           "FLVE",
			"provider_reference": resp.SessionID,
			"updated_at":         now,
		}).Error; err != nil {
		logger.Errorf("StartFLVESession: failed to update DB with session ID: %v", err)
		return "", err
	}

	logger.Infof("FLVE session started: kyc_id=%s session_id=%s", kycID, resp.SessionID)
	return resp.SessionID, nil
}

// ── DB helpers ────────────────────────────────────────────────────────────────

func (s *KYCService) create(ctx context.Context, k *kycv1.KYCVerification) error {
	return s.db.WithContext(ctx).Exec(
		`insert into authn_schema.kyc_verifications
			(verification_id, type, entity_type, entity_id, method, provider, provider_reference, documents, status, verification_result, rejection_reason, verified_by, verified_at, expires_at, audit_info)
		 values (?, ?, ?, ?, ?, ?, ?, ?::jsonb, ?, ?::jsonb, ?, ?, ?, ?, '{}'::jsonb)`,
		k.Id,
		strings.TrimPrefix(k.Type.String(), "VERIFICATION_TYPE_"),
		k.EntityType,
		k.EntityId,
		strings.TrimPrefix(k.Method.String(), "VERIFICATION_METHOD_"),
		nullStr(k.Provider),
		nullStr(k.ProviderReference),
		"null",
		strings.TrimPrefix(k.Status.String(), "VERIFICATION_STATUS_"),
		"null",
		nullStr(k.RejectionReason),
		nullStr(k.VerifiedBy),
		nil,
		nil,
	).Error
}

func (s *KYCService) getByID(ctx context.Context, id string) (*kycv1.KYCVerification, error) {
	cols := `verification_id, type, entity_type, entity_id, method, provider, provider_reference, status, rejection_reason, verified_by, verified_at, expires_at`
	q := `select ` + cols + ` from authn_schema.kyc_verifications where verification_id = ? limit 1`
	row := s.db.WithContext(ctx).Raw(q, id).Row()
	if err := row.Err(); err != nil {
		return nil, err
	}
	return scanKYC(row)
}

func (s *KYCService) appendVerificationResult(ctx context.Context, id, jsonSnapshot string) error {
	return s.db.WithContext(ctx).Exec(
		`update authn_schema.kyc_verifications
		 set verification_result = coalesce(verification_result, '[]'::jsonb) || ?::jsonb
		 where verification_id = ?`,
		"["+jsonSnapshot+"]", id,
	).Error
}

func (s *KYCService) advanceStatusIfPending(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Exec(
		`update authn_schema.kyc_verifications set status = 'IN_PROGRESS' where verification_id = ? and status = 'PENDING'`,
		id,
	).Error
}

// ── scanner ──────────────────────────────────────────────────────────────────

type rowScanner interface {
	Scan(...any) error
}

func scanKYC(row rowScanner) (*kycv1.KYCVerification, error) {
	var k kycv1.KYCVerification
	var typeStr, methodStr, statusStr string
	var provider, providerRef, rejectionReason, verifiedBy *string
	var verifiedAt, expiresAt *time.Time

	if err := row.Scan(&k.Id, &typeStr, &k.EntityType, &k.EntityId, &methodStr,
		&provider, &providerRef, &statusStr, &rejectionReason, &verifiedBy,
		&verifiedAt, &expiresAt); err != nil {
		return nil, err
	}
	if provider != nil {
		k.Provider = *provider
	}
	if providerRef != nil {
		k.ProviderReference = *providerRef
	}
	if rejectionReason != nil {
		k.RejectionReason = *rejectionReason
	}
	if verifiedBy != nil {
		k.VerifiedBy = *verifiedBy
	}
	k.Type = verificationTypeFromString(typeStr)
	k.Method = verificationMethodFromString(methodStr)
	k.Status = verificationStatusFromString(statusStr)
	return &k, nil
}

// ── enum mappers ──────────────────────────────────────────────────────────────

func mapVerificationMethod(s string) kycv1.VerificationMethod {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "PORICHOY":
		return kycv1.VerificationMethod_VERIFICATION_METHOD_PORICHOY
	case "NID":
		return kycv1.VerificationMethod_VERIFICATION_METHOD_NID
	case "PASSPORT":
		return kycv1.VerificationMethod_VERIFICATION_METHOD_PASSPORT
	case "TRADE_LICENSE":
		return kycv1.VerificationMethod_VERIFICATION_METHOD_TRADE_LICENSE
	case "FLVE_EKYC", "FLVE":
		return kycv1.VerificationMethod_VERIFICATION_METHOD_FLVE_EKYC
	default:
		return kycv1.VerificationMethod_VERIFICATION_METHOD_MANUAL
	}
}

func mapVerificationType(s string) kycv1.VerificationType {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "KYB":
		return kycv1.VerificationType_VERIFICATION_TYPE_KYB
	default:
		return kycv1.VerificationType_VERIFICATION_TYPE_KYC
	}
}

func verificationTypeFromString(s string) kycv1.VerificationType {
	s = strings.ToUpper(strings.TrimPrefix(strings.TrimSpace(s), "VERIFICATION_TYPE_"))
	if v, ok := kycv1.VerificationType_value["VERIFICATION_TYPE_"+s]; ok {
		return kycv1.VerificationType(v)
	}
	return kycv1.VerificationType_VERIFICATION_TYPE_UNSPECIFIED
}

func verificationMethodFromString(s string) kycv1.VerificationMethod {
	s = strings.ToUpper(strings.TrimPrefix(strings.TrimSpace(s), "VERIFICATION_METHOD_"))
	if v, ok := kycv1.VerificationMethod_value["VERIFICATION_METHOD_"+s]; ok {
		return kycv1.VerificationMethod(v)
	}
	return kycv1.VerificationMethod_VERIFICATION_METHOD_UNSPECIFIED
}

func verificationStatusFromString(s string) kycv1.VerificationStatus {
	s = strings.ToUpper(strings.TrimPrefix(strings.TrimSpace(s), "VERIFICATION_STATUS_"))
	if v, ok := kycv1.VerificationStatus_value["VERIFICATION_STATUS_"+s]; ok {
		return kycv1.VerificationStatus(v)
	}
	return kycv1.VerificationStatus_VERIFICATION_STATUS_UNSPECIFIED
}

func nullStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
