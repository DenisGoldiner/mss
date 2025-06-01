package mss

import "errors"

// The list of supported errors.
var (
	ErrEngineNotFound      = errors.New("engine not found")
	ErrGlobalQueueOverflow = errors.New("global queue is full")
	ErrInvalidTask         = errors.New("the task does not satisfy the system requirements")
)
