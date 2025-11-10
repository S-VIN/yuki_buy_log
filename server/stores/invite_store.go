package stores

import (
	"sync"
	"time"
	"yuki_buy_log/database"
	"yuki_buy_log/models"
)

type InviteStore struct {
	data  []models.Invite
	mutex sync.RWMutex
}

var (
	inviteStoreInstance *InviteStore
	inviteStoreOnce     sync.Once
)

func GetInviteStore() *InviteStore {
	inviteStoreOnce.Do(func() {
		invites, err := database.GetAllInvites()
		if err != nil {
			invites = []models.Invite{}
		}
		inviteStoreInstance = &InviteStore{
			data: invites,
		}
	})
	return inviteStoreInstance
}

func (s *InviteStore) GetInviteById(id models.InviteId) *models.Invite {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, invite := range s.data {
		if invite.Id == id {
			return &invite
		}
	}
	return nil
}

func (s *InviteStore) GetInvitesFromUser(fromUserId models.UserId) []models.Invite {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var result []models.Invite
	for _, invite := range s.data {
		if invite.FromUserId == fromUserId {
			result = append(result, invite)
		}
	}
	return result
}

func (s *InviteStore) GetInvitesToUser(toUserId models.UserId) []models.Invite {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var result []models.Invite
	for _, invite := range s.data {
		if invite.ToUserId == toUserId {
			result = append(result, invite)
		}
	}
	return result
}

func (s *InviteStore) GetInvite(fromUserId, toUserId models.UserId) *models.Invite {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, invite := range s.data {
		if invite.FromUserId == fromUserId && invite.ToUserId == toUserId {
			return &invite
		}
	}
	return nil
}

// CreateInvite добавляет новое приглашение
func (s *InviteStore) CreateInvite(invite *models.Invite) error {
	// Добавляем в БД
	inviteId, err := database.CreateInvite(invite.FromUserId, invite.ToUserId)
	if err != nil {
		return err
	}

	// Обновляем ID в инвайте
	invite.Id = inviteId

	// Добавляем в локальный стор
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data = append(s.data, *invite)
	return nil
}

// DeleteInvite удаляет приглашение по ID
func (s *InviteStore) DeleteInvite(id models.InviteId) error {
	// Удаляем из БД
	err := database.DeleteInvite(id)
	if err != nil {
		return err
	}

	// Удаляем из локального стора
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var newData []models.Invite
	for _, invite := range s.data {
		if invite.Id != id {
			newData = append(newData, invite)
		}
	}
	s.data = newData

	return nil
}

// DeleteInviteByUsers удаляет приглашения между пользователями
func (s *InviteStore) DeleteInviteByUsers(fromUserId, toUserId models.UserId) error {
	// Удаляем из БД
	err := database.DeleteInvitesBetweenUsers(fromUserId, toUserId)
	if err != nil {
		return err
	}

	// Удаляем из локального стора
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var newData []models.Invite
	for _, invite := range s.data {
		if !((invite.FromUserId == fromUserId && invite.ToUserId == toUserId) ||
			(invite.FromUserId == toUserId && invite.ToUserId == fromUserId)) {
			newData = append(newData, invite)
		}
	}
	s.data = newData

	return nil
}

func (s *InviteStore) DeleteOldInvites(cutoffTime time.Time) (int64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Удаляем из БД
	rowsAffected, err := database.DeleteOldInvites(cutoffTime)
	if err != nil {
		return 0, err
	}

	// Удаляем из локального стора
	var newData []models.Invite
	for _, invite := range s.data {
		if invite.CreatedAt.After(cutoffTime) || invite.CreatedAt.Equal(cutoffTime) {
			newData = append(newData, invite)
		}
	}
	s.data = newData

	return rowsAffected, nil
}
