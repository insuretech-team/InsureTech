package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	mediav1 "github.com/newage-saint/insuretech/gen/go/insuretech/media/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProcessingJobRepository handles database operations for processing jobs.
type ProcessingJobRepository struct {
	db *sqlx.DB
}

var (
	ErrJobNotFound = errors.New("processing job not found")
)

// NewProcessingJobRepository creates a new processing job repository.
func NewProcessingJobRepository(db *sqlx.DB) *ProcessingJobRepository {
	return &ProcessingJobRepository{db: db}
}

// Create stores a new processing job.
func (r *ProcessingJobRepository) Create(ctx context.Context, job *mediav1.ProcessingJob) (*mediav1.ProcessingJob, error) {
	query := `
		INSERT INTO media_schema.processing_jobs (
			job_id, media_id, processing_type, status, priority,
			retry_count, max_retries, started_at, completed_at,
			error_message, result_data, audit_info, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW())
		RETURNING created_at, updated_at
	`

	var createdAt, updatedAt sql.NullTime
	var startedAt, completedAt sql.NullTime

	if job.StartedAt != nil {
		startedAt = sql.NullTime{Time: job.StartedAt.AsTime(), Valid: true}
	}
	if job.CompletedAt != nil {
		completedAt = sql.NullTime{Time: job.CompletedAt.AsTime(), Valid: true}
	}

	err := r.db.QueryRowContext(ctx, query,
		job.Id,
		job.MediaId,
		job.ProcessingType.String(),
		job.Status.String(),
		job.Priority,
		job.RetryCount,
		job.MaxRetries,
		startedAt,
		completedAt,
		nullString(job.ErrorMessage),
		nullString(job.ResultData),
		"{}",
	).Scan(&createdAt, &updatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create processing job: %w", err)
	}

	if job.AuditInfo == nil {
		job.AuditInfo = &commonv1.AuditInfo{}
	}
	if createdAt.Valid {
		job.AuditInfo.CreatedAt = timestamppb.New(createdAt.Time)
	}
	if updatedAt.Valid {
		job.AuditInfo.UpdatedAt = timestamppb.New(updatedAt.Time)
	}

	return job, nil
}

// GetByID retrieves a processing job by ID. If tenantID is provided, access is tenant-scoped.
func (r *ProcessingJobRepository) GetByID(ctx context.Context, tenantID, jobID string) (*mediav1.ProcessingJob, error) {
	query := `
		SELECT
			pj.job_id, pj.media_id, pj.processing_type, pj.status, pj.priority,
			pj.retry_count, pj.max_retries, pj.started_at, pj.completed_at,
			pj.error_message, pj.result_data, pj.created_at, pj.updated_at
		FROM media_schema.processing_jobs pj
		JOIN media_schema.media_files mf ON mf.media_id = pj.media_id
		WHERE pj.job_id = $1`
	args := []any{jobID}
	if strings.TrimSpace(tenantID) != "" {
		query += " AND mf.tenant_id = $2"
		args = append(args, tenantID)
	}

	job, err := scanProcessingJob(func(dest ...any) error {
		return r.db.QueryRowContext(ctx, query, args...).Scan(dest...)
	})
	if err == sql.ErrNoRows {
		return nil, ErrJobNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get processing job: %w", err)
	}
	return job, nil
}

// ListByMediaID retrieves all processing jobs for a media file.
func (r *ProcessingJobRepository) ListByMediaID(ctx context.Context, tenantID, mediaID string) ([]*mediav1.ProcessingJob, error) {
	jobs, _, err := r.List(ctx, tenantID, mediaID, nil, nil, 1000, 0)
	return jobs, err
}

// List retrieves processing jobs with optional filters and pagination.
func (r *ProcessingJobRepository) List(
	ctx context.Context,
	tenantID, mediaID string,
	processingType *mediav1.ProcessingType,
	status *mediav1.ProcessingStatus,
	limit, offset int,
) ([]*mediav1.ProcessingJob, int, error) {
	conditions := []string{"1=1"}
	args := make([]any, 0, 6)

	if strings.TrimSpace(tenantID) != "" {
		conditions = append(conditions, fmt.Sprintf("mf.tenant_id = $%d", len(args)+1))
		args = append(args, tenantID)
	}
	if strings.TrimSpace(mediaID) != "" {
		conditions = append(conditions, fmt.Sprintf("pj.media_id = $%d", len(args)+1))
		args = append(args, mediaID)
	}
	if processingType != nil {
		conditions = append(conditions, fmt.Sprintf("pj.processing_type = $%d", len(args)+1))
		args = append(args, processingType.String())
	}
	if status != nil {
		conditions = append(conditions, fmt.Sprintf("pj.status = $%d", len(args)+1))
		args = append(args, status.String())
	}

	whereClause := strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM media_schema.processing_jobs pj
		JOIN media_schema.media_files mf ON mf.media_id = pj.media_id
		WHERE %s
	`, whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count processing jobs: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT
			pj.job_id, pj.media_id, pj.processing_type, pj.status, pj.priority,
			pj.retry_count, pj.max_retries, pj.started_at, pj.completed_at,
			pj.error_message, pj.result_data, pj.created_at, pj.updated_at
		FROM media_schema.processing_jobs pj
		JOIN media_schema.media_files mf ON mf.media_id = pj.media_id
		WHERE %s
		ORDER BY pj.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, len(args)+1, len(args)+2)
	queryArgs := append(append([]any{}, args...), limit, offset)

	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list processing jobs: %w", err)
	}
	defer rows.Close()

	jobs := make([]*mediav1.ProcessingJob, 0)
	for rows.Next() {
		job, scanErr := scanProcessingJob(rows.Scan)
		if scanErr != nil {
			return nil, 0, fmt.Errorf("failed to scan processing job: %w", scanErr)
		}
		jobs = append(jobs, job)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating processing jobs: %w", err)
	}

	return jobs, total, nil
}

// GetNextPendingJob retrieves the next pending job by priority.
func (r *ProcessingJobRepository) GetNextPendingJob(ctx context.Context, processingType mediav1.ProcessingType) (*mediav1.ProcessingJob, error) {
	query := `
		SELECT
			job_id, media_id, processing_type, status, priority,
			retry_count, max_retries, started_at, completed_at,
			error_message, result_data, created_at, updated_at
		FROM media_schema.processing_jobs
		WHERE status = 'PENDING' AND processing_type = $1
		ORDER BY priority DESC, created_at ASC
		LIMIT 1
		FOR UPDATE SKIP LOCKED
	`

	job, err := scanProcessingJob(func(dest ...any) error {
		return r.db.QueryRowContext(ctx, query, processingType.String()).Scan(dest...)
	})
	if err == sql.ErrNoRows {
		return nil, nil // No pending jobs.
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get next pending job: %w", err)
	}
	return job, nil
}

// UpdateStatus updates the status of a processing job.
func (r *ProcessingJobRepository) UpdateStatus(ctx context.Context, jobID string, status mediav1.ProcessingStatus) error {
	query := `
		UPDATE media_schema.processing_jobs
		SET status = $1, updated_at = NOW()
		WHERE job_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, status.String(), jobID)
	if err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrJobNotFound
	}

	return nil
}

// MarkAsStarted marks a job as in progress.
func (r *ProcessingJobRepository) MarkAsStarted(ctx context.Context, jobID string) error {
	query := `
		UPDATE media_schema.processing_jobs
		SET status = 'IN_PROGRESS', started_at = NOW(), updated_at = NOW()
		WHERE job_id = $1
	`

	result, err := r.db.ExecContext(ctx, query, jobID)
	if err != nil {
		return fmt.Errorf("failed to mark job as started: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrJobNotFound
	}

	return nil
}

// MarkAsCompleted marks a job as completed with result data.
func (r *ProcessingJobRepository) MarkAsCompleted(ctx context.Context, jobID string, resultData string) error {
	query := `
		UPDATE media_schema.processing_jobs
		SET status = 'COMPLETED', completed_at = NOW(), result_data = $1, updated_at = NOW()
		WHERE job_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, resultData, jobID)
	if err != nil {
		return fmt.Errorf("failed to mark job as completed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrJobNotFound
	}

	return nil
}

// MarkAsFailed marks a job as failed with error message and increments retry count.
func (r *ProcessingJobRepository) MarkAsFailed(ctx context.Context, jobID string, errorMsg string) error {
	query := `
		UPDATE media_schema.processing_jobs
		SET status = 'FAILED', error_message = $1, retry_count = retry_count + 1, updated_at = NOW()
		WHERE job_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, errorMsg, jobID)
	if err != nil {
		return fmt.Errorf("failed to mark job as failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrJobNotFound
	}

	return nil
}

func scanProcessingJob(scan func(dest ...any) error) (*mediav1.ProcessingJob, error) {
	var job mediav1.ProcessingJob
	var processingType, status string
	var startedAt, completedAt, createdAt, updatedAt sql.NullTime
	var errorMessage, resultData sql.NullString

	err := scan(
		&job.Id,
		&job.MediaId,
		&processingType,
		&status,
		&job.Priority,
		&job.RetryCount,
		&job.MaxRetries,
		&startedAt,
		&completedAt,
		&errorMessage,
		&resultData,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	job.ProcessingType = parseProcessingType(processingType)
	job.Status = parseProcessingStatus(status)

	if startedAt.Valid {
		job.StartedAt = timestamppb.New(startedAt.Time)
	}
	if completedAt.Valid {
		job.CompletedAt = timestamppb.New(completedAt.Time)
	}
	if errorMessage.Valid {
		job.ErrorMessage = errorMessage.String
	}
	if resultData.Valid {
		job.ResultData = resultData.String
	}

	job.AuditInfo = &commonv1.AuditInfo{}
	if createdAt.Valid {
		job.AuditInfo.CreatedAt = timestamppb.New(createdAt.Time)
	}
	if updatedAt.Valid {
		job.AuditInfo.UpdatedAt = timestamppb.New(updatedAt.Time)
	}

	return &job, nil
}

// Helper functions.
func parseProcessingType(s string) mediav1.ProcessingType {
	if val, ok := mediav1.ProcessingType_value[s]; ok {
		return mediav1.ProcessingType(val)
	}
	return mediav1.ProcessingType_PROCESSING_TYPE_UNSPECIFIED
}

func parseProcessingStatus(s string) mediav1.ProcessingStatus {
	if val, ok := mediav1.ProcessingStatus_value[s]; ok {
		return mediav1.ProcessingStatus(val)
	}
	return mediav1.ProcessingStatus_PROCESSING_STATUS_UNSPECIFIED
}

// GetPendingJobs retrieves pending jobs for worker processing
// Implements the JobUpdater interface required by ProcessingWorker
type ProcessingJobRecord struct {
	ID       string
	MediaID  string
	TenantID string
	JobType  string // THUMBNAIL, OPTIMIZATION, OCR, VIRUS_SCAN
	Priority int
}

// GetPendingJobs retrieves a batch of pending jobs
func (r *ProcessingJobRepository) GetPendingJobs(ctx context.Context, limit int) ([]*ProcessingJobRecord, error) {
	query := `
		SELECT
			pj.job_id, pj.media_id, mf.tenant_id, pj.processing_type, pj.priority
		FROM media_schema.processing_jobs pj
		JOIN media_schema.media_files mf ON mf.media_id = pj.media_id
		WHERE pj.status = 'PENDING'
		ORDER BY pj.priority DESC, pj.created_at ASC
		LIMIT $1
		FOR UPDATE SKIP LOCKED
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending jobs: %w", err)
	}
	defer rows.Close()

	var jobs []*ProcessingJobRecord
	for rows.Next() {
		var job ProcessingJobRecord
		var jobType string
		if err := rows.Scan(&job.ID, &job.MediaID, &job.TenantID, &jobType, &job.Priority); err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}
		job.JobType = jobType
		jobs = append(jobs, &job)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating pending jobs: %w", err)
	}

	return jobs, nil
}

// MarkJobStarted marks a job as started (wrapper for worker compatibility)
func (r *ProcessingJobRepository) MarkJobStarted(ctx context.Context, jobID string) error {
	return r.MarkAsStarted(ctx, jobID)
}

// MarkJobCompleted marks a job as completed (wrapper for worker compatibility)
func (r *ProcessingJobRepository) MarkJobCompleted(ctx context.Context, jobID string, result string) error {
	return r.MarkAsCompleted(ctx, jobID, result)
}

// MarkJobFailed marks a job as failed (wrapper for worker compatibility)
func (r *ProcessingJobRepository) MarkJobFailed(ctx context.Context, jobID string, errMsg string) error {
	return r.MarkAsFailed(ctx, jobID, errMsg)
}
