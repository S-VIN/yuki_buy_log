package database

import (
	"database/sql"
	"fmt"
	"log"
	"yuki_buy_log/models"
)

func (d *DatabaseManager) GetUserById(id *models.UserId) (user *models.User, err error) {
	user = &models.User{}
	err = d.db.QueryRow(`SELECT id, login, password_hash FROM users WHERE id = $1`, id).Scan(&user.Id, &user.Login, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("Cant find user with id: %d in db, err: %e", id, err)
	}
	return user, nil
}

func (d *DatabaseManager) GetUserByLogin(login string) (user *models.User, err error) {
	err = d.db.QueryRow(`SELECT id, password_hash FROM users WHERE login=$1`, login).Scan(&user.Id, &user.Password)
	if err != nil {
		log.Printf("User not found or database error for login %s: %v", login, err)
		return nil, fmt.Errorf("User not found or database error for login %s", login)
	}
	return user, nil
}

func (d *DatabaseManager) GetUsersByGroupId(id *models.GroupId) (users []models.User, err error) {
	err = d.db.QueryRow(`SELECT id FROM groups WHERE user_id = $1`, id).Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			return users, fmt.Errorf("Cant find users in group id %d, err: %e", id, err)
		}
		return users, nil
	}

	rows, err := d.db.Query(`SELECT user_id, login, password_hash FROM groups JOIN users u on groups.user_id = u.id WHERE groups.id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("Cant find user with id: %d in db, err: %e", id, err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Login, &user.Password)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (d *DatabaseManager) AddUser(user *models.User) (err error) {
	err = d.db.QueryRow(`INSERT INTO users (login, password_hash) VALUES ($1,$2) RETURNING id`, user.Login, user.Password).Scan(&user.Id)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return err
	}
	return nil
}

func (d *DatabaseManager) UpdateUser(user *models.User) error {
	_, err := d.db.Exec(`UPDATE users SET login = $1, password_hash = $2 WHERE id = $3`, user.Login, user.Password, user.Id)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		return err
	}
	return nil
}

func (d *DatabaseManager) DeleteUser(userId models.UserId) error {
	_, err := d.db.Exec(`DELETE FROM users WHERE id = $1`, userId)
	if err != nil {
		log.Printf("Failed to delete user: %v", err)
		return err
	}
	return nil
}

func (d *DatabaseManager) GetAllUsers() ([]models.User, error) {
	rows, err := d.db.Query(`SELECT id, login, password_hash FROM users`)
	if err != nil {
		return nil, fmt.Errorf("Failed to get all users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Login, &user.Password)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}
