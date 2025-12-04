package asynq

import "errors"

// Asynq task queue infrastructure errors
var (
	ErrTaskEnqueue   = errors.New("failed to enqueue task")
	ErrTaskProcess   = errors.New("failed to process task")
	ErrTaskTimeout   = errors.New("task processing timeout")
	ErrTaskRetry     = errors.New("task retry limit exceeded")
	ErrTaskDuplicate = errors.New("duplicate task detected")
	ErrWorkerStart   = errors.New("failed to start worker")
	ErrWorkerStop    = errors.New("failed to stop worker")
)
