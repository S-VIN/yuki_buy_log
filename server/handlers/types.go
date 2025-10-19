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

// Вспомогательная функция для получения ID пользователя из контекста
func userID(r *http.Request) (int64, bool) {
	// Эта функция должна извлекать user ID из контекста запроса
	// Предполагаем, что middleware аутентификации добавляет это в контекст
	if userID := r.Context().Value("user_id"); userID != nil {
		if id, ok := userID.(int64); ok {
			return id, true
		}
	}
	return 0, false
}