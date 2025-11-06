package handlers

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"yuki_buy_log/database"
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
		u.Password = string(hash)

		err = database.AddUser(&u)
		if err != nil {
			log.Printf("Failed to register user: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
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

		// Проверка HTTP метода
		if r.Method != http.MethodPost {
			log.Printf("Method not allowed for login: %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Декодирование учетных данных из запроса
		// Структура заполняется не полностью, а только login и hash
		var credentials models.User
		if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			log.Printf("Failed to decode login JSON: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Login attempt for user: %s", credentials.Login)

		// Получение пользователя из базы данных
		user, err := database.GetUserByLogin(credentials.Login)
		if err != nil {
			log.Printf("Failed to get user %s: %v", credentials.Login, err)
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		// Проверка пароля
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
			log.Printf("Invalid password for user: %s", credentials.Login)
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		// Генерация токена аутентификации
		token, err := deps.Auth.GenerateToken(user.Id)
		if err != nil {
			log.Printf("Failed to generate token for user %d: %v", user.Id, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Успешный ответ
		log.Printf("Successfully logged in user %s with ID: %d", credentials.Login, user.Id)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}
