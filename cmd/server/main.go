package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/arrase21/mobileapi/internal/repository"
	"github.com/arrase21/mobileapi/internal/service"
	transporthttp "github.com/arrase21/mobileapi/internal/transport/http"
)

func main() {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") != "" {
		log.Println("🔥 Running with Firestore emulator")
	} else {
		if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
			os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "serviceAccountKey.json")
		}
	}

	projectID := os.Getenv("GCP_PROJECT")
	if projectID == "" {
		if os.Getenv("FIRESTORE_EMULATOR_HOST") != "" {
			projectID = "demo-no-project"
		} else {
			log.Fatal("GCP_PROJECT environment variable is required")
		}
	}

	client, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatalf("failed to create firestore client: %v", err)
	}
	defer client.Close()

	userRepo := repository.NewFireUserRepo(client)
	roleRepo := repository.NewFireRoleRepo(client)
	assessmentRepo := repository.NewFireAssessmentRepo(client)
	skinfoldRepo := repository.NewFireSkinfoldRepo(client)
	tenantRepo := repository.NewFireTenantRepo(client)
	permissionRepo := repository.NewFirePermissionRepo(client)
	permissionActionRepo := repository.NewFirePermissionActionRepo(client)

	userSvc := service.NewUserService(userRepo, tenantRepo)
	roleSvc := service.NewRoleService(roleRepo)
	assessmentSvc := service.NewAssessmentService(assessmentRepo, userRepo)
	skinfoldSvc := service.NewSkinfoldService(skinfoldRepo, assessmentRepo)
	tenantSvc := service.NewTenantService(tenantRepo)
	permissionSvc := service.NewPermissionService(permissionRepo)
	permissionActionSvc := service.NewPermissionActionService(permissionActionRepo)

	router := transporthttp.NewRouter(
		userSvc,
		roleSvc,
		assessmentSvc,
		skinfoldSvc,
		tenantSvc,
		permissionSvc,
		permissionActionSvc,
		userRepo,
		roleRepo,
		projectID,
	)

	addr := ":8000"
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("🚀 Server running on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("⏳ Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("✨ Server stopped gracefully")
}
