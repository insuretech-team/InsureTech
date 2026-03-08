package domain

import (
	"fmt"
	"regexp"
	"strings"
)

// ActionMatches evaluates request vs policy action safely.
// Supported policy formats:
// - "*" wildcard (match all)
// - exact verb (case-insensitive), e.g. GET/POST
// - regex pattern
//
// Invalid regex patterns fail closed (false), never panic.
func ActionMatches(requestAction, policyAction string) bool {
	requestAction = strings.TrimSpace(requestAction)
	policyAction = strings.TrimSpace(policyAction)
	if requestAction == "" || policyAction == "" {
		return false
	}
	if policyAction == "*" || strings.EqualFold(policyAction, requestAction) {
		return true
	}
	matched, err := regexp.MatchString(policyAction, requestAction)
	return err == nil && matched
}

// ActionMatchExpressionFunc adapts ActionMatches for Casbin expression funcs.
func ActionMatchExpressionFunc(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return false, nil
	}
	requestAction := strings.TrimSpace(fmt.Sprint(args[0]))
	policyAction := strings.TrimSpace(fmt.Sprint(args[1]))
	return ActionMatches(requestAction, policyAction), nil
}
