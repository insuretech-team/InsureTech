package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/respond"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"go.uber.org/zap"
)

// Recovery middleware recovers from panics and logs the error
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.String("remote_addr", r.RemoteAddr),
					zap.ByteString("stack", debug.Stack()),
				)

				// Return 500 error to client
				respond.Error(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
