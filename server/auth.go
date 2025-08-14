package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// contextKey is used to store values in request context.
type contextKey string

// UserLoginKey is the context key for the authenticated user's login.
const UserLoginKey contextKey = "userLogin"

// Authenticator handles token generation and verification.
type Authenticator struct {
	secret []byte
}

func NewAuthenticator(secret []byte) *Authenticator {
	return &Authenticator{secret: secret}
}

// GenerateToken creates a signed JWT for the given user login.
func (a *Authenticator) GenerateToken(login string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   login,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.secret)
}

// Middleware verifies the Authorization header and adds user login to the context.
func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return a.secret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		login, ok := claims["sub"].(string)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserLoginKey, login)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// userLogin retrieves the authenticated user login from the request context.
func userLogin(r *http.Request) (string, bool) {
	login, ok := r.Context().Value(UserLoginKey).(string)
	return login, ok
}
