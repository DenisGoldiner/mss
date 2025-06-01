package mss

import (
	"log/slog"
)

const (
	logKeyResult = "result"
	logKeyError  = "error"
)

// LogPostprocessor is an implementation of the postprocessor interface.
// It only logs if the execution was successful or not.
type LogPostprocessor[Out any] struct {
	logger *slog.Logger
}

// NewLogPostprocessor is a constructor for the LogPostprocessor.
func NewLogPostprocessor[Out any](logger *slog.Logger) LogPostprocessor[Out] {
	return LogPostprocessor[Out]{logger: logger}
}

func (p LogPostprocessor[Out]) ok(resVal Out) {
	p.logger.Info("Successfully processed task", logKeyResult, resVal)
}

func (p LogPostprocessor[Out]) err(resErr error) {
	p.logger.Info("Failed to process task", logKeyError, resErr)
}
