package mediaserver

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/pion/webrtc/v4"
	"github.com/pion/interceptor"
)

// SFU (Selective Forwarding Unit) manages media routing for conferences
type SFU struct {
	rooms  sync.Map // roomID (string) -> *Room
	config *Config
	api    *webrtc.API
}

// Config holds SFU configuration
type Config struct {
	ICEServers []webrtc.ICEServer
	// STUN/TURN servers for NAT traversal
	
	// Media settings
	MaxBitrate      uint64 // Maximum bitrate per track in bps
	EnableSimulcast bool   // Enable multi-quality simulcast
	
	// Logging
	Debug bool
}

// Room represents a conference room with multiple peers
type Room struct {
	id    string
	peers sync.Map // peerID (string) -> *Peer
	mu    sync.RWMutex
}

// Peer represents a WebRTC peer connection
type Peer struct {
	id   string
	pc   *webrtc.PeerConnection
	room *Room
	
	// Tracks
	localTracks  []*webrtc.TrackLocalStaticRTP
	remoteTracks []*webrtc.TrackRemote
	
	// State
	connected bool
	mu        sync.RWMutex
}

// NewSFU creates a new SFU instance
func NewSFU(config *Config) (*SFU, error) {
	if config == nil {
		config = DefaultConfig()
	}
	
	// Create media engine
	m := &webrtc.MediaEngine{}
	
	// Register default codecs
	if err := m.RegisterDefaultCodecs(); err != nil {
		return nil, fmt.Errorf("failed to register codecs: %w", err)
	}
	
	// Create interceptor registry for RTCP handling
	i := &interceptor.Registry{}
	if err := webrtc.RegisterDefaultInterceptors(m, i); err != nil {
		return nil, fmt.Errorf("failed to register interceptors: %w", err)
	}
	
	// Create WebRTC API
	api := webrtc.NewAPI(
		webrtc.WithMediaEngine(m),
		webrtc.WithInterceptorRegistry(i),
	)
	
	return &SFU{
		config: config,
		api:    api,
	}, nil
}

// DefaultConfig returns default SFU configuration
func DefaultConfig() *Config {
	return &Config{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		},
		MaxBitrate:      2_500_000, // 2.5 Mbps
		EnableSimulcast: true,
		Debug:           false,
	}
}

// CreateRoom creates a new conference room
func (s *SFU) CreateRoom(roomID string) (*Room, error) {
	room := &Room{
		id: roomID,
	}
	
	s.rooms.Store(roomID, room)
	return room, nil
}

// GetRoom retrieves an existing room
func (s *SFU) GetRoom(roomID string) (*Room, bool) {
	val, ok := s.rooms.Load(roomID)
	if !ok {
		return nil, false
	}
	return val.(*Room), true
}

// DeleteRoom removes a room and disconnects all peers
func (s *SFU) DeleteRoom(roomID string) error {
	val, ok := s.rooms.Load(roomID)
	if !ok {
		return fmt.Errorf("room not found: %s", roomID)
	}
	
	room := val.(*Room)
	
	// Disconnect all peers
	room.peers.Range(func(key, value interface{}) bool {
		peer := value.(*Peer)
		peer.Close()
		return true
	})
	
	s.rooms.Delete(roomID)
	return nil
}

// AddPeer adds a new peer to a room
func (s *SFU) AddPeer(ctx context.Context, roomID, peerID string) (*Peer, error) {
	// Get or create room
	var room *Room
	val, ok := s.rooms.Load(roomID)
	if !ok {
		var err error
		room, err = s.CreateRoom(roomID)
		if err != nil {
			return nil, err
		}
	} else {
		room = val.(*Room)
	}
	
	// Create peer connection configuration
	config := webrtc.Configuration{
		ICEServers: s.config.ICEServers,
	}
	
	// Create peer connection
	pc, err := s.api.NewPeerConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create peer connection: %w", err)
	}
	
	peer := &Peer{
		id:        peerID,
		pc:        pc,
		room:      room,
		connected: false,
	}
	
	// Set up peer connection handlers
	s.setupPeerHandlers(peer)
	
	// Add peer to room
	room.peers.Store(peerID, peer)
	
	return peer, nil
}

// setupPeerHandlers configures event handlers for a peer connection
func (s *SFU) setupPeerHandlers(peer *Peer) {
	// Handle incoming tracks
	peer.pc.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		if s.config.Debug {
			fmt.Printf("Track received: %s, type: %s\n", track.ID(), track.Kind())
		}
		
		peer.mu.Lock()
		peer.remoteTracks = append(peer.remoteTracks, track)
		peer.mu.Unlock()
		
		// Forward track to all other peers in room
		go s.forwardTrack(peer, track)
	})
	
	// Handle ICE connection state
	peer.pc.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		if s.config.Debug {
			fmt.Printf("Peer %s ICE state: %s\n", peer.id, state.String())
		}
		
		if state == webrtc.ICEConnectionStateConnected {
			peer.mu.Lock()
			peer.connected = true
			peer.mu.Unlock()
		} else if state == webrtc.ICEConnectionStateFailed || 
			      state == webrtc.ICEConnectionStateDisconnected {
			peer.mu.Lock()
			peer.connected = false
			peer.mu.Unlock()
		}
	})
	
	// Handle peer connection state
	peer.pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		if s.config.Debug {
			fmt.Printf("Peer %s connection state: %s\n", peer.id, state.String())
		}
		
		if state == webrtc.PeerConnectionStateFailed || 
		   state == webrtc.PeerConnectionStateClosed {
			peer.Close()
		}
	})
}

// forwardTrack forwards a remote track to all other peers in the room
func (s *SFU) forwardTrack(sourcePeer *Peer, remoteTrack *webrtc.TrackRemote) {
	// Create local track for forwarding
	localTrack, err := webrtc.NewTrackLocalStaticRTP(
		remoteTrack.Codec().RTPCodecCapability,
		remoteTrack.ID(),
		remoteTrack.StreamID(),
	)
	if err != nil {
		if s.config.Debug {
			fmt.Printf("Failed to create local track: %v\n", err)
		}
		return
	}
	
	sourcePeer.mu.Lock()
	sourcePeer.localTracks = append(sourcePeer.localTracks, localTrack)
	sourcePeer.mu.Unlock()
	
	// Add track to all other peers in room
	sourcePeer.room.peers.Range(func(key, value interface{}) bool {
		targetPeer := value.(*Peer)
		
		// Don't send track back to source
		if targetPeer.id == sourcePeer.id {
			return true
		}
		
		// Add track to target peer
		if _, err := targetPeer.pc.AddTrack(localTrack); err != nil {
			if s.config.Debug {
				fmt.Printf("Failed to add track to peer %s: %v\n", targetPeer.id, err)
			}
			return true
		}
		
		// Trigger renegotiation
		go s.renegotiate(targetPeer)
		
		return true
	})
	
	// Read RTP packets and forward
	buf := make([]byte, 1500)
	for {
		n, _, err := remoteTrack.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			}
			if s.config.Debug {
				fmt.Printf("Error reading track: %v\n", err)
			}
			return
		}
		
		// Write to local track (which forwards to all subscribed peers)
		if _, err = localTrack.Write(buf[:n]); err != nil {
			if s.config.Debug {
				fmt.Printf("Error writing to local track: %v\n", err)
			}
			return
		}
	}
}

// renegotiate triggers SDP renegotiation for a peer
func (s *SFU) renegotiate(peer *Peer) {
	// Create new offer
	offer, err := peer.pc.CreateOffer(nil)
	if err != nil {
		if s.config.Debug {
			fmt.Printf("Failed to create offer: %v\n", err)
		}
		return
	}
	
	// Set local description
	if err = peer.pc.SetLocalDescription(offer); err != nil {
		if s.config.Debug {
			fmt.Printf("Failed to set local description: %v\n", err)
		}
		return
	}
	
	// Note: In production, send offer to client via signaling
	// This requires integration with SignalingService
}

// RemovePeer removes a peer from its room
func (s *SFU) RemovePeer(roomID, peerID string) error {
	val, ok := s.rooms.Load(roomID)
	if !ok {
		return fmt.Errorf("room not found: %s", roomID)
	}
	
	room := val.(*Room)
	
	val, ok = room.peers.Load(peerID)
	if !ok {
		return fmt.Errorf("peer not found: %s", peerID)
	}
	
	peer := val.(*Peer)
	peer.Close()
	
	room.peers.Delete(peerID)
	
	return nil
}

// Close closes a peer connection and cleans up resources
func (p *Peer) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.pc != nil {
		return p.pc.Close()
	}
	
	return nil
}

// GetPeerConnection returns the underlying peer connection
func (p *Peer) GetPeerConnection() *webrtc.PeerConnection {
	return p.pc
}

// IsConnected returns true if the peer is connected
func (p *Peer) IsConnected() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.connected
}

// GetPeerCount returns the number of peers in a room
func (r *Room) GetPeerCount() int {
	count := 0
	r.peers.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}
