// Package webrtc provides production-grade WebRTC conferencing services
// for the LifePlusCore telemedicine platform.
//
// Architecture:
//   - core: Configuration and service management
//   - services: gRPC service implementations
//   - repository: Data access layer
//   - mediaserver: SFU media routing
//   - turn: TURN server integration
//
// Usage:
//   import "github.com/fahara02/lifepluscore/lpc/pkg/webrtc"
//   import "github.com/fahara02/lifepluscore/lpc/pkg/webrtc/core"
//   
//   config := core.LoadFromEnv()
//   manager := webrtc.NewServiceManager(db, config)
package webrtc

import (
	"gorm.io/gorm"
	
	"github.com/newage-saint/insuretech/backend/inscore/pkg/webrtc/core"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/webrtc/repository"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/webrtc/services"
)

// ServiceManager manages all WebRTC services
type ServiceManager struct {
	Repository       *repository.Repository
	RoomService      *services.RoomService
	PeerService      *services.PeerService
	TrackService     *services.TrackService
	SignalingService *services.SignalingService
	StatsService     *services.StatsService
}

// NewServiceManager creates a new WebRTC service manager with all services initialized
func NewServiceManager(db *gorm.DB) *ServiceManager {
	repo := repository.NewRepository(db)

	return &ServiceManager{
		Repository:       repo,
		RoomService:      services.NewRoomService(repo),
		PeerService:      services.NewPeerService(repo),
		TrackService:     services.NewTrackService(repo),
		SignalingService: services.NewSignalingService(repo),
		StatsService:     services.NewStatsService(repo),
	}
}

// GetRepository returns the repository instance
func (sm *ServiceManager) GetRepository() *repository.Repository {
	return sm.Repository
}

// GetRoomService returns the room service instance
func (sm *ServiceManager) GetRoomService() *services.RoomService {
	return sm.RoomService
}

// GetPeerService returns the peer service instance
func (sm *ServiceManager) GetPeerService() *services.PeerService {
	return sm.PeerService
}

// GetTrackService returns the track service instance
func (sm *ServiceManager) GetTrackService() *services.TrackService {
	return sm.TrackService
}

// GetSignalingService returns the signaling service instance
func (sm *ServiceManager) GetSignalingService() *services.SignalingService {
	return sm.SignalingService
}

// GetStatsService returns the stats service instance
func (sm *ServiceManager) GetStatsService() *services.StatsService {
	return sm.StatsService
}

// Re-export commonly used types and functions from core for convenience
type (
	// Config types
	ProductionConfig   = core.ProductionConfig
	ICEServerConfig    = core.ICEServerConfig
	TLSConfig          = core.TLSConfig
	SecurityConfig     = core.SecurityConfig
	MediaConfig        = core.MediaConfig
	TelemedicineConfig = core.TelemedicineConfig
	
	// State types
	ConnectionState = core.ConnectionState
	RoomState       = core.RoomState
	TrackKind       = core.TrackKind
	TrackState      = core.TrackState
)

// Re-export config functions
var (
	DefaultProductionConfig = core.DefaultProductionConfig
	DevelopmentConfig       = core.DevelopmentConfig
	LoadFromEnv             = core.LoadFromEnv
)

// Re-export state constants
const (
	ConnectionStateNew          = core.ConnectionStateNew
	ConnectionStateConnecting   = core.ConnectionStateConnecting
	ConnectionStateConnected    = core.ConnectionStateConnected
	ConnectionStateDisconnected = core.ConnectionStateDisconnected
	ConnectionStateFailed       = core.ConnectionStateFailed
	ConnectionStateClosed       = core.ConnectionStateClosed
	
	RoomStateActive = core.RoomStateActive
	RoomStateIdle   = core.RoomStateIdle
	RoomStateClosed = core.RoomStateClosed
	
	TrackKindAudio  = core.TrackKindAudio
	TrackKindVideo  = core.TrackKindVideo
	TrackKindScreen = core.TrackKindScreen
	
	TrackStateActive = core.TrackStateActive
	TrackStateMuted  = core.TrackStateMuted
	TrackStateEnded  = core.TrackStateEnded
)
