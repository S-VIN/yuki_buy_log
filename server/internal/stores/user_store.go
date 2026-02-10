package stores

import (
	"sync"
	"yuki_buy_log/internal/database"
	"yuki_buy_log/internal/domain"
)

var (
	userStoreInstance *UserStore
	userStoreLock     sync.Once
)

type UserStore struct {
	data  map[domain.UserId]domain.User
	mutex sync.RWMutex
	db    database.DatabaseManager
}

func GetUserStore() *UserStore {
	userStoreLock.Do(func() {
		var db, _ = database.GetDBManager()
		users, err := db.GetAllUsers()
		if err != nil {
			users = []domain.User{}
		}

		// Преобразуем список пользователей в map[UserId]User
		userMap := make(map[domain.UserId]domain.User)
		for _, user := range users {
			userMap[user.Id] = user
		}

		userStoreInstance = &UserStore{
			data: userMap,
			db:   *db,
		}
	})
	return userStoreInstance
}

// GetUserById возвращает пользователя по ID
func (s *UserStore) GetUserById(id domain.UserId) *domain.User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if user, ok := s.data[id]; ok {
		// Возвращаем копию, чтобы избежать модификации извне
		userCopy := user
		return &userCopy
	}
	return nil
}

// GetUserByLogin возвращает пользователя по логину
func (s *UserStore) GetUserByLogin(login string) *domain.User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, user := range s.data {
		if user.Login == login {
			// Возвращаем копию
			userCopy := user
			return &userCopy
		}
	}
	return nil
}

// AddUser добавляет нового пользователя
func (s *UserStore) AddUser(user *domain.User) error {
	// Обновляем локальный стор
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Добавляем в БД
	err := s.db.AddUser(user)
	if err != nil {
		return err
	}

	s.data[user.Id] = *user
	return nil
}

// UpdateUser обновляет данные пользователя
func (s *UserStore) UpdateUser(user *domain.User) error {
	// Обновляем в БД
	err := s.db.UpdateUser(user)
	if err != nil {
		return err
	}

	// Обновляем локальный стор
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[user.Id] = *user
	return nil
}

// DeleteUser удаляет пользователя
func (s *UserStore) DeleteUser(userId domain.UserId) error {
	// Удаляем из БД
	err := s.db.DeleteUser(userId)
	if err != nil {
		return err
	}

	// Удаляем из локального стора
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.data, userId)
	return nil
}

// GetUsersByGroupId возвращает всех пользователей группы
func (s *UserStore) GetUsersByGroupId(groupId domain.GroupId) []domain.User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Получаем группу
	var groupStore = GetGroupStore()
	group := groupStore.GetGroupById(groupId)
	if group == nil {
		return nil
	}

	// Собираем пользователей
	var users []domain.User
	for _, member := range group.Members {
		if user, ok := s.data[member.UserId]; ok {
			users = append(users, user)
		}
	}

	return users
}

// GetAllUsers возвращает всех пользователей
func (s *UserStore) GetAllUsers() []domain.User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	users := make([]domain.User, 0, len(s.data))
	for _, user := range s.data {
		users = append(users, user)
	}
	return users
}
