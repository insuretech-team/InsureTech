package domain

import (
	"fmt"
	"path"
	"strings"
)

// ObjMatch evaluates whether a request object matches a policy object.
//
// Unlike Casbin's built-in keyMatch2, this function does NOT interpret the
// colon in "svc:<name>" as a named-parameter placeholder. keyMatch2 treats
// ":name" as a regex wildcard ([^/]+), which causes "svc:claim/*" to match
// ANY "svc:*/…" object — a critical over-permission bug.
//
// Matching rules (evaluated in order):
//  1. Exact string equality.
//  2. If the policy contains no wildcards (* ? [), no match.
//  3. "svc:*" — matches any object starting with "svc:".
//  4. Policy ending in "/*" — prefix match (any depth below prefix).
//  5. Otherwise, path.Match (shell glob; * matches one segment only).
func ObjMatch(requestObj, policyObj string) bool {
	if requestObj == policyObj {
		return true
	}
	if !strings.ContainsAny(policyObj, "*?[") {
		return false
	}
	// "svc:*" → service-level wildcard (super_admin / readonly).
	if policyObj == "svc:*" {
		return strings.HasPrefix(requestObj, "svc:")
	}
	// "svc:xxx/*" → any resource under the service prefix (any depth).
	if strings.HasSuffix(policyObj, "/*") {
		prefix := policyObj[:len(policyObj)-2]
		return strings.HasPrefix(requestObj, prefix+"/")
	}
	// Fall back to shell-glob (single-segment * matching).
	matched, err := path.Match(policyObj, requestObj)
	return err == nil && matched
}

// ObjMatchExpressionFunc adapts ObjMatch for Casbin expression evaluation.
func ObjMatchExpressionFunc(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return false, nil
	}
	requestObj := strings.TrimSpace(fmt.Sprint(args[0]))
	policyObj := strings.TrimSpace(fmt.Sprint(args[1]))
	return ObjMatch(requestObj, policyObj), nil
}
