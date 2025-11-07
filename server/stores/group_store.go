package stores

import (
	"sync"
	"yuki_buy_log/database"
	"yuki_buy_log/models"
)

type GroupStore struct {
	data  map[models.GroupId][]models.GroupMember
	mutex sync.RWMutex
}

var (
	groupStoreInstance *GroupStore
	groupStoreOnce     sync.Once
)

func GetGroupStore() *GroupStore {
	groupStoreOnce.Do(func() {
		members, err := database.GetAllGroups()
		if err != nil {
			members = []models.GroupMember{}
		}

		// Преобразуем список членов в map[GroupId][]GroupMember
		groupMap := make(map[models.GroupId][]models.GroupMember)
		for _, member := range members {
			groupMap[member.GroupId] = append(groupMap[member.GroupId], member)
		}

		groupStoreInstance = &GroupStore{
			data: groupMap,
		}
	})
	return groupStoreInstance
}

// GetGroupById возвращает всех участников группы по ID группы
func (s *GroupStore) GetGroupById(id models.GroupId) []models.GroupMember {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if members, ok := s.data[id]; ok {
		// Возвращаем копию, чтобы избежать модификации извне
		result := make([]models.GroupMember, len(members))
		copy(result, members)
		return result
	}
	return nil
}

// GetGroupByUserId возвращает группу, в которой состоит пользователь
func (s *GroupStore) GetGroupByUserId(userId models.UserId) []models.GroupMember {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, members := range s.data {
		for _, member := range members {
			if member.UserId == userId {
				// Возвращаем копию
				result := make([]models.GroupMember, len(members))
				copy(result, members)
				return result
			}
		}
	}
	return nil
}

// IsUserInGroup проверяет, состоит ли пользователь в какой-либо группе
func (s *GroupStore) IsUserInGroup(userId models.UserId) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, members := range s.data {
		for _, member := range members {
			if member.UserId == userId {
				return true
			}
		}
	}
	return false
}

// CreateNewGroup создает новую группу с пользователем в качестве первого участника
func (s *GroupStore) CreateNewGroup(userId models.UserId) (*models.GroupId, error) {
	// Создаем группу в БД
	groupId, err := database.CreateNewGroup(userId)
	if err != nil {
		return nil, err
	}

	// Обновляем локальный store
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Получаем данные о пользователе из БД для создания полного объекта GroupMember
	members, err := database.GetGroupById(groupId)
	if err != nil {
		return &groupId, err
	}

	s.data[groupId] = members
	return &groupId, nil
}

// AddUserToGroup добавляет пользователя в группу
func (s *GroupStore) AddUserToGroup(groupId models.GroupId, userId models.UserId, memberNumber int) error {
	// Добавляем в БД
	err := database.AddUserToGroup(groupId, userId, memberNumber)
	if err != nil {
		return err
	}

	// Обновляем локальный store
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Получаем обновленные данные группы из БД
	members, err := database.GetGroupById(groupId)
	if err != nil {
		return err
	}

	s.data[groupId] = members
	return nil
}

// DeleteUserFromGroup удаляет пользователя из группы
func (s *GroupStore) DeleteUserFromGroup(userId models.UserId) error {
	// Сначала находим группу пользователя
	group := s.GetGroupByUserId(userId)
	if len(group) == 0 {
		return nil // Пользователь не в группе
	}
	groupId := group[0].GroupId

	// Удаляем из БД
	err := database.DeleteUserFromGroup(userId)
	if err != nil {
		return err
	}

	// Обновляем локальный store
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Проверяем, осталась ли группа после удаления
	count, err := database.GetGroupUserCount(groupId)
	if err != nil {
		return err
	}

	if count == 0 {
		// Группа пустая, удаляем её из store
		delete(s.data, groupId)
	} else {
		// Обновляем данные группы
		members, err := database.GetGroupById(groupId)
		if err != nil {
			return err
		}
		s.data[groupId] = members
	}

	return nil
}

// DeleteGroupById удаляет всю группу
func (s *GroupStore) DeleteGroupById(id models.GroupId) error {
	// Удаляем из БД
	err := database.DeleteGroupById(id)
	if err != nil {
		return err
	}

	// Удаляем из локального store
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.data, id)
	return nil
}

// RenumberGroupMembers перенумеровывает участников группы
func (s *GroupStore) RenumberGroupMembers(groupId models.GroupId) error {
	// Перенумеровываем в БД
	err := database.RenumberGroupMembers(groupId)
	if err != nil {
		return err
	}

	// Обновляем локальный store
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Получаем обновленные данные группы
	members, err := database.GetGroupById(groupId)
	if err != nil {
		return err
	}

	s.data[groupId] = members
	return nil
}
