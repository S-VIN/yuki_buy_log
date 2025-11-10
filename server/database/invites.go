package database

import (
	"database/sql"
	"time"
	"yuki_buy_log/models"
)

func scanRowsToInvites(rows *sql.Rows) ([]models.Invite, error) {
	var invites []models.Invite
	for rows.Next() {
		var inv models.Invite
		if err := rows.Scan(&inv.Id, &inv.FromUserId, &inv.ToUserId, &inv.FromLogin, &inv.ToLogin, &inv.CreatedAt); err != nil {
			return nil, err
		}
		invites = append(invites, inv)
	}
	return invites, nil
}

func GetIncomingInvites(userId models.UserId) ([]models.Invite, error) {
	rows, err := db.Query(`
		SELECT i.id, i.from_user_id, i.to_user_id, u_from.login, u_to.login, i.created_at
		FROM invites i
		JOIN users u_from ON i.from_user_id = u_from.id
		JOIN users u_to ON i.to_user_id = u_to.id
		WHERE i.to_user_id = $1`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowsToInvites(rows)
}

func GetInvite(fromUserId, toUserId models.UserId) (models.Invite, error) {
	var invite models.Invite
	err := db.QueryRow(`
		SELECT i.id, i.from_user_id, i.to_user_id, u_from.login, u_to.login, i.created_at
		FROM invites i
		JOIN users u_from ON i.from_user_id = u_from.id
		JOIN users u_to ON i.to_user_id = u_to.id
		WHERE i.from_user_id = $1 AND i.to_user_id = $2`, fromUserId, toUserId).Scan(&invite.Id, &invite.FromUserId, &invite.ToUserId, &invite.FromLogin, &invite.ToLogin, &invite.CreatedAt)
	return invite, err
}

func GetInvitesFromUser(fromUserId models.InviteId) ([]models.Invite, error) {
	rows, err := db.Query(`
		SELECT i.id, i.from_user_id, i.to_user_id, u_from.login, u_to.login, i.created_at
		FROM invites i
		JOIN users u_from ON i.from_user_id = u_from.id
		JOIN users u_to ON i.to_user_id = u_to.id
		WHERE i.from_user_id = $1`, fromUserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowsToInvites(rows)
}

func GetInvitesToUser(toUserId models.UserId) ([]models.Invite, error) {
	rows, err := db.Query(`
		SELECT i.id, i.from_user_id, i.to_user_id, u_from.login, u_to.login, i.created_at
		FROM invites i
		JOIN users u_from ON i.from_user_id = u_from.id
		JOIN users u_to ON i.to_user_id = u_to.id
		WHERE i.to_user_id = $1`, toUserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowsToInvites(rows)
}

func DeleteInvitesBetweenUsers(userId1, userId2 models.UserId) error {
	_, err := db.Exec(`DELETE FROM invites WHERE (from_user_id = $1 AND to_user_id = $2) OR (from_user_id = $2 AND to_user_id = $1)`, userId1, userId2)
	return err
}

func GetAllInvites() ([]models.Invite, error) {
	rows, err := db.Query(`
		SELECT i.id, i.from_user_id, i.to_user_id, u_from.login, u_to.login, i.created_at
		FROM invites i
		JOIN users u_from ON i.from_user_id = u_from.id
		JOIN users u_to ON i.to_user_id = u_to.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowsToInvites(rows)
}

func CreateInvite(fromUserId, toUserId models.UserId) (models.InviteId, error) {
	var inviteId models.InviteId
	err := db.QueryRow(`INSERT INTO invites (from_user_id, to_user_id) VALUES ($1, $2) RETURNING id`, fromUserId, toUserId).Scan(&inviteId)
	return inviteId, err
}

func DeleteInvite(inviteId models.InviteId) error {
	_, err := db.Exec(`DELETE FROM invites WHERE id = $1`, inviteId)
	return err
}

func DeleteOldInvites(cutoffTime time.Time) (int64, error) {
	result, err := db.Exec(`DELETE FROM invites WHERE created_at < $1`, cutoffTime)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
