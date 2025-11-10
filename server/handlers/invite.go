package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"yuki_buy_log/models"
)

func InviteHandler(auth Authenticator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Invite handler called: %s %s", r.Method, r.URL.Path)
		switch r.Method {
		case http.MethodGet:
			getIncomingInvites(w, r)
		case http.MethodPost:
			sendInvite(w, r)
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
	inviteStore := GetInviteStore()
	invites := inviteStore.GetInvitesToUser(user.Id)
	log.Printf("Successfully fetched %d incoming invites for user %d", len(invites), user.Id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"invites": invites})
}

func canMergeUsersToGroups(firstUserId models.UserId, secondUserId models.UserId) bool {
	groupStore := GetGroupStore()

	// Check if users are already in groups
	firstUserGroup := groupStore.GetGroupByUserId(firstUserId)
	secondUserGroup := groupStore.GetGroupByUserId(secondUserId)
	if (firstUserGroup != nil) && (secondUserGroup == nil) {
		return false
	}
	return true
}

// // Соединяет двух пользователей в группу. Если нельзя, то возвращает ошибку
func mergeUsersToGroups(firstUserId models.UserId, secondUserId models.UserId) error {
	groupStore := GetGroupStore()

	// Check if users are already in groups
	firstUserGroup := groupStore.GetGroupByUserId(firstUserId)
	secondUserGroup := groupStore.GetGroupByUserId(secondUserId)

	if (firstUserGroup == nil) && (secondUserGroup == nil) {
		// Создаем нового пользователя с member_number = 1
		newGroupId, err := groupStore.CreateNewGroup(firstUserId)
		if err != nil {
			log.Printf("Failed to create new group: %v", err)
			return err
		}

		err = groupStore.AddUserToGroup(*newGroupId, secondUserId)
		if err != nil {
			log.Printf("Failed to add new group to user: %v", err)
			return err
		}
	}

	// Первый юзер в группе, а второй не в группе
	// Добавляем второго пользователя в группу к первому
	if firstUserGroup != nil {
		// Добавляем второго пользователя
		err := groupStore.AddUserToGroup(firstUserGroup.Id, secondUserId)
		if err != nil {
			log.Printf("Failed to add new group to user: %v", err)
			return err
		}
	}

	// Тоже самое только наоборот
	if secondUserGroup != nil {
		// Добавляем первого пользователя
		err := groupStore.AddUserToGroup(secondUserGroup.Id, firstUserId)
		if err != nil {
			log.Printf("Failed to add new group to user: %v", err)
			return err
		}
	}

	return nil
}

func sendInvite(w http.ResponseWriter, r *http.Request) {
	userStore := GetUserStore()
	inviteStore := GetInviteStore()

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
	targetUser := userStore.GetUserByLogin(req.Login)
	if targetUser == nil {
		log.Printf("User with login '%s' not found", req.Login)
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// Проверяем был ли уже точно такой инвайт
	invite := inviteStore.GetInvite(user.Id, targetUser.Id)
	if invite != nil {
		http.Error(w, "Invite already exists", http.StatusInternalServerError)
		return
	}

	// Не можем объединить пользователей в группы
	if !canMergeUsersToGroups(user.Id, targetUser.Id) {
		log.Printf("Cannot invite: user %d and user %d are in different groups", user.Id, targetUser.Id)
		http.Error(w, "cannot invite users who are in different groups", http.StatusBadRequest)
		return
	}

	// Если уже есть инвайт от противоположного пользователя - соединяем в группу
	oppositeInvite := inviteStore.GetInvite(targetUser.Id, user.Id)
	if oppositeInvite != nil {
		err = mergeUsersToGroups(user.Id, targetUser.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = inviteStore.DeleteInviteByUsers(user.Id, targetUser.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Successfully created group from user %d to user %d", targetUser.Id, user.Id)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "group created", "userId": user.Id, "targetUser": targetUser.Id})
		return
	}

	// Остался вариант, что мы можем отправить инвайт, но объединить в группу не можем пока, потому что нет обратного инвайта
	var newInvite = models.Invite{
		FromUserId: user.Id,
		ToUserId:   targetUser.Id,
		FromLogin:  user.Login,
		ToLogin:    targetUser.Login,
		CreatedAt:  time.Now(),
	}
	err = inviteStore.CreateInvite(&newInvite)
	if err != nil {
		log.Printf("Cannot send invite")
		http.Error(w, "Cannot invite users", http.StatusBadRequest)
		return
	}
}
