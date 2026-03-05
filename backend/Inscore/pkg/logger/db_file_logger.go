package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// DatabaseLogger handles logging for database operations
type DatabaseLogger struct {
	logDir string
}

// NewDatabaseLogger creates a new database logger instance
func NewDatabaseLogger() *DatabaseLogger {
	logDir := filepath.Join("internal", "db", "logs")
	// Ensure log directory exists (best-effort)
	_ = os.MkdirAll(logDir, 0o755)

	return &DatabaseLogger{logDir: logDir}
}

// LogServerStart logs when the server starts with database initialization
func (dl *DatabaseLogger) LogServerStart(dbType string) {
	dl.writeLog("server.log", fmt.Sprintf("SERVER_START: Database type: %s", dbType))
}

// LogMigrationStart logs when migrations begin
func (dl *DatabaseLogger) LogMigrationStart(dbType string) {
	dl.writeLog("migrations.log", fmt.Sprintf("MIGRATION_START: Database: %s", dbType))
}

// LogMigrationFile logs individual migration file execution
func (dl *DatabaseLogger) LogMigrationFile(filename, dbType string, success bool) {
	status := "SUCCESS"
	if !success {
		status = "FAILED"
	}
	dl.writeLog("migrations.log", fmt.Sprintf("MIGRATION_FILE: %s on %s - %s", filename, dbType, status))
}

// LogMigrationComplete logs when all migrations are complete
func (dl *DatabaseLogger) LogMigrationComplete(dbType string, count int) {
	dl.writeLog("migrations.log", fmt.Sprintf("MIGRATION_COMPLETE: Database: %s, Files processed: %d", dbType, count))
}

// LogSeederStart logs when seeders begin
func (dl *DatabaseLogger) LogSeederStart(dbType string) {
	dl.writeLog("seeders.log", fmt.Sprintf("SEEDER_START: Database: %s", dbType))
}

// LogSeederFile logs individual seeder file execution
func (dl *DatabaseLogger) LogSeederFile(filename, dbType string, success bool) {
	status := "SUCCESS"
	if !success {
		status = "FAILED"
	}
	dl.writeLog("seeders.log", fmt.Sprintf("SEEDER_FILE: %s on %s - %s", filename, dbType, status))
}

// LogSeederComplete logs when all seeders are complete
func (dl *DatabaseLogger) LogSeederComplete(dbType string, count int) {
	dl.writeLog("seeders.log", fmt.Sprintf("SEEDER_COMPLETE: Database: %s, Files processed: %d", dbType, count))
}

// LogGormSeed logs GORM-based seeding operations
func (dl *DatabaseLogger) LogGormSeed(seedType, dbType string, success bool, recordCount int) {
	status := "SUCCESS"
	if !success {
		status = "FAILED"
	}
	dl.writeLog("gorm_seeds.log", fmt.Sprintf("GORM_SEED: %s on %s - %s, Records: %d", seedType, dbType, status, recordCount))
}

// LogDatabaseConnection logs database connection events
func (dl *DatabaseLogger) LogDatabaseConnection(dbType string, success bool, dsn string) {
	status := "SUCCESS"
	if !success {
		status = "FAILED"
	}
	// Mask sensitive information in DSN
	maskedDSN := dl.maskSensitiveInfo(dsn)
	dl.writeLog("connections.log", fmt.Sprintf("CONNECTION: %s - %s, DSN: %s", dbType, status, maskedDSN))
}

// LogFailover logs database failover events
func (dl *DatabaseLogger) LogFailover(fromDB, toDB string) {
	dl.writeLog("failover.log", fmt.Sprintf("FAILOVER: From %s to %s", fromDB, toDB))
}

// LogSync logs database synchronization events
func (dl *DatabaseLogger) LogSync(operation string, recordCount int, success bool) {
	status := "SUCCESS"
	if !success {
		status = "FAILED"
	}
	dl.writeLog("sync.log", fmt.Sprintf("SYNC: %s - %s, Records: %d", operation, status, recordCount))
}

// LogError logs database errors
func (dl *DatabaseLogger) LogError(operation, dbType string, err error) {
	dl.writeLog("errors.log", fmt.Sprintf("ERROR: %s on %s - %v", operation, dbType, err))
}

// LogDoubleExecution logs when double execution is detected and prevented
func (dl *DatabaseLogger) LogDoubleExecution(operation, dbType string, prevented bool) {
	status := "DETECTED"
	if prevented {
		status = "PREVENTED"
	}
	dl.writeLog("double_execution.log", fmt.Sprintf("DOUBLE_EXECUTION: %s on %s - %s", operation, dbType, status))
}

// writeLog writes a timestamped entry to the specified log file
func (dl *DatabaseLogger) writeLog(filename, message string) {
	logPath := filepath.Join(dl.logDir, filename)

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return // Fail silently to avoid disrupting database operations
	}
	defer file.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	logEntry := fmt.Sprintf("[%s] %s\n", timestamp, message)
	_, _ = file.WriteString(logEntry)
}

// maskSensitiveInfo masks passwords and sensitive information in DSN strings
func (dl *DatabaseLogger) maskSensitiveInfo(dsn string) string {
	// Simple masking - replace a middle part with ***MASKED*** to avoid leaking secrets
	if len(dsn) > 50 {
		return dsn[:20] + "***MASKED***" + dsn[len(dsn)-10:]
	}
	return "***MASKED***"
}

// GetLogSummary returns a summary of recent database operations
func (dl *DatabaseLogger) GetLogSummary() map[string]int {
	summary := make(map[string]int)

	logFiles := []string{
		"server.log",
		"migrations.log",
		"seeders.log",
		"gorm_seeds.log",
		"connections.log",
		"failover.log",
		"sync.log",
		"errors.log",
		"double_execution.log",
	}

	for _, logFile := range logFiles {
		logPath := filepath.Join(dl.logDir, logFile)
		if info, err := os.Stat(logPath); err == nil {
			summary[logFile] = int(info.Size())
		}
	}

	return summary
}

// Global database logger instance
var dbLogger *DatabaseLogger

// InitDatabaseLogger initializes the global database logger
func InitDatabaseLogger() { dbLogger = NewDatabaseLogger() }

// GetDatabaseLogger returns the global database logger instance
func GetDatabaseLogger() *DatabaseLogger {
	if dbLogger == nil {
		InitDatabaseLogger()
	}
	return dbLogger
}
