package service

// session_limiter.go — Concurrent session enforcer (Sprint 1.11).
//
// SessionLimiter tracks active sessions per user in a Redis sorted set and
// evicts the oldest sessions when the limit is exceeded.
//
// Redis key:  sessions:active:<userID>
// Member:     sessionID (string UUID)
// Score:      creation unix timestamp (int64)
//
// Algorithm on TrackSession:
//   1. ZADD  key score=now.Unix() member=sessionID
//   2. ZCARD key                           — count active sessions
//   3. If count > maxSessions: ZPOPMIN key (count-maxSessions) → evicted IDs
//   4. Return evicted IDs so the caller can revoke them.
//
// Uses creation time as the score so ZPOPMIN always evicts the oldest-created
// session regardless of session type (JWT vs SERVER_SIDE) or TTL differences.

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
// userID and evicts the oldest-created sessions if the per-user limit is
// exceeded.
//
// Score = creation timestamp (not expiry), so ZPOPMIN always evicts the
// oldest-created session regardless of session type or TTL differences.
//
// Returns the list of evicted session IDs.  The caller is responsible for
// revoking those sessions (e.g. via TokenService.RevokeSession).
func (sl *SessionLimiter) TrackSession(ctx context.Context, userID, sessionID string, _ time.Time) (evicted []string, err error) {
	if sl.rdb == nil {
		// No Redis — limiter is a no-op (single-instance deployments use DB-level revocation).
		return nil, nil
	}

	k := sl.key(userID)
	now := time.Now().UTC()

	// 1. ZADD key score=now.Unix() member=sessionID
	//    Score is creation time so the oldest-created session has the lowest
	//    score and gets evicted first — never the brand-new session.
	if err := sl.rdb.ZAdd(ctx, k, redis.Z{
		Score:  float64(now.Unix()),
		Member: sessionID,
	}).Err(); err != nil {
		logger.Errorf("session_limiter ZADD: %v", err)
		return nil, errors.New("session_limiter ZADD")
	}

	// 2. Count active sessions
	count, err := sl.rdb.ZCard(ctx, k).Result()
	if err != nil {
		logger.Errorf("session_limiter ZCARD: %v", err)
		return nil, errors.New("session_limiter ZCARD")
	}

	// 3. Evict oldest-created sessions if over the limit
	if count > int64(sl.maxSessions) {
		overflow := count - int64(sl.maxSessions)

		// ZPOPMIN returns members with the lowest scores (oldest creation time)
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
