package asynq

import (
	"github.com/hibiken/asynq"
)

// NewClient creates a new Asynq client for enqueueing tasks
func NewClient(redisAddr string) *asynq.Client {
	redisOpt := asynq.RedisClientOpt{
		Addr: redisAddr,
	}
	return asynq.NewClient(redisOpt)
}

// NewServer creates a new Asynq server for processing tasks
func NewServer(redisAddr string, concurrency int) *asynq.Server {
	redisOpt := asynq.RedisClientOpt{
		Addr: redisAddr,
	}

	return asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: concurrency,
			Queues: map[string]int{
				"critical": 6, // processed 60% of the time
				"default":  3, // processed 30% of the time
				"low":      1, // processed 10% of the time
			},
		},
	)
}
