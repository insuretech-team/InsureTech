package core

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// ICEServerConfig represents a STUN/TURN server configuration
type ICEServerConfig struct {
	URLs       []string `json:"urls" yaml:"urls"`
	Username   string   `json:"username,omitempty" yaml:"username,omitempty"`
	Credential string   `json:"credential,omitempty" yaml:"credential,omitempty"`
}

// TLSConfig represents TLS/HTTPS configuration
type TLSConfig struct {
	Enabled      bool   `json:"enabled" yaml:"enabled"`
	CertFile     string `json:"cert_file" yaml:"cert_file"`
	KeyFile      string `json:"key_file" yaml:"key_file"`
	AutoGenerate bool   `json:"auto_generate" yaml:"auto_generate"`
}

// SecurityConfig represents security settings for production
type SecurityConfig struct {
	JWTEnabled       bool     `json:"jwt_enabled" yaml:"jwt_enabled"`
	JWTSecret        string   `json:"jwt_secret" yaml:"jwt_secret"`
	JWTExpiration    int      `json:"jwt_expiration" yaml:"jwt_expiration"`
	RateLimitEnabled bool     `json:"rate_limit_enabled" yaml:"rate_limit_enabled"`
	MaxRoomsPerUser  int      `json:"max_rooms_per_user" yaml:"max_rooms_per_user"`
	MaxPeersPerRoom  int      `json:"max_peers_per_room" yaml:"max_peers_per_room"`
	AllowedOrigins   []string `json:"allowed_origins" yaml:"allowed_origins"`
	AllowedIPs       []string `json:"allowed_ips" yaml:"allowed_ips"`
	BlockedIPs       []string `json:"blocked_ips" yaml:"blocked_ips"`
}

// MediaConfig represents media quality settings
type MediaConfig struct {
	MaxVideoBitrate  int  `json:"max_video_bitrate" yaml:"max_video_bitrate"`
	MaxVideoWidth    int  `json:"max_video_width" yaml:"max_video_width"`
	MaxVideoHeight   int  `json:"max_video_height" yaml:"max_video_height"`
	MaxFrameRate     int  `json:"max_frame_rate" yaml:"max_frame_rate"`
	MaxAudioBitrate  int  `json:"max_audio_bitrate" yaml:"max_audio_bitrate"`
	EchoCancellation bool `json:"echo_cancellation" yaml:"echo_cancellation"`
	NoiseSuppression bool `json:"noise_suppression" yaml:"noise_suppression"`
	RecordingEnabled bool `json:"recording_enabled" yaml:"recording_enabled"`
	RecordingPath    string `json:"recording_path" yaml:"recording_path"`
}

// TelemedicineConfig represents telemedicine-specific features
type TelemedicineConfig struct {
	Enabled                bool     `json:"enabled" yaml:"enabled"`
	WaitingRoomEnabled     bool     `json:"waiting_room_enabled" yaml:"waiting_room_enabled"`
	MandatoryRecording     bool     `json:"mandatory_recording" yaml:"mandatory_recording"`
	ScreenSharingEnabled   bool     `json:"screen_sharing_enabled" yaml:"screen_sharing_enabled"`
	FileTransferEnabled    bool     `json:"file_transfer_enabled" yaml:"file_transfer_enabled"`
	MaxFileSize            int64    `json:"max_file_size" yaml:"max_file_size"`
	AllowedFileTypes       []string `json:"allowed_file_types" yaml:"allowed_file_types"`
	MaxConsultationDuration int     `json:"max_consultation_duration" yaml:"max_consultation_duration"`
	HIPAAComplianceEnabled bool     `json:"hipaa_compliance_enabled" yaml:"hipaa_compliance_enabled"`
	AuditLoggingEnabled    bool     `json:"audit_logging_enabled" yaml:"audit_logging_enabled"`
	EncryptionRequired     bool     `json:"encryption_required" yaml:"encryption_required"`
}

// ProductionConfig represents complete production configuration
type ProductionConfig struct {
	Host                  string             `json:"host" yaml:"host"`
	Port                  int                `json:"port" yaml:"port"`
	TLS                   TLSConfig          `json:"tls" yaml:"tls"`
	ICEServers            []ICEServerConfig  `json:"ice_servers" yaml:"ice_servers"`
	Security              SecurityConfig     `json:"security" yaml:"security"`
	Media                 MediaConfig        `json:"media" yaml:"media"`
	Telemedicine          TelemedicineConfig `json:"telemedicine" yaml:"telemedicine"`
	WebSocketReadTimeout  time.Duration      `json:"websocket_read_timeout" yaml:"websocket_read_timeout"`
	WebSocketWriteTimeout time.Duration      `json:"websocket_write_timeout" yaml:"websocket_write_timeout"`
	DatabaseEnabled       bool               `json:"database_enabled" yaml:"database_enabled"`
	MetricsEnabled        bool               `json:"metrics_enabled" yaml:"metrics_enabled"`
	MetricsPort           int                `json:"metrics_port" yaml:"metrics_port"`
	LogLevel              string             `json:"log_level" yaml:"log_level"`
}

// DefaultProductionConfig returns a secure production-ready configuration
func DefaultProductionConfig() *ProductionConfig {
	return &ProductionConfig{
		Host: "0.0.0.0",
		Port: 8443,
		TLS: TLSConfig{
			Enabled:      true,
			CertFile:     "/etc/inscore/certs/server.crt",
			KeyFile:      "/etc/inscore/certs/server.key",
			AutoGenerate: false,
		},
		ICEServers: []ICEServerConfig{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
			{URLs: []string{"stun:stun1.l.google.com:19302"}},
		},
		Security: SecurityConfig{
			JWTEnabled:       true,
			JWTSecret:        os.Getenv("JWT_SECRET"),
			JWTExpiration:    24,
			RateLimitEnabled: true,
			MaxRoomsPerUser:  10,
			MaxPeersPerRoom:  50,
			AllowedOrigins:   []string{"https://yourdomain.com"},
			AllowedIPs:       []string{},
			BlockedIPs:       []string{},
		},
		Media: MediaConfig{
			MaxVideoBitrate:  2500,
			MaxVideoWidth:    1920,
			MaxVideoHeight:   1080,
			MaxFrameRate:     30,
			MaxAudioBitrate:  128,
			EchoCancellation: true,
			NoiseSuppression: true,
			RecordingEnabled: false,
			RecordingPath:    "/var/inscore/recordings",
		},
		Telemedicine: TelemedicineConfig{
			Enabled:                 true,
			WaitingRoomEnabled:      true,
			MandatoryRecording:      false,
			ScreenSharingEnabled:    true,
			FileTransferEnabled:     true,
			MaxFileSize:             10 * 1024 * 1024,
			AllowedFileTypes:        []string{".pdf", ".jpg", ".png", ".dcm"},
			MaxConsultationDuration: 60,
			HIPAAComplianceEnabled:  true,
			AuditLoggingEnabled:     true,
			EncryptionRequired:      true,
		},
		WebSocketReadTimeout:  60 * time.Second,
		WebSocketWriteTimeout: 10 * time.Second,
		DatabaseEnabled:       true,
		MetricsEnabled:        true,
		MetricsPort:           9090,
		LogLevel:              "info",
	}
}

// DevelopmentConfig returns a configuration for local development
func DevelopmentConfig() *ProductionConfig {
	config := DefaultProductionConfig()
	config.Port = 8080
	config.TLS.Enabled = false
	config.Security.JWTEnabled = false
	config.Security.RateLimitEnabled = false
	config.Security.AllowedOrigins = []string{"*"}
	config.Telemedicine.HIPAAComplianceEnabled = false
	config.LogLevel = "debug"
	return config
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() *ProductionConfig {
	config := DefaultProductionConfig()
	
	if host := os.Getenv("WEBRTC_HOST"); host != "" {
		config.Host = host
	}
	
	if portStr := os.Getenv("WEBRTC_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			config.Port = port
		}
	}
	
	if tlsEnabled := os.Getenv("WEBRTC_TLS_ENABLED"); tlsEnabled != "" {
		config.TLS.Enabled = tlsEnabled == "true"
	}
	
	if certFile := os.Getenv("WEBRTC_TLS_CERT"); certFile != "" {
		config.TLS.CertFile = certFile
	}
	
	if keyFile := os.Getenv("WEBRTC_TLS_KEY"); keyFile != "" {
		config.TLS.KeyFile = keyFile
	}
	
	if turnURL := os.Getenv("TURN_URL"); turnURL != "" {
		config.ICEServers = append(config.ICEServers, ICEServerConfig{
			URLs:       []string{turnURL},
			Username:   os.Getenv("TURN_USERNAME"),
			Credential: os.Getenv("TURN_PASSWORD"),
		})
	}
	
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		config.Security.JWTSecret = jwtSecret
		config.Security.JWTEnabled = true
	}
	
	if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
		config.Security.AllowedOrigins = strings.Split(origins, ",")
	}
	
	if env := os.Getenv("ENVIRONMENT"); env == "development" {
		return DevelopmentConfig()
	}
	
	return config
}

// Validate checks if the configuration is valid
func (c *ProductionConfig) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("invalid port: %d", c.Port)
	}
	
	if c.TLS.Enabled {
		if c.TLS.CertFile == "" || c.TLS.KeyFile == "" {
			if !c.TLS.AutoGenerate {
				return fmt.Errorf("TLS enabled but cert/key files not specified and auto-generate is false")
			}
		}
	}
	
	if c.Security.JWTEnabled && c.Security.JWTSecret == "" {
		return fmt.Errorf("JWT enabled but secret not set")
	}
	
	if c.Telemedicine.HIPAAComplianceEnabled && !c.TLS.Enabled {
		return fmt.Errorf("HIPAA compliance requires TLS/HTTPS")
	}
	
	if len(c.ICEServers) == 0 {
		return fmt.Errorf("at least one ICE server (STUN/TURN) is required")
	}
	
	return nil
}

// GetTLSConfig returns a tls.Config for the server
func (c *ProductionConfig) GetTLSConfig() (*tls.Config, error) {
	if !c.TLS.Enabled {
		return nil, nil
	}
	
	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
		PreferServerCipherSuites: true,
	}, nil
}

// Address returns the server address (host:port)
func (c *ProductionConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// IsProduction returns true if running in production mode
func (c *ProductionConfig) IsProduction() bool {
	return c.TLS.Enabled && c.Security.JWTEnabled
}

// GetICEServersJSON returns ICE servers in JSON format for client
func (c *ProductionConfig) GetICEServersJSON() []map[string]interface{} {
	servers := make([]map[string]interface{}, 0, len(c.ICEServers))
	
	for _, server := range c.ICEServers {
		s := map[string]interface{}{
			"urls": server.URLs,
		}
		
		if server.Username != "" {
			s["username"] = server.Username
		}
		if server.Credential != "" {
			s["credential"] = server.Credential
		}
		
		servers = append(servers, s)
	}
	
	return servers
}
