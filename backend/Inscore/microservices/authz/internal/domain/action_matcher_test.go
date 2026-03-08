package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestActionMatches(t *testing.T) {
	tests := []struct {
		name          string
		requestAction string
		policyAction  string
		want          bool
	}{
		{name: "wildcard", requestAction: "POST", policyAction: "*", want: true},
		{name: "exact_case_insensitive", requestAction: "post", policyAction: "POST", want: true},
		{name: "regex", requestAction: "PATCH", policyAction: "P(UT|ATCH)", want: true},
		{name: "invalid_regex_fails_closed", requestAction: "GET", policyAction: "[", want: false},
		{name: "blank_request", requestAction: "", policyAction: "GET", want: false},
		{name: "blank_policy", requestAction: "GET", policyAction: "", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, ActionMatches(tc.requestAction, tc.policyAction))
		})
	}
}

func TestActionMatchExpressionFunc(t *testing.T) {
	result, err := ActionMatchExpressionFunc("GET", "*")
	require.NoError(t, err)
	require.Equal(t, true, result)

	result, err = ActionMatchExpressionFunc("GET", "[")
	require.NoError(t, err)
	require.Equal(t, false, result)

	result, err = ActionMatchExpressionFunc("GET")
	require.NoError(t, err)
	require.Equal(t, false, result)
}
