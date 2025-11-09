package stores

import (
	"errors"
	"sort"
	"sync"
	"yuki_buy_log/database"
	"yuki_buy_log/models"
)

type GroupStore struct {
	groupIdByUserId map[models.UserId]models.GroupId
	groupById       map[models.GroupId]models.Group
	mutex           sync.RWMutex
}

var (
	groupStoreInstance *GroupStore
	groupStoreOnce     sync.Once
)

var (
	ErrMaxMembersInGroup = errors.New("max members in group")
	ErrNotFound          = errors.New("not found")
)

func (s *GroupStore) getMaxMemberNumber(members []models.GroupMember) int {
	maxMemberNumber := 0
	for _, value := range members {
		if value.MemberNumber > maxMemberNumber {
			maxMemberNumber = value.MemberNumber
		}
	}
	return maxMemberNumber
}

func (s *GroupStore) needRenumberMembers(members []models.GroupMember) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	sort.Slice(members, func(i, j int) bool { return members[i].MemberNumber < members[j].MemberNumber })
	for index, value := range members {
		if index != value.MemberNumber {
			return true
		}
	}
	return false
}

func (s *GroupStore) renumberMembers(members []models.GroupMember) {
	if len(members) == 0 {
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	sort.Slice(members, func(i, j int) bool { return members[i].MemberNumber < members[j].MemberNumber })
	for index := range members {
		if index != members[index].MemberNumber {
			members[index].MemberNumber = index
			database.UpdateGroupMember(&members[index])
		}
	}
	group := s.groupById[members[0].GroupId]
	group.Members = members
	s.groupById[members[0].GroupId] = group
	return
}

func GetGroupStore() *GroupStore {
	groupStoreOnce.Do(func() {
		members, err := database.GetAllGroupMembers()
		if err != nil {
			members = []models.GroupMember{}
		}

		groupStoreInstance = &GroupStore{}
		for _, member := range members {
			if value, ok := groupStoreInstance.groupById[member.GroupId]; ok {
				value.Members = append(value.Members, member)
			} else {
				value = models.Group{
					Id:      member.GroupId,
					Members: []models.GroupMember{member},
				}
			}
		}
	})
	return groupStoreInstance
}

// GetGroupById возвращает всех участников группы по ID группы
func (s *GroupStore) GetGroupById(id models.GroupId) *models.Group {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if group, ok := s.groupById[id]; ok {
		// Возвращаем ссылку на копию, чтобы избежать модификации извне
		result := group
		return &result
	}
	return nil
}

// GetGroupIdByUserId возвращает ID группы, в которой состоит пользователь
func (s *GroupStore) GetGroupIdByUserId(userId models.UserId) *models.GroupId {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if groupId, ok := s.groupIdByUserId[userId]; ok {
		result := groupId
		return &result
	}
	return nil
}

// GetGroupByUserId возвращает группу, в которой состоит пользователь
func (s *GroupStore) GetGroupByUserId(userId models.UserId) *models.Group {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	groupId := s.GetGroupIdByUserId(userId)
	if groupId == nil {
		return nil
	}

	if group, ok := s.groupById[*groupId]; ok {
		result := group
		return &result
	}

	return nil
}

// GetGroupUserCount возвращает количество участников в группе
func (s *GroupStore) GetGroupUserCount(groupId models.GroupId) int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	group := s.GetGroupById(groupId)
	if group == nil {
		return 0
	}
	return len(group.Members)
}

// IsUserInGroup проверяет, состоит ли пользователь в какой-либо группе
func (s *GroupStore) IsUserInGroup(userId models.UserId) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	group := s.GetGroupByUserId(userId)
	if group == nil {
		return false
	}

	for _, member := range group.Members {
		if member.UserId == userId {
			return true
		}
	}
	return false
}

// CreateNewGroup создает новую группу с пользователем в качестве первого участника
func (s *GroupStore) CreateNewGroup(userId models.UserId) (*models.GroupId, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Создаем группу в БД
	groupId, err := database.CreateNewGroup(userId)
	if err != nil {
		return nil, err
	}

	// Получаем данные о пользователе из БД для создания полного объекта GroupMember
	members, err := database.GetGroupMembersByGroupId(groupId)
	if err != nil {
		return &groupId, err
	}

	s.groupById[groupId] = models.Group{Id: groupId, Members: members}
	s.groupIdByUserId[userId] = groupId
	return &groupId, nil
}

// AddUserToGroup добавляет пользователя в группу
func (s *GroupStore) AddUserToGroup(groupId models.GroupId, userId models.UserId) error {
	// Обновляем локальный store
	s.mutex.Lock()
	defer s.mutex.Unlock()

	group := s.GetGroupById(groupId)
	if group == nil {
		return nil
	}

	maxMemberNumber := s.getMaxMemberNumber(group.Members)
	if maxMemberNumber >= 5 {
		return ErrMaxMembersInGroup
	}

	group.Members = append(group.Members, models.GroupMember{GroupId: groupId, UserId: userId, MemberNumber: maxMemberNumber + 1})

	// Добавляем в БД
	err := database.AddUserToGroup(groupId, userId, maxMemberNumber+1)
	if err != nil {
		return err
	}

	s.groupById[groupId] = *group
	return nil
}

// DeleteUserFromGroup удаляет пользователя из группы
func (s *GroupStore) DeleteUserFromGroup(userId models.UserId) error {
	// Обновляем локальный store
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Сначала находим группу пользователя
	group := s.GetGroupByUserId(userId)
	if group == nil {
		return ErrNotFound
	}
	if len(group.Members) == 0 {
		return nil // Пользователь не в группе
	}

	// удаляем из внутренних структур
	for i, member := range group.Members {
		if member.UserId == userId {
			group.Members[i] = group.Members[len(group.Members)-1]
			group.Members = group.Members[:len(group.Members)-1] // уменьшаем slice
			break
		}
	}

	// Удаляем из БД
	err := database.DeleteUserFromGroup(userId)
	if err != nil {
		return err
	}

	// Если остался 1 или 0 членов группы, то группа распалась - удаляем
	if len(group.Members) < 2 {
		s.DeleteGroupById(group.Id)
	}

	return nil
}

// DeleteGroupById удаляет всю группу
func (s *GroupStore) DeleteGroupById(id models.GroupId) error {
	// Удаляем из локального store
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Удаляем из БД
	err := database.DeleteGroupMembersByGroupId(id)
	if err != nil {
		return err
	}

	delete(s.groupById, id)
	for key, value := range s.groupIdByUserId {
		if value == id {
			delete(s.groupIdByUserId, key)
			break
		}
	}

	return nil
}
