package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gieart87/gohexaclean/internal/bootstrap"
	pb "github.com/gieart87/gohexaclean/api/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load configuration
	configPath := getConfigPath()

	// Initialize container
	container, err := bootstrap.NewContainer(configPath)
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer container.Close()

	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(1024 * 1024 * 10), // 10MB
		grpc.MaxSendMsgSize(1024 * 1024 * 10), // 10MB
	)

	// Register services
	pb.RegisterUserServiceServer(grpcServer, container.UserGRPCHandler)

	// Register reflection service for gRPC tools (e.g., grpcurl)
	reflection.Register(grpcServer)

	// Start server
	port := fmt.Sprintf(":%d", container.Config.Server.GRPC.Port)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		container.Logger.Fatal(fmt.Sprintf("Failed to listen: %v", err))
	}

	container.Logger.Info(fmt.Sprintf("gRPC Server starting on port %d", container.Config.Server.GRPC.Port))

	// Graceful shutdown
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			container.Logger.Fatal(fmt.Sprintf("Failed to serve: %v", err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	container.Logger.Info("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	container.Logger.Info("gRPC Server exited")
}

// getConfigPath returns the configuration file path
func getConfigPath() string {
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}
	return "config/app.yaml"
}
