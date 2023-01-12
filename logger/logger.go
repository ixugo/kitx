package logger

import (
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New constructs a Sugared Logger that writes to stdout and
// provides human-readable timestamps.
// func New(service string) (*zap.SugaredLogger, error) {
// 	config := zap.NewProductionConfig()
// 	config.OutputPaths = []string{"stdout"}
// 	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
// 	config.DisableStacktrace = true
// 	config.InitialFields = map[string]any{
// 		"service": service,
// 	}

// 	log, err := config.Build()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return log.Sugar(), nil
// }

// InitJSONLogger ...
func InitJSONLogger(logPath string, debug bool) (*zap.SugaredLogger, error) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	config.NameKey = ""
	l, err := rotatelog(logPath)
	if err != nil {
		return nil, err
	}
	mulitWriteSyncer := []zapcore.WriteSyncer{
		zapcore.AddSync(l),
	}
	if debug {
		mulitWriteSyncer = append(mulitWriteSyncer, zapcore.AddSync(os.Stdout))
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		zapcore.NewMultiWriteSyncer(mulitWriteSyncer...),
		zap.InfoLevel,
	)
	return zap.New(core, zap.AddCaller()).Sugar(), nil
}

// newCustomEncoderConfig 文本型日志
func newTextConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// NewTextLogger 文本日志
func NewTextLogger(logPath string, debug bool) (*zap.SugaredLogger, error) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	config.NameKey = ""
	l, err := rotatelog(logPath)
	if err != nil {
		return nil, err
	}
	mulitWriteSyncer := []zapcore.WriteSyncer{
		zapcore.AddSync(l),
	}
	if debug {
		mulitWriteSyncer = append(mulitWriteSyncer, zapcore.AddSync(os.Stdout))
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(newTextConfig()),
		zapcore.NewMultiWriteSyncer(mulitWriteSyncer...),
		zap.InfoLevel,
	)
	return zap.New(core, zap.AddCaller()).Sugar(), nil
}

func rotatelog(logPath string) (*rotatelogs.RotateLogs, error) {
	return rotatelogs.New(
		filepath.Join(logPath, "%Y%m%d_%H_%M_%S.log"),
		rotatelogs.WithMaxAge(2*7*24*time.Hour),
		rotatelogs.WithRotationTime(12*time.Hour),
	)
}
