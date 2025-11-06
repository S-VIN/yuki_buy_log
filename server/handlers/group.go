package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"yuki_buy_log/database"
	"yuki_buy_log/models"
)

func GroupHandler(d *Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Group handler called: %s %s", r.Method, r.URL.Path)
		switch r.Method {
		case http.MethodGet:
			getGroupMembers(d, w, r)
		case http.MethodDelete:
			leaveGroup(d, w, r)
		default:
			log.Printf("Method not allowed for group: %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func getGroupMembers(d *Dependencies, w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching group members from database")
	user, err := getUser(r)
	if err != nil {
		log.Println("Unauthorized access attempt to group")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("Fetching group members for user ID: %d", user.Id)

	// Get the group_id for the current user
	groupId, err := database.GetGroupIdByUserId(user.Id)
	if err != nil {
		log.Printf("User %d is not in any group", user.Id)
		// Return empty list if user is not in a group
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"members":         []models.GroupMember{},
			"current_user_id": user.Id,
		})
		return
	}

	// Get all members of the same group
	members, err := database.GetGroupById(groupId)
	log.Printf("Successfully fetched %d group members for user %d", len(members), user.Id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"members":         members,
		"current_user_id": user.Id,
	})
}

func leaveGroup(d *Dependencies, w http.ResponseWriter, r *http.Request) {
	log.Println("User leaving group")
	user, err := getUser(r)
	if err != nil {
		log.Println("Unauthorized access attempt to leave group")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the group_id for the current user
	groupId, err := database.GetGroupIdByUserId(user.Id)

	// Remove user from group
	err = database.DeleteUserFromGroup(user.Id)
	if err != nil {
		log.Printf("Failed to remove user from group: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check remaining members count
	count, err := database.GetGroupUserCount(groupId)
	if err != nil {
		log.Printf("Failed to count remaining group members: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If only 1 member remains, delete the entire group
	if count == 1 {
		log.Printf("Only 1 member remains in group %d, deleting group", groupId)
		err = database.DeleteGroupById(groupId)
		if err != nil {
			log.Printf("Failed to delete group: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if count > 1 {
		// Renumber remaining members to fill gaps
		err = database.RenumberGroupMembers(groupId)
		if err != nil {
			log.Printf("Failed to renumber group members: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	log.Printf("User %d successfully left group %d", user.Id, groupId)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "left group successfully"})
}
