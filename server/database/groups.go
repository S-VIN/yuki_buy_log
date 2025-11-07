package database

import (
	"log"
	"yuki_buy_log/models"
)

func GetAllGroups() (result []models.GroupMember, err error) {
	// Получение всех участников всех групп
	rows, err := db.Query(`
		SELECT g.id, g.user_id, u.login, g.member_number
		FROM groups g
		JOIN users u ON g.user_id = u.id
		ORDER BY g.id, g.member_number`)
	if err != nil {
		log.Printf("Failed to query all groups: %v", err)
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

func GetGroupById(id models.GroupId) (result []models.GroupMember, err error) {
	// Получение всех участников группы по ID группы
	rows, err := db.Query(`
		SELECT g.id, g.user_id, u.login, g.member_number
		FROM groups g
		JOIN users u ON g.user_id = u.id
		WHERE g.id = $1
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

func GetGroupUserCount(id models.GroupId) (count int64, err error) {
	err = db.QueryRow(`SELECT COUNT(*) FROM groups WHERE id = $1`, id).Scan(&count)
	if err != nil {
		log.Printf("Failed to count remaining group members: %v", err)
		return 0, err
	}
	return count, nil
}

func GetGroupIdByUserId(userId models.UserId) (groupId models.GroupId, err error) {
	// Get the group_id for the current user
	err = db.QueryRow(`SELECT id FROM groups WHERE user_id = $1`, userId).Scan(&groupId)
	if err != nil {
		log.Printf("User %d is not in any group", userId)
		// Return empty list if user is not in a group
		return groupId, err
	}
	return groupId, nil
}

func DeleteGroupById(id models.GroupId) error {
	_, err := db.Exec(`DELETE FROM groups WHERE id = $1`, id)
	if err != nil {
		log.Printf("Failed to delete group %d: %v", id, err)
		return err
	}
	return nil
}

func DeleteUserFromGroup(userId models.UserId) (err error) {
	// Remove user from group
	_, err = db.Exec(`DELETE FROM groups WHERE user_id = $1`, userId)
	if err != nil {
		log.Printf("Failed to remove user from group: %v", err)
		return err
	}
	return nil
}

// renumberGroupMembers reassigns member numbers to group members sequentially (1, 2, 3, ...)
// to eliminate gaps after a member leaves
func RenumberGroupMembers(groupId models.GroupId) error {
	// Get all members ordered by their current member_number
	rows, err := db.Query(`
		SELECT user_id
		FROM groups
		WHERE id = $1
		ORDER BY member_number`, groupId)
	if err != nil {
		return err
	}
	defer rows.Close()

	var userIDs []int64
	for rows.Next() {
		var userID int64
		if err := rows.Scan(&userID); err != nil {
			return err
		}
		userIDs = append(userIDs, userID)
	}

	// Update each member with new sequential number
	for i, userID := range userIDs {
		newNumber := i + 1
		_, err = db.Exec(`
			UPDATE groups
			SET member_number = $1
			WHERE id = $2 AND user_id = $3`, newNumber, groupId, userID)
		if err != nil {
			return err
		}
	}

	return nil
}

func AddUserToGroup(groupId models.GroupId, userId models.UserId, memberNumber int) error {
	_, err := db.Exec(`INSERT INTO groups (id, user_id, member_number) VALUES ($1, $2, $3)`, groupId, userId, memberNumber)
	return err
}

func CreateNewGroup(userId models.UserId) (groupId models.GroupId, err error) {
	err = db.QueryRow(`INSERT INTO groups (user_id, member_number) VALUES ($1, $2) RETURNING id`, userId, 1).Scan(&groupId)
	return groupId, err
}
