package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/firestore"
	"github.com/arrase21/mobileapi/internal/repository"
	"github.com/arrase21/mobileapi/internal/service"
	"github.com/arrase21/mobileapi/internal/transport/http"
)

func main() {
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "serviceAccountKey.json")
	}

	projectID := os.Getenv("GCP_PROJECT")
	if projectID == "" {
		log.Fatal("GCP_PROJECT environment variable is required")
	}

	client, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatalf("failed to create firestore client: %v", err)
	}
	defer client.Close()

	userRepo := repository.NewFireUserRepo(client)
	userSvc := service.NewUserService(userRepo)
	router := http.NewRouter(userSvc)

	go func() {
		addr := ":8000"
		log.Printf("🚀 Server running on %s", addr)
		if err := router.Run(addr); err != nil {
			log.Fatalf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("✨ Server stopped gracefully")
}
