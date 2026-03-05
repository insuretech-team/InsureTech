package db

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// DatabaseConfig represents the complete database configuration
type DatabaseConfig struct {
	Database struct {
		Primary        DatabaseConnection `yaml:"primary"`
		Backup         DatabaseConnection `yaml:"backup"`
		Failover       FailoverConfig     `yaml:"failover"`
		Sync           SyncConfig         `yaml:"sync"`
		BackupSettings BackupConfig       `yaml:"backup_settings"`
		Monitoring     MonitoringConfig   `yaml:"monitoring"`
	} `yaml:"database"`
}

// DatabaseConnection represents a single database connection configuration
type DatabaseConnection struct {
	Provider string `yaml:"provider"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"ssl_mode"`

	// Pool connection configuration (optional)
	PoolPort     string `yaml:"pool_port"`
	PoolDatabase string `yaml:"pool_database"`

	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

// FailoverConfig represents failover configuration
type FailoverConfig struct {
	Enabled             bool          `yaml:"enabled"`
	HealthCheckInterval time.Duration `yaml:"health_check_interval"`
	MaxRetryAttempts    int           `yaml:"max_retry_attempts"`
	RetryDelay          time.Duration `yaml:"retry_delay"`
	AutoSwitchBack      bool          `yaml:"auto_switch_back"`
	SwitchBackDelay     time.Duration `yaml:"switch_back_delay"`
}

// SyncConfig represents synchronization configuration
type SyncConfig struct {
	Enabled       bool          `yaml:"enabled"`
	Interval      time.Duration `yaml:"interval"`
	BatchSize     int           `yaml:"batch_size"`
	TablesToSync  []string      `yaml:"tables_to_sync"`
	ExcludeTables []string      `yaml:"exclude_tables"`
}

// BackupConfig represents backup configuration
type BackupConfig struct {
	AutoBackup     bool          `yaml:"auto_backup"`
	BackupInterval time.Duration `yaml:"backup_interval"`
	RetentionDays  int           `yaml:"retention_days"`
	Compression    bool          `yaml:"compression"`
	BackupPath     string        `yaml:"backup_path"`
}

// MonitoringConfig represents monitoring configuration
type MonitoringConfig struct {
	Enabled              bool          `yaml:"enabled"`
	MetricsInterval      time.Duration `yaml:"metrics_interval"`
	ConnectionPoolAlerts bool          `yaml:"connection_pool_alerts"`
}

// LoadDatabaseConfig loads database configuration from YAML file with environment variable substitution
func LoadDatabaseConfig(configPath string) (*DatabaseConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		// Fallback: build config from environment variables (Neon/PG*), used when no YAML is provided
		cfg, envErr := loadConfigFromEnv()
		if envErr == nil {
			return cfg, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Expand environment variables
	expandedData := os.ExpandEnv(string(data))

	var config DatabaseConfig
	if err := yaml.Unmarshal([]byte(expandedData), &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	// Validate and set defaults
	if err := validateAndSetDefaults(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// loadConfigFromEnv constructs DatabaseConfig from environment variables.
// Expected vars for primary: PGHOST, PGPORT, PGDATABASE, PGUSER, PGPASSWORD, PGSSLMODE
// For backup: PGHOST2, PGPORT2 (optional; falls back to PGPORT), PGDATABASE2, PGUSER2, PGPASSWORD2, PGSSLMODE2
func loadConfigFromEnv() (*DatabaseConfig, error) {
	primary := DatabaseConnection{
		Provider: "digitalocean",
		Host:     os.Getenv("PGHOST"),
		Port:     firstNonEmpty(os.Getenv("PGPORT"), "25060"),
		Database: os.Getenv("PGDATABASE"),
		Username: os.Getenv("PGUSER"),
		Password: os.Getenv("PGPASSWORD"),
		SSLMode:  firstNonEmpty(os.Getenv("PGSSLMODE"), "require"),
	}

	backup := DatabaseConnection{
		Provider: "neon",
		Host:     os.Getenv("PGHOST2"),
		Port:     firstNonEmpty(os.Getenv("PGPORT2"), os.Getenv("PGPORT"), "5432"),
		Database: os.Getenv("PGDATABASE2"),
		Username: os.Getenv("PGUSER2"),
		Password: os.Getenv("PGPASSWORD2"),
		SSLMode:  firstNonEmpty(os.Getenv("PGSSLMODE2"), os.Getenv("PGSSLMODE"), "require"),
	}

	cfg := &DatabaseConfig{}
	cfg.Database.Primary = primary
	cfg.Database.Backup = backup

	// Validate/defaults
	if err := validateAndSetDefaults(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// firstNonEmpty returns the first non-empty value from the provided args.
func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

// validateAndSetDefaults validates the configuration and sets default values
func validateAndSetDefaults(config *DatabaseConfig) error {
	// Validate primary database
	if err := validateConnection(&config.Database.Primary, "primary"); err != nil {
		return err
	}

	// Validate backup database
	if err := validateConnection(&config.Database.Backup, "backup"); err != nil {
		return err
	}

	// Set defaults for failover
	if config.Database.Failover.HealthCheckInterval == 0 {
		config.Database.Failover.HealthCheckInterval = 30 * time.Second
	}
	if config.Database.Failover.MaxRetryAttempts == 0 {
		config.Database.Failover.MaxRetryAttempts = 3
	}
	if config.Database.Failover.RetryDelay == 0 {
		config.Database.Failover.RetryDelay = 5 * time.Second
	}

	// Set defaults for sync
	if config.Database.Sync.Interval == 0 {
		config.Database.Sync.Interval = 1 * time.Hour
	}
	if config.Database.Sync.BatchSize == 0 {
		config.Database.Sync.BatchSize = 1000
	}

	// Set defaults for backup
	if config.Database.BackupSettings.BackupInterval == 0 {
		config.Database.BackupSettings.BackupInterval = 6 * time.Hour
	}
	if config.Database.BackupSettings.RetentionDays == 0 {
		config.Database.BackupSettings.RetentionDays = 30
	}
	if config.Database.BackupSettings.BackupPath == "" {
		config.Database.BackupSettings.BackupPath = "./backups"
	}

	return nil
}

// validateConnection validates a database connection configuration
func validateConnection(conn *DatabaseConnection, name string) error {
	if conn.Host == "" {
		return fmt.Errorf("%s database host is required", name)
	}
	if conn.Database == "" {
		return fmt.Errorf("%s database name is required", name)
	}
	if conn.Username == "" {
		return fmt.Errorf("%s database username is required", name)
	}
	if conn.Password == "" {
		return fmt.Errorf("%s database password is required", name)
	}

	// Set defaults
	if conn.Port == "" {
		// DigitalOcean managed Postgres uses 25060 by default
		if strings.EqualFold(conn.Provider, "digitalocean") {
			conn.Port = "25060"
		} else {
			conn.Port = "5432"
		}
	}
	if conn.SSLMode == "" {
		conn.SSLMode = "require"
	}
	if conn.MaxOpenConns == 0 {
		conn.MaxOpenConns = 25
	}
	if conn.MaxIdleConns == 0 {
		conn.MaxIdleConns = 10
	}
	if conn.ConnMaxLifetime == 0 {
		conn.ConnMaxLifetime = 1 * time.Hour
	}

	return nil
}

// BuildDSN constructs the database connection string
func (conn *DatabaseConnection) BuildDSN() string {
	port := conn.Port
	if port == "" {
		if strings.EqualFold(conn.Provider, "digitalocean") {
			port = "25060"
		} else {
			port = "5432"
		}
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s connect_timeout=30",
		conn.Host, port, conn.Username, conn.Password, conn.Database, conn.SSLMode)

	// Don't add statement_timeout for connection pool port 25061
	if port != "25061" {
		dsn += " statement_timeout=60000"
	}

	// For DigitalOcean managed databases: Use proper SSL configuration
	if conn.Provider == "digitalocean" {
		// Smart cert path detection - try Docker path first, fallback to local if file doesn't exist
		var certPath string
		dockerPath := os.Getenv("DO_DB_SSL_ROOT_CERT")
		localPath := os.Getenv("DO_DB_SSL_ROOT_CERT_LOCAL")
		
		// Try Docker path first (if set and file exists)
		if dockerPath != "" {
			if _, err := os.Stat(dockerPath); err == nil {
				certPath = dockerPath
			}
		}
		// Fallback to local path if Docker path not found
		if certPath == "" && localPath != "" {
			if _, err := os.Stat(localPath); err == nil {
				certPath = localPath
			}
		}
		// Final fallback to generic DB_SSL_ROOT_CERT
		if certPath == "" {
			certPath = os.Getenv("DB_SSL_ROOT_CERT")
		}
		
		if certPath != "" {
			// If provided path is relative, resolve to absolute based on current working directory
			if !filepath.IsAbs(certPath) {
				if abs, err := filepath.Abs(certPath); err == nil {
					certPath = abs
				}
			}
			// Only add sslrootcert if cert file is explicitly provided
			dsn += " sslrootcert=" + certPath
		}
		// If no cert path provided, use sslmode=require without explicit cert
		// This allows system SSL verification to handle the connection
	} else if conn.Provider == "neon" {
		// Neon handles SSL internally
	}

	// Note: channel_binding is not a connection parameter - it's handled internally by the driver
	// Removed PGCHANNELBINDING env var usage as it causes "unsupported startup parameter" error

	return dsn
}

// BuildPoolDSN constructs the pool connection string (for application traffic)
func (conn *DatabaseConnection) BuildPoolDSN() string {
	if conn.PoolPort == "" || conn.PoolDatabase == "" {
		return "" // No pool configuration
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s connect_timeout=30",
		conn.Host, conn.PoolPort, conn.Username, conn.Password, conn.PoolDatabase, conn.SSLMode)

	// Include DO sslrootcert for pool connections when provided
	if conn.Provider == "digitalocean" {
		// Smart cert path detection - try Docker path first, fallback to local if file doesn't exist
		var certPath string
		dockerPath := os.Getenv("DO_DB_SSL_ROOT_CERT")
		localPath := os.Getenv("DO_DB_SSL_ROOT_CERT_LOCAL")
		
		// Try Docker path first (if set and file exists)
		if dockerPath != "" {
			if _, err := os.Stat(dockerPath); err == nil {
				certPath = dockerPath
			}
		}
		// Fallback to local path if Docker path not found
		if certPath == "" && localPath != "" {
			if _, err := os.Stat(localPath); err == nil {
				certPath = localPath
			}
		}
		// Final fallback to generic DB_SSL_ROOT_CERT
		if certPath == "" {
			certPath = os.Getenv("DB_SSL_ROOT_CERT")
		}
		
		if certPath != "" {
			if !filepath.IsAbs(certPath) {
				if abs, err := filepath.Abs(certPath); err == nil {
					certPath = abs
				}
			}
			dsn += " sslrootcert=" + certPath
		}
	}

	// Optional channel_binding for pool as well
	if cb := os.Getenv("PGCHANNELBINDING"); cb != "" {
		dsn += " channel_binding=" + cb
	} else if cb := os.Getenv("DB_CHANNEL_BINDING"); cb != "" {
		dsn += " channel_binding=" + cb
	}

	// Pool connections don't support statement_timeout
	return dsn
}

// GetConnectionDSN returns the appropriate DSN based on operation type
func (conn *DatabaseConnection) GetConnectionDSN(usePool bool) string {
	if usePool && conn.PoolPort != "" && conn.PoolDatabase != "" {
		return conn.BuildPoolDSN()
	}
	return conn.BuildDSN()
}

// GetEnvWithDefault gets environment variable with default value
func GetEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvAsInt gets environment variable as integer with default value
func GetEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvAsBool gets environment variable as boolean with default value
func GetEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
