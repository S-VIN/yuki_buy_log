package database

import (
	"database/sql"
	"time"
	"yuki_buy_log/internal/domain"
)

func scanRowsToInvites(rows *sql.Rows) ([]domain.Invite, error) {
	var invites []domain.Invite
	for rows.Next() {
		var inv domain.Invite
		if err := rows.Scan(&inv.Id, &inv.FromUserId, &inv.ToUserId, &inv.FromLogin, &inv.ToLogin, &inv.CreatedAt); err != nil {
			return nil, err
		}
		invites = append(invites, inv)
	}
	return invites, nil
}

func (d *DatabaseManager) GetIncomingInvites(userId domain.UserId) ([]domain.Invite, error) {
	rows, err := d.db.Query(`
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

func (d *DatabaseManager) GetInvite(fromUserId, toUserId domain.UserId) (domain.Invite, error) {
	var invite domain.Invite
	err := d.db.QueryRow(`
		SELECT i.id, i.from_user_id, i.to_user_id, u_from.login, u_to.login, i.created_at
		FROM invites i
		JOIN users u_from ON i.from_user_id = u_from.id
		JOIN users u_to ON i.to_user_id = u_to.id
		WHERE i.from_user_id = $1 AND i.to_user_id = $2`, fromUserId, toUserId).Scan(&invite.Id, &invite.FromUserId, &invite.ToUserId, &invite.FromLogin, &invite.ToLogin, &invite.CreatedAt)
	return invite, err
}

func (d *DatabaseManager) GetInvitesFromUser(fromUserId domain.InviteId) ([]domain.Invite, error) {
	rows, err := d.db.Query(`
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

func (d *DatabaseManager) GetInvitesToUser(toUserId domain.UserId) ([]domain.Invite, error) {
	rows, err := d.db.Query(`
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

func (d *DatabaseManager) DeleteInvitesBetweenUsers(userId1, userId2 domain.UserId) error {
	_, err := d.db.Exec(`DELETE FROM invites WHERE (from_user_id = $1 AND to_user_id = $2) OR (from_user_id = $2 AND to_user_id = $1)`, userId1, userId2)
	return err
}

func (d *DatabaseManager) GetAllInvites() ([]domain.Invite, error) {
	rows, err := d.db.Query(`
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

func (d *DatabaseManager) CreateInvite(fromUserId, toUserId domain.UserId) (domain.InviteId, error) {
	var inviteId domain.InviteId
	err := d.db.QueryRow(`INSERT INTO invites (from_user_id, to_user_id) VALUES ($1, $2) RETURNING id`, fromUserId, toUserId).Scan(&inviteId)
	return inviteId, err
}

func (d *DatabaseManager) DeleteOldInvites(cutoffTime time.Time) (int64, error) {
	result, err := d.db.Exec(`DELETE FROM invites WHERE created_at < $1`, cutoffTime)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
