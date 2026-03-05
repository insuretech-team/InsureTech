package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/webrtc/mediaserver"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/webrtc/repository"
	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/entity"
	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/events"
	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/service"
	"github.com/pion/webrtc/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// SignalingServiceWithSFU integrates signaling with the SFU media server
type SignalingServiceWithSFU struct {
	service.UnimplementedSignalingServiceServer
	repo      *repository.Repository
	sfu       *mediaserver.SFU
	streams   sync.Map // Map of peer_id to stream channel
	roomPeers sync.Map // Map of room_id to set of peer_ids
	peerConns sync.Map // Map of peer_id to *webrtc.PeerConnection
	mu        sync.RWMutex
}

// NewSignalingServiceWithSFU creates a new signaling service with SFU integration
func NewSignalingServiceWithSFU(repo *repository.Repository, sfu *mediaserver.SFU) *SignalingServiceWithSFU {
	return &SignalingServiceWithSFU{
		repo: repo,
		sfu:  sfu,
	}
}

// Connect establishes a bidirectional signaling connection with full SFU integration
func (s *SignalingServiceWithSFU) Connect(stream service.SignalingService_ConnectServer) error {
	ctx := stream.Context()
	var peerID string
	var roomID string
	var sfuPeer *mediaserver.Peer

	// Create a channel for this peer's outgoing messages
	outgoing := make(chan *service.SignalResponse, 100)
	defer func() {
		close(outgoing)
		if peerID != "" {
			s.cleanupPeerConnection(ctx, peerID, roomID)
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

				// Verify peer exists in database
				peer, err := s.repo.GetPeer(ctx, peerID)
				if err != nil || peer == nil {
					return status.Error(codes.NotFound, "peer not found")
				}

				// Create SFU peer connection
				sfuPeer, err = s.sfu.AddPeer(ctx, roomID, peerID)
				if err != nil {
					return status.Errorf(codes.Internal, "failed to create SFU peer: %v", err)
				}

				// Setup peer connection callbacks
				s.setupPeerConnectionCallbacks(sfuPeer, peerID, roomID, outgoing)

				// Store stream info
				info := &streamInfo{
					channel: outgoing,
					roomID:  roomID,
					peerID:  peerID,
				}
				s.streams.Store(peerID, info)
				s.peerConns.Store(peerID, sfuPeer.GetPeerConnection())

				// Add peer to room's peer list
				s.addPeerToRoom(roomID, peerID)

				// Notify other peers in room about new peer
				if err := s.broadcastPeerJoined(ctx, roomID, peerID); err != nil {
					return err
				}

				// Update peer state to connecting
				peer.State = entity.PeerConnectionState_PEER_CONNECTION_STATE_CONNECTING
				_ = s.repo.UpdatePeerState(ctx, peerID, peer.State)
			}

			// Handle different signal types
			switch payload := req.Payload.(type) {
			case *service.SignalRequest_Offer:
				if err := s.handleOfferWithSFU(ctx, peerID, roomID, payload.Offer, outgoing, sfuPeer); err != nil {
					return err
				}
			case *service.SignalRequest_Answer:
				if err := s.handleAnswerWithSFU(ctx, peerID, roomID, payload.Answer, outgoing, sfuPeer); err != nil {
					return err
				}
			case *service.SignalRequest_IceCandidate:
				if err := s.handleICECandidateWithSFU(ctx, peerID, payload.IceCandidate, sfuPeer); err != nil {
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

// setupPeerConnectionCallbacks sets up WebRTC peer connection event handlers
func (s *SignalingServiceWithSFU) setupPeerConnectionCallbacks(
	sfuPeer *mediaserver.Peer,
	peerID, roomID string,
	outgoing chan *service.SignalResponse,
) {
	pc := sfuPeer.GetPeerConnection()

	// Handle ICE candidates
	pc.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			return
		}

		// Send ICE candidate to client
		candidateInit := candidate.ToJSON()
		select {
		case outgoing <- &service.SignalResponse{
			PeerId:    peerID,
			RoomId:    roomID,
			Timestamp: timestamppb.Now(),
			Payload: &service.SignalResponse_IceCandidateReceived{
				IceCandidateReceived: &events.ICECandidateReceivedEvent{
					RoomId: roomID,
					Candidate: &entity.ICECandidate{
						FromPeerId:    "server",
						ToPeerId:      peerID,
						Candidate:     candidateInit.Candidate,
						SdpMid:        *candidateInit.SDPMid,
						SdpMLineIndex: int32(*candidateInit.SDPMLineIndex),
					},
					ReceivedAt: timestamppb.Now(),
				},
			},
		}:
		default:
			// Channel full
		}
	})

	// Handle connection state changes
	pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		// Update peer state in database
		var dbState entity.PeerConnectionState
		switch state {
		case webrtc.PeerConnectionStateNew:
			dbState = entity.PeerConnectionState_PEER_CONNECTION_STATE_NEW
		case webrtc.PeerConnectionStateConnecting:
			dbState = entity.PeerConnectionState_PEER_CONNECTION_STATE_CONNECTING
		case webrtc.PeerConnectionStateConnected:
			dbState = entity.PeerConnectionState_PEER_CONNECTION_STATE_CONNECTED
		case webrtc.PeerConnectionStateDisconnected:
			dbState = entity.PeerConnectionState_PEER_CONNECTION_STATE_DISCONNECTED
		case webrtc.PeerConnectionStateFailed:
			dbState = entity.PeerConnectionState_PEER_CONNECTION_STATE_FAILED
		case webrtc.PeerConnectionStateClosed:
			dbState = entity.PeerConnectionState_PEER_CONNECTION_STATE_CLOSED
		}

		ctx := context.Background()
		_ = s.repo.UpdatePeerState(ctx, peerID, dbState)

		// Notify peers about state change
		s.broadcastPeerStateChanged(ctx, roomID, peerID, dbState)
	})

	// Handle track events
	pc.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		ctx := context.Background()

		// Determine track type
		var trackType entity.TrackType
		switch track.Kind() {
		case webrtc.RTPCodecTypeAudio:
			trackType = entity.TrackType_TRACK_TYPE_AUDIO
		case webrtc.RTPCodecTypeVideo:
			trackType = entity.TrackType_TRACK_TYPE_VIDEO
		}

		// Create track entity in database
		trackEntity := &entity.Track{
			TrackId: track.ID(),
			PeerId:  peerID,
			Type:    trackType,
			Label:   track.StreamID(), // Use StreamID instead of Label
			Muted:   false,
			State:   entity.TrackState_TRACK_STATE_ACTIVE,
		}
		_ = s.repo.PublishTrack(ctx, trackEntity)

		// Broadcast track published event
		s.broadcastTrackPublished(ctx, roomID, peerID, trackEntity)
	})

	// Handle negotiation needed
	pc.OnNegotiationNeeded(func() {
		// Create and send new offer
		offer, err := pc.CreateOffer(nil)
		if err != nil {
			return
		}

		if err := pc.SetLocalDescription(offer); err != nil {
			return
		}

		// Send renegotiation required event
		select {
		case outgoing <- &service.SignalResponse{
			PeerId:    peerID,
			RoomId:    roomID,
			Timestamp: timestamppb.Now(),
			Payload: &service.SignalResponse_OfferReceived{
				OfferReceived: &events.OfferReceivedEvent{
					RoomId: roomID,
					Offer: &entity.SDPOffer{
						FromPeerId: "server",
						ToPeerId:   peerID,
						Sdp:        offer.SDP,
						Type:       entity.SDPType_SDP_TYPE_OFFER,
					},
					ReceivedAt: timestamppb.Now(),
				},
			},
		}:
		default:
		}
	})
}

// handleOfferWithSFU handles SDP offer from client
func (s *SignalingServiceWithSFU) handleOfferWithSFU(
	ctx context.Context,
	peerID, roomID string,
	req *service.SendOfferRequest,
	outgoing chan *service.SignalResponse,
	sfuPeer *mediaserver.Peer,
) error {
	if sfuPeer == nil {
		return status.Error(codes.Internal, "SFU peer not initialized")
	}

	pc := sfuPeer.GetPeerConnection()

	// Set remote description (client's offer)
	offer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  req.Sdp,
	}

	if err := pc.SetRemoteDescription(offer); err != nil {
		return status.Errorf(codes.InvalidArgument, "failed to set remote description: %v", err)
	}

	// Create answer
	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to create answer: %v", err)
	}

	// Set local description (server's answer)
	if err := pc.SetLocalDescription(answer); err != nil {
		return status.Errorf(codes.Internal, "failed to set local description: %v", err)
	}

	// Send answer back to client
	select {
	case outgoing <- &service.SignalResponse{
		PeerId:    peerID,
		RoomId:    roomID,
		Timestamp: timestamppb.Now(),
		Payload: &service.SignalResponse_AnswerReceived{
			AnswerReceived: &events.AnswerReceivedEvent{
				RoomId: roomID,
				Answer: &entity.SDPAnswer{
					FromPeerId: "server",
					ToPeerId:   peerID,
					Sdp:        answer.SDP,
					Type:       entity.SDPType_SDP_TYPE_ANSWER,
				},
				ReceivedAt: timestamppb.Now(),
			},
		},
	}:
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

// handleAnswerWithSFU handles SDP answer from client
func (s *SignalingServiceWithSFU) handleAnswerWithSFU(
	ctx context.Context,
	peerID, roomID string,
	req *service.SendAnswerRequest,
	outgoing chan *service.SignalResponse,
	sfuPeer *mediaserver.Peer,
) error {
	if sfuPeer == nil {
		return status.Error(codes.Internal, "SFU peer not initialized")
	}

	pc := sfuPeer.GetPeerConnection()

	// Set remote description (client's answer)
	answer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer,
		SDP:  req.Sdp,
	}

	if err := pc.SetRemoteDescription(answer); err != nil {
		return status.Errorf(codes.InvalidArgument, "failed to set remote description: %v", err)
	}

	return nil
}

// handleICECandidateWithSFU handles ICE candidate from client
func (s *SignalingServiceWithSFU) handleICECandidateWithSFU(
	ctx context.Context,
	peerID string,
	req *service.SendICECandidateRequest,
	sfuPeer *mediaserver.Peer,
) error {
	if sfuPeer == nil {
		return status.Error(codes.Internal, "SFU peer not initialized")
	}

	pc := sfuPeer.GetPeerConnection()

	// Add ICE candidate
	candidate := webrtc.ICECandidateInit{
		Candidate:     req.Candidate,
		SDPMid:        &req.SdpMid,
		SDPMLineIndex: func() *uint16 { v := uint16(req.SdpMLineIndex); return &v }(),
	}

	if err := pc.AddICECandidate(candidate); err != nil {
		return status.Errorf(codes.InvalidArgument, "failed to add ICE candidate: %v", err)
	}

	return nil
}

// SendOffer sends an SDP offer (peer-to-peer signaling, not used in SFU mode)
func (s *SignalingServiceWithSFU) SendOffer(ctx context.Context, req *service.SendOfferRequest) (*service.SendOfferResponse, error) {
	// In SFU mode, offers are handled through the Connect stream
	return &service.SendOfferResponse{Success: true}, nil
}

// SendAnswer sends an SDP answer (peer-to-peer signaling, not used in SFU mode)
func (s *SignalingServiceWithSFU) SendAnswer(ctx context.Context, req *service.SendAnswerRequest) (*service.SendAnswerResponse, error) {
	// In SFU mode, answers are handled through the Connect stream
	return &service.SendAnswerResponse{Success: true}, nil
}

// SendICECandidate sends an ICE candidate (peer-to-peer signaling, not used in SFU mode)
func (s *SignalingServiceWithSFU) SendICECandidate(ctx context.Context, req *service.SendICECandidateRequest) (*service.SendICECandidateResponse, error) {
	// In SFU mode, ICE candidates are handled through the Connect stream
	return &service.SendICECandidateResponse{Success: true}, nil
}

// Helper methods

func (s *SignalingServiceWithSFU) addPeerToRoom(roomID, peerID string) {
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

func (s *SignalingServiceWithSFU) removePeerFromRoom(roomID, peerID string) {
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

func (s *SignalingServiceWithSFU) cleanupPeerConnection(ctx context.Context, peerID, roomID string) {
	s.streams.Delete(peerID)
	s.peerConns.Delete(peerID)

	if roomID != "" {
		// Remove from SFU
		_ = s.sfu.RemovePeer(roomID, peerID)

		// Remove from room
		s.removePeerFromRoom(roomID, peerID)

		// Notify other peers
		_ = s.broadcastPeerLeft(ctx, roomID, peerID)
	}
}

func (s *SignalingServiceWithSFU) broadcastPeerJoined(ctx context.Context, roomID, newPeerID string) error {
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

func (s *SignalingServiceWithSFU) broadcastPeerLeft(ctx context.Context, roomID, leftPeerID string) error {
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

func (s *SignalingServiceWithSFU) broadcastPeerStateChanged(ctx context.Context, roomID, peerID string, state entity.PeerConnectionState) {
	msg := &service.SignalResponse{
		PeerId:    peerID,
		RoomId:    roomID,
		Timestamp: timestamppb.Now(),
		Payload: &service.SignalResponse_PeerStateChanged{
			PeerStateChanged: &events.PeerStateChangedEvent{
				RoomId:    roomID,
				PeerId:    peerID,
				OldState:  entity.PeerConnectionState_PEER_CONNECTION_STATE_UNSPECIFIED,
				NewState:  state,
				ChangedAt: timestamppb.Now(),
			},
		},
	}

	s.broadcastToRoom(roomID, peerID, msg)
}

func (s *SignalingServiceWithSFU) broadcastTrackPublished(ctx context.Context, roomID, peerID string, track *entity.Track) {
	msg := &service.SignalResponse{
		PeerId:    peerID,
		RoomId:    roomID,
		Timestamp: timestamppb.Now(),
		Payload: &service.SignalResponse_TrackPublished{
			TrackPublished: &events.TrackPublishedEvent{
				RoomId:      roomID,
				PeerId:      peerID,
				Track:       track,
				PublishedAt: timestamppb.Now(),
			},
		},
	}

	s.broadcastToRoom(roomID, "", msg)
}

func (s *SignalingServiceWithSFU) broadcastToRoom(roomID, senderPeerID string, msg *service.SignalResponse) error {
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

func (s *SignalingServiceWithSFU) sendToPeer(peerID string, msg *service.SignalResponse) error {
	val, ok := s.streams.Load(peerID)
	if !ok {
		return fmt.Errorf("peer not connected: %s", peerID)
	}

	info := val.(*streamInfo)
	select {
	case info.channel <- msg:
		return nil
	default:
		return fmt.Errorf("peer message queue full: %s", peerID)
	}
}
