package logger

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	initOnce    sync.Once
	base        *zap.Logger
	initialized bool
	mu          sync.RWMutex
)

type Config struct {
	Level   string
	Format  string
	Output  string
	Verbose bool
}

func NoFileConfig() Config {
	return Config{Level: "info", Format: "text", Output: "console", Verbose: false}
}

const (
	colorReset  = "\033[0m"
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorPurple = "\033[35m"
	colorGray   = "\033[90m"
	colorBlue   = "\033[34m"
)

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
	enc.AppendString(fmt.Sprintf("%s[%s]%s\n   ", colorCyan, caller.TrimmedPath(), colorReset))
}

func Initialize(cfg Config) error {
	mu.Lock()
	defer mu.Unlock()

	enableVirtualTerminalProcessing()

	var zcfg zap.Config
	if cfg.Format == "json" {
		zcfg = zap.NewProductionConfig()
	} else {
		zcfg = zap.NewDevelopmentConfig()
		zcfg.Encoding = "console"
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
	if err != nil {
		return err
	}
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

func Infof(format string, args ...interface{})  { GetLogger().Sugar().Infof(format, args...) }
func Warnf(format string, args ...interface{})  { GetLogger().Sugar().Warnf(format, args...) }
func Errorf(format string, args ...interface{}) { GetLogger().Sugar().Errorf(format, args...) }

func SearchHeader(w io.Writer, title, root string) {
	fmt.Fprintf(w, "%s%s%s\n", colorGreen, title, colorReset)
	fmt.Fprintf(w, "%sRoot:%s %s\n\n", colorGray, colorReset, root)
}

func SearchFile(w io.Writer, path string) {
	fmt.Fprintf(w, "%s%s%s\n", colorCyan, path, colorReset)
}

func SearchHit(w io.Writer, lineNo int, frequency int, text string, includeFrequency bool, highlight string) {
	text = Highlight(text, highlight)
	if includeFrequency {
		fmt.Fprintf(w, "  %sL%-6d%s %shits=%-3d%s %s\n", colorYellow, lineNo, colorReset, colorPurple, frequency, colorReset, text)
		return
	}
	fmt.Fprintf(w, "  %sL%-6d%s %s\n", colorYellow, lineNo, colorReset, text)
}

func SearchSummary(w io.Writer, matches int) {
	fmt.Fprintf(w, "\n%sMatches shown:%s %d\n", colorGray, colorReset, matches)
}

func SearchContext(w io.Writer, lineNo int, text string) {
	fmt.Fprintf(w, "  %sL%-6d%s %s%s%s\n", colorGray, lineNo, colorReset, colorGray, text, colorReset)
}

func Prompt(w io.Writer, text string) {
	fmt.Fprintf(w, "%s%s%s", colorBlue, text, colorReset)
}

func Highlight(text, needle string) string {
	if needle == "" {
		return text
	}

	lowerText := strings.ToLower(text)
	lowerNeedle := strings.ToLower(needle)
	if lowerNeedle == "" {
		return text
	}

	var builder strings.Builder
	start := 0
	for {
		idx := strings.Index(lowerText[start:], lowerNeedle)
		if idx < 0 {
			builder.WriteString(text[start:])
			break
		}
		idx += start
		builder.WriteString(text[start:idx])
		builder.WriteString(colorYellow)
		builder.WriteString(text[idx : idx+len(needle)])
		builder.WriteString(colorReset)
		start = idx + len(needle)
	}

	return builder.String()
}

func HighlightRegex(text string, re *regexp.Regexp) string {
	if re == nil {
		return text
	}

	indexes := re.FindAllStringIndex(text, -1)
	if len(indexes) == 0 {
		return text
	}

	var builder strings.Builder
	last := 0
	for _, idx := range indexes {
		start, end := idx[0], idx[1]
		if start < last {
			continue
		}
		builder.WriteString(text[last:start])
		builder.WriteString(colorYellow)
		builder.WriteString(text[start:end])
		builder.WriteString(colorReset)
		last = end
	}
	builder.WriteString(text[last:])

	return builder.String()
}
