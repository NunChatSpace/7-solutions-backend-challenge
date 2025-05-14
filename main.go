package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/database"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/database/mongo/repositories"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/grpc"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/http"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/config"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/di"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/domain"
	"github.com/savsgio/atreugo/v11"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	appConfig := atreugo.Config{
		Addr:             cfg.App.Port,
		GracefulShutdown: true,
	}
	deps := di.NewDependency(cfg)

	repositories.ProvideRepositories(deps)
	services.ProvideServices(deps)

	server := http.NewServer(deps, appConfig)
	grpcServer := grpc.NewGRPCServer(deps)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Background user logger
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Shutting down background user logger")
				return
			case <-ticker.C:
				userRepo := di.Get[database.IUserRepository](deps)
				users, err := userRepo.Search(domain.User{})
				if err != nil {
					log.Printf("Error fetching users: %v\n", err)
					continue
				}
				log.Printf("User count: %d\n", len(users))
			}
		}
	}()

	// Run HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("HTTP server error: %+v", err)
		}
		log.Println("HTTP server stopped")
	}()

	// Run gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %+v", err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %+v", err)
		}
		log.Println("gRPC server stopped")
	}()

	// Handle shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	log.Println("Shutdown signal received")
	cancel()

	// Stop HTTP server
	if err := server.Shutdown(); err != nil {
		log.Printf("HTTP Server Shutdown Failed: %+v", err)
	}

	// Stop gRPC server gracefully
	grpcServer.GracefulStop()

	wg.Wait()
	log.Println("All services gracefully stopped")
}
