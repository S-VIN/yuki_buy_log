package stores

import (
	"sync"
	"time"
	"yuki_buy_log/internal/database"
	"yuki_buy_log/internal/domain"
)

type InviteStore struct {
	data  []domain.Invite
	mutex sync.RWMutex
	db    database.DatabaseManager
}

var (
	inviteStoreInstance *InviteStore
	inviteStoreLock     sync.Once
)

func GetInviteStore() *InviteStore {
	inviteStoreLock.Do(func() {
		var db, _ = database.GetDBManager()

		invites, err := db.GetAllInvites()
		if err != nil {
			invites = []domain.Invite{}
		}

		inviteStoreInstance = &InviteStore{
			data: invites,
			db:   *db,
		}
	})
	return inviteStoreInstance
}

func (s *InviteStore) GetInviteById(id domain.InviteId) *domain.Invite {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, invite := range s.data {
		if invite.Id == id {
			return &invite
		}
	}
	return nil
}

func (s *InviteStore) GetInvitesFromUser(fromUserId domain.UserId) []domain.Invite {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var result []domain.Invite
	for _, invite := range s.data {
		if invite.FromUserId == fromUserId {
			result = append(result, invite)
		}
	}
	return result
}

func (s *InviteStore) GetInvitesToUser(toUserId domain.UserId) []domain.Invite {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var result []domain.Invite
	for _, invite := range s.data {
		if invite.ToUserId == toUserId {
			result = append(result, invite)
		}
	}
	return result
}

func (s *InviteStore) GetInvite(fromUserId, toUserId domain.UserId) *domain.Invite {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, invite := range s.data {
		if invite.FromUserId == fromUserId && invite.ToUserId == toUserId {
			return &invite
		}
	}
	return nil
}

func (s *InviteStore) AddInvite(invite domain.Invite) (domain.InviteId, error) {
	// Добавляем в БД
	//_, err := database.CreateInvite(fromUserId, toUserId)
	//if err != nil {
	//	return 0, err
	//}

	// Добавляем в локальный стор
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var maxId = domain.InviteId(0)
	for _, item := range s.data {
		if item.Id > maxId {
			maxId = item.Id
		}
	}

	// При добавлении нового инвайта добавляем его с Id больше на 1
	invite.Id = maxId + 1
	s.data = append(s.data, invite)

	return invite.Id, nil
}

func (s *InviteStore) DeleteInvites(fromUserId, toUserId domain.UserId) error {
	// Удаляем из локального стора
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Удаляем из БД
	err := s.db.DeleteInvitesBetweenUsers(fromUserId, toUserId)
	if err != nil {
		return err
	}

	var newData []domain.Invite
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
	rowsAffected, err := s.db.DeleteOldInvites(cutoffTime)
	if err != nil {
		return 0, err
	}

	// Удаляем из локального стора
	var newData []domain.Invite
	for _, invite := range s.data {
		if invite.CreatedAt.After(cutoffTime) || invite.CreatedAt.Equal(cutoffTime) {
			newData = append(newData, invite)
		}
	}
	s.data = newData

	return rowsAffected, nil
}
