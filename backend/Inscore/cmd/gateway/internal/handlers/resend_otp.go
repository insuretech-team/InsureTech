package handlers

import (
	"context"
	"net/http"

	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	"google.golang.org/protobuf/proto"
)

func (h *AuthnHandler) ResendOTP(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req authnservicev1.ResendOTPRequest
		if err := protoUnmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.ResendOTP(ctx, &req)
	})
}
