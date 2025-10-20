package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"yuki_buy_log/models"
)

func InviteHandler(deps *Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Invite handler called: %s %s", r.Method, r.URL.Path)
		switch r.Method {
		case http.MethodGet:
			getIncomingInvites(deps, w, r)
		case http.MethodPost:
			sendInvite(deps, w, r)
		default:
			log.Printf("Method not allowed for invite: %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func getIncomingInvites(deps *Dependencies, w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching incoming invites from database")
	uid, ok := userID(r)
	if !ok {
		log.Println("Unauthorized access attempt to invites")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("Fetching incoming invites for user ID: %d", uid)

	rows, err := deps.DB.Query(`
		SELECT i.id, i.from_user_id, i.to_user_id, u_from.login, u_to.login, i.created_at
		FROM invites i
		JOIN users u_from ON i.from_user_id = u_from.id
		JOIN users u_to ON i.to_user_id = u_to.id
		WHERE i.to_user_id = $1`, uid)
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
	log.Printf("Successfully fetched %d incoming invites for user %d", len(invites), uid)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"invites": invites})
}

func sendInvite(deps *Dependencies, w http.ResponseWriter, r *http.Request) {
	log.Println("Sending invite")
	var req struct {
		Login string `json:"login"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode invite JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uid, ok := userID(r)
	if !ok {
		log.Println("Unauthorized access attempt to send invite")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get target user ID by login
	var targetUserID int64
	err := deps.DB.QueryRow(`SELECT id FROM users WHERE login = $1`, req.Login).Scan(&targetUserID)
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

	// Check if users are already in the same group
	var currentUserGroupID, targetUserGroupID sql.NullInt64
	deps.DB.QueryRow(`SELECT id FROM groups WHERE user_id = $1`, uid).Scan(&currentUserGroupID)
	deps.DB.QueryRow(`SELECT id FROM groups WHERE user_id = $1`, targetUserID).Scan(&targetUserGroupID)

	if currentUserGroupID.Valid && targetUserGroupID.Valid && currentUserGroupID.Int64 == targetUserGroupID.Int64 {
		log.Printf("Users %d and %d are already in the same group", uid, targetUserID)
		http.Error(w, "users are already in the same group", http.StatusBadRequest)
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

	// Check if reverse invite exists (mutual invite)
	var reverseInviteID int64
	err = tx.QueryRow(`SELECT id FROM invites WHERE from_user_id = $1 AND to_user_id = $2`, targetUserID, uid).Scan(&reverseInviteID)
	mutualInvite := err == nil

	if mutualInvite {
		log.Printf("Mutual invite detected between users %d and %d, creating group", uid, targetUserID)

		// Delete both invites
		_, err = tx.Exec(`DELETE FROM invites WHERE (from_user_id = $1 AND to_user_id = $2) OR (from_user_id = $2 AND to_user_id = $1)`, uid, targetUserID)
		if err != nil {
			log.Printf("Failed to delete invites: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create or merge into group
		if currentUserGroupID.Valid {
			// Current user already has a group, add target user to it
			_, err = tx.Exec(`INSERT INTO groups (id, user_id) VALUES ($1, $2)`, currentUserGroupID.Int64, targetUserID)
			if err != nil {
				log.Printf("Failed to add user to existing group: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("Added user %d to existing group %d", targetUserID, currentUserGroupID.Int64)
		} else if targetUserGroupID.Valid {
			// Target user already has a group, add current user to it
			_, err = tx.Exec(`INSERT INTO groups (id, user_id) VALUES ($1, $2)`, targetUserGroupID.Int64, uid)
			if err != nil {
				log.Printf("Failed to add user to existing group: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("Added user %d to existing group %d", uid, targetUserGroupID.Int64)
		} else {
			// Neither user has a group, create a new one
			var newGroupID int64
			err = tx.QueryRow(`INSERT INTO groups (user_id) VALUES ($1) RETURNING id`, uid).Scan(&newGroupID)
			if err != nil {
				log.Printf("Failed to create new group: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = tx.Exec(`INSERT INTO groups (id, user_id) VALUES ($1, $2)`, newGroupID, targetUserID)
			if err != nil {
				log.Printf("Failed to add second user to new group: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("Created new group %d with users %d and %d", newGroupID, uid, targetUserID)
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
		err = tx.QueryRow(`INSERT INTO invites (from_user_id, to_user_id) VALUES ($1, $2) RETURNING id`, uid, targetUserID).Scan(&inviteID)
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

		log.Printf("Successfully created invite %d from user %d to user %d", inviteID, uid, targetUserID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "invite sent", "invite_id": inviteID})
	}
}
