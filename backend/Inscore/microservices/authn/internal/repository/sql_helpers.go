package repository

import (
	"database/sql"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func nilOrTime(ts *timestamppb.Timestamp) any {
	if ts == nil {
		return nil
	}
	return ts.AsTime()
}

func requireTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		// Caller must ensure required timestamps are present.
		return time.Time{}
	}
	return ts.AsTime()
}

func nilIfZero(v int32) any {
	if v == 0 {
		return nil
	}
	return v
}

func nullableString(s string) any {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

func nullableJSON(s string) any {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

// sqlNullTime helper for scans
func sqlNullTime(t time.Time) sql.NullTime { return sql.NullTime{Time: t, Valid: true} }

func splitCSV(s sql.NullString) []string {
	if !s.Valid {
		return nil
	}
	trim := strings.TrimSpace(s.String)
	if trim == "" {
		return nil
	}
	parts := strings.Split(trim, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}
