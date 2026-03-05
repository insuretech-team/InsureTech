package core

import (
	"gorm.io/gorm"
)

// ServiceManager manages all WebRTC services
// This is the main entry point for WebRTC functionality
type ServiceManager struct {
	db               *gorm.DB
	Repository       Repository
	RoomService      RoomService
	PeerService      PeerService
	TrackService     TrackService
	SignalingService SignalingService
	StatsService     StatsService
}

// Repository interface defines data access methods
type Repository interface {
	// Room operations
	CreateRoom(room interface{}) error
	GetRoom(roomID string) (interface{}, error)
	UpdateRoom(room interface{}) error
	DeleteRoom(roomID string) error
	ListRooms(filters interface{}) ([]interface{}, error)
	
	// Peer operations
	CreatePeer(peer interface{}) error
	GetPeer(peerID string) (interface{}, error)
	UpdatePeer(peer interface{}) error
	DeletePeer(peerID string) error
	ListPeersInRoom(roomID string) ([]interface{}, error)
	
	// Track operations
	CreateTrack(track interface{}) error
	GetTrack(trackID string) (interface{}, error)
	UpdateTrack(track interface{}) error
	DeleteTrack(trackID string) error
	
	// Session operations
	CreateSession(session interface{}) error
	GetSession(sessionID string) (interface{}, error)
	EndSession(sessionID string) error
}

// RoomService interface defines room management operations
type RoomService interface {
	CreateRoom(name string, config interface{}) (interface{}, error)
	GetRoom(roomID string) (interface{}, error)
	UpdateRoom(roomID string, config interface{}) error
	CloseRoom(roomID string) error
	ListRooms() ([]interface{}, error)
}

// PeerService interface defines peer management operations
type PeerService interface {
	JoinRoom(roomID, peerID, displayName string) (interface{}, error)
	LeaveRoom(roomID, peerID string) error
	GetPeer(peerID string) (interface{}, error)
	ListPeersInRoom(roomID string) ([]interface{}, error)
	UpdatePeerState(peerID string, state interface{}) error
}

// TrackService interface defines media track operations
type TrackService interface {
	PublishTrack(peerID string, trackInfo interface{}) (interface{}, error)
	UnpublishTrack(trackID string) error
	MuteTrack(trackID string) error
	UnmuteTrack(trackID string) error
	GetTrack(trackID string) (interface{}, error)
}

// SignalingService interface defines signaling operations
type SignalingService interface {
	// Bidirectional signaling stream
	Connect(stream interface{}) error
	
	// SDP exchange
	SendOffer(fromPeerID, toPeerID, sdp string) error
	SendAnswer(fromPeerID, toPeerID, sdp string) error
	
	// ICE candidate exchange
	SendICECandidate(fromPeerID, toPeerID, candidate string) error
}

// StatsService interface defines statistics operations
type StatsService interface {
	RecordPeerStats(peerID string, stats interface{}) error
	GetPeerStats(peerID string) (interface{}, error)
	GetRoomStats(roomID string) (interface{}, error)
	GetSystemStats() (interface{}, error)
}

// NewServiceManager creates a new WebRTC service manager
func NewServiceManager(
	db *gorm.DB,
	repo Repository,
	roomSvc RoomService,
	peerSvc PeerService,
	trackSvc TrackService,
	signalingSvc SignalingService,
	statsSvc StatsService,
) *ServiceManager {
	return &ServiceManager{
		db:               db,
		Repository:       repo,
		RoomService:      roomSvc,
		PeerService:      peerSvc,
		TrackService:     trackSvc,
		SignalingService: signalingSvc,
		StatsService:     statsSvc,
	}
}

// GetDB returns the database instance
func (sm *ServiceManager) GetDB() *gorm.DB {
	return sm.db
}
