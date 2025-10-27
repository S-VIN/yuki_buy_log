package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"yuki_buy_log/models"
)

func InviteHandler(d *Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Invite handler called: %s %s", r.Method, r.URL.Path)
		switch r.Method {
		case http.MethodGet:
			getIncomingInvites(d, w, r)
		case http.MethodPost:
			sendInvite(d, w, r)
		default:
			log.Printf("Method not allowed for invite: %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func getIncomingInvites(d *Dependencies, w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching incoming invites from database")
	user, err := getUser(d, r)
	if err != nil {
		log.Println("Unauthorized access attempt to invites")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("Fetching incoming invites for user ID: %d", user.Id)

	rows, err := d.DB.Query(`
		SELECT i.id, i.from_user_id, i.to_user_id, u_from.login, u_to.login, i.created_at
		FROM invites i
		JOIN users u_from ON i.from_user_id = u_from.id
		JOIN users u_to ON i.to_user_id = u_to.id
		WHERE i.to_user_id = $1`, user.Id)
	if err != nil {
		log.Printf("Failed to query invites: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	invites := []models.Invite{}
	for rows.Next() {
		var inv models.Invite
		if err := rows.Scan(&inv.Id, &inv.FromUserId, &inv.ToUserId, &inv.FromLogin, &inv.ToLogin, &inv.CreatedAt); err != nil {
			log.Printf("Failed to scan invite row: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		invites = append(invites, inv)
	}
	log.Printf("Successfully fetched %d incoming invites for user %d", len(invites), user.Id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"invites": invites})
}

func sendInvite(d *Dependencies, w http.ResponseWriter, r *http.Request) {
	log.Println("Sending invite")
	var req struct {
		Login string `json:"login"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode invite JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := getUser(d, r)
	if err != nil {
		log.Println("Unauthorized access attempt to send invite")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get target user ID by login
	var targetUserID int64
	err = d.DB.QueryRow(`SELECT id FROM users WHERE login = $1`, req.Login).Scan(&targetUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("User with login '%s' not found", req.Login)
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		log.Printf("Failed to query user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if users are already in groups
	var currentUserGroupID, targetUserGroupID sql.NullInt64
	d.DB.QueryRow(`SELECT id FROM groups WHERE user_id = $1`, user.Id).Scan(&currentUserGroupID)
	d.DB.QueryRow(`SELECT id FROM groups WHERE user_id = $1`, targetUserID).Scan(&targetUserGroupID)

	// Cannot send invite if both users are already in groups
	if currentUserGroupID.Valid && targetUserGroupID.Valid {
		// Check if they are in the same group
		if currentUserGroupID.Int64 == targetUserGroupID.Int64 {
			log.Printf("Both users %d and %d are already in the same group %d", user.Id, targetUserID, currentUserGroupID.Int64)
			http.Error(w, "users are already in the same group", http.StatusBadRequest)
			return
		}
		// Both users in different groups - cannot invite
		log.Printf("Cannot invite: user %d in group %d, user %d in group %d", user.Id, currentUserGroupID.Int64, targetUserID, targetUserGroupID.Int64)
		http.Error(w, "cannot invite users who are in different groups", http.StatusBadRequest)
		return
	}

	// Check group size limit for the group that will receive the new member
	var groupToCheck sql.NullInt64
	if currentUserGroupID.Valid {
		groupToCheck = currentUserGroupID
	} else if targetUserGroupID.Valid {
		groupToCheck = targetUserGroupID
	}

	if groupToCheck.Valid {
		var groupSize int
		err := d.DB.QueryRow(`SELECT COUNT(*) FROM groups WHERE id = $1`, groupToCheck.Int64).Scan(&groupSize)
		if err != nil {
			log.Printf("Failed to count group members: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if groupSize >= 5 {
			log.Printf("Group %d has reached maximum size of 5 members", groupToCheck.Int64)
			http.Error(w, "group has reached maximum size of 5 members", http.StatusBadRequest)
			return
		}
	}

	// Start transaction
	tx, err := d.DB.Begin()
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Check if reverse invite exists (mutual invite)
	var reverseInviteID int64
	err = tx.QueryRow(`SELECT id FROM invites WHERE from_user_id = $1 AND to_user_id = $2`, targetUserID, user.Id).Scan(&reverseInviteID)
	mutualInvite := err == nil

	if mutualInvite {
		log.Printf("Mutual invite detected between users %d and %d, creating group", user.Id, targetUserID)

		// Delete both invites
		_, err = tx.Exec(`DELETE FROM invites WHERE (from_user_id = $1 AND to_user_id = $2) OR (from_user_id = $2 AND to_user_id = $1)`, user.Id, targetUserID)
		if err != nil {
			log.Printf("Failed to delete invites: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Handle group creation based on current group status
		// We need to refresh group IDs within the transaction to ensure consistency
		var txCurrentUserGroupID, txTargetUserGroupID sql.NullInt64
		tx.QueryRow(`SELECT id FROM groups WHERE user_id = $1`, user.Id).Scan(&txCurrentUserGroupID)
		tx.QueryRow(`SELECT id FROM groups WHERE user_id = $1`, targetUserID).Scan(&txTargetUserGroupID)

		if txCurrentUserGroupID.Valid && txTargetUserGroupID.Valid {
			// Both users already in groups - this shouldn't happen due to earlier checks
			log.Printf("Error: Both users already in groups during mutual invite")
			http.Error(w, "both users are already in groups", http.StatusBadRequest)
			return
		} else if txCurrentUserGroupID.Valid {
			// Current user already has a group, add target user to it
			// Check group size limit again within transaction
			var groupSize int
			err = tx.QueryRow(`SELECT COUNT(*) FROM groups WHERE id = $1`, txCurrentUserGroupID.Int64).Scan(&groupSize)
			if err != nil {
				log.Printf("Failed to count group members: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if groupSize >= 5 {
				log.Printf("Group %d has reached maximum size", txCurrentUserGroupID.Int64)
				http.Error(w, "group has reached maximum size of 5 members", http.StatusBadRequest)
				return
			}

			// Add target user with next available member number
			nextMemberNumber := groupSize + 1
			_, err = tx.Exec(`INSERT INTO groups (id, user_id, member_number) VALUES ($1, $2, $3)`, txCurrentUserGroupID.Int64, targetUserID, nextMemberNumber)
			if err != nil {
				log.Printf("Failed to add user to existing group: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("Added user %d to existing group %d with member number %d", targetUserID, txCurrentUserGroupID.Int64, nextMemberNumber)
		} else if txTargetUserGroupID.Valid {
			// Target user already has a group, add current user to it
			var groupSize int
			err = tx.QueryRow(`SELECT COUNT(*) FROM groups WHERE id = $1`, txTargetUserGroupID.Int64).Scan(&groupSize)
			if err != nil {
				log.Printf("Failed to count group members: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if groupSize >= 5 {
				log.Printf("Group %d has reached maximum size", txTargetUserGroupID.Int64)
				http.Error(w, "group has reached maximum size of 5 members", http.StatusBadRequest)
				return
			}

			// Add current user with next available member number
			nextMemberNumber := groupSize + 1
			_, err = tx.Exec(`INSERT INTO groups (id, user_id, member_number) VALUES ($1, $2, $3)`, txTargetUserGroupID.Int64, user.Id, nextMemberNumber)
			if err != nil {
				log.Printf("Failed to add user to existing group: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("Added user %d to existing group %d with member number %d", user.Id, txTargetUserGroupID.Int64, nextMemberNumber)
		} else {
			// Neither user has a group, create a new one
			var newGroupID int64
			err = tx.QueryRow(`INSERT INTO groups (user_id, member_number) VALUES ($1, $2) RETURNING id`, user.Id, 1).Scan(&newGroupID)
			if err != nil {
				log.Printf("Failed to create new group: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = tx.Exec(`INSERT INTO groups (id, user_id, member_number) VALUES ($1, $2, $3)`, newGroupID, targetUserID, 2)
			if err != nil {
				log.Printf("Failed to add second user to new group: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("Created new group %d with users %d (member 1) and %d (member 2)", newGroupID, user.Id, targetUserID)
		}

		if err = tx.Commit(); err != nil {
			log.Printf("Failed to commit transaction: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "group created", "mutual_invite": true})
	} else {
		// No mutual invite, just create the invite
		var inviteID int64
		err = tx.QueryRow(`INSERT INTO invites (from_user_id, to_user_id) VALUES ($1, $2) RETURNING id`, user.Id, targetUserID).Scan(&inviteID)
		if err != nil {
			log.Printf("Failed to create invite: %v", err)
			http.Error(w, "failed to create invite (may already exist)", http.StatusBadRequest)
			return
		}

		if err = tx.Commit(); err != nil {
			log.Printf("Failed to commit transaction: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Successfully created invite %d from user %d to user %d", inviteID, user.Id, targetUserID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "invite sent", "invite_id": inviteID})
	}
}
