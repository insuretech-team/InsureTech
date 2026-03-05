package logger

import (
    "fmt"
    "sync"
    "strings"
    "time"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var (
    initOnce     sync.Once
    base         *zap.Logger
    initialized  bool
    mu           sync.RWMutex
)

// Config is a minimal logger configuration for CLI tools
type Config struct {
    Level   string // debug|info|warn|error
    Format  string // json|text
    Output  string // console (only for now)
    Verbose bool
}

// NoFileConfig returns a default console-oriented config used by dbmanager
func NoFileConfig() Config {
    return Config{Level: "info", Format: "text", Output: "console", Verbose: false}
}

// ANSI color codes
const (
    colorReset  = "\033[0m"
    colorCyan   = "\033[36m"
    colorGreen  = "\033[32m"
    colorYellow = "\033[33m"
    colorRed    = "\033[31m"
    colorPurple = "\033[35m"
    colorGray   = "\033[90m"
)

// Custom encoders for beautiful console output with colors
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
    enc.AppendString(fmt.Sprintf("%s[%s]%s", colorGray, t.Format("15:04"), colorReset))
}

func customLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
    var color string
    switch level {
    case zapcore.DebugLevel:
        color = colorPurple
    case zapcore.InfoLevel:
        color = colorGreen
    case zapcore.WarnLevel:
        color = colorYellow
    case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
        color = colorRed
    default:
        color = colorReset
    }
    enc.AppendString(fmt.Sprintf("%s[%s]%s", color, level.CapitalString(), colorReset))
}

func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
    enc.AppendString(fmt.Sprintf("%s[%s]%s\n  ", colorCyan, caller.TrimmedPath(), colorReset))
}

// Initialize configures the global logger according to the provided Config
// This will override any previous initialization
func Initialize(cfg Config) error {
    mu.Lock()
    defer mu.Unlock()
    
    var zcfg zap.Config
    if cfg.Format == "json" {
        zcfg = zap.NewProductionConfig()
    } else {
        zcfg = zap.NewDevelopmentConfig()
        zcfg.Encoding = "console"
        
        // Custom encoder config for beautiful console output with brackets
        zcfg.EncoderConfig.TimeKey = "T"
        zcfg.EncoderConfig.LevelKey = "L"
        zcfg.EncoderConfig.NameKey = "N"
        zcfg.EncoderConfig.CallerKey = "C"
        zcfg.EncoderConfig.MessageKey = "M"
        zcfg.EncoderConfig.StacktraceKey = "S"
        
        zcfg.EncoderConfig.EncodeLevel = customLevelEncoder
        zcfg.EncoderConfig.EncodeTime = customTimeEncoder
        zcfg.EncoderConfig.EncodeCaller = customCallerEncoder
        zcfg.EncoderConfig.EncodeName = zapcore.FullNameEncoder
        zcfg.EncoderConfig.LineEnding = zapcore.DefaultLineEnding
        zcfg.EncoderConfig.ConsoleSeparator = " "
        
        zcfg.DisableStacktrace = true
        zcfg.DisableCaller = false
    }
    // Level
    switch strings.ToLower(cfg.Level) {
    case "debug":
        zcfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
    case "warn":
        zcfg.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
    case "error":
        zcfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
    default:
        zcfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
    }
    l, err := zcfg.Build(zap.AddCallerSkip(1))
    if err != nil { return err }
    base = l
    initialized = true
    return nil
}

func initLogger() {
    cfg := zap.NewProductionConfig()
    cfg.EncoderConfig.TimeKey = "ts"
    cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    l, _ := cfg.Build()
    base = l
}

func GetLogger() *zap.Logger {
	mu.RLock()
	if initialized && base != nil {
		mu.RUnlock()
		return base
	}
	mu.RUnlock()
	
	initOnce.Do(initLogger)
	return base
}

// Convenience top-level logging
func Info(msg string, fields ...zap.Field)  { GetLogger().Info(msg, fields...) }
func Warn(msg string, fields ...zap.Field)  { GetLogger().Warn(msg, fields...) }
func Error(msg string, fields ...zap.Field) { GetLogger().Error(msg, fields...) }
func Debug(msg string, fields ...zap.Field) { GetLogger().Debug(msg, fields...) }

// Formatted helpers used by dbmanager
func Infof(format string, args ...interface{})  { GetLogger().Sugar().Infof(format, args...) }
func Warnf(format string, args ...interface{})  { GetLogger().Sugar().Warnf(format, args...) }
func Errorf(format string, args ...interface{}) { GetLogger().Sugar().Errorf(format, args...) }
func Fatal(msg string, fields ...zap.Field)     { GetLogger().Fatal(msg, fields...) }
func Fatalf(format string, args ...interface{}) { GetLogger().Sugar().Fatalf(format, args...) }

// Entry is a small wrapper for structured logging with pre-bound fields
type Entry struct {
	l      *zap.Logger
	fields []zap.Field
}

func WithError(err error) *Entry {
	return &Entry{l: GetLogger(), fields: []zap.Field{zap.Error(err)}}
}

func WithField(key string, val interface{}) *Entry {
	return &Entry{l: GetLogger(), fields: []zap.Field{zap.Any(key, val)}}
}

func (e *Entry) With(arg interface{}, vals ...interface{}) *Entry {
	// Flexible helper to support With(key, val) and With(zap.Field)
	switch v := arg.(type) {
	case zap.Field:
		return &Entry{l: e.l, fields: append(e.fields, v)}
	case string:
		var any interface{}
		if len(vals) > 0 {
			any = vals[0]
		}
		return &Entry{l: e.l, fields: append(e.fields, zap.Any(v, any))}
	default:
		return &Entry{l: e.l, fields: append(e.fields, zap.Any("field", v))}
	}
}

func (e *Entry) Info(msg string)  { e.l.With(e.fields...).Info(msg) }
func (e *Entry) Warn(msg string)  { e.l.With(e.fields...).Warn(msg) }
func (e *Entry) Error(msg string) { e.l.With(e.fields...).Error(msg) }
func (e *Entry) Debug(msg string) { e.l.With(e.fields...).Debug(msg) }

// Domain-specific named loggers
func SyncLogger() *zap.Logger                 { return GetLogger().Named("sync") }
func GetDatabaseOperationLogger() *DBOpLogger { return &DBOpLogger{l: GetLogger().Named("dbops")} }

// DBOpLogger provides helper methods for DB operations logging used by existing code
type DBOpLogger struct{ l *zap.Logger }

func (d *DBOpLogger) LogServerStart(mode string) {
	d.l.Info("server start", zap.String("mode", mode))
}

func (d *DBOpLogger) LogDatabaseConnection(kind string, success bool, host string) {
	lvl := zap.InfoLevel
	if !success {
		lvl = zap.ErrorLevel
	}
	ce := d.l.Check(lvl, "database connection")
	if ce != nil {
		ce.Write(zap.String("type", kind), zap.Bool("success", success), zap.String("host", host))
	}
}

func (d *DBOpLogger) LogMigrationStart(target string) {
	d.l.Info("migrations start", zap.String("target", target))
}
func (d *DBOpLogger) LogMigrationFile(name, target string, success bool) {
	lvl := zap.InfoLevel
	if !success {
		lvl = zap.ErrorLevel
	}
	ce := d.l.Check(lvl, "migration file")
	if ce != nil {
		ce.Write(zap.String("file", name), zap.String("target", target), zap.Bool("success", success))
	}
}
func (d *DBOpLogger) LogAllMigrationsComplete(target string, count int) {
	d.l.Info("migrations complete", zap.String("target", target), zap.Int("applied", count))
}

func (d *DBOpLogger) LogSeederStart(target string) {
	d.l.Info("seeders start", zap.String("target", target))
}
func (d *DBOpLogger) LogSeederFile(name, target string, success bool) {
	lvl := zap.InfoLevel
	if !success {
		lvl = zap.ErrorLevel
	}
	ce := d.l.Check(lvl, "seeder file")
	if ce != nil {
		ce.Write(zap.String("file", name), zap.String("target", target), zap.Bool("success", success))
	}
}
func (d *DBOpLogger) LogSeederComplete(target string, count int) {
	d.l.Info("seeders complete", zap.String("target", target), zap.Int("applied", count))
}

func (d *DBOpLogger) LogDoubleExecution(kind, target string, skipped bool) {
	d.l.Info("skip double execution", zap.String("kind", kind), zap.String("target", target), zap.Bool("skipped", skipped))
}

func (d *DBOpLogger) LogGormSeed(name, dbType string, success bool, count int) {
	lvl := zap.InfoLevel
	if !success {
		lvl = zap.ErrorLevel
	}
	ce := d.l.Check(lvl, "gorm seed")
	if ce != nil {
		ce.Write(zap.String("name", name), zap.String("db_type", dbType), zap.Bool("success", success), zap.Int("count", count))
	}
}

func (d *DBOpLogger) LogFailover(from, to string) {
	d.l.Warn("failover", zap.String("from", from), zap.String("to", to))
}
