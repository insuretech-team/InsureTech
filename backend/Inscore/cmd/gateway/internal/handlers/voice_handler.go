package handlers

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	roomservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/service"
)

// VoiceHandler proxies voice session requests to the WebRTC/room gRPC service.
type VoiceHandler struct {
	client roomservicev1.RoomServiceClient
}

// NewVoiceHandler creates a VoiceHandler from a gRPC connection to the WebRTC service.
func NewVoiceHandler(conn *grpc.ClientConn) *VoiceHandler {
	return &VoiceHandler{client: roomservicev1.NewRoomServiceClient(conn)}
}

// Create creates a new voice/video session (room).
// POST /v1/voice/sessions
func (h *VoiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req roomservicev1.CreateRoomRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.CreateRoom(ctx, &req)
	})
}

// Get retrieves a voice session (room) by ID.
// GET /v1/voice/sessions/{session_id}
func (h *VoiceHandler) Get(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("session_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetRoom(ctx, &roomservicev1.GetRoomRequest{
			RoomId: sessionID,
		})
	})
}
