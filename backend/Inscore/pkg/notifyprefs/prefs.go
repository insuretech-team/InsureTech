package notifyprefs

import (
	"fmt"
	"slices"
	"sort"
	"strings"
	"time"

	notificationv1 "github.com/newage-saint/insuretech/gen/go/insuretech/notification/entity/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var channelOrder = []notificationv1.NotificationChannel{
	notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
	notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS,
	notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_PUSH,
	notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
	notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_WHATSAPP,
}

func Default(userID string) *notificationv1.NotificationPreference {
	return &notificationv1.NotificationPreference{
		UserId: userID,
		ChannelPreferences: []*notificationv1.ChannelPreference{
			{
				Channel: notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL,
				Enabled: true,
			},
			{
				Channel: notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS,
				Enabled: true,
			},
			{
				Channel: notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP,
				Enabled: true,
			},
		},
		MarketingOptIn:     false,
		TransactionalOptIn: true,
		UpdatedAt:          timestamppb.New(time.Now()),
	}
}

// Parse converts either JSON or legacy compact preference strings into proto form.
func Parse(userID, raw string) (*notificationv1.NotificationPreference, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return Default(userID), nil
	}

	if strings.HasPrefix(raw, "{") {
		prefs := &notificationv1.NotificationPreference{}
		if err := protojson.Unmarshal([]byte(raw), prefs); err != nil {
			return nil, fmt.Errorf("decode notification preferences: %w", err)
		}
		if prefs.GetUserId() == "" {
			prefs.UserId = userID
		}
		if prefs.GetUpdatedAt() == nil {
			prefs.UpdatedAt = timestamppb.New(time.Now())
		}
		return prefs, nil
	}

	return parseLegacy(userID, raw), nil
}

// Compact encodes notification preferences into the shared compact string format used by authn rows.
func Compact(preferences *notificationv1.NotificationPreference) string {
	if preferences == nil {
		return "ALL"
	}

	tokens := make([]string, 0, len(preferences.GetChannelPreferences()))
	for _, pref := range preferences.GetChannelPreferences() {
		if pref == nil || !pref.GetEnabled() {
			continue
		}
		token, ok := tokenForChannel(pref.GetChannel())
		if !ok {
			continue
		}
		tokens = append(tokens, token)
	}
	if len(tokens) == 0 {
		return fmt.Sprintf("NONE;MKT=%d;TXN=%d", boolToInt(preferences.GetMarketingOptIn()), boolToInt(preferences.GetTransactionalOptIn()))
	}

	sort.Strings(tokens)
	tokens = slices.Compact(tokens)
	return fmt.Sprintf("%s;MKT=%d;TXN=%d", strings.Join(tokens, ","), boolToInt(preferences.GetMarketingOptIn()), boolToInt(preferences.GetTransactionalOptIn()))
}

func parseLegacy(userID, raw string) *notificationv1.NotificationPreference {
	normalized := strings.ToUpper(strings.TrimSpace(raw))
	switch normalized {
	case "", "DEFAULT":
		return Default(userID)
	case "NONE":
		return build(userID, map[notificationv1.NotificationChannel]bool{}, false, true)
	case "ALL":
		allChannels := make(map[notificationv1.NotificationChannel]bool, len(channelOrder))
		for _, channel := range channelOrder {
			allChannels[channel] = true
		}
		return build(userID, allChannels, false, true)
	}

	enabled := map[notificationv1.NotificationChannel]bool{}
	marketingOptIn := false
	transactionalOptIn := true

	segments := strings.Split(normalized, ";")
	for _, token := range strings.FieldsFunc(segments[0], func(r rune) bool {
		return r == ',' || r == '|' || r == ' '
	}) {
		if channel, ok := channelFromToken(token); ok {
			enabled[channel] = true
		}
	}
	for _, segment := range segments[1:] {
		switch strings.TrimSpace(segment) {
		case "MKT=1":
			marketingOptIn = true
		case "MKT=0":
			marketingOptIn = false
		case "TXN=1":
			transactionalOptIn = true
		case "TXN=0":
			transactionalOptIn = false
		}
	}

	if len(enabled) == 0 {
		prefs := Default(userID)
		prefs.MarketingOptIn = marketingOptIn
		prefs.TransactionalOptIn = transactionalOptIn
		return prefs
	}

	return build(userID, enabled, marketingOptIn, transactionalOptIn)
}

func build(userID string, enabled map[notificationv1.NotificationChannel]bool, marketingOptIn, transactionalOptIn bool) *notificationv1.NotificationPreference {
	prefs := &notificationv1.NotificationPreference{
		UserId:             userID,
		ChannelPreferences: make([]*notificationv1.ChannelPreference, 0, len(channelOrder)),
		MarketingOptIn:     marketingOptIn,
		TransactionalOptIn: transactionalOptIn,
		UpdatedAt:          timestamppb.New(time.Now()),
	}
	for _, channel := range channelOrder {
		prefs.ChannelPreferences = append(prefs.ChannelPreferences, &notificationv1.ChannelPreference{
			Channel: channel,
			Enabled: enabled[channel],
		})
	}
	return prefs
}

func channelFromToken(token string) (notificationv1.NotificationChannel, bool) {
	switch strings.ToUpper(strings.TrimSpace(token)) {
	case "SMS":
		return notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS, true
	case "EMAIL":
		return notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL, true
	case "PUSH":
		return notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_PUSH, true
	case "IN_APP":
		return notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP, true
	case "WHATSAPP":
		return notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_WHATSAPP, true
	default:
		return notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_UNSPECIFIED, false
	}
}

func tokenForChannel(channel notificationv1.NotificationChannel) (string, bool) {
	switch channel {
	case notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_SMS:
		return "SMS", true
	case notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL:
		return "EMAIL", true
	case notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_PUSH:
		return "PUSH", true
	case notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP:
		return "IN_APP", true
	case notificationv1.NotificationChannel_NOTIFICATION_CHANNEL_WHATSAPP:
		return "WHATSAPP", true
	default:
		return "", false
	}
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
