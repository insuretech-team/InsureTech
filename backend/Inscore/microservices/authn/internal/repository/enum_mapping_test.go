package repository

import (
	"testing"

	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// These tests verify that the proto-generated enum value maps cover all values
// used by the auth service. Since the proto_enum GORM serializer uses the
// enum's String() / value map for DB round-trips, these tests act as a
// contract between the proto definition and the auth service logic.
// ---------------------------------------------------------------------------

func TestUserType_ProtoValueMap_CoversAllExpectedValues(t *testing.T) {
	cases := []struct {
		name  string
		value authnentityv1.UserType
		str   string
	}{
		{"B2C_CUSTOMER", authnentityv1.UserType_USER_TYPE_B2C_CUSTOMER, "USER_TYPE_B2C_CUSTOMER"},
		{"AGENT", authnentityv1.UserType_USER_TYPE_AGENT, "USER_TYPE_AGENT"},
		{"BUSINESS_BENEFICIARY", authnentityv1.UserType_USER_TYPE_BUSINESS_BENEFICIARY, "USER_TYPE_BUSINESS_BENEFICIARY"},
		{"SYSTEM_USER", authnentityv1.UserType_USER_TYPE_SYSTEM_USER, "USER_TYPE_SYSTEM_USER"},
		{"UNSPECIFIED", authnentityv1.UserType_USER_TYPE_UNSPECIFIED, "USER_TYPE_UNSPECIFIED"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// String() → proto name
			require.Equal(t, tc.str, tc.value.String())
			// value map round-trip: name → int32 → enum
			v, ok := authnentityv1.UserType_value[tc.str]
			require.True(t, ok, "value map missing key %s", tc.str)
			require.Equal(t, tc.value, authnentityv1.UserType(v))
			// name map round-trip: int32 → name
			n, ok := authnentityv1.UserType_name[int32(tc.value)]
			require.True(t, ok, "name map missing value %d", tc.value)
			require.Equal(t, tc.str, n)
		})
	}
}

func TestUserStatus_ProtoValueMap_CoversAllExpectedValues(t *testing.T) {
	cases := []struct {
		name  string
		value authnentityv1.UserStatus
		str   string
	}{
		{"UNSPECIFIED", authnentityv1.UserStatus_USER_STATUS_UNSPECIFIED, "USER_STATUS_UNSPECIFIED"},
		{"PENDING_VERIFICATION", authnentityv1.UserStatus_USER_STATUS_PENDING_VERIFICATION, "USER_STATUS_PENDING_VERIFICATION"},
		{"ACTIVE", authnentityv1.UserStatus_USER_STATUS_ACTIVE, "USER_STATUS_ACTIVE"},
		{"SUSPENDED", authnentityv1.UserStatus_USER_STATUS_SUSPENDED, "USER_STATUS_SUSPENDED"},
		{"LOCKED", authnentityv1.UserStatus_USER_STATUS_LOCKED, "USER_STATUS_LOCKED"},
		{"DELETED", authnentityv1.UserStatus_USER_STATUS_DELETED, "USER_STATUS_DELETED"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.str, tc.value.String())
			v, ok := authnentityv1.UserStatus_value[tc.str]
			require.True(t, ok, "value map missing key %s", tc.str)
			require.Equal(t, tc.value, authnentityv1.UserStatus(v))
			n, ok := authnentityv1.UserStatus_name[int32(tc.value)]
			require.True(t, ok, "name map missing value %d", tc.value)
			require.Equal(t, tc.str, n)
		})
	}
}

func TestSessionType_ProtoValueMap_CoversAllExpectedValues(t *testing.T) {
	cases := []struct {
		value authnentityv1.SessionType
		str   string
	}{
		{authnentityv1.SessionType_SESSION_TYPE_UNSPECIFIED, "SESSION_TYPE_UNSPECIFIED"},
		{authnentityv1.SessionType_SESSION_TYPE_SERVER_SIDE, "SESSION_TYPE_SERVER_SIDE"},
		{authnentityv1.SessionType_SESSION_TYPE_JWT, "SESSION_TYPE_JWT"},
	}
	for _, tc := range cases {
		t.Run(tc.str, func(t *testing.T) {
			require.Equal(t, tc.str, tc.value.String())
			v, ok := authnentityv1.SessionType_value[tc.str]
			require.True(t, ok)
			require.Equal(t, tc.value, authnentityv1.SessionType(v))
		})
	}
}

func TestDeviceType_ProtoValueMap_CoversAllExpectedValues(t *testing.T) {
	cases := []struct {
		value authnentityv1.DeviceType
		str   string
	}{
		{authnentityv1.DeviceType_DEVICE_TYPE_UNSPECIFIED, "DEVICE_TYPE_UNSPECIFIED"},
		{authnentityv1.DeviceType_DEVICE_TYPE_WEB, "DEVICE_TYPE_WEB"},
		{authnentityv1.DeviceType_DEVICE_TYPE_MOBILE_ANDROID, "DEVICE_TYPE_MOBILE_ANDROID"},
		{authnentityv1.DeviceType_DEVICE_TYPE_MOBILE_IOS, "DEVICE_TYPE_MOBILE_IOS"},
		{authnentityv1.DeviceType_DEVICE_TYPE_DESKTOP, "DEVICE_TYPE_DESKTOP"},
		{authnentityv1.DeviceType_DEVICE_TYPE_API, "DEVICE_TYPE_API"},
	}
	for _, tc := range cases {
		t.Run(tc.str, func(t *testing.T) {
			require.Equal(t, tc.str, tc.value.String())
			v, ok := authnentityv1.DeviceType_value[tc.str]
			require.True(t, ok)
			require.Equal(t, tc.value, authnentityv1.DeviceType(v))
		})
	}
}
