package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/arrase21/mobileapi/internal/domain"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

type FirebaseAuth struct {
	authClient *auth.Client
}

func NewFirebaseAuth(projectID string) *FirebaseAuth {
	var opts []option.ClientOption

	if projectID == "" {
		projectID = os.Getenv("GCP_PROJECT")
	}
	if projectID == "" {
		projectID = os.Getenv("GCLOUD_PROJECT")
	}

	emulatorHost := os.Getenv("FIREBASE_AUTH_EMULATOR_HOST")
	if emulatorHost == "" && os.Getenv("FIRESTORE_EMULATOR_HOST") != "" {
		emulatorHost = "localhost:9099"
		os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", emulatorHost)
	}

	if emulatorHost != "" {
		log.Printf("🔥 Firebase Auth: using emulator at %s (project: %s)", emulatorHost, projectID)
	} else {
		credPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		if credPath == "" {
			credPath = "serviceAccountKey.json"
		}
		opts = append(opts, option.WithCredentialsFile(credPath))
	}

	config := &firebase.Config{
		ProjectID: projectID,
	}

	app, err := firebase.NewApp(context.Background(), config, opts...)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v", err)
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting firebase auth client: %v", err)
	}

	return &FirebaseAuth{authClient: authClient}
}

func (m *FirebaseAuth) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := m.authClient.VerifyIDToken(c.Request.Context(), idToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("firebase_uid", token.UID)

		if tenant, ok := token.Claims["tenant_id"]; ok {
			c.Set("tenant_id", tenant)
		}

		c.Next()
	}
}

func RequirePermission(userRepo domain.UserRepo, roleRepo domain.RoleRepo, permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		firebaseUID, exists := c.Get("firebase_uid")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing firebase uid"})
			c.Abort()
			return
		}

		var tenantID string
		if tid, ok := c.Get("tenant_id"); ok {
			tenantID = tid.(string)
		}

		headerTenantID := c.GetHeader("X-Tenant-ID")

		if tenantID == "" && headerTenantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant id"})
			c.Abort()
			return
		}

		if tenantID == "" && headerTenantID != "" {
			tenantID = headerTenantID
		}

		var user *domain.User
		var searchErr error

		if headerTenantID != "" && headerTenantID != tenantID {
			user, searchErr = userRepo.GetByFirebaseUID(c.Request.Context(), headerTenantID, firebaseUID.(string))
			if searchErr != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "user not found in this tenant"})
				c.Abort()
				return
			}
			if !user.IsSuperAdmin {
				c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions: super admin required for cross-tenant access"})
				c.Abort()
				return
			}
			tenantID = headerTenantID
		} else {
			user, searchErr = userRepo.GetByFirebaseUID(c.Request.Context(), tenantID, firebaseUID.(string))
			if searchErr != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "user not found in this tenant"})
				c.Abort()
				return
			}
		}

		c.Set("tenant_id", tenantID)
		c.Set("user_id", user.ID)

		if user.IsSuperAdmin {
			c.Next()
			return
		}

		role, err := roleRepo.GetByID(c.Request.Context(), tenantID, user.RoleID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "user role not found"})
			c.Abort()
			return
		}

		for _, required := range permissions {
			if slices.Contains(role.Permissions, required) {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		c.Abort()
	}
}
