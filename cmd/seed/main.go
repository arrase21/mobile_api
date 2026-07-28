package main

import (
	"context"
	"log/slog"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/arrase21/mobileapi/internal/logger"
)

func main() {
	logger.Init("info", "")
	ctx := context.Background()

	projectID := "demo-no-project"

	os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", "localhost:9099")
	os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8080")

	config := &firebase.Config{
		ProjectID: projectID,
	}

	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		slog.Error("error initializing app", "error", err)
		os.Exit(1)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		slog.Error("error getting auth client", "error", err)
		os.Exit(1)
	}

	firestoreClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		slog.Error("error creating firestore client", "error", err)
		os.Exit(1)
	}
	defer firestoreClient.Close()

	tenantID := "tenant-1"
	email := "admin@test.com"
	password := "test1234"

	slog.Info("creating tenant in Firestore", "step", 1)
	_, err = firestoreClient.Collection("tenants").Doc(tenantID).Set(ctx, map[string]interface{}{
		"id":         tenantID,
		"name":       "Mi Empresa",
		"plan":       "premium",
		"is_active":  true,
		"max_users":  100,
		"created_at": firestore.ServerTimestamp,
		"updated_at": firestore.ServerTimestamp,
	}, firestore.MergeAll)
	if err != nil {
		slog.Error("error creating tenant", "error", err)
		os.Exit(1)
	}
	slog.Info("tenant created", "tenant_id", tenantID)

	slog.Info("creating admin role", "step", 2)
	permissions := []string{"user.create", "user.read", "user.update", "user.delete", "user.restore"}
	roleID := "role-admin"
	_, err = firestoreClient.Collection("tenants").Doc(tenantID).Collection("roles").Doc(roleID).Set(ctx, map[string]interface{}{
		"id":          roleID,
		"name":        "Admin",
		"description": "Administrador del tenant",
		"permissions": permissions,
		"tenant_id":   tenantID,
		"is_deleted":  false,
		"created_at":  firestore.ServerTimestamp,
		"updated_at":  firestore.ServerTimestamp,
	}, firestore.MergeAll)
	if err != nil {
		slog.Error("error creating role", "error", err)
		os.Exit(1)
	}
	slog.Info("role created", "role_id", roleID, "permissions", permissions)

	slog.Info("creating Firebase Auth user", "step", 3)
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password).
		EmailVerified(false)

	fbUser, err := authClient.GetUserByEmail(ctx, email)
	if err != nil {
		fbUser, err = authClient.CreateUser(ctx, params)
		if err != nil {
			slog.Error("error creating firebase user", "error", err)
			os.Exit(1)
		}
		slog.Info("firebase user created", "email", email, "uid", fbUser.UID)
	} else {
		_, err = authClient.UpdateUser(ctx, fbUser.UID, (&auth.UserToUpdate{}).Password(password))
		if err != nil {
			slog.Error("error updating firebase user", "error", err)
			os.Exit(1)
		}
		slog.Info("firebase user already exists, password updated", "email", email, "uid", fbUser.UID)
	}

	slog.Info("setting custom claims", "step", 4)
	err = authClient.SetCustomUserClaims(ctx, fbUser.UID, map[string]interface{}{
		"tenant_id":      tenantID,
		"is_super_admin": true,
	})
	if err != nil {
		slog.Error("error setting custom claims", "error", err)
		os.Exit(1)
	}
	slog.Info("custom claims set", "claims", []string{"tenant_id", "is_super_admin"})

	slog.Info("creating user in Firestore", "step", 5)
	_, err = firestoreClient.Collection("tenants").Doc(tenantID).Collection("users").Doc(fbUser.UID).Set(ctx, map[string]interface{}{
		"id":             fbUser.UID,
		"firebase_uid":   fbUser.UID,
		"first_name":     "Admin",
		"last_name":      "Test",
		"dni":            "00000000",
		"gender":         "male",
		"phone":          "1234567890",
		"email":          email,
		"nickname":       "admin",
		"is_super_admin": true,
		"role_id":        roleID,
		"tenant_id":      tenantID,
		"created_by":     "",
		"is_deleted":     false,
		"created_at":     firestore.ServerTimestamp,
		"updated_at":     firestore.ServerTimestamp,
	}, firestore.MergeAll)
	if err != nil {
		slog.Error("error creating user in firestore", "error", err)
		os.Exit(1)
	}
	slog.Info("user created in Firestore", "tenant_id", tenantID, "user_id", fbUser.UID)

	slog.Info("seed completed",
		"token_url", "http://localhost:9099/identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=dummy",
		"email", email,
	)
}
