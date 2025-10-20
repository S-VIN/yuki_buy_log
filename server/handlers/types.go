package handlers

import (
	"database/sql"
	"net/http"
	"yuki_buy_log/validators"
)

// Тип Validator из пакета validators
type Validator = validators.Validator

type Authenticator interface {
	GenerateToken(userID int64) (string, error)
}

// Структура с зависимостями для handlers
type Dependencies struct {
	DB        *sql.DB
	Validator Validator
	Auth      Authenticator
}

// contextKey is used to store values in request context.
type contextKey string

// UserIDKey is the context key for the authenticated user's id.
const UserIDKey contextKey = "userID"

// Вспомогательная функция для получения ID пользователя из контекста
func userID(r *http.Request) (int64, bool) {
	// Эта функция должна извлекать user ID из контекста запроса
	// Предполагаем, что middleware аутентификации добавляет это в контекст
	id, ok := r.Context().Value(UserIDKey).(int64)
	return id, ok
}