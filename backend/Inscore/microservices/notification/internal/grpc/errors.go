package grpc

import (
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toGRPCError(err error) error {
	if err == nil {
		return nil
	}
	if isGRPCStatus(err) {
		return err
	}

	msg := err.Error()
	lower := strings.ToLower(msg)

	switch {
	case contains(lower, "not found", "record not found", "no rows"):
		return status.Error(codes.NotFound, msg)
	case contains(lower, "already exists", "duplicate"):
		return status.Error(codes.AlreadyExists, msg)
	case contains(lower, "required", "invalid", "unsupported", "must be", "no enabled notification channels"):
		return status.Error(codes.InvalidArgument, msg)
	case contains(lower, "disabled by user preferences", "not active", "not configured"):
		return status.Error(codes.FailedPrecondition, msg)
	case contains(lower, "email sendmail failed", "sms failed", "provider did not accept"):
		return status.Error(codes.Unavailable, msg)
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}

func contains(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func isGRPCStatus(err error) bool {
	_, ok := status.FromError(err)
	return ok && status.Code(err) != codes.OK && status.Code(err) != codes.Unknown
}
