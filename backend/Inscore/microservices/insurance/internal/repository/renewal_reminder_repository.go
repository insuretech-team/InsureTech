package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"

	renewalv1 "github.com/newage-saint/insuretech/gen/go/insuretech/renewal/entity/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
)

type RenewalReminderRepository struct {
	db *gorm.DB
}

func NewRenewalReminderRepository(db *gorm.DB) *RenewalReminderRepository {
	return &RenewalReminderRepository{db: db}
}

func (r *RenewalReminderRepository) Create(ctx context.Context, reminder *renewalv1.RenewalReminder) (*renewalv1.RenewalReminder, error) {
	if reminder.Id == "" {
		return nil, fmt.Errorf("reminder_id is required")
	}

	// Handle timestamps
	var scheduledAt time.Time
	if reminder.ScheduledAt != nil {
		scheduledAt = reminder.ScheduledAt.AsTime()
	}

	var sentAt sql.NullTime
	if reminder.SentAt != nil {
		sentAt = sql.NullTime{Time: reminder.SentAt.AsTime(), Valid: true}
	}

	var notificationID sql.NullString
	if reminder.NotificationId != "" {
		notificationID = sql.NullString{String: reminder.NotificationId, Valid: true}
	}

	var auditInfo interface{}
	if reminder.AuditInfo != nil {
		auditInfo = "{}"
	}

	err := r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.renewal_reminders
			(reminder_id, renewal_schedule_id, days_before_renewal, channel, status,
			 scheduled_at, sent_at, notification_id, audit_info)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		reminder.Id,
		reminder.RenewalScheduleId,
		reminder.DaysBeforeRenewal,
		strings.ToUpper(reminder.Channel.String()),
		strings.ToUpper(reminder.Status.String()),
		scheduledAt,
		sentAt,
		notificationID,
		auditInfo,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert renewal reminder: %w", err)
	}

	return r.GetByID(ctx, reminder.Id)
}

func (r *RenewalReminderRepository) GetByID(ctx context.Context, reminderID string) (*renewalv1.RenewalReminder, error) {
	var (
		rem            renewalv1.RenewalReminder
		channelStr     sql.NullString
		statusStr      sql.NullString
		scheduledAt    time.Time
		sentAt         sql.NullTime
		notificationID sql.NullString
		auditInfo      sql.NullString
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT reminder_id, renewal_schedule_id, days_before_renewal, channel, status,
		       scheduled_at, sent_at, notification_id, audit_info
		FROM insurance_schema.renewal_reminders
		WHERE reminder_id = $1`,
		reminderID,
	).Row().Scan(
		&rem.Id,
		&rem.RenewalScheduleId,
		&rem.DaysBeforeRenewal,
		&channelStr,
		&statusStr,
		&scheduledAt,
		&sentAt,
		&notificationID,
		&auditInfo,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get renewal reminder: %w", err)
	}

	// Set optional fields
	if notificationID.Valid {
		rem.NotificationId = notificationID.String
	}

	// Parse enums
	if channelStr.Valid {
		k := strings.ToUpper(channelStr.String)
		if v, ok := renewalv1.ReminderChannel_value[k]; ok {
			rem.Channel = renewalv1.ReminderChannel(v)
		}
	}
	if statusStr.Valid {
		k := strings.ToUpper(statusStr.String)
		if v, ok := renewalv1.ReminderStatus_value[k]; ok {
			rem.Status = renewalv1.ReminderStatus(v)
		}
	}

	// Set timestamps
	if !scheduledAt.IsZero() {
		rem.ScheduledAt = timestamppb.New(scheduledAt)
	}
	if sentAt.Valid {
		rem.SentAt = timestamppb.New(sentAt.Time)
	}

	// Set audit info
	if auditInfo.Valid {
		rem.AuditInfo = &commonv1.AuditInfo{}
	}

	return &rem, nil
}

func (r *RenewalReminderRepository) Update(ctx context.Context, reminder *renewalv1.RenewalReminder) (*renewalv1.RenewalReminder, error) {
	// Handle timestamps
	var scheduledAt time.Time
	if reminder.ScheduledAt != nil {
		scheduledAt = reminder.ScheduledAt.AsTime()
	}

	var sentAt sql.NullTime
	if reminder.SentAt != nil {
		sentAt = sql.NullTime{Time: reminder.SentAt.AsTime(), Valid: true}
	}

	var notificationID sql.NullString
	if reminder.NotificationId != "" {
		notificationID = sql.NullString{String: reminder.NotificationId, Valid: true}
	}

	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.renewal_reminders
		SET renewal_schedule_id = $2,
		    days_before_renewal = $3,
		    channel = $4,
		    status = $5,
		    scheduled_at = $6,
		    sent_at = $7,
		    notification_id = $8
		WHERE reminder_id = $1`,
		reminder.Id,
		reminder.RenewalScheduleId,
		reminder.DaysBeforeRenewal,
		strings.ToUpper(reminder.Channel.String()),
		strings.ToUpper(reminder.Status.String()),
		scheduledAt,
		sentAt,
		notificationID,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update renewal reminder: %w", err)
	}

	return r.GetByID(ctx, reminder.Id)
}

func (r *RenewalReminderRepository) Delete(ctx context.Context, reminderID string) error {
	err := r.db.WithContext(ctx).Exec(`
		DELETE FROM insurance_schema.renewal_reminders
		WHERE reminder_id = $1`,
		reminderID,
	).Error

	if err != nil {
		return fmt.Errorf("failed to delete renewal reminder: %w", err)
	}

	return nil
}

func (r *RenewalReminderRepository) ListByScheduleID(ctx context.Context, scheduleID string) ([]*renewalv1.RenewalReminder, error) {
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT reminder_id, renewal_schedule_id, days_before_renewal, channel, status,
		       scheduled_at, sent_at, notification_id, audit_info
		FROM insurance_schema.renewal_reminders
		WHERE renewal_schedule_id = $1
		ORDER BY scheduled_at DESC`,
		scheduleID,
	).Rows()

	if err != nil {
		return nil, fmt.Errorf("failed to list renewal reminders: %w", err)
	}
	defer rows.Close()

	reminders := make([]*renewalv1.RenewalReminder, 0)
	for rows.Next() {
		var (
			rem            renewalv1.RenewalReminder
			channelStr     sql.NullString
			statusStr      sql.NullString
			scheduledAt    time.Time
			sentAt         sql.NullTime
			notificationID sql.NullString
			auditInfo      sql.NullString
		)

		err := rows.Scan(
			&rem.Id,
			&rem.RenewalScheduleId,
			&rem.DaysBeforeRenewal,
			&channelStr,
			&statusStr,
			&scheduledAt,
			&sentAt,
			&notificationID,
			&auditInfo,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan renewal reminder: %w", err)
		}

		// Set optional fields
		if notificationID.Valid {
			rem.NotificationId = notificationID.String
		}

		// Parse enums
		if channelStr.Valid {
			k := strings.ToUpper(channelStr.String)
			if v, ok := renewalv1.ReminderChannel_value[k]; ok {
				rem.Channel = renewalv1.ReminderChannel(v)
			}
		}
		if statusStr.Valid {
			k := strings.ToUpper(statusStr.String)
			if v, ok := renewalv1.ReminderStatus_value[k]; ok {
				rem.Status = renewalv1.ReminderStatus(v)
			}
		}

		// Set timestamps
		if !scheduledAt.IsZero() {
			rem.ScheduledAt = timestamppb.New(scheduledAt)
		}
		if sentAt.Valid {
			rem.SentAt = timestamppb.New(sentAt.Time)
		}

		// Set audit info
		if auditInfo.Valid {
			rem.AuditInfo = &commonv1.AuditInfo{}
		}

		reminders = append(reminders, &rem)
	}

	return reminders, nil
}
