package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UploadUserDocument validates the document type and creates a document record.
func (s *AuthService) UploadUserDocument(ctx context.Context, req *authnservicev1.UploadUserDocumentRequest) (*authnservicev1.UploadUserDocumentResponse, error) {
	if s.userDocumentRepo == nil || s.documentTypeRepo == nil {
		return nil, errors.New("document repositories not configured")
	}

	// Validate document type exists.
	if _, err := s.documentTypeRepo.GetByID(ctx, req.DocumentTypeId); err != nil {
		logger.Errorf("invalid document_type_id: %v", err)
		return nil, errors.New("invalid document_type_id")
	}

	// Basic file_url validation.
	if !strings.HasPrefix(req.FileUrl, "https://") && !strings.HasPrefix(req.FileUrl, "http://") {
		return nil, errors.New("file_url must be a valid URL")
	}

	doc := &authnentityv1.UserDocument{
		UserDocumentId:     uuid.NewString(),
		UserId:             req.UserId,
		DocumentTypeId:     req.DocumentTypeId,
		FileUrl:            req.FileUrl,
		PolicyId:           req.PolicyId,
		VerificationStatus: "PENDING",
		CreatedAt:          timestamppb.Now(),
		UpdatedAt:          timestamppb.Now(),
	}

	if err := s.userDocumentRepo.Create(ctx, doc); err != nil {
		appLogger.Errorf("UploadUserDocument: failed for user %s: %v", req.UserId, err)
		logger.Errorf("failed to save document: %v", err)
		return nil, errors.New("failed to save document")
	}

	appLogger.Infof("UploadUserDocument: created document %s for user %s", doc.UserDocumentId, req.UserId)

	return &authnservicev1.UploadUserDocumentResponse{
		Document: doc,
		Message:  "Document uploaded successfully",
	}, nil
}

// ListUserDocuments lists all documents for a user with optional type filter.
func (s *AuthService) ListUserDocuments(ctx context.Context, req *authnservicev1.ListUserDocumentsRequest) (*authnservicev1.ListUserDocumentsResponse, error) {
	if s.userDocumentRepo == nil {
		return nil, errors.New("document repository not configured")
	}

	docs, err := s.userDocumentRepo.ListByUser(ctx, req.UserId)
	if err != nil {
		logger.Errorf("failed to list documents: %v", err)
		return nil, errors.New("failed to list documents")
	}

	// Optional filter by document_type_id.
	if req.DocumentTypeId != "" {
		filtered := make([]*authnentityv1.UserDocument, 0, len(docs))
		for _, d := range docs {
			if d.DocumentTypeId == req.DocumentTypeId {
				filtered = append(filtered, d)
			}
		}
		docs = filtered
	}

	return &authnservicev1.ListUserDocumentsResponse{
		Documents:  docs,
		TotalCount: int32(len(docs)),
	}, nil
}

// GetUserDocument retrieves a single document by ID.
func (s *AuthService) GetUserDocument(ctx context.Context, req *authnservicev1.GetUserDocumentRequest) (*authnservicev1.GetUserDocumentResponse, error) {
	if s.userDocumentRepo == nil {
		return nil, errors.New("document repository not configured")
	}

	doc, err := s.userDocumentRepo.GetByID(ctx, req.UserDocumentId)
	if err != nil {
		logger.Errorf("document not found: %v", err)
		return nil, errors.New("document not found")
	}

	return &authnservicev1.GetUserDocumentResponse{Document: doc}, nil
}

// UpdateUserDocument updates mutable fields for a user document.
func (s *AuthService) UpdateUserDocument(ctx context.Context, req *authnservicev1.UpdateUserDocumentRequest) (*authnservicev1.UpdateUserDocumentResponse, error) {
	if s.userDocumentRepo == nil {
		return nil, errors.New("document repository not configured")
	}

	updates := map[string]any{}

	if req.DocumentTypeId != "" {
		if s.documentTypeRepo == nil {
			return nil, errors.New("document type repository not configured")
		}
		if _, err := s.documentTypeRepo.GetByID(ctx, req.DocumentTypeId); err != nil {
			logger.Errorf("invalid document_type_id: %v", err)
			return nil, errors.New("invalid document_type_id")
		}
		updates["document_type_id"] = req.DocumentTypeId
	}

	if req.FileUrl != "" {
		if !strings.HasPrefix(req.FileUrl, "https://") && !strings.HasPrefix(req.FileUrl, "http://") {
			return nil, errors.New("file_url must be a valid URL")
		}
		updates["file_url"] = req.FileUrl
	}

	if req.PolicyId != "" {
		updates["policy_id"] = req.PolicyId
	}

	if len(updates) == 0 {
		return nil, errors.New("no updatable fields provided")
	}

	if err := s.userDocumentRepo.Update(ctx, req.UserDocumentId, updates); err != nil {
		logger.Errorf("failed to update document: %v", err)
		return nil, errors.New("failed to update document")
	}

	doc, err := s.userDocumentRepo.GetByID(ctx, req.UserDocumentId)
	if err != nil {
		return nil, errors.New("failed to fetch updated document")
	}

	return &authnservicev1.UpdateUserDocumentResponse{
		Document: doc,
		Message:  "Document updated successfully",
	}, nil
}

// DeleteUserDocument soft-deletes a document.
func (s *AuthService) DeleteUserDocument(ctx context.Context, req *authnservicev1.DeleteUserDocumentRequest) (*authnservicev1.DeleteUserDocumentResponse, error) {
	if s.userDocumentRepo == nil {
		return nil, errors.New("document repository not configured")
	}

	if err := s.userDocumentRepo.Delete(ctx, req.UserDocumentId); err != nil {
		logger.Errorf("failed to delete document: %v", err)
		return nil, errors.New("failed to delete document")
	}

	appLogger.Infof("DeleteUserDocument: deleted document %s", req.UserDocumentId)

	return &authnservicev1.DeleteUserDocumentResponse{
		Message: "Document deleted successfully",
	}, nil
}

// ListDocumentTypes returns all (or active-only) document types.
func (s *AuthService) ListDocumentTypes(ctx context.Context, req *authnservicev1.ListDocumentTypesRequest) (*authnservicev1.ListDocumentTypesResponse, error) {
	if s.documentTypeRepo == nil {
		return nil, errors.New("document type repository not configured")
	}

	types, err := s.documentTypeRepo.ListActive(ctx)
	if err != nil {
		logger.Errorf("failed to list document types: %v", err)
		return nil, errors.New("failed to list document types")
	}

	return &authnservicev1.ListDocumentTypesResponse{Types: types}, nil
}
