package lib

import (
	"go-clean-arch/src/constants"
	"io"
	"time"

	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gormlogger "gorm.io/gorm/logger"
)

var (
	globalLogger *Logger
	zapLogger    *zap.Logger
)

// Logger structure
type Logger struct {
	*zap.SugaredLogger
}

func newSugaredLogger(logger *zap.Logger) *Logger {
	return &Logger{
		SugaredLogger: logger.Sugar(),
	}
}

// newLogger sets up logger
func newLogger(env *Env) Logger {

	config := zap.NewDevelopmentConfig()
	logOutput := env.LogOutput

	if (env.Environment == constants.Local) || (env.Environment == constants.Development) {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else if env.Environment == constants.Production && logOutput != "" {
		config.OutputPaths = []string{logOutput}
	}

	logLevel := env.LogLevel
	level := zap.PanicLevel
	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "fatal":
		level = zapcore.FatalLevel
	default:
		level = zap.PanicLevel
	}
	config.Level.SetLevel(level)

	zapLogger, _ = config.Build()
	logger := newSugaredLogger(zapLogger)

	return *logger
}

// GetLogger gets the global instance of the logger
func GetLogger() *Logger {
	env := GetEnv()
	if globalLogger == nil {
		logger := newLogger(env)
		globalLogger = &logger
	}
	return globalLogger
}

// GetGormLogger build gorm logger from zap logger (sub-logger)
func (l Logger) GetGormLogger() gormlogger.Interface {
	logger := zapLogger.WithOptions(
		zap.AddCaller(),
		zap.AddCallerSkip(3),
	)

	ignoreRecordNotFoundError := false
	colorful := true
	if globalEnv != nil {
		if globalEnv.Environment == constants.Production {
			ignoreRecordNotFoundError = true
			colorful = false
		}
	}

	result := &GormLogger{
		Logger: newSugaredLogger(logger),
		Config: gormlogger.Config{
			SlowThreshold:             250 * time.Millisecond,
			LogLevel:                  gormlogger.Warn,
			IgnoreRecordNotFoundError: ignoreRecordNotFoundError,
			Colorful:                  colorful,
		},
	}
	result.setup()
	return result
}

// GetFxLogger gets logger for go-fx
func (l Logger) GetFxLogger() fxevent.Logger {
	logger := zapLogger.WithOptions(
		zap.WithCaller(false),
	)
	result := &FxLogger{Logger: newSugaredLogger(logger)}
	return result
}

// GetGinLogger gets logger for gin framework debugging
func (l Logger) GetGinLogger() io.Writer {
	logger := zapLogger.WithOptions(
		zap.WithCaller(false),
	)
	result := &GinLogger{
		Logger: newSugaredLogger(logger),
	}
	return result
}
