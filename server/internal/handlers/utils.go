package handlers

import (
	"fmt"
	"net/http"
	"yuki_buy_log/internal/domain"
	"yuki_buy_log/internal/stores"
)

type Authenticator interface {
	GenerateToken(userId domain.UserId) (string, error)
}

// Вспомогательная функция для получения пользователя из контекста
func getUser(r *http.Request) (user *domain.User, err error) {
	// Используем простую строку как ключ контекста
	userId := r.Context().Value("userId")
	if userId == nil {
		return nil, fmt.Errorf("Cannot get userId from context")
	}

	// Приводим к типу domain.UserId
	userIdTyped, ok := userId.(domain.UserId)
	if !ok {
		return nil, fmt.Errorf("Invalid userId type in context")
	}

	var userStore = stores.GetUserStore()
	user = userStore.GetUserById(userIdTyped)
	if user == nil {
		return nil, fmt.Errorf("User not found")
	}
	return user, nil
}
