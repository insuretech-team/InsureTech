package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	notificationv1 "github.com/newage-saint/insuretech/gen/go/insuretech/notification/entity/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestNotificationRepository_LiveDB_CRUDAndTransitions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testNotificationDB(t)

	userID := uuid.NewString()
	mobile := genValidMobile()
	insertUserMinimal(t, dbConn, userID, mobile, "notif_live_"+uuid.NewString()[:8]+"@example.com", "hash", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupNotificationUser(ctx, t, dbConn, userID) })

	repo := NewNotificationRepository(dbConn.Table("notification_schema.notifications"))

	now := time.Now().UTC()
	queuedID := uuid.NewString()
	sentID := uuid.NewString()

	require.NoError(t, repo.Create(ctx, &notificationv1.Notification{
		NotificationId: queuedID,
		RecipientId:    userID,
		Type:           notificationv1.NotificationType_NOTIFICATION_TYPE_PAYMENT_CONFIRMATION,
		Channel:        notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
		Subject:        "Payment received",
		Message:        "Your premium payment has been received.",
		TemplateData: map[string]string{
			"policy_number": "PL-TEST-001",
			"amount":        "1250.00",
		},
		Priority:    notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH,
		Status:      notificationv1.NotificationStatus_NOTIFICATION_STATUS_QUEUED,
		ScheduledAt: timestamppb.New(now.Add(-1 * time.Minute)),
		RetryCount:  0,
	}))

	require.NoError(t, repo.Create(ctx, &notificationv1.Notification{
		NotificationId: sentID,
		RecipientId:    userID,
		Type:           notificationv1.NotificationType_NOTIFICATION_TYPE_POLICY_ISSUED,
		Channel:        notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
		Subject:        "Policy issued",
		Message:        "Your policy is active now.",
		TemplateData: map[string]string{
			"policy_number": "PL-TEST-002",
		},
		Priority:   notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_NORMAL,
		Status:     notificationv1.NotificationStatus_NOTIFICATION_STATUS_SENT,
		SentAt:     timestamppb.New(now),
		RetryCount: 0,
	}))

	stored, err := repo.GetByID(ctx, queuedID)
	require.NoError(t, err)
	require.Equal(t, userID, stored.GetRecipientId())
	require.Equal(t, notificationv1.NotificationType_NOTIFICATION_TYPE_PAYMENT_CONFIRMATION, stored.GetType())
	require.Equal(t, "1250.00", stored.GetTemplateData()["amount"])

	due, err := repo.ListDue(ctx, now.Add(30*time.Second), 10)
	require.NoError(t, err)
	require.NotEmpty(t, due)

	foundQueued := false
	for _, item := range due {
		if item.GetNotificationId() == queuedID {
			foundQueued = true
			break
		}
	}
	require.True(t, foundQueued, "expected queued notification in due list")

	require.NoError(t, repo.MarkSending(ctx, queuedID))
	require.NoError(t, repo.MarkSent(ctx, queuedID, now.Add(10*time.Second)))
	require.NoError(t, repo.MarkDelivered(ctx, queuedID, now.Add(20*time.Second)))
	require.NoError(t, repo.MarkAsRead(ctx, []string{queuedID, sentID}))

	afterRead, err := repo.GetByID(ctx, queuedID)
	require.NoError(t, err)
	require.Equal(t, notificationv1.NotificationStatus_NOTIFICATION_STATUS_READ, afterRead.GetStatus())
	require.NotNil(t, afterRead.GetReadAt())

	require.NoError(t, repo.ScheduleRetry(ctx, queuedID, 1, now.Add(2*time.Minute), "temporary provider issue"))
	retried, err := repo.GetByID(ctx, queuedID)
	require.NoError(t, err)
	require.Equal(t, notificationv1.NotificationStatus_NOTIFICATION_STATUS_QUEUED, retried.GetStatus())
	require.EqualValues(t, 1, retried.GetRetryCount())

	require.NoError(t, repo.MarkFailed(ctx, queuedID, 3, "provider permanently rejected"))
	failed, err := repo.GetByID(ctx, queuedID)
	require.NoError(t, err)
	require.Equal(t, notificationv1.NotificationStatus_NOTIFICATION_STATUS_FAILED, failed.GetStatus())
	require.EqualValues(t, 3, failed.GetRetryCount())
	require.Contains(t, failed.GetErrorMessage(), "provider permanently rejected")

	list, total, unread, err := repo.ListByRecipient(ctx, userID, false, 20, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, total, int32(2))
	require.Len(t, list, 2)
	require.GreaterOrEqual(t, unread, int32(0))
}
