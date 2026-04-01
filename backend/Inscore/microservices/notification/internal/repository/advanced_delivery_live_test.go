package repository

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPushTokenRepository_LiveDB(t *testing.T) {
	dbConn := testNotificationDB(t)
	ctx := context.Background()

	userID := uuid.NewString()
	cleanupNotificationUser(ctx, t, dbConn, userID)
	insertUserMinimal(t, dbConn, userID, genValidMobile(), "push-"+userID[:8]+"@example.com", "hash", 1)
	defer cleanupNotificationUser(ctx, t, dbConn, userID)

	repo := NewPushTokenRepository(dbConn)
	lastSeen := time.Now().UTC().Add(-time.Minute)
	require.NoError(t, repo.Upsert(ctx, &PushDeviceToken{
		UserID:      userID,
		Provider:    "FCM",
		Platform:    "android",
		DeviceToken: "token-" + uuid.NewString(),
		DeviceID:    "device-" + uuid.NewString()[:8],
		AppID:       "customer-app",
		IsActive:    true,
		LastSeenAt:  &lastSeen,
	}))

	tokens, err := repo.ListActiveByUserID(ctx, userID)
	require.NoError(t, err)
	require.Len(t, tokens, 1)
	require.Equal(t, "FCM", tokens[0].Provider)
	require.Equal(t, "android", tokens[0].Platform)

	require.NoError(t, repo.DeactivateByDeviceTokens(ctx, []string{tokens[0].DeviceToken}))

	tokens, err = repo.ListActiveByUserID(ctx, userID)
	require.NoError(t, err)
	require.Empty(t, tokens)
}

func TestWebhookRepository_LiveDB(t *testing.T) {
	dbConn := testNotificationDB(t)
	ctx := context.Background()

	repo := NewWebhookRepository(dbConn)
	subscriptionID := uuid.NewString()
	defer cleanupWebhookSubscription(ctx, t, dbConn, subscriptionID)

	require.NoError(t, repo.CreateSubscription(ctx, &WebhookSubscription{
		SubscriptionID: subscriptionID,
		SubscriberName: "partner-sync",
		TargetURL:      "https://example.com/hooks/notifications",
		Secret:         "secret-123",
		EventTypes:     []string{"notification.sent"},
		TopicGroups:    []string{"customer_identity"},
		Topics:         []string{"authn.user.registered"},
		Channels:       []string{"email"},
		TimeoutSeconds: 9,
		MaxAttempts:    4,
		IsActive:       true,
	}))

	subscription, err := repo.GetSubscriptionByID(ctx, subscriptionID)
	require.NoError(t, err)
	require.Equal(t, "partner-sync", subscription.SubscriberName)

	matches, err := repo.ListMatchingSubscriptions(ctx, "notification.sent", "authn.user.registered", "customer_identity", "EMAIL")
	require.NoError(t, err)
	require.Len(t, matches, 1)

	payload := json.RawMessage(`{"notification_id":"n-1","event_type":"NOTIFICATION.SENT"}`)
	attemptID := uuid.NewString()
	require.NoError(t, repo.EnqueueAttempt(ctx, &WebhookDeliveryAttempt{
		AttemptID:      attemptID,
		SubscriptionID: subscriptionID,
		NotificationID: "",
		LifecycleEvent: "notification.sent",
		SourceTopic:    "authn.user.registered",
		Payload:        payload,
		Status:         "QUEUED",
	}))

	attempts, err := repo.ListDueAttempts(ctx, time.Now().UTC().Add(time.Second), 10)
	require.NoError(t, err)
	found := false
	for _, attempt := range attempts {
		if attempt != nil && attempt.AttemptID == attemptID {
			found = true
			break
		}
	}
	require.True(t, found)

	require.NoError(t, repo.MarkAttemptSending(ctx, attemptID))
	require.NoError(t, repo.ScheduleAttemptRetry(ctx, attemptID, 1, time.Now().UTC().Add(time.Minute), 503, "retry later", "temporary upstream error"))

	attempt, err := repo.GetAttemptByID(ctx, attemptID)
	require.NoError(t, err)
	require.Equal(t, 1, attempt.RetryCount)
	require.Equal(t, "QUEUED", attempt.Status)
	require.Equal(t, 503, attempt.ResponseStatus)

	require.NoError(t, repo.MarkAttemptFailed(ctx, attemptID, 4, 400, "bad request", "permanent failure", time.Now().UTC()))
	attempt, err = repo.GetAttemptByID(ctx, attemptID)
	require.NoError(t, err)
	require.Equal(t, "FAILED", attempt.Status)
	require.Equal(t, 4, attempt.RetryCount)
	require.Equal(t, "permanent failure", attempt.ErrorMessage)
}
