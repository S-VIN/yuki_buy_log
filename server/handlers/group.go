package handlers

import (
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
		SELECT g.id, g.user_id, u.login
		FROM groups g
		JOIN users u ON g.user_id = u.id
		WHERE g.id = $1`, groupID)
	if err != nil {
		log.Printf("Failed to query group members: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	members := []models.GroupMember{}
	for rows.Next() {
		var m models.GroupMember
		if err := rows.Scan(&m.GroupId, &m.UserId, &m.Login); err != nil {
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
