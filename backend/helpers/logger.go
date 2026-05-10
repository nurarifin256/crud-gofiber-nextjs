package helpers

import (
	"go.uber.org/zap"
)

// LogError logs error with structured fields for common database operations
func LogError(message string, err error, fields ...zap.Field) {
	allFields := append([]zap.Field{zap.Error(err)}, fields...)
	zap.L().Error(message, allFields...)
}

// LogInfo logs informational messages with structured fields
func LogInfo(message string, fields ...zap.Field) {
	zap.L().Info(message, fields...)
}

// LogWarn logs warning messages with structured fields
func LogWarn(message string, fields ...zap.Field) {
	zap.L().Warn(message, fields...)
}
