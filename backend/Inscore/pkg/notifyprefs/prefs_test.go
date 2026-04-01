package notifyprefs

import (
	"testing"

	notificationv1 "github.com/newage-saint/insuretech/gen/go/insuretech/notification/entity/v1"
)

func TestParseDefaultAndCompact(t *testing.T) {
	prefs, err := Parse("user-1", "")
	if err != nil {
		t.Fatalf("Parse(default): %v", err)
	}
	if !prefs.GetTransactionalOptIn() || prefs.GetMarketingOptIn() {
		t.Fatalf("unexpected default flags: %+v", prefs)
	}

	encoded := Compact(&notificationv1.NotificationPreference{
		ChannelPreferences: []*notificationv1.ChannelPreference{
			{Channel: notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL, Enabled: true},
			{Channel: notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_PUSH, Enabled: true},
			{Channel: notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL, Enabled: true},
		},
		MarketingOptIn:     true,
		TransactionalOptIn: false,
	})
	if encoded != "EMAIL,PUSH;MKT=1;TXN=0" {
		t.Fatalf("Compact() = %q", encoded)
	}
}

func TestParseLegacy(t *testing.T) {
	prefs, err := Parse("user-1", "EMAIL,PUSH;MKT=1;TXN=0")
	if err != nil {
		t.Fatalf("Parse(legacy): %v", err)
	}
	if !prefs.GetMarketingOptIn() || prefs.GetTransactionalOptIn() {
		t.Fatalf("unexpected flags: %+v", prefs)
	}
	var emailEnabled, pushEnabled bool
	for _, pref := range prefs.GetChannelPreferences() {
		switch pref.GetChannel() {
		case notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL:
			emailEnabled = pref.GetEnabled()
		case notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_PUSH:
			pushEnabled = pref.GetEnabled()
		}
	}
	if !emailEnabled || !pushEnabled {
		t.Fatalf("expected email and push enabled: %+v", prefs)
	}
}
