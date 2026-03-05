package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// DLRHandler handles delivery report (DLR) webhooks from SSLWireless.
// POST /v1/internal/sms/dlr
type DLRHandler struct {
	client authnservicev1.AuthServiceClient
}

// NewDLRHandler creates a DLRHandler from a gRPC connection.
func NewDLRHandler(conn *grpc.ClientConn) *DLRHandler {
	return &DLRHandler{
		client: authnservicev1.NewAuthServiceClient(conn),
	}
}

// dlrBody is the JSON payload sent by SSLWireless.
type dlrBody struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
	ErrorCode string `json:"error_code,omitempty"`
}

// ServeHTTP processes the DLR webhook.
func (h *DLRHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Authenticate via shared secret header.
	secret := os.Getenv("DLR_WEBHOOK_SECRET")
	if secret == "" || r.Header.Get("X-DLR-Secret") != secret {
		logger.Warn("DLR webhook: unauthorized request",
			zap.String("remote_addr", r.RemoteAddr),
		)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var body dlrBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logger.Warn("DLR webhook: failed to decode body", zap.Error(err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	logger.Info("DLR webhook received",
		zap.String("provider_message_id", body.MessageID),
		zap.String("status", body.Status),
		zap.String("error_code", body.ErrorCode),
	)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, err := h.client.UpdateDLRStatus(ctx, &authnservicev1.UpdateDLRStatusRequest{
		ProviderMessageId: body.MessageID,
		Status:            body.Status,
		ErrorCode:         body.ErrorCode,
	})
	if err != nil {
		logger.Warn("DLR webhook: UpdateDLRStatus gRPC call failed",
			zap.String("provider_message_id", body.MessageID),
			zap.Error(err),
		)
		// Still return 200 to prevent SSLWireless from retrying indefinitely.
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
