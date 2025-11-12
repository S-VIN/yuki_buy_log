package database

import (
	"log"
	"yuki_buy_log/models"
)

func (d *DatabaseManager) GetAllGroupMembers() (result []models.GroupMember, err error) {
	// Получение всех участников всех групп
	rows, err := d.db.Query(`
		SELECT g.group_id, g.user_id, u.login, g.member_number
		FROM group_members g
		JOIN users u ON g.user_id = u.id
		ORDER BY g.group_id, g.member_number`)
	if err != nil {
		log.Printf("Failed to query all group_members: %v", err)
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var member models.GroupMember
		if err := rows.Scan(&member.GroupId, &member.UserId, &member.Login, &member.MemberNumber); err != nil {
			log.Printf("Failed to scan group member row: %v", err)
			return result, err
		}
		result = append(result, member)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v", err)
		return result, err
	}

	log.Printf("Successfully retrieved %d group members", len(result))
	return result, nil
}

func (d *DatabaseManager) GetGroupMembersByGroupId(id models.GroupId) (result []models.GroupMember, err error) {
	// Получение всех участников группы по ID группы
	rows, err := d.db.Query(`
		SELECT g.group_id, g.user_id, u.login, g.member_number
		FROM group_members g
		JOIN users u ON g.user_id = u.id
		WHERE g.group_id = $1
		ORDER BY g.member_number`, id)
	if err != nil {
		log.Printf("Failed to query group members for group %d: %v", id, err)
		return result, err
	}
	defer rows.Close()

	// Обработка всех строк из результата запроса
	for rows.Next() {
		var member models.GroupMember
		if err := rows.Scan(&member.GroupId, &member.UserId, &member.Login, &member.MemberNumber); err != nil {
			log.Printf("Failed to scan group member row: %v", err)
			return result, err
		}
		result = append(result, member)
	}

	// Проверка на ошибки итерации
	if err = rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v", err)
		return result, err
	}

	log.Printf("Successfully retrieved %d members for group %d", len(result), id)
	return result, nil
}

func (d *DatabaseManager) DeleteGroupMembersByGroupId(id models.GroupId) error {
	_, err := d.db.Exec(`DELETE FROM group_members WHERE group_id = $1`, id)
	if err != nil {
		log.Printf("Failed to delete group %d: %v", id, err)
		return err
	}
	return nil
}

func (d *DatabaseManager) DeleteUserFromGroup(userId models.UserId) (err error) {
	// Remove user from group
	_, err = d.db.Exec(`DELETE FROM group_members WHERE user_id = $1`, userId)
	if err != nil {
		log.Printf("Failed to remove user from group: %v", err)
		return err
	}
	return nil
}

func (d *DatabaseManager) AddUserToGroup(groupId models.GroupId, userId models.UserId, memberNumber int) error {
	_, err := d.db.Exec(`INSERT INTO group_members (group_id, user_id, member_number) VALUES ($1, $2, $3)`, groupId, userId, memberNumber)
	return err
}

// Обновляет member number как в groupMember
func (d *DatabaseManager) UpdateGroupMember(groupMember *models.GroupMember) error {
	_, err := d.db.Exec(`UPDATE group_members SET member_number = $1 WHERE user_id = $2`,
		groupMember.MemberNumber, groupMember.UserId)
	if err != nil {
		log.Printf("Failed to update group member for user %d: %v", groupMember.UserId, err)
		return err
	}
	return nil
}

func (d *DatabaseManager) CreateNewGroup(userId models.UserId) (groupId models.GroupId, err error) {
	err = d.db.QueryRow(`INSERT INTO group_members (group_id, user_id, member_number) VALUES (
		(SELECT COALESCE(MAX(group_id), 0) + 1 FROM group_members), $1, $2) RETURNING group_id`).Scan(&groupId, &userId, 0)
	return groupId, err
}
