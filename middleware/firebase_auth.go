package middleware

import (
	"context"
	"encoding/json"
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

// NewFirebaseAuthClient builds a Firebase Auth client from the project id,
// client email and private key of a service account.
func NewFirebaseAuthClient(cfg *config.Config) (*auth.Client, error) {
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

// RequireFirebaseAuth verifies the Firebase ID token sent as a
// `Authorization: Bearer <token>` header.
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			return
		}

		c.Set("uid", decoded.UID)
		c.Set("email", decoded.Claims["email"])
		c.Next()
	}
}
