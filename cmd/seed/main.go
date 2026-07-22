package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

func main() {
	ctx := context.Background()

	projectID := "demo-no-project"

	os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", "localhost:9099")
	os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8080")

	config := &firebase.Config{
		ProjectID: projectID,
	}

	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting auth client: %v", err)
	}

	firestoreClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("error creating firestore client: %v", err)
	}
	defer firestoreClient.Close()

	tenantID := "tenant-1"
	email := "admin@test.com"
	password := "test1234"

	log.Println("1. Creando tenant en Firestore...")
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
		log.Fatalf("error creating tenant: %v", err)
	}
	log.Printf("   Tenant '%s' creado\n", tenantID)

	log.Println("2. Creando rol admin...")
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
		log.Fatalf("error creating role: %v", err)
	}
	log.Printf("   Rol '%s' creado con permisos: %v\n", roleID, permissions)

	log.Println("3. Creando usuario en Firebase Auth...")
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password).
		EmailVerified(false)

	fbUser, err := authClient.GetUserByEmail(ctx, email)
	if err != nil {
		fbUser, err = authClient.CreateUser(ctx, params)
		if err != nil {
			log.Fatalf("error creating firebase user: %v", err)
		}
		log.Printf("   Usuario Firebase creado: %s (uid: %s)\n", email, fbUser.UID)
	} else {
		// Update password in case it changed
		_, err = authClient.UpdateUser(ctx, fbUser.UID, (&auth.UserToUpdate{}).Password(password))
		if err != nil {
			log.Fatalf("error updating firebase user: %v", err)
		}
		log.Printf("   Usuario Firebase ya existía: %s (uid: %s), password actualizado\n", email, fbUser.UID)
	}

	log.Println("4. Seteando custom claims...")
	err = authClient.SetCustomUserClaims(ctx, fbUser.UID, map[string]interface{}{
		"tenant_id":      tenantID,
		"is_super_admin": true,
	})
	if err != nil {
		log.Fatalf("error setting custom claims: %v", err)
	}
	log.Println("   Custom claims seteados (tenant_id + is_super_admin)")

	log.Println("5. Creando usuario en Firestore...")
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
		log.Fatalf("error creating user in firestore: %v", err)
	}
	log.Printf("   Usuario Firestore creado en %s/users/%s\n", tenantID, fbUser.UID)

	log.Println("")
	log.Println("═══════════════════════════════════════════════════")
	log.Println("  Seed completado!")
	log.Println("")
	log.Println("  Para obtener el token:")
	log.Printf("  curl -X POST 'http://localhost:9099/identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=dummy' \\\n")
	log.Println("    -H 'Content-Type: application/json' \\\n")
	log.Printf("    -d '{\"email\":\"%s\",\"password\":\"%s\",\"returnSecureToken\":true}'", email, password)
	log.Println("")
	log.Println("═══════════════════════════════════════════════════")
}
