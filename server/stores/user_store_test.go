package stores

import (
	"errors"
	"testing"
	"yuki_buy_log/mocks"
	"yuki_buy_log/models"

	"go.uber.org/mock/gomock"
)

func TestNewUserStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockIDataBaseManager(ctrl)

	t.Run("успешное создание store с пользователями", func(t *testing.T) {
		users := []models.User{
			{Id: 1, Login: "user1", Password: "pass1"},
			{Id: 2, Login: "user2", Password: "pass2"},
		}

		mockDB.EXPECT().GetAllUsers().Return(users, nil)

		store, err := NewUserStore(mockDB)
		if err != nil {
			t.Fatalf("ожидали успешное создание store, получили ошибку: %v", err)
		}

		if store == nil {
			t.Fatal("store не должен быть nil")
		}

		if len(store.data) != 2 {
			t.Errorf("ожидали 2 пользователей в store, получили %d", len(store.data))
		}

		if store.data[1].Login != "user1" {
			t.Errorf("ожидали login 'user1', получили '%s'", store.data[1].Login)
		}
	})

	t.Run("создание store при ошибке получения пользователей", func(t *testing.T) {
		mockDB.EXPECT().GetAllUsers().Return(nil, errors.New("database error"))

		store, err := NewUserStore(mockDB)
		if err != nil {
			t.Fatalf("не ожидали ошибку при создании store: %v", err)
		}

		if store == nil {
			t.Fatal("store не должен быть nil")
		}

		if len(store.data) != 0 {
			t.Errorf("ожидали пустой store, получили %d пользователей", len(store.data))
		}
	})
}

func TestGetUserById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockIDataBaseManager(ctrl)
	users := []models.User{
		{Id: 1, Login: "user1", Password: "pass1"},
		{Id: 2, Login: "user2", Password: "pass2"},
	}

	mockDB.EXPECT().GetAllUsers().Return(users, nil)
	store, _ := NewUserStore(mockDB)

	t.Run("получение существующего пользователя", func(t *testing.T) {
		user := store.GetUserById(1)
		if user == nil {
			t.Fatal("ожидали найти пользователя, получили nil")
		}

		if user.Login != "user1" {
			t.Errorf("ожидали login 'user1', получили '%s'", user.Login)
		}
	})

	t.Run("получение несуществующего пользователя", func(t *testing.T) {
		user := store.GetUserById(999)
		if user != nil {
			t.Error("ожидали nil для несуществующего пользователя")
		}
	})

	t.Run("возвращается копия пользователя", func(t *testing.T) {
		user := store.GetUserById(1)
		if user == nil {
			t.Fatal("пользователь не найден")
		}

		originalLogin := user.Login
		user.Login = "modified"

		user2 := store.GetUserById(1)
		if user2.Login != originalLogin {
			t.Error("изменение возвращенного пользователя повлияло на store")
		}
	})
}

func TestGetUserByLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockIDataBaseManager(ctrl)
	users := []models.User{
		{Id: 1, Login: "user1", Password: "pass1"},
		{Id: 2, Login: "user2", Password: "pass2"},
	}

	mockDB.EXPECT().GetAllUsers().Return(users, nil)
	store, _ := NewUserStore(mockDB)

	t.Run("получение пользователя по существующему login", func(t *testing.T) {
		user := store.GetUserByLogin("user1")
		if user == nil {
			t.Fatal("ожидали найти пользователя, получили nil")
		}

		if user.Id != 1 {
			t.Errorf("ожидали ID 1, получили %d", user.Id)
		}
	})

	t.Run("получение пользователя по несуществующему login", func(t *testing.T) {
		user := store.GetUserByLogin("nonexistent")
		if user != nil {
			t.Error("ожидали nil для несуществующего login")
		}
	})

	t.Run("возвращается копия пользователя", func(t *testing.T) {
		user := store.GetUserByLogin("user1")
		if user == nil {
			t.Fatal("пользователь не найден")
		}

		originalLogin := user.Login
		user.Login = "modified"

		user2 := store.GetUserByLogin("user1")
		if user2.Login != originalLogin {
			t.Error("изменение возвращенного пользователя повлияло на store")
		}
	})
}

func TestAddUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockIDataBaseManager(ctrl)

	t.Run("успешное добавление пользователя", func(t *testing.T) {
		mockDB.EXPECT().GetAllUsers().Return([]models.User{}, nil)
		store, _ := NewUserStore(mockDB)

		newUser := &models.User{Id: 1, Login: "newuser", Password: "newpass"}
		mockDB.EXPECT().AddUser(newUser).Return(nil)

		err := store.AddUser(newUser)
		if err != nil {
			t.Fatalf("не ожидали ошибку при добавлении пользователя: %v", err)
		}

		user := store.GetUserById(1)
		if user == nil {
			t.Fatal("пользователь должен быть добавлен в store")
		}

		if user.Login != "newuser" {
			t.Errorf("ожидали login 'newuser', получили '%s'", user.Login)
		}
	})

	t.Run("ошибка при добавлении пользователя в БД", func(t *testing.T) {
		mockDB.EXPECT().GetAllUsers().Return([]models.User{}, nil)
		store, _ := NewUserStore(mockDB)

		newUser := &models.User{Id: 2, Login: "erroruser", Password: "pass"}
		mockDB.EXPECT().AddUser(newUser).Return(errors.New("database error"))

		err := store.AddUser(newUser)
		if err == nil {
			t.Error("ожидали ошибку при добавлении пользователя")
		}

		user := store.GetUserById(2)
		if user != nil {
			t.Error("пользователь не должен быть добавлен в store при ошибке БД")
		}
	})
}

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockIDataBaseManager(ctrl)

	t.Run("успешное обновление пользователя", func(t *testing.T) {
		users := []models.User{
			{Id: 1, Login: "user1", Password: "pass1"},
		}
		mockDB.EXPECT().GetAllUsers().Return(users, nil)
		store, _ := NewUserStore(mockDB)

		updatedUser := &models.User{Id: 1, Login: "updated", Password: "newpass"}
		mockDB.EXPECT().UpdateUser(updatedUser).Return(nil)

		err := store.UpdateUser(updatedUser)
		if err != nil {
			t.Fatalf("не ожидали ошибку при обновлении пользователя: %v", err)
		}

		user := store.GetUserById(1)
		if user.Login != "updated" {
			t.Errorf("ожидали обновленный login 'updated', получили '%s'", user.Login)
		}
	})

	t.Run("ошибка при обновлении пользователя в БД", func(t *testing.T) {
		users := []models.User{
			{Id: 2, Login: "user2", Password: "pass2"},
		}
		mockDB.EXPECT().GetAllUsers().Return(users, nil)
		store, _ := NewUserStore(mockDB)

		updatedUser := &models.User{Id: 2, Login: "updated", Password: "newpass"}
		mockDB.EXPECT().UpdateUser(updatedUser).Return(errors.New("database error"))

		err := store.UpdateUser(updatedUser)
		if err == nil {
			t.Error("ожидали ошибку при обновлении пользователя")
		}

		user := store.GetUserById(2)
		if user.Login == "updated" {
			t.Error("пользователь не должен быть обновлен в store при ошибке БД")
		}
	})
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockIDataBaseManager(ctrl)

	t.Run("успешное удаление пользователя", func(t *testing.T) {
		users := []models.User{
			{Id: 1, Login: "user1", Password: "pass1"},
		}
		mockDB.EXPECT().GetAllUsers().Return(users, nil)
		store, _ := NewUserStore(mockDB)

		mockDB.EXPECT().DeleteUser(models.UserId(1)).Return(nil)

		err := store.DeleteUser(1)
		if err != nil {
			t.Fatalf("не ожидали ошибку при удалении пользователя: %v", err)
		}

		user := store.GetUserById(1)
		if user != nil {
			t.Error("пользователь должен быть удален из store")
		}
	})

	t.Run("ошибка при удалении пользователя из БД", func(t *testing.T) {
		users := []models.User{
			{Id: 2, Login: "user2", Password: "pass2"},
		}
		mockDB.EXPECT().GetAllUsers().Return(users, nil)
		store, _ := NewUserStore(mockDB)

		mockDB.EXPECT().DeleteUser(models.UserId(2)).Return(errors.New("database error"))

		err := store.DeleteUser(2)
		if err == nil {
			t.Error("ожидали ошибку при удалении пользователя")
		}

		user := store.GetUserById(2)
		if user == nil {
			t.Error("пользователь не должен быть удален из store при ошибке БД")
		}
	})
}

func TestGetAllUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockIDataBaseManager(ctrl)

	t.Run("получение всех пользователей", func(t *testing.T) {
		users := []models.User{
			{Id: 1, Login: "user1", Password: "pass1"},
			{Id: 2, Login: "user2", Password: "pass2"},
			{Id: 3, Login: "user3", Password: "pass3"},
		}
		mockDB.EXPECT().GetAllUsers().Return(users, nil)
		store, _ := NewUserStore(mockDB)

		allUsers := store.GetAllUsers()
		if len(allUsers) != 3 {
			t.Errorf("ожидали 3 пользователей, получили %d", len(allUsers))
		}
	})

	t.Run("получение пустого списка", func(t *testing.T) {
		mockDB.EXPECT().GetAllUsers().Return([]models.User{}, nil)
		store, _ := NewUserStore(mockDB)

		allUsers := store.GetAllUsers()
		if len(allUsers) != 0 {
			t.Errorf("ожидали пустой список, получили %d пользователей", len(allUsers))
		}
	})
}

func TestConcurrency(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockIDataBaseManager(ctrl)
	users := []models.User{
		{Id: 1, Login: "user1", Password: "pass1"},
	}
	mockDB.EXPECT().GetAllUsers().Return(users, nil)
	store, _ := NewUserStore(mockDB)

	t.Run("конкурентное чтение безопасно", func(t *testing.T) {
		done := make(chan bool)

		for i := 0; i < 10; i++ {
			go func() {
				user := store.GetUserById(1)
				if user == nil {
					t.Error("пользователь должен существовать")
				}
				done <- true
			}()
		}

		for i := 0; i < 10; i++ {
			<-done
		}
	})
}
