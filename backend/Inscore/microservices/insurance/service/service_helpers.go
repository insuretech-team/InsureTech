package service

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func newID() string {
	return uuid.NewString()
}

func nowTS() *timestamppb.Timestamp {
	return timestamppb.New(time.Now().UTC())
}

func pageOffset(pageToken string) int {
	if strings.TrimSpace(pageToken) == "" {
		return 0
	}

	offset, err := strconv.Atoi(pageToken)
	if err != nil || offset < 0 {
		return 0
	}

	return offset
}

func nextToken(offset, count, total int) string {
	next := offset + count
	if next >= total {
		return ""
	}
	return strconv.Itoa(next)
}

func errorResponse(code, message string, httpStatus int32) *commonv1.Error {
	return &commonv1.Error{
		Code:           code,
		Message:        message,
		HttpStatusCode: httpStatus,
	}
}

func ensureMetadata(metadata map[string]string) map[string]string {
	if metadata == nil {
		return map[string]string{}
	}
	return metadata
}

func marshalJSON(v any) string {
	if v == nil {
		return ""
	}

	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}

	return string(b)
}

func unmarshalJSON[T any](raw string, target *T) bool {
	if strings.TrimSpace(raw) == "" || target == nil {
		return false
	}

	if err := json.Unmarshal([]byte(raw), target); err != nil {
		return false
	}

	return true
}

func money(amount int64, currency string) *commonv1.Money {
	if strings.TrimSpace(currency) == "" {
		currency = "BDT"
	}

	return &commonv1.Money{
		Amount:   amount,
		Currency: currency,
	}
}
