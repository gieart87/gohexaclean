package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gieart87/gohexaclean/internal/infra/asynq"
	"github.com/gieart87/gohexaclean/internal/infra/asynq/tasks"
	asynqlib "github.com/hibiken/asynq"
)

func main() {
	// Get Redis address from environment or use default
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	// Get concurrency from environment or use default
	concurrency := 10

	// Create Asynq server
	srv := asynq.NewServer(redisAddr, concurrency)

	// Create task mux (router)
	mux := asynqlib.NewServeMux()

	// Register task handlers
	mux.HandleFunc(tasks.TypeEmailWelcome, tasks.HandleEmailWelcomeTask)

	// Setup graceful shutdown
	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatalf("Could not run asynq server: %v", err)
		}
	}()

	log.Printf("Asynq worker started (Redis: %s, Concurrency: %d)", redisAddr, concurrency)

	// Wait for interrupt signal to gracefully shutdown the worker
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")
	srv.Shutdown()
	log.Println("Worker stopped")
}
