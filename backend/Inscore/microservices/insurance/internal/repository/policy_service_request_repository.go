package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"

	policyv1 "github.com/newage-saint/insuretech/gen/go/insuretech/policy/entity/v1"
)

type PolicyServiceRequestRepository struct {
	db *gorm.DB
}

func NewPolicyServiceRequestRepository(db *gorm.DB) *PolicyServiceRequestRepository {
	return &PolicyServiceRequestRepository{db: db}
}

func (r *PolicyServiceRequestRepository) Create(ctx context.Context, request *policyv1.PolicyServiceRequest) (*policyv1.PolicyServiceRequest, error) {
	if request.RequestId == "" {
		return nil, fmt.Errorf("request_id is required")
	}

	var processedBy sql.NullString
	if request.ProcessedBy != "" {
		processedBy = sql.NullString{String: request.ProcessedBy, Valid: true}
	}

	var processedAt sql.NullTime
	if request.ProcessedAt != nil {
		processedAt = sql.NullTime{Time: request.ProcessedAt.AsTime(), Valid: true}
	}

	err := r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.policy_service_requests
			(request_id, policy_id, customer_id, request_type, request_data, status, processed_by, processed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		request.RequestId,
		request.PolicyId,
		request.CustomerId,
		strings.ToUpper(request.RequestType.String()),
		request.RequestData,
		strings.ToUpper(request.Status.String()),
		processedBy,
		processedAt,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert policy service request: %w", err)
	}

	return r.GetByID(ctx, request.RequestId)
}

func (r *PolicyServiceRequestRepository) GetByID(ctx context.Context, requestID string) (*policyv1.PolicyServiceRequest, error) {
	var (
		req           policyv1.PolicyServiceRequest
		requestTypeStr sql.NullString
		statusStr     sql.NullString
		requestData   sql.NullString
		processedBy   sql.NullString
		processedAt   sql.NullTime
		createdAt     time.Time
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT request_id, policy_id, customer_id, request_type, request_data, status, processed_by, processed_at, created_at
		FROM insurance_schema.policy_service_requests
		WHERE request_id = $1`,
		requestID,
	).Row().Scan(
		&req.RequestId,
		&req.PolicyId,
		&req.CustomerId,
		&requestTypeStr,
		&requestData,
		&statusStr,
		&processedBy,
		&processedAt,
		&createdAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get policy service request: %w", err)
	}

	if requestTypeStr.Valid {
		k := strings.ToUpper(requestTypeStr.String)
		if v, ok := policyv1.ServiceRequestType_value[k]; ok {
			req.RequestType = policyv1.ServiceRequestType(v)
		}
	}

	if statusStr.Valid {
		k := strings.ToUpper(statusStr.String)
		if v, ok := policyv1.ServiceRequestStatus_value[k]; ok {
			req.Status = policyv1.ServiceRequestStatus(v)
		}
	}

	if requestData.Valid {
		req.RequestData = requestData.String
	}

	if processedBy.Valid {
		req.ProcessedBy = processedBy.String
	}

	if processedAt.Valid {
		req.ProcessedAt = timestamppb.New(processedAt.Time)
	}

	req.CreatedAt = timestamppb.New(createdAt)

	return &req, nil
}

func (r *PolicyServiceRequestRepository) Update(ctx context.Context, request *policyv1.PolicyServiceRequest) (*policyv1.PolicyServiceRequest, error) {
	var processedBy sql.NullString
	if request.ProcessedBy != "" {
		processedBy = sql.NullString{String: request.ProcessedBy, Valid: true}
	}

	var processedAt sql.NullTime
	if request.ProcessedAt != nil {
		processedAt = sql.NullTime{Time: request.ProcessedAt.AsTime(), Valid: true}
	}

	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.policy_service_requests
		SET policy_id = $2,
		    customer_id = $3,
		    request_type = $4,
		    request_data = $5,
		    status = $6,
		    processed_by = $7,
		    processed_at = $8
		WHERE request_id = $1`,
		request.RequestId,
		request.PolicyId,
		request.CustomerId,
		strings.ToUpper(request.RequestType.String()),
		request.RequestData,
		strings.ToUpper(request.Status.String()),
		processedBy,
		processedAt,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update policy service request: %w", err)
	}

	return r.GetByID(ctx, request.RequestId)
}

func (r *PolicyServiceRequestRepository) Delete(ctx context.Context, requestID string) error {
	err := r.db.WithContext(ctx).Exec(`
		DELETE FROM insurance_schema.policy_service_requests
		WHERE request_id = $1`,
		requestID,
	).Error

	if err != nil {
		return fmt.Errorf("failed to delete policy service request: %w", err)
	}

	return nil
}

func (r *PolicyServiceRequestRepository) ListByPolicyID(ctx context.Context, policyID string) ([]*policyv1.PolicyServiceRequest, error) {
	query := `
		SELECT request_id, policy_id, customer_id, request_type, request_data, status, processed_by, processed_at, created_at
		FROM insurance_schema.policy_service_requests
		WHERE policy_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.WithContext(ctx).Raw(query, policyID).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to list policy service requests: %w", err)
	}
	defer rows.Close()

	requests := make([]*policyv1.PolicyServiceRequest, 0)
	for rows.Next() {
		var (
			req            policyv1.PolicyServiceRequest
			requestTypeStr sql.NullString
			statusStr      sql.NullString
			requestData    sql.NullString
			processedBy    sql.NullString
			processedAt    sql.NullTime
			createdAt      time.Time
		)

		err := rows.Scan(
			&req.RequestId,
			&req.PolicyId,
			&req.CustomerId,
			&requestTypeStr,
			&requestData,
			&statusStr,
			&processedBy,
			&processedAt,
			&createdAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan policy service request: %w", err)
		}

		if requestTypeStr.Valid {
			k := strings.ToUpper(requestTypeStr.String)
			if v, ok := policyv1.ServiceRequestType_value[k]; ok {
				req.RequestType = policyv1.ServiceRequestType(v)
			}
		}

		if statusStr.Valid {
			k := strings.ToUpper(statusStr.String)
			if v, ok := policyv1.ServiceRequestStatus_value[k]; ok {
				req.Status = policyv1.ServiceRequestStatus(v)
			}
		}

		if requestData.Valid {
			req.RequestData = requestData.String
		}

		if processedBy.Valid {
			req.ProcessedBy = processedBy.String
		}

		if processedAt.Valid {
			req.ProcessedAt = timestamppb.New(processedAt.Time)
		}

		req.CreatedAt = timestamppb.New(createdAt)

		requests = append(requests, &req)
	}

	return requests, nil
}
