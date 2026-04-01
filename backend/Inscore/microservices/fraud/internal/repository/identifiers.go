package repository

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

func newSequenceNumber(prefix string, now time.Time) string {
	suffix := strings.ToUpper(strings.ReplaceAll(uuid.NewString(), "-", ""))
	if len(suffix) > 8 {
		suffix = suffix[:8]
	}
	return prefix + "-" + now.UTC().Format("20060102-150405") + "-" + suffix
}
