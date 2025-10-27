package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"yuki_buy_log/models"
	"yuki_buy_log/validators"
)

// Вспомогательная функция для получения ID пользователя из контекста
func getUser(d *Dependencies, r *http.Request) (user *models.User, err error) {
	// Используем простую строку как ключ контекста
	userId := r.Context().Value("userId")
	if userId == nil {
		return nil, fmt.Errorf("Cannot get userId from context")
	}

	user = &models.User{}
	err = d.DB.QueryRow(`SELECT id, login, password_hash FROM users WHERE id = $1`, userId).Scan(&user.Id, &user.Login, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("Cant find user with id: %d in db, err: %e", userId, err)
	}
	return user, nil
}

// получаем всех пользователей группы, в которой находится текущий пользователь
func getUsersInGroupForUser(d *Dependencies, userId *models.UserId) ([]models.User, error) {
	var groupId int64
	var result []models.User
	err := d.DB.QueryRow(`SELECT id FROM groups WHERE user_id = $1`, userId).Scan(&groupId)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("Cant find users in group id %d, err: %e", userId, err)
		}
		// Если пользователь не в группе, то возвращаем только текущего пользователя
		result = append(result, models.User{})
		return result, nil
	}

	rows, err := d.DB.Query(`SELECT user_id, login, password_hash FROM groups JOIN users u on groups.user_id = u.id WHERE groups.id = $1`, groupId)
	if err != nil {
		return nil, fmt.Errorf("Cant find user with id: %d in db, err: %e", userId, err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Login, &user.Password)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		result = append(result, user)
	}
	return result, nil
}

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
