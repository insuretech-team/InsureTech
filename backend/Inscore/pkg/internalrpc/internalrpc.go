package internalrpc

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/metadata"
)

const (
	HeaderService   = "x-internal-service"
	HeaderTimestamp = "x-internal-timestamp"
	HeaderSignature = "x-internal-signature"
	secretEnvKey    = "INTERNAL_RPC_SHARED_SECRET"
	maxClockSkew    = 2 * time.Minute
)

var (
	ErrMissingMetadata = errors.New("missing internal rpc metadata")
	ErrMissingSecret   = errors.New("missing INTERNAL_RPC_SHARED_SECRET")
	ErrInvalidService  = errors.New("invalid internal service")
	ErrInvalidTime     = errors.New("invalid internal rpc timestamp")
	ErrInvalidSig      = errors.New("invalid internal rpc signature")
)

func OutgoingContext(ctx context.Context, serviceName string) context.Context {
	serviceName = strings.TrimSpace(serviceName)
	if serviceName == "" {
		return ctx
	}

	md := metadata.Pairs(HeaderService, serviceName)
	if secret := strings.TrimSpace(os.Getenv(secretEnvKey)); secret != "" {
		ts := strconv.FormatInt(time.Now().UTC().Unix(), 10)
		md.Set(HeaderTimestamp, ts)
		md.Set(HeaderSignature, sign(serviceName, ts, secret))
	}

	if existing, ok := metadata.FromOutgoingContext(ctx); ok {
		cloned := existing.Copy()
		for key, values := range md {
			cloned.Set(key, values...)
		}
		return metadata.NewOutgoingContext(ctx, cloned)
	}
	return metadata.NewOutgoingContext(ctx, md)
}

func ValidateIncoming(ctx context.Context, trusted map[string]struct{}) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrMissingMetadata
	}

	serviceName := strings.ToLower(strings.TrimSpace(first(md, HeaderService)))
	if serviceName == "" {
		return "", ErrMissingMetadata
	}
	if _, ok := trusted[serviceName]; !ok {
		return "", ErrInvalidService
	}

	secret := strings.TrimSpace(os.Getenv(secretEnvKey))
	if secret == "" {
		return "", ErrMissingSecret
	}

	timestamp := strings.TrimSpace(first(md, HeaderTimestamp))
	if timestamp == "" {
		return "", ErrInvalidTime
	}
	unixSeconds, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return "", ErrInvalidTime
	}
	now := time.Now().UTC()
	signedAt := time.Unix(unixSeconds, 0).UTC()
	if signedAt.Before(now.Add(-maxClockSkew)) || signedAt.After(now.Add(maxClockSkew)) {
		return "", ErrInvalidTime
	}

	signature := strings.TrimSpace(first(md, HeaderSignature))
	if signature == "" {
		return "", ErrInvalidSig
	}
	expected := sign(serviceName, timestamp, secret)
	if !hmac.Equal([]byte(signature), []byte(expected)) {
		return "", ErrInvalidSig
	}

	return serviceName, nil
}

func first(md metadata.MD, key string) string {
	for _, value := range md.Get(key) {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}

func sign(serviceName, timestamp, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(serviceName))
	_, _ = mac.Write([]byte{'\n'})
	_, _ = mac.Write([]byte(timestamp))
	return hex.EncodeToString(mac.Sum(nil))
}
