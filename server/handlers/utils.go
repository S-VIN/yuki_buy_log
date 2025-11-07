package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"yuki_buy_log/models"
	"yuki_buy_log/stores"
	"yuki_buy_log/validators"
)

// Тип Validator из пакета validators
type Validator = validators.Validator

type Authenticator interface {
	GenerateToken(userId models.UserId) (string, error)
}

// Структура с зависимостями для handlers
type Dependencies struct {
	DB        *sql.DB
	Validator Validator
	Auth      Authenticator
}

// Вспомогательная функция для получения ID пользователя из контекста
func getUser(r *http.Request) (user *models.User, err error) {
	// Используем простую строку как ключ контекста
	userId := r.Context().Value("userId")
	if userId == nil {
		return nil, fmt.Errorf("Cannot get userId from context")
	}

	// Приводим к типу models.UserId
	userIdTyped, ok := userId.(models.UserId)
	if !ok {
		return nil, fmt.Errorf("Invalid userId type in context")
	}

	userStore := stores.GetUserStore()
	user = userStore.GetUserById(userIdTyped)
	if user == nil {
		return nil, fmt.Errorf("User not found")
	}
	return user, nil
}
