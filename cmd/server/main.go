package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/arrase21/mobileapi/internal/logger"
	"github.com/arrase21/mobileapi/internal/repository"
	"github.com/arrase21/mobileapi/internal/service"
	transporthttp "github.com/arrase21/mobileapi/internal/transport/http"
)

func main() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	logFile := os.Getenv("LOG_FILE")
	if _, ok := os.LookupEnv("LOG_FILE"); !ok {
		logFile = "logs/app.log"
	}
	logger.Init(logLevel, logFile)

	if os.Getenv("FIRESTORE_EMULATOR_HOST") != "" {
		slog.Info("running with Firestore emulator")
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
			slog.Error("GCP_PROJECT environment variable is required")
			os.Exit(1)
		}
	}

	client, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		slog.Error("failed to create firestore client", "error", err)
		os.Exit(1)
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	addr := ":" + port
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("server starting", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped gracefully")
}
