package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"yuki_buy_log/models"
)

func GroupHandler(auth Authenticator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Group handler called: %s %s", r.Method, r.URL.Path)
		switch r.Method {
		case http.MethodGet:
			getGroupMembers(w, r)
		case http.MethodDelete:
			leaveGroup(w, r)
		default:
			log.Printf("Method not allowed for group: %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func getGroupMembers(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching group members from store")
	user, err := getUser(r)
	if err != nil {
		log.Println("Unauthorized access attempt to group")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("Fetching group members for user ID: %d", user.Id)

	groupStore := getGroupStore()

	// Get the group_id for the current user
	group := groupStore.GetGroupByUserId(user.Id)
	if group == nil {
		log.Printf("User %d is not in any group", user.Id)
		// Return empty list if user is not in a group
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"members":         []models.GroupMember{},
			"current_user_id": user.Id,
		})
		return
	}

	log.Printf("Successfully fetched %d group members for user %d", len(group.Members), user.Id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"members":         group.Members,
		"current_user_id": user.Id,
	})
}

func leaveGroup(w http.ResponseWriter, r *http.Request) {
	log.Println("User leaving group")
	user, err := getUser(r)
	if err != nil {
		log.Println("Unauthorized access attempt to leave group")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	groupStore := getGroupStore()

	// Get the group_id for the current user
	group := groupStore.GetGroupByUserId(user.Id)
	if group == nil {
		log.Printf("User %d is not in any group: %v", user.Id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Remove user from group
	err = groupStore.DeleteUserFromGroup(user.Id)
	if err != nil {
		log.Printf("Failed to remove user from group: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("User %d successfully left group %d", user.Id, group.Id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "left group successfully"})
}
