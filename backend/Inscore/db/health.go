package db

import (
	"context"
	"time"

	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"gorm.io/gorm"
)

// startHealthCheck starts the health check routine for database connections
func (dm *DatabaseManager) startHealthCheck() {
	ticker := time.NewTicker(dm.config.Database.Failover.HealthCheckInterval)
	defer ticker.Stop()

	appLogger.Infof("Starting health check with interval: %v", dm.config.Database.Failover.HealthCheckInterval)

	for {
		select {
		case <-ticker.C:
			dm.performHealthCheck()
		case <-dm.healthCheckStop:
			appLogger.Info("Health check stopped")
			return
		}
	}
}

// performHealthCheck checks the health of both database connections
func (dm *DatabaseManager) performHealthCheck() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check primary database health
	primaryHealthy := dm.checkDatabaseHealth(dm.primaryDB, "primary", ctx)

	// Check backup database health
	backupHealthy := dm.checkDatabaseHealth(dm.backupDB, "backup", ctx)

	dm.mu.Lock()
	if primaryHealthy {
		dm.primaryStatus = StatusHealthy
	} else {
		dm.primaryStatus = StatusUnhealthy
	}

	if backupHealthy {
		dm.backupStatus = StatusHealthy
	} else {
		dm.backupStatus = StatusUnhealthy
	}
	dm.mu.Unlock()

	// Handle failover logic
	dm.handleFailover(primaryHealthy, backupHealthy)
}

// checkDatabaseHealth checks if a database connection is healthy
func (dm *DatabaseManager) checkDatabaseHealth(db *gorm.DB, dbType string, ctx context.Context) bool {
	if db == nil {
		return false
	}

	sqlDB, err := db.DB()
	if err != nil {
		appLogger.Warnf("Failed to get SQL DB for %s: %v", dbType, err)
		return false
	}

	// Ping with context timeout
	if err := sqlDB.PingContext(ctx); err != nil {
		appLogger.Warnf("Health check failed for %s database: %v", dbType, err)
		return false
	}

	// Check connection pool stats
	stats := sqlDB.Stats()
	if stats.OpenConnections > stats.MaxOpenConnections {
		appLogger.Warnf("%s database connection pool exhausted", dbType)
	}

	return true
}

// handleFailover handles database failover logic
func (dm *DatabaseManager) handleFailover(primaryHealthy, backupHealthy bool) {
	currentType := dm.GetCurrentType()

	// If currently using primary and it becomes unhealthy, switch to backup
	if currentType == Primary && !primaryHealthy && backupHealthy {
		appLogger.Info("Primary database unhealthy, attempting failover to backup")
		if err := dm.switchToBackup(); err != nil {
			appLogger.Errorf("Failed to switch to backup database: %v", err)
		} else {
			appLogger.Info("Successfully failed over to backup database")
		}
		return
	}

	// If currently using backup and primary becomes healthy again, switch back
	if currentType == Backup && primaryHealthy && dm.config.Database.Failover.AutoSwitchBack {
		// Wait for switch back delay to ensure primary is stable
		go dm.delayedSwitchBack()
		return
	}

	// If currently using backup and it becomes unhealthy, try primary
	if currentType == Backup && !backupHealthy && primaryHealthy {
		appLogger.Info("Backup database unhealthy, switching back to primary")
		if err := dm.switchToPrimary(); err != nil {
			appLogger.Errorf("Failed to switch to primary database: %v", err)
		} else {
			appLogger.Info("Successfully switched back to primary database")
		}
		return
	}

	// If both databases are unhealthy, log critical error
	if !primaryHealthy && !backupHealthy {
		appLogger.Error("CRITICAL: Both primary and backup databases are unhealthy!")
	}
}

// delayedSwitchBack waits for the configured delay before switching back to primary
func (dm *DatabaseManager) delayedSwitchBack() {
	time.Sleep(dm.config.Database.Failover.SwitchBackDelay)

	// Re-check primary health before switching
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if dm.checkDatabaseHealth(dm.primaryDB, "primary", ctx) {
		if err := dm.switchToPrimary(); err != nil {
			appLogger.Errorf("Failed to switch back to primary database: %v", err)
		} else {
			appLogger.Info("Successfully switched back to primary database after delay")
		}
	} else {
		appLogger.Info("Primary database still unhealthy, staying on backup")
	}
}

// ForceFailover manually triggers a failover to backup database
func (dm *DatabaseManager) ForceFailover() error {
	appLogger.Info("Manual failover triggered")
	return dm.switchToBackup()
}

// ForceSwitchBack manually switches back to primary database
func (dm *DatabaseManager) ForceSwitchBack() error {
	appLogger.Info("Manual switch back triggered")
	return dm.switchToPrimary()
}

// IsHealthy returns true if the current database is healthy
func (dm *DatabaseManager) IsHealthy() bool {
	primaryStatus, backupStatus := dm.GetStatus()
	currentType := dm.GetCurrentType()

	if currentType == Primary {
		return primaryStatus == StatusHealthy
	}
	return backupStatus == StatusHealthy
}

// GetDetailedStatus returns detailed status information
func (dm *DatabaseManager) GetDetailedStatus() map[string]interface{} {
	primaryStatus, backupStatus := dm.GetStatus()
	currentType := dm.GetCurrentType()
	metrics := dm.GetMetrics()

	return map[string]interface{}{
		"current_database":    string(currentType),
		"primary_status":      string(primaryStatus),
		"backup_status":       string(backupStatus),
		"primary_connections": metrics.PrimaryConnections,
		"backup_connections":  metrics.BackupConnections,
		"failover_count":      metrics.FailoverCount,
		"last_failover_time":  metrics.LastFailoverTime,
		"last_sync_time":      metrics.LastSyncTime,
		"last_backup_time":    metrics.LastBackupTime,
		"slow_queries":        metrics.SlowQueries,
	}
}
