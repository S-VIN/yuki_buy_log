package stores

import (
	"errors"
	"sort"
	"sync"
	"yuki_buy_log/internal/database"
	"yuki_buy_log/internal/domain"
)

type GroupStore struct {
	groupIdByUserId map[domain.UserId]domain.GroupId
	groupById       map[domain.GroupId]domain.Group
	mutex           sync.RWMutex
	db              database.DatabaseManager
}

var (
	ErrMaxMembersInGroup = errors.New("max members in group")
	ErrNotFound          = errors.New("not found")
)

var (
	groupStoreInstance *GroupStore
	groupStoreLock     sync.Once
)

func (s *GroupStore) getMaxMemberNumber(members []domain.GroupMember) int {
	maxMemberNumber := 0
	for _, value := range members {
		if value.MemberNumber > maxMemberNumber {
			maxMemberNumber = value.MemberNumber
		}
	}
	return maxMemberNumber
}

func (s *GroupStore) needRenumberMembers(members []domain.GroupMember) bool {
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

func (s *GroupStore) renumberMembers(members []domain.GroupMember) {
	if len(members) == 0 {
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	sort.Slice(members, func(i, j int) bool { return members[i].MemberNumber < members[j].MemberNumber })
	for index := range members {
		if index != members[index].MemberNumber {
			members[index].MemberNumber = index
			s.db.UpdateGroupMember(&members[index])
		}
	}
	group := s.groupById[members[0].GroupId]
	group.Members = members
	s.groupById[members[0].GroupId] = group
	return
}

func GetGroupStore() *GroupStore {
	groupStoreLock.Do(func() {
		var db, _ = database.GetDBManager()
		members, err := db.GetAllGroupMembers()
		if err != nil {
			members = []domain.GroupMember{}
		}

		groupStoreInstance = &GroupStore{
			groupIdByUserId: make(map[domain.UserId]domain.GroupId),
			groupById:       make(map[domain.GroupId]domain.Group),
			db:              *db,
		}
		for _, member := range members {
			if value, ok := groupStoreInstance.groupById[member.GroupId]; ok {
				value.Members = append(value.Members, member)
				groupStoreInstance.groupById[member.GroupId] = value
			} else {
				value = domain.Group{
					Id:      member.GroupId,
					Members: []domain.GroupMember{member},
				}
				groupStoreInstance.groupById[member.GroupId] = value
			}
			groupStoreInstance.groupIdByUserId[member.UserId] = member.GroupId
		}
	})
	return groupStoreInstance
}

// GetGroupById возвращает всех участников группы по ID группы
func (s *GroupStore) GetGroupById(id domain.GroupId) *domain.Group {
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
func (s *GroupStore) GetGroupIdByUserId(userId domain.UserId) *domain.GroupId {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if groupId, ok := s.groupIdByUserId[userId]; ok {
		result := groupId
		return &result
	}
	return nil
}

// GetGroupByUserId возвращает группу, в которой состоит пользователь
func (s *GroupStore) GetGroupByUserId(userId domain.UserId) *domain.Group {
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
func (s *GroupStore) GetGroupUserCount(groupId domain.GroupId) int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	group := s.GetGroupById(groupId)
	if group == nil {
		return 0
	}
	return len(group.Members)
}

// IsUserInGroup проверяет, состоит ли пользователь в какой-либо группе
func (s *GroupStore) IsUserInGroup(userId domain.UserId) bool {
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
func (s *GroupStore) CreateNewGroup(userId domain.UserId) (*domain.GroupId, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Создаем группу в БД
	groupId, err := s.db.CreateNewGroup(userId)
	if err != nil {
		return nil, err
	}

	// Получаем данные о пользователе из БД для создания полного объекта GroupMember
	members, err := s.db.GetGroupMembersByGroupId(groupId)
	if err != nil {
		return &groupId, err
	}

	s.groupById[groupId] = domain.Group{Id: groupId, Members: members}
	s.groupIdByUserId[userId] = groupId
	return &groupId, nil
}

// AddUserToGroup добавляет пользователя в группу
func (s *GroupStore) AddUserToGroup(groupId domain.GroupId, userId domain.UserId) error {
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

	group.Members = append(group.Members, domain.GroupMember{GroupId: groupId, UserId: userId, MemberNumber: maxMemberNumber + 1})

	// Добавляем в БД
	err := s.db.AddUserToGroup(groupId, userId, maxMemberNumber+1)
	if err != nil {
		return err
	}

	s.groupById[groupId] = *group
	return nil
}

// DeleteUserFromGroup удаляет пользователя из группы
func (s *GroupStore) DeleteUserFromGroup(userId domain.UserId) error {
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
	err := s.db.DeleteUserFromGroup(userId)
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
func (s *GroupStore) DeleteGroupById(id domain.GroupId) error {
	// Удаляем из локального store
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Удаляем из БД
	err := s.db.DeleteGroupMembersByGroupId(id)
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
