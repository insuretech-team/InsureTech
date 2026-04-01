package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestObjMatch(t *testing.T) {
	tests := []struct {
		name       string
		requestObj string
		policyObj  string
		want       bool
	}{
		// ── Exact match ──────────────────────────────────────
		{name: "exact_match", requestObj: "svc:authz/roles", policyObj: "svc:authz/roles", want: true},
		{name: "exact_no_match", requestObj: "svc:authz/policies", policyObj: "svc:authz/roles", want: false},

		// ── svc:* super-wildcard ─────────────────────────────
		{name: "svc_star_matches_any", requestObj: "svc:authz/roles", policyObj: "svc:*", want: true},
		{name: "svc_star_matches_deep", requestObj: "svc:policy/my/list", policyObj: "svc:*", want: true},
		{name: "svc_star_no_match_nonsvc", requestObj: "other:foo", policyObj: "svc:*", want: false},

		// ── Prefix wildcard (/*) ─────────────────────────────
		{name: "slash_star_single", requestObj: "svc:claim/list", policyObj: "svc:claim/*", want: true},
		{name: "slash_star_deep", requestObj: "svc:claim/my/detail", policyObj: "svc:claim/*", want: true},
		{name: "slash_star_wrong_svc", requestObj: "svc:authz/list", policyObj: "svc:claim/*", want: false},
		{name: "slash_star_portals", requestObj: "svc:authz/portals/b2c", policyObj: "svc:authz/portals/*", want: true},
		{name: "slash_star_portals_wrong", requestObj: "svc:authz/roles", policyObj: "svc:authz/portals/*", want: false},
		{name: "slash_star_nested", requestObj: "svc:policy/my/list", policyObj: "svc:policy/my/*", want: true},

		// ── THE BUG: keyMatch2 would incorrectly match these ─
		{name: "bug_cross_svc_claim", requestObj: "svc:authz/authz/policies", policyObj: "svc:claim/*", want: false},
		{name: "bug_cross_svc_policy", requestObj: "svc:authz/authz/roles", policyObj: "svc:policy/*", want: false},
		{name: "bug_cross_svc_product", requestObj: "svc:authz/authz/audits", policyObj: "svc:product/*", want: false},
		{name: "bug_cross_svc_document", requestObj: "svc:authz/authz/policies", policyObj: "svc:document/*", want: false},

		// ── No wildcards → no match if not exact ─────────────
		{name: "no_wildcard_different", requestObj: "svc:authz/authz/check", policyObj: "svc:authz/check", want: false},
		{name: "no_wildcard_exact", requestObj: "svc:authz/check", policyObj: "svc:authz/check", want: true},

		// ── Edge cases ───────────────────────────────────────
		{name: "empty_request", requestObj: "", policyObj: "svc:authz/*", want: false},
		{name: "empty_policy", requestObj: "svc:authz/roles", policyObj: "", want: false},
		{name: "both_empty", requestObj: "", policyObj: "", want: true},
		{name: "glob_question_mark", requestObj: "svc:authz/a", policyObj: "svc:authz/?", want: true},
		{name: "glob_question_no_match", requestObj: "svc:authz/ab", policyObj: "svc:authz/?", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, ObjMatch(tc.requestObj, tc.policyObj), "ObjMatch(%q, %q)", tc.requestObj, tc.policyObj)
		})
	}
}

func TestObjMatchExpressionFunc(t *testing.T) {
	result, err := ObjMatchExpressionFunc("svc:claim/list", "svc:claim/*")
	require.NoError(t, err)
	require.Equal(t, true, result)

	result, err = ObjMatchExpressionFunc("svc:authz/authz/roles", "svc:claim/*")
	require.NoError(t, err)
	require.Equal(t, false, result)

	// Wrong arg count → false, no error
	result, err = ObjMatchExpressionFunc("only-one")
	require.NoError(t, err)
	require.Equal(t, false, result)
}
