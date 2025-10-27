package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"yuki_buy_log/models"
)

func GroupHandler(deps *Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Group handler called: %s %s", r.Method, r.URL.Path)
		switch r.Method {
		case http.MethodGet:
			getGroupMembers(deps, w, r)
		case http.MethodDelete:
			leaveGroup(deps, w, r)
		default:
			log.Printf("Method not allowed for group: %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func getGroupMembers(deps *Dependencies, w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching group members from database")
	uid, ok := userID(r)
	if !ok {
		log.Println("Unauthorized access attempt to group")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("Fetching group members for user ID: %d", uid)

	// Get the group_id for the current user
	var groupID int64
	err := deps.DB.QueryRow(`SELECT id FROM groups WHERE user_id = $1`, uid).Scan(&groupID)
	if err != nil {
		log.Printf("User %d is not in any group", uid)
		// Return empty list if user is not in a group
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"members": []models.GroupMember{}})
		return
	}

	// Get all members of the same group
	rows, err := deps.DB.Query(`
		SELECT g.id, g.user_id, u.login, g.member_number
		FROM groups g
		JOIN users u ON g.user_id = u.id
		WHERE g.id = $1
		ORDER BY g.member_number`, groupID)
	if err != nil {
		log.Printf("Failed to query group members: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	members := []models.GroupMember{}
	for rows.Next() {
		var m models.GroupMember
		if err := rows.Scan(&m.GroupId, &m.UserId, &m.Login, &m.MemberNumber); err != nil {
			log.Printf("Failed to scan group member row: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		members = append(members, m)
	}
	log.Printf("Successfully fetched %d group members for user %d", len(members), uid)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"members": members})
}

func leaveGroup(deps *Dependencies, w http.ResponseWriter, r *http.Request) {
	log.Println("User leaving group")
	uid, ok := userID(r)
	if !ok {
		log.Println("Unauthorized access attempt to leave group")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the group_id for the current user
	var groupID int64
	err := deps.DB.QueryRow(`SELECT id FROM groups WHERE user_id = $1`, uid).Scan(&groupID)
	if err != nil {
		log.Printf("User %d is not in any group", uid)
		http.Error(w, "you are not in a group", http.StatusBadRequest)
		return
	}

	// Start transaction
	tx, err := deps.DB.Begin()
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Remove user from group
	_, err = tx.Exec(`DELETE FROM groups WHERE user_id = $1`, uid)
	if err != nil {
		log.Printf("Failed to remove user from group: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check remaining members count
	var remainingCount int
	err = tx.QueryRow(`SELECT COUNT(*) FROM groups WHERE id = $1`, groupID).Scan(&remainingCount)
	if err != nil {
		log.Printf("Failed to count remaining group members: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If only 1 member remains, delete the entire group
	if remainingCount == 1 {
		log.Printf("Only 1 member remains in group %d, deleting group", groupID)
		_, err = tx.Exec(`DELETE FROM groups WHERE id = $1`, groupID)
		if err != nil {
			log.Printf("Failed to delete group: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if remainingCount > 1 {
		// Renumber remaining members to fill gaps
		err = renumberGroupMembers(tx, groupID)
		if err != nil {
			log.Printf("Failed to renumber group members: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("User %d successfully left group %d", uid, groupID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "left group successfully"})
}

// renumberGroupMembers reassigns member numbers to group members sequentially (1, 2, 3, ...)
// to eliminate gaps after a member leaves
func renumberGroupMembers(tx *sql.Tx, groupID int64) error {
	// Get all members ordered by their current member_number
	rows, err := tx.Query(`
		SELECT user_id
		FROM groups
		WHERE id = $1
		ORDER BY member_number`, groupID)
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
		_, err = tx.Exec(`
			UPDATE groups
			SET member_number = $1
			WHERE id = $2 AND user_id = $3`, newNumber, groupID, userID)
		if err != nil {
			return err
		}
	}

	return nil
}
