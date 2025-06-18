package logger

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func NewLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.InitialFields = map[string]interface{}{
		"correlation_id": uuid.New().String(),
	}
	return config.Build()
}