package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/raulito1500/merkadapp/config"
	"google.golang.org/api/option"
)

type serviceAccountKey struct {
	Type        string `json:"type"`
	ProjectID   string `json:"project_id"`
	PrivateKey  string `json:"private_key"`
	ClientEmail string `json:"client_email"`
	TokenURI    string `json:"token_uri"`
}

func NewFirebaseAuthClient(cfg *config.Config) (*auth.Client, error) {
	if cfg.FirebaseProjectID == "" || cfg.FirebaseClientEmail == "" || cfg.FirebasePrivateKey == "" {
		return nil, errors.New("FIREBASE_PROJECT_ID, FIREBASE_CLIENT_EMAIL and FIREBASE_PRIVATE_KEY must be set")
	}

	privateKey := strings.ReplaceAll(cfg.FirebasePrivateKey, "\\n", "\n")

	credentialsJSON, err := json.Marshal(serviceAccountKey{
		Type:        "service_account",
		ProjectID:   cfg.FirebaseProjectID,
		PrivateKey:  privateKey,
		ClientEmail: cfg.FirebaseClientEmail,
		TokenURI:    "https://oauth2.googleapis.com/token",
	})
	if err != nil {
		return nil, err
	}

	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON(credentialsJSON))
	if err != nil {
		return nil, err
	}

	return app.Auth(context.Background())
}

func RequireFirebaseAuth(authClient *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Missing bearer token"})
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")

		decoded, err := authClient.VerifyIDToken(c.Request.Context(), token)
		if err != nil {
			log.Printf("firebase auth: token verification failed: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			return
		}

		c.Set("uid", decoded.UID)
		c.Set("email", decoded.Claims["email"])
		c.Next()
	}
}
