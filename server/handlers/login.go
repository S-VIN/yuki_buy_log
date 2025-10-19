package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"yuki_buy_log/models"
)

func RegisterHandler(deps *Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Register handler called: %s %s", r.Method, r.URL.Path)
		if r.Method != http.MethodPost {
			log.Printf("Method not allowed for register: %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var u models.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			log.Printf("Failed to decode user JSON: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Registering new user: %s", u.Login)
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = deps.DB.QueryRow(`INSERT INTO users (login, password_hash) VALUES ($1,$2) RETURNING id`, u.Login, string(hash)).Scan(&u.Id)
		if err != nil {
			log.Printf("Failed to insert user: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		token, err := deps.Auth.GenerateToken(u.Id)
		if err != nil {
			log.Printf("Failed to generate token for user %d: %v", u.Id, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Successfully registered user %s with ID: %d", u.Login, u.Id)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}

func LoginHandler(deps *Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Login handler called: %s %s", r.Method, r.URL.Path)
		if r.Method != http.MethodPost {
			log.Printf("Method not allowed for login: %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var creds models.User
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			log.Printf("Failed to decode login JSON: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Login attempt for user: %s", creds.Login)
		var id int64
		var hash string
		err := deps.DB.QueryRow(`SELECT id, password_hash FROM users WHERE login=$1`, creds.Login).Scan(&id, &hash)
		if err != nil {
			log.Printf("User not found or database error for login %s: %v", creds.Login, err)
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		if bcrypt.CompareHashAndPassword([]byte(hash), []byte(creds.Password)) != nil {
			log.Printf("Invalid password for user: %s", creds.Login)
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		token, err := deps.Auth.GenerateToken(id)
		if err != nil {
			log.Printf("Failed to generate token for user %d: %v", id, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Successfully logged in user %s with ID: %d", creds.Login, id)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}
