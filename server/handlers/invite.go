package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"yuki_buy_log/database"
	"yuki_buy_log/models"
)

func InviteHandler(d *Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Invite handler called: %s %s", r.Method, r.URL.Path)
		switch r.Method {
		case http.MethodGet:
			getIncomingInvites(w, r)
		case http.MethodPost:
			sendInvite(d, w, r)
		default:
			log.Printf("Method not allowed for invite: %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func getIncomingInvites(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching incoming invites from database")
	user, err := getUser(r)
	if err != nil {
		log.Println("Unauthorized access attempt to invites")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("Fetching incoming invites for user ID: %d", user.Id)

	invites, err := database.GetIncomingInvites(user.Id)
	if err != nil {
		log.Printf("Failed to query invites: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully fetched %d incoming invites for user %d", len(invites), user.Id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"invites": invites})
}

// Соединяет двух пользователей в группу. Если нельзя, то возвращает ошибку
func mergeUsersToGroups(firstUser models.User, secondUser models.User) error {
	// Check if users are already in groups
	firstUserGroupId, err := database.GetGroupIdByUserId(firstUser.Id)
	if err != nil {
		return fmt.Errorf()
	}
	secondUserGroupId, err := database.GetGroupIdByUserId(secondUser.Id)
	if err != nil {
		log.Printf("Failed to get target user group: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}


func sendInvite(w http.ResponseWriter, r *http.Request) {
	log.Println("Sending invite")
	var req struct {
		Login string `json:"login"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode invite JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := getUser(r)
	if err != nil {
		log.Println("Unauthorized access attempt to send invite")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get target user ID by login
	targetUser, err := database.GetUserByLogin(req.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("User with login '%s' not found", req.Login)
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		log.Printf("Failed to query user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = database.GetInvite(user.Id, targetUser.Id)
	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		log.Printf("Failed to query user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err == nil {
		http.Error(w, "invite already exists", http.StatusConflict)
		return
	}

	

	// Cannot send invite if both users are already in groups
	if (userGroupId != 0) && (targetUserGroupId != 0) {
		// Both users in different groups - cannot invite
		log.Printf("Cannot invite: user %d in group %d, user %d in group %d", user.Id, userGroupId, targetUser.Id, targetUserGroupId)
		http.Error(w, "cannot invite users who are in different groups", http.StatusBadRequest)
		return
	}

	// Никто не в группе. Два пользователя соединяются в группу
	if (userGroupId == 0) && (targetUserGroupId == 0) {
		// Создаем нового пользователя с member_number = 1
		newGroupId, err := database.CreateNewGroup(user.Id)
		if err != nil {
			log.Printf("Failed to create new group: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Добавляем второго пользователя с member_number = 2
		err = database.AddUserToGroup(newGroupId, targetUser.Id, 2)
		if err != nil {
			log.Printf("Failed to add new group to user: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Delete both invites
		err = database.DeleteInvitesBetweenUsers(user.Id, targetUserGroupId)
		if err != nil {
			log.Printf("Failed to delete invites: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}



	// Check if reverse invite exists (mutual invite)
	reverseInviteID, err := database.GetInvite(targetUserID, user.Id)
	mutualInvite := err == nil


		log.Printf("Mutual invite detected between users %d and %d, creating group", user.Id, targetUserID)



		// Handle group creation based on current group status
		if currentUserGroupID.Valid && targetUserGroupID.Valid {
			// Both users already in groups - this shouldn't happen due to earlier checks
			log.Printf("Error: Both users already in groups during mutual invite")
			http.Error(w, "both users are already in groups", http.StatusBadRequest)
			return
		} else if currentUserGroupID.Valid {
			// Current user already has a group, add target user to it
			// Check group size limit again
			groupSize, err := database.CountGroupMembers(currentUserGroupID.Int64)
			if err != nil {
				log.Printf("Failed to count group members: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if groupSize >= 5 {
				log.Printf("Group %d has reached maximum size", currentUserGroupID.Int64)
				http.Error(w, "group has reached maximum size of 5 members", http.StatusBadRequest)
				return
			}

			// Add target user with next available member number
			nextMemberNumber := groupSize + 1
			err = database.AddUserToGroup(currentUserGroupID.Int64, targetUserID, nextMemberNumber)
			if err != nil {
				log.Printf("Failed to add user to existing group: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("Added user %d to existing group %d with member number %d", targetUserID, currentUserGroupID.Int64, nextMemberNumber)
		} else if targetUserGroupID.Valid {
			// Target user already has a group, add current user to it
			groupSize, err := database.CountGroupMembers(targetUserGroupID.Int64)
			if err != nil {
				log.Printf("Failed to count group members: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if groupSize >= 5 {
				log.Printf("Group %d has reached maximum size", targetUserGroupID.Int64)
				http.Error(w, "group has reached maximum size of 5 members", http.StatusBadRequest)
				return
			}

			// Add current user with next available member number
			nextMemberNumber := groupSize + 1
			err = database.AddUserToGroup(targetUserGroupID.Int64, user.Id, nextMemberNumber)
			if err != nil {
				log.Printf("Failed to add user to existing group: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("Added user %d to existing group %d with member number %d", user.Id, targetUserGroupID.Int64, nextMemberNumber)
		} else {
			// Neither user has a group, create a new one
			newGroupID, err := database.CreateNewGroup(user.Id, 1)
			if err != nil {
				log.Printf("Failed to create new group: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = database.AddUserToGroup(newGroupID, targetUserID, 2)
			if err != nil {
				log.Printf("Failed to add second user to new group: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("Created new group %d with users %d (member 1) and %d (member 2)", newGroupID, user.Id, targetUserID)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "group created", "mutual_invite": true})
	} else {
		// No mutual invite, just create the invite
		inviteID, err := database.CreateInvite(user.Id, targetUserID)
		if err != nil {
			log.Printf("Failed to create invite: %v", err)
			http.Error(w, "failed to create invite (may already exist)", http.StatusBadRequest)
			return
		}

		log.Printf("Successfully created invite %d from user %d to user %d", inviteID, user.Id, targetUserID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "invite sent", "invite_id": inviteID})
	}
}
