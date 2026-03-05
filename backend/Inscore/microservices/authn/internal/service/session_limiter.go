package service

// session_limiter.go — Concurrent session enforcer (Sprint 1.11).
//
// SessionLimiter tracks active sessions per user in a Redis sorted set and
// evicts the oldest sessions when the limit is exceeded.
//
// Redis key:  sessions:active:<userID>
// Member:     sessionID (string UUID)
// Score:      expiry unix timestamp (int64)
//
// Algorithm on TrackSession:
//   1. ZADD  key score=expiry.Unix() member=sessionID
//   2. ZREMRANGEBYSCORE key -inf <now>     — purge expired entries
//   3. ZCARD key                           — count remaining active sessions
//   4. If count > maxSessions: ZPOPMIN key (count-maxSessions) → evicted IDs
//   5. Return evicted IDs so the caller can revoke them.

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"

	"github.com/redis/go-redis/v9"
)

const defaultMaxSessions = 5

// SessionLimiter enforces per-user concurrent session limits using a Redis
// sorted set.  Pass maxSessions ≤ 0 to use the default (5).
type SessionLimiter struct {
	rdb         redis.UniversalClient
	maxSessions int
	mu          sync.Mutex // guards in-memory fallback only
}

// NewSessionLimiter creates a SessionLimiter backed by the provided Redis
// client.  maxSessions ≤ 0 defaults to 5.
func NewSessionLimiter(rdb redis.UniversalClient, maxSessions int) *SessionLimiter {
	if maxSessions <= 0 {
		maxSessions = defaultMaxSessions
	}
	return &SessionLimiter{
		rdb:         rdb,
		maxSessions: maxSessions,
	}
}

func (sl *SessionLimiter) key(userID string) string {
	return "sessions:active:" + userID
}

// TrackSession registers sessionID in the active-session sorted set for
// userID, removes expired entries, and evicts the oldest sessions if the
// per-user limit is exceeded.
//
// Returns the list of evicted session IDs.  The caller is responsible for
// revoking those sessions (e.g. via TokenService.RevokeSession).
func (sl *SessionLimiter) TrackSession(ctx context.Context, userID, sessionID string, expiry time.Time) (evicted []string, err error) {
	if sl.rdb == nil {
		// No Redis — limiter is a no-op (single-instance deployments use DB-level revocation).
		return nil, nil
	}

	k := sl.key(userID)
	now := time.Now().UTC()

	// 1. ZADD key score=expiry.Unix() member=sessionID
	if err := sl.rdb.ZAdd(ctx, k, redis.Z{
		Score:  float64(expiry.Unix()),
		Member: sessionID,
	}).Err(); err != nil {
		logger.Errorf("session_limiter ZADD: %v", err)
		return nil, errors.New("session_limiter ZADD")
	}

	// 2. Remove expired sessions (score < now)
	if err := sl.rdb.ZRemRangeByScore(ctx, k, "-inf", strconv.FormatInt(now.Unix(), 10)).Err(); err != nil {
		// Non-fatal: continue even if cleanup fails.
		_ = err
	}

	// 3. Count active sessions
	count, err := sl.rdb.ZCard(ctx, k).Result()
	if err != nil {
		logger.Errorf("session_limiter ZCARD: %v", err)
		return nil, errors.New("session_limiter ZCARD")
	}

	// 4. Evict oldest sessions if over the limit
	if count > int64(sl.maxSessions) {
		overflow := count - int64(sl.maxSessions)

		// ZPOPMIN returns members with the lowest scores (oldest expiry = oldest sessions)
		result, err := sl.rdb.ZPopMin(ctx, k, overflow).Result()
		if err != nil {
			logger.Errorf("session_limiter ZPOPMIN: %v", err)
			return nil, errors.New("session_limiter ZPOPMIN")
		}
		for _, z := range result {
			if id, ok := z.Member.(string); ok && id != "" {
				evicted = append(evicted, id)
			}
		}
	}

	return evicted, nil
}

// RemoveSession removes a single session from the active-session sorted set.
// Call this on explicit logout / session revocation so the slot is freed
// immediately rather than waiting for TTL expiry.
func (sl *SessionLimiter) RemoveSession(ctx context.Context, userID, sessionID string) error {
	if sl.rdb == nil {
		return nil
	}
	return sl.rdb.ZRem(ctx, sl.key(userID), sessionID).Err()
}

// ActiveCount returns the number of non-expired active sessions for userID.
// Expired sessions are pruned before counting.
func (sl *SessionLimiter) ActiveCount(ctx context.Context, userID string) (int64, error) {
	if sl.rdb == nil {
		return 0, nil
	}
	k := sl.key(userID)
	now := time.Now().UTC()

	// Prune expired entries first.
	_ = sl.rdb.ZRemRangeByScore(ctx, k, "-inf", strconv.FormatInt(now.Unix(), 10))

	return sl.rdb.ZCard(ctx, k).Result()
}
