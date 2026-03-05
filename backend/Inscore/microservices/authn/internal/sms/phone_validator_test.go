package sms

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizePhoneNumber(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		out     string
		wantErr error
	}{
		{name: "local", in: "01712345678", out: "8801712345678"},
		{name: "with plus", in: "+8801712345678", out: "8801712345678"},
		{name: "with 00 prefix", in: "008801712345678", out: "8801712345678"},
		{name: "ten digits", in: "1712345678", out: "8801712345678"},
		{name: "invalid empty", in: "", wantErr: ErrInvalidPhoneNumber},
		{name: "invalid operator", in: "01112345678", wantErr: ErrInvalidPhoneNumber},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NormalizePhoneNumber(tc.in)
			if tc.wantErr != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tc.wantErr))
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.out, got)
		})
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	require.True(t, ValidatePhoneNumber("01712345678"))
	require.True(t, ValidatePhoneNumber("+8801712345678"))
	require.False(t, ValidatePhoneNumber("01112345678"))
}

func TestGetOperator(t *testing.T) {
	op, err := GetOperator("01712345678")
	require.NoError(t, err)
	require.Equal(t, OperatorGrameenphone, op)

	op, err = GetOperator("01812345678")
	require.NoError(t, err)
	require.Equal(t, OperatorRobi, op)

	op, err = GetOperator("01912345678")
	require.NoError(t, err)
	require.Equal(t, OperatorBanglalink, op)

	op, err = GetOperator("01612345678")
	require.NoError(t, err)
	require.Equal(t, OperatorTeletalk, op)

	op, err = GetOperator("01112345678")
	require.Error(t, err)
	require.Equal(t, OperatorUnknown, op)
}

func TestMaskPhoneNumberAndFormatForDisplay(t *testing.T) {
	require.Equal(t, "880171***5678", MaskPhoneNumber("8801712345678"))
	require.Equal(t, "12345", MaskPhoneNumber("12345"))

	require.Equal(t, "+880 171 234 5678", FormatForDisplay("01712345678"))
	require.Equal(t, "invalid", FormatForDisplay("invalid"))
}
