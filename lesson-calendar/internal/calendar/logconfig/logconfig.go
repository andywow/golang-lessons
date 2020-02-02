package logconfig

import (
	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// GetLoggerForConfig get logger for config
func GetLoggerForConfig(cfg *config.Config) (*zap.Logger, error) {

	return zap.Config{
		Encoding: "json",

		// set log level
		Level: func(logLevel string) zap.AtomicLevel {
			switch logLevel {
			case "debug":
				return zap.NewAtomicLevelAt(zap.DebugLevel)
			case "error":
				return zap.NewAtomicLevelAt(zap.ErrorLevel)
			case "info":
				return zap.NewAtomicLevelAt(zap.InfoLevel)
			case "warn":
				return zap.NewAtomicLevelAt(zap.WarnLevel)
			default:
				return zap.NewAtomicLevelAt(zap.InfoLevel)
			}
		}(cfg.LogLevel),

		// set log paths
		OutputPaths: func(logFile string, logConsole bool) []string {
			result := []string{}
			if logConsole {
				result = append(result, "stdout")
			}
			if logFile != "" {
				result = append(result, logFile)
			}
			return result
		}(cfg.LogFile, cfg.LogStdout),

		// set error log paths
		ErrorOutputPaths: func(logFile string, logConsole bool) []string {
			result := []string{}
			if logConsole {
				result = append(result, "stderr")
			}
			if logFile != "" {
				result = append(result, logFile)
			}
			return result
		}(cfg.LogFile, cfg.LogStdout),

		// log fields
		EncoderConfig: zapcore.EncoderConfig{
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			TimeKey:      "time",
			LevelKey:     "level",
			CallerKey:    "caller",
			MessageKey:   "message",
		},
	}.Build()

}
