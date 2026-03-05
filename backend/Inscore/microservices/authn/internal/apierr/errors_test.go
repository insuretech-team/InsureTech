package apierr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDomainError_ConstructorsAndIs(t *testing.T) {
	e := NotFound("missing", errors.New("db"))
	require.Contains(t, e.Error(), "missing")
	require.ErrorIs(t, e, ErrNotFound)

	require.ErrorIs(t, AlreadyExists("dup", nil), ErrAlreadyExists)
	require.ErrorIs(t, InvalidCredentials("bad"), ErrInvalidCredentials)
	require.ErrorIs(t, InvalidArgument("invalid"), ErrInvalidArgument)
	require.ErrorIs(t, Expired("expired"), ErrExpired)
	require.ErrorIs(t, RateLimited("slow"), ErrRateLimited)
	require.ErrorIs(t, PermissionDenied("deny"), ErrPermissionDenied)
	require.ErrorIs(t, Internal("oops", errors.New("x")), ErrInternal)
	require.ErrorIs(t, Unauthenticated("unauth"), ErrUnauthenticated)
}
