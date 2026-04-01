package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	notificationv1 "github.com/newage-saint/insuretech/gen/go/insuretech/notification/entity/v1"
	"github.com/stretchr/testify/require"
)

func TestTemplateRepository_LiveDB_CRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testNotificationDB(t)
	repo := NewTemplateRepository(dbConn.Table("notification_schema.notification_templates"))

	templateID := uuid.NewString()
	templateName := "live_template_" + uuid.NewString()[:8]
	t.Cleanup(func() { cleanupNotificationTemplate(ctx, t, dbConn, templateID) })

	require.NoError(t, repo.Create(ctx, &notificationv1.NotificationTemplate{
		TemplateId:      templateID,
		TemplateName:    templateName,
		Type:            notificationv1.NotificationType_NOTIFICATION_TYPE_RENEWAL_REMINDER,
		Channel:         notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
		SubjectTemplate: "Renewal due soon",
		BodyTemplate:    "Dear {{name}}, your policy {{policy_number}} will expire soon.",
		Language:        "en",
		IsActive:        true,
	}))

	stored, err := repo.GetByID(ctx, templateID)
	require.NoError(t, err)
	require.Equal(t, templateName, stored.GetTemplateName())
	require.Equal(t, notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL, stored.GetChannel())

	require.NoError(t, repo.Update(ctx, templateID, templateName+"_v2", "Renewal due", "Updated body"))
	updated, err := repo.GetByID(ctx, templateID)
	require.NoError(t, err)
	require.Equal(t, templateName+"_v2", updated.GetTemplateName())
	require.Equal(t, "Renewal due", updated.GetSubjectTemplate())
	require.Equal(t, "Updated body", updated.GetBodyTemplate())

	require.NoError(t, repo.Deactivate(ctx, templateID))
	deactivated, err := repo.GetByID(ctx, templateID)
	require.NoError(t, err)
	require.False(t, deactivated.GetIsActive())
}

func TestUserRepository_LiveDB_Preferences(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testNotificationDB(t)

	userID := uuid.NewString()
	mobile := genValidMobile()
	insertUserMinimal(t, dbConn, userID, mobile, "notif_pref_"+uuid.NewString()[:8]+"@example.com", "hash", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupNotificationUser(ctx, t, dbConn, userID) })

	repo := NewUserRepository(dbConn.Table("authn_schema.users"))

	defaults, err := repo.GetPreferences(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, userID, defaults.GetUserId())
	require.True(t, defaults.GetTransactionalOptIn())
	require.False(t, defaults.GetMarketingOptIn())
	require.True(t, channelEnabled(defaults, notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL))
	require.True(t, channelEnabled(defaults, notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS))
	require.True(t, channelEnabled(defaults, notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP))

	require.NoError(t, dbConn.Exec(`
		UPDATE authn_schema.users
		SET notification_preference = $2,
		    updated_at = NOW()
		WHERE user_id = $1 AND deleted_at IS NULL`,
		userID, "EMAIL,IN_APP,PUSH;MKT=1;TXN=1",
	).Error)

	updated, err := repo.GetPreferences(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, userID, updated.GetUserId())
	require.True(t, updated.GetMarketingOptIn())
	require.True(t, channelEnabled(updated, notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL))
	require.False(t, channelEnabled(updated, notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS))
	require.True(t, channelEnabled(updated, notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_PUSH))
	require.True(t, channelEnabled(updated, notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP))

	user, err := repo.GetByID(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, userID, user.GetUserId())
	require.Equal(t, mobile, user.GetMobileNumber())
}

func channelEnabled(prefs *notificationv1.NotificationPreference, channel notificationv1.NotificationChannel) bool {
	for _, pref := range prefs.GetChannelPreferences() {
		if pref.GetChannel() == channel {
			return pref.GetEnabled()
		}
	}
	return false
}
