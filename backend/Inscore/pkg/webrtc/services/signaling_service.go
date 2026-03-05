package services

import (
	"context"
	"sync"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/webrtc/repository"
	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/entity"
	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/events"
	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// SignalingService implements the WebRTC SignalingService gRPC API
type SignalingService struct {
	service.UnimplementedSignalingServiceServer
	repo      *repository.Repository
	streams   sync.Map // Map of peer_id to stream channel
	roomPeers sync.Map // Map of room_id to set of peer_ids
	mu        sync.RWMutex
}

// streamInfo holds information about a peer's stream
type streamInfo struct {
	channel chan *service.SignalResponse
	roomID  string
	peerID  string
}

// NewSignalingService creates a new SignalingService
func NewSignalingService(repo *repository.Repository) *SignalingService {
	return &SignalingService{
		repo: repo,
	}
}

// Connect establishes a bidirectional signaling connection
func (s *SignalingService) Connect(stream service.SignalingService_ConnectServer) error {
	ctx := stream.Context()
	var peerID string
	var roomID string

	// Create a channel for this peer's outgoing messages
	outgoing := make(chan *service.SignalResponse, 100)
	defer func() {
		close(outgoing)
		if peerID != "" {
			s.cleanupPeerConnection(peerID, roomID)
		}
	}()

	// Goroutine to send messages to client
	errChan := make(chan error, 1)
	sendDone := make(chan struct{})
	go func() {
		defer close(sendDone)
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-outgoing:
				if !ok {
					return
				}
				if err := stream.Send(msg); err != nil {
					errChan <- err
					return
				}
			}
		}
	}()

	// Receive messages from client
	for {
		select {
		case err := <-errChan:
			return err
		case <-ctx.Done():
			return ctx.Err()
		default:
			req, err := stream.Recv()
			if err != nil {
				return err
			}

			// Store peer stream on first message
			if peerID == "" && req.PeerId != "" && req.RoomId != "" {
				peerID = req.PeerId
				roomID = req.RoomId

				// Store stream info
				info := &streamInfo{
					channel: outgoing,
					roomID:  roomID,
					peerID:  peerID,
				}
				s.streams.Store(peerID, info)

				// Add peer to room's peer list
				s.addPeerToRoom(roomID, peerID)

				// Notify other peers in room about new peer
				if err := s.broadcastPeerJoined(ctx, roomID, peerID); err != nil {
					return err
				}
			}

			// Handle different signal types
			switch payload := req.Payload.(type) {
			case *service.SignalRequest_Offer:
				if err := s.handleOffer(ctx, req.PeerId, req.RoomId, payload.Offer, outgoing); err != nil {
					return err
				}
			case *service.SignalRequest_Answer:
				if err := s.handleAnswer(ctx, req.PeerId, req.RoomId, payload.Answer, outgoing); err != nil {
					return err
				}
			case *service.SignalRequest_IceCandidate:
				if err := s.handleICECandidate(ctx, req.PeerId, req.RoomId, payload.IceCandidate, outgoing); err != nil {
					return err
				}
			case *service.SignalRequest_Ping:
				// Respond with pong
				select {
				case outgoing <- &service.SignalResponse{
					PeerId:    req.PeerId,
					RoomId:    req.RoomId,
					Timestamp: timestamppb.Now(),
					Payload: &service.SignalResponse_Pong{
						Pong: &service.PongResponse{
							Timestamp: timestamppb.Now(),
						},
					},
				}:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
	}
}

// SendOffer sends an SDP offer to another peer
func (s *SignalingService) SendOffer(ctx context.Context, req *service.SendOfferRequest) (*service.SendOfferResponse, error) {
	if req.FromPeerId == "" || req.ToPeerId == "" {
		return nil, status.Error(codes.InvalidArgument, "from_peer_id and to_peer_id are required")
	}

	if req.Sdp == "" {
		return nil, status.Error(codes.InvalidArgument, "sdp is required")
	}

	// Verify both peers exist
	fromPeer, err := s.repo.GetPeer(ctx, req.FromPeerId)
	if err != nil || fromPeer == nil {
		return nil, status.Error(codes.NotFound, "from_peer not found")
	}

	toPeer, err := s.repo.GetPeer(ctx, req.ToPeerId)
	if err != nil || toPeer == nil {
		return nil, status.Error(codes.NotFound, "to_peer not found")
	}

	// Verify peers are in the same room
	// Note: entity.Peer might still use wrappers, check access
	if fromPeer.RoomId != toPeer.RoomId {
		return nil, status.Error(codes.PermissionDenied, "peers are not in the same room")
	}

	// Forward offer to target peer's stream
	msg := &service.SignalResponse{
		PeerId:    req.ToPeerId,
		RoomId:    toPeer.RoomId,
		Timestamp: timestamppb.Now(),
		Payload: &service.SignalResponse_OfferReceived{
			OfferReceived: &events.OfferReceivedEvent{
				RoomId: toPeer.RoomId,
				Offer: &entity.SDPOffer{
					FromPeerId: req.FromPeerId,
					ToPeerId:   req.ToPeerId,
					Sdp:        req.Sdp,
					Type:       entity.SDPType_SDP_TYPE_OFFER,
				},
				ReceivedAt: timestamppb.Now(),
			},
		},
	}

	if err := s.sendToPeer(req.ToPeerId, msg); err != nil {
		return nil, err
	}

	return &service.SendOfferResponse{Success: true}, nil
}

// SendAnswer sends an SDP answer to another peer
func (s *SignalingService) SendAnswer(ctx context.Context, req *service.SendAnswerRequest) (*service.SendAnswerResponse, error) {
	if req.FromPeerId == "" || req.ToPeerId == "" {
		return nil, status.Error(codes.InvalidArgument, "from_peer_id and to_peer_id are required")
	}

	if req.Sdp == "" {
		return nil, status.Error(codes.InvalidArgument, "sdp is required")
	}

	// Verify both peers exist
	fromPeer, err := s.repo.GetPeer(ctx, req.FromPeerId)
	if err != nil || fromPeer == nil {
		return nil, status.Error(codes.NotFound, "from_peer not found")
	}

	toPeer, err := s.repo.GetPeer(ctx, req.ToPeerId)
	if err != nil || toPeer == nil {
		return nil, status.Error(codes.NotFound, "to_peer not found")
	}

	// Verify peers are in the same room
	if fromPeer.RoomId != toPeer.RoomId {
		return nil, status.Error(codes.PermissionDenied, "peers are not in the same room")
	}

	// Forward answer to target peer's stream
	msg := &service.SignalResponse{
		PeerId:    req.ToPeerId,
		RoomId:    toPeer.RoomId,
		Timestamp: timestamppb.Now(),
		Payload: &service.SignalResponse_AnswerReceived{
			AnswerReceived: &events.AnswerReceivedEvent{
				RoomId: toPeer.RoomId,
				Answer: &entity.SDPAnswer{
					FromPeerId: req.FromPeerId,
					ToPeerId:   req.ToPeerId,
					Sdp:        req.Sdp,
					Type:       entity.SDPType_SDP_TYPE_ANSWER,
				},
				ReceivedAt: timestamppb.Now(),
			},
		},
	}

	if err := s.sendToPeer(req.ToPeerId, msg); err != nil {
		return nil, err
	}

	return &service.SendAnswerResponse{Success: true}, nil
}

// SendICECandidate sends an ICE candidate to another peer
func (s *SignalingService) SendICECandidate(ctx context.Context, req *service.SendICECandidateRequest) (*service.SendICECandidateResponse, error) {
	if req.FromPeerId == "" || req.ToPeerId == "" {
		return nil, status.Error(codes.InvalidArgument, "from_peer_id and to_peer_id are required")
	}

	if req.Candidate == "" {
		return nil, status.Error(codes.InvalidArgument, "candidate is required")
	}

	// Verify both peers exist
	fromPeer, err := s.repo.GetPeer(ctx, req.FromPeerId)
	if err != nil || fromPeer == nil {
		return nil, status.Error(codes.NotFound, "from_peer not found")
	}

	toPeer, err := s.repo.GetPeer(ctx, req.ToPeerId)
	if err != nil || toPeer == nil {
		return nil, status.Error(codes.NotFound, "to_peer not found")
	}

	// Verify peers are in the same room
	if fromPeer.RoomId != toPeer.RoomId {
		return nil, status.Error(codes.PermissionDenied, "peers are not in the same room")
	}

	// Forward ICE candidate to target peer's stream
	msg := &service.SignalResponse{
		PeerId:    req.ToPeerId,
		RoomId:    toPeer.RoomId,
		Timestamp: timestamppb.Now(),
		Payload: &service.SignalResponse_IceCandidateReceived{
			IceCandidateReceived: &events.ICECandidateReceivedEvent{
				RoomId: toPeer.RoomId,
				Candidate: &entity.ICECandidate{
					FromPeerId:       req.FromPeerId,
					ToPeerId:         req.ToPeerId,
					Candidate:        req.Candidate,
					SdpMid:           req.SdpMid,
					SdpMLineIndex:    req.SdpMLineIndex,
					UsernameFragment: "",
				},
				ReceivedAt: timestamppb.Now(),
			},
		},
	}

	if err := s.sendToPeer(req.ToPeerId, msg); err != nil {
		return nil, err
	}

	return &service.SendICECandidateResponse{Success: true}, nil
}

// Helper methods

func (s *SignalingService) handleOffer(ctx context.Context, peerID string, roomID string, req *service.SendOfferRequest, outgoing chan *service.SignalResponse) error {
	_, err := s.SendOffer(ctx, req)
	return err
}

func (s *SignalingService) handleAnswer(ctx context.Context, peerID string, roomID string, req *service.SendAnswerRequest, outgoing chan *service.SignalResponse) error {
	_, err := s.SendAnswer(ctx, req)
	return err
}

func (s *SignalingService) handleICECandidate(ctx context.Context, peerID string, roomID string, req *service.SendICECandidateRequest, outgoing chan *service.SignalResponse) error {
	_, err := s.SendICECandidate(ctx, req)
	return err
}

// addPeerToRoom adds a peer to a room's peer list
func (s *SignalingService) addPeerToRoom(roomID, peerID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var peerSet map[string]bool
	if val, ok := s.roomPeers.Load(roomID); ok {
		peerSet = val.(map[string]bool)
	} else {
		peerSet = make(map[string]bool)
	}
	peerSet[peerID] = true
	s.roomPeers.Store(roomID, peerSet)
}

// removePeerFromRoom removes a peer from a room's peer list
func (s *SignalingService) removePeerFromRoom(roomID, peerID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if val, ok := s.roomPeers.Load(roomID); ok {
		peerSet := val.(map[string]bool)
		delete(peerSet, peerID)
		if len(peerSet) == 0 {
			s.roomPeers.Delete(roomID)
		} else {
			s.roomPeers.Store(roomID, peerSet)
		}
	}
}

// cleanupPeerConnection cleans up a peer's connection
func (s *SignalingService) cleanupPeerConnection(peerID, roomID string) {
	s.streams.Delete(peerID)
	if roomID != "" {
		s.removePeerFromRoom(roomID, peerID)
		// Notify other peers about peer leaving
		ctx := context.Background()
		_ = s.broadcastPeerLeft(ctx, roomID, peerID)
	}
}

// broadcastPeerJoined notifies all peers in a room about a new peer joining
func (s *SignalingService) broadcastPeerJoined(ctx context.Context, roomID, newPeerID string) error {
	peer, err := s.repo.GetPeer(ctx, newPeerID)
	if err != nil || peer == nil {
		return err
	}

	msg := &service.SignalResponse{
		PeerId:    newPeerID,
		RoomId:    roomID,
		Timestamp: timestamppb.Now(),
		Payload: &service.SignalResponse_PeerJoined{
			PeerJoined: &events.PeerJoinedEvent{
				RoomId:   roomID,
				Peer:     peer,
				JoinedAt: timestamppb.Now(),
			},
		},
	}

	return s.broadcastToRoom(roomID, newPeerID, msg)
}

// broadcastPeerLeft notifies all peers in a room about a peer leaving
func (s *SignalingService) broadcastPeerLeft(ctx context.Context, roomID, leftPeerID string) error {
	msg := &service.SignalResponse{
		PeerId:    leftPeerID,
		RoomId:    roomID,
		Timestamp: timestamppb.Now(),
		Payload: &service.SignalResponse_PeerLeft{
			PeerLeft: &events.PeerLeftEvent{
				RoomId: roomID,
				PeerId: leftPeerID,
				LeftAt: timestamppb.Now(),
				Reason: "disconnected",
			},
		},
	}

	return s.broadcastToRoom(roomID, leftPeerID, msg)
}

// broadcastToRoom sends a message to all peers in a room except the sender
func (s *SignalingService) broadcastToRoom(roomID, senderPeerID string, msg *service.SignalResponse) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.roomPeers.Load(roomID)
	if !ok {
		return nil
	}

	peerSet := val.(map[string]bool)
	for peerID := range peerSet {
		if peerID == senderPeerID {
			continue
		}

		if val, ok := s.streams.Load(peerID); ok {
			info := val.(*streamInfo)
			select {
			case info.channel <- msg:
			default:
				// Channel full, skip this peer
			}
		}
	}

	return nil
}

// sendToPeer sends a message to a specific peer
func (s *SignalingService) sendToPeer(peerID string, msg *service.SignalResponse) error {
	val, ok := s.streams.Load(peerID)
	if !ok {
		return status.Error(codes.NotFound, "peer not connected")
	}

	info := val.(*streamInfo)
	select {
	case info.channel <- msg:
		return nil
	default:
		return status.Error(codes.ResourceExhausted, "peer message queue full")
	}
}
