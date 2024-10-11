package middleware

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"gotempl/views"
	"gotempl/views/handlers"
	"net/http"
	"os"
	"strings"
	"time"

	clerkjwt "github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

type ClerkPublicAuthMiddleware struct {
	JwtPublicSigningKey string
}

func (c *ClerkPublicAuthMiddleware) Init() error {

	err := godotenv.Load()
	if err != nil {
		return err
	}

	keyPath := os.Getenv("JWT_PUBLIC_KEY_PATH")
	if keyPath != "" {
		key, err := os.ReadFile(keyPath)
		if err == nil {
			jwtSigningKey = key
			return nil
		}
		return err
	}
	return errors.New("unable to load JWT_PUBLIC_KEY_PATH from .env file")
}

// Middleware for authentication
func (c *ClerkPublicAuthMiddleware) ClerkAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("__session")
		if err != nil || cookie == "" {
			handlers.Render(c, http.StatusForbidden, views.Error500("Access denied: authentication is needed"))
			c.Abort()
			return
		}

		// Extract session token from the cookie
		sessionToken := strings.TrimSpace(cookie)

		ret, err := verifyTokenLocally(sessionToken)
		if err != nil {
			handlers.Render(c, http.StatusForbidden, views.Error500(fmt.Sprintf("Access denied: %s", err.Error())))
			c.Abort()
			return
		}
		if ret.Valid {
			c.Next()
			return

		}

		// Verify the session
		claims, err := clerkjwt.Verify(c.Request.Context(), &clerkjwt.VerifyParams{
			Token:  sessionToken,
			Leeway: 10 * time.Second,
		})
		if err != nil {
			handlers.Render(c, http.StatusForbidden, views.Error500(fmt.Sprintf("Access denied: %v", err)))

			c.Abort()
			return
		}

		// Get user information
		usr, err := user.Get(c.Request.Context(), claims.Subject)
		if err != nil {
			handlers.Render(c, http.StatusInternalServerError, views.Error500("Failed to get user information"))
			c.Abort()
			return
		}

		// Check if the user is banned
		if usr.Banned {
			handlers.Render(c, http.StatusForbidden, views.Error500("Access denied: user is banned"))
			c.Abort()
			return
		}

		// Set user info in context for use in other handlers
		c.Set("user", usr)

		c.Next()
	}
}

var jwtSigningKey = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzDqLjNpG813z8Oxdpz0f
9BQ/8P18DvcPA27se8m9gl9jyTh2wEWqhKAV/ydtp47e/bQfhe1A4+NkXfYx6+Np
JobfBLAEJKDy+puhsqOMJkpcExNkzh6Wu/TBdWCQDvLl7lAZQzRP9XNyQ2CkOBAe
jSllN390M2eogQHa/tuF1JjjtJW3G1ywaZwNHGcl/2FI1a0GMMCAPX8yWu55VzYS
iXR61Eti7/+iErKFiSvpwTc9hkmTP2x2+Zl6CVqUC+sEzGgPsbRdcoxwIS+tPB3V
tLf0iPvQGNUalMi32M1jnus2iwGUUeD+IcL/wAcTlQW6iK+OaFW37ct8YlqY92sT
cQIDAQAB
-----END PUBLIC KEY-----
`)

func verifyTokenLocally(tokenString string) (*jwt.Token, error) {

	publicKey, err := parseRSAPublicKey(jwtSigningKey)
	if err != nil {
		fmt.Errorf("Failed to parse public key: %v", err)
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	}

	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		// Check if the error is because the token is expired
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, fmt.Errorf("token is expired")
			}
		}
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	return token, nil
}

// Parse the RSA public key from the PEM-encoded data
func parseRSAPublicKey(pemKey []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemKey)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	// Assert the parsed key is of type *rsa.PublicKey
	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key is not of type *rsa.PublicKey")
	}

	return rsaPubKey, nil
}
