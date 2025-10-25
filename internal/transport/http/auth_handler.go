package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

type FirebaseAuth struct {
	authClient *auth.Client
}

func NewFirebaseAuth() *FirebaseAuth {
	credPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credPath == "" {
		credPath = "service.json"
	}
	opt := option.WithCredentialsFile(credPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}
	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error gettung auth client: %v", err)
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ivalid token"})
			c.Abort()
			return
		}
		c.Set("firebase_uid", token.UID)
		c.Set("claims", token.Claims)
		if tenant, ok := token.Claims["tenant_id"]; ok {
			c.Set("tenant_id", tenant)
		}
		c.Next()
	}
}
