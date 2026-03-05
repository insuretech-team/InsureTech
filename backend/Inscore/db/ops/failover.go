package ops

import (
	"github.com/newage-saint/insuretech/backend/inscore/db"
)

// ForceFailover manually triggers a failover to the backup database
func (api *ManagementAPI) ForceFailover() error {
	return api.manager.ForceFailover()
}

// ForceSwitchBack manually switches back to the primary database
func (api *ManagementAPI) ForceSwitchBack() error {
	return api.manager.ForceSwitchBack()
}

// RebuildBackupSchema rebuilds the backup database schema to match primary
func (api *ManagementAPI) RebuildBackupSchema() error {
	return api.manager.RebuildBackupSchema()
}

// ValidateSchemaConsistency checks if backup database schema matches primary
func (api *ManagementAPI) ValidateSchemaConsistency() ([]string, error) {
	return api.manager.ValidateSchemaConsistency()
}

// IsHealthy returns whether the database system is healthy
func (api *ManagementAPI) IsHealthy() bool {
	primaryStatus, backupStatus := api.manager.GetStatus()
	return primaryStatus == db.StatusHealthy || backupStatus == db.StatusHealthy
}

// GetDetailedStatus returns detailed status information for both databases
func (api *ManagementAPI) GetDetailedStatus() map[string]interface{} {
	primaryStatus, backupStatus := api.manager.GetStatus()
	metrics := api.manager.GetMetrics()

	return map[string]interface{}{
		"primary_status":      string(primaryStatus),
		"backup_status":       string(backupStatus),
		"current_db":          string(api.manager.GetCurrentType()),
		"primary_connections": metrics.PrimaryConnections,
		"backup_connections":  metrics.BackupConnections,
		"failover_count":      metrics.FailoverCount,
		"last_failover":       metrics.LastFailoverTime,
	}
}

