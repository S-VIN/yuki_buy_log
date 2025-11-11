package stores

import (
	"sync"
	"yuki_buy_log/database"
	"yuki_buy_log/models"
)

type UserStore struct {
	db    database.Database
	data  map[models.UserId]models.User
	mutex sync.RWMutex
}

var (
	userStoreInstance *UserStore
	userStoreOnce     sync.Once
)

func GetUserStore(db database.Database) *UserStore {
	userStoreOnce.Do(func() {
		users, err := db.GetAllUsers()
		if err != nil {
			users = []models.User{}
		}

		// Преобразуем список пользователей в map[UserId]User
		userMap := make(map[models.UserId]models.User)
		for _, user := range users {
			userMap[user.Id] = user
		}

		userStoreInstance = &UserStore{
			db:   db,
			data: userMap,
		}
	})
	return userStoreInstance
}

// GetUserById возвращает пользователя по ID
func (s *UserStore) GetUserById(id models.UserId) *models.User {
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
func (s *UserStore) GetUserByLogin(login string) *models.User {
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
func (s *UserStore) AddUser(user *models.User) error {
	// Добавляем в БД
	err := s.db.AddUser(user)
	if err != nil {
		return err
	}

	// Обновляем локальный стор
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[user.Id] = *user
	return nil
}

// UpdateUser обновляет данные пользователя
func (s *UserStore) UpdateUser(user *models.User) error {
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
func (s *UserStore) DeleteUser(userId models.UserId) error {
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
func (s *UserStore) GetUsersByGroupId(groupId models.GroupId) []models.User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Получаем группу
	groupStore := GetGroupStore(s.db)
	group := groupStore.GetGroupById(groupId)
	if group == nil {
		return nil
	}

	// Собираем пользователей
	var users []models.User
	for _, member := range group.Members {
		if user, ok := s.data[member.UserId]; ok {
			users = append(users, user)
		}
	}

	return users
}

// GetAllUsers возвращает всех пользователей
func (s *UserStore) GetAllUsers() []models.User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	users := make([]models.User, 0, len(s.data))
	for _, user := range s.data {
		users = append(users, user)
	}
	return users
}
