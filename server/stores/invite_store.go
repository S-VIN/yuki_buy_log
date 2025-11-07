package stores

import (
	"sync"
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

func (s *InviteStore) AddInvite(invite models.Invite) error {
	// Добавляем в БД
	//_, err := database.CreateInvite(fromUserId, toUserId)
	//if err != nil {
	//	return err
	//}

	// Добавляем в локальный стор
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var maxId = models.InviteId(0)
	for _, item := range s.data {
		if item.Id > maxId {
			maxId = item.Id
		}
	}

	// При добавлении нового инвайта добавляем его с Id больше на 1
	invite.Id = maxId + 1
	s.data = append(s.data, invite)

	return nil
}

func (s *InviteStore) DeleteInvites(fromUserId, toUserId models.UserId) error {
	// Удаляем из локального стора
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Удаляем из БД
	err := database.DeleteInvitesBetweenUsers(fromUserId, toUserId)
	if err != nil {
		return err
	}

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
