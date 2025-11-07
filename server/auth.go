package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"yuki_buy_log/models"

	"github.com/golang-jwt/jwt/v5"
)

// Authenticator handles token generation and verification.
type Authenticator struct {
	secret []byte
}

func NewAuthenticator(secret []byte) *Authenticator {
	return &Authenticator{secret: secret}
}

// GenerateToken creates a signed JWT for the given user id.
func (a *Authenticator) GenerateToken(userId models.UserId) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(int64(userId), 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.secret)
}

// Middleware verifies the Authorization header and adds user id to the context.
func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Auth middleware processing request: %s %s", r.Method, r.URL.Path)
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			log.Printf("Missing or invalid Authorization header for %s", r.URL.Path)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return a.secret, nil
		})
		if err != nil || !token.Valid {
			log.Printf("Invalid token for %s: %v", r.URL.Path, err)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Printf("Invalid token claims for %s", r.URL.Path)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		sub, ok := claims["sub"].(string)
		if !ok {
			log.Printf("Missing subject in token claims for %s", r.URL.Path)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		id, err := strconv.ParseInt(sub, 10, 64)
		if err != nil {
			log.Printf("Invalid user ID in token for %s: %v", r.URL.Path, err)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		log.Printf("Successfully authenticated user %d for %s", id, r.URL.Path)
		ctx := context.WithValue(r.Context(), "userId", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
