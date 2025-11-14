package tests

import (
	"sync"
	"testing"
	"yuki_buy_log/models"
	"yuki_buy_log/stores"

	"github.com/stretchr/testify/assert"
)

// Эти тесты проверяют только публичный API UserStore
// Они не требуют подключения к базе данных и используют только экспортированные методы

func TestGetUserStore(t *testing.T) {
	t.Run("успешное создание store", func(t *testing.T) {
		store := stores.GetUserStore()
		assert.NotNil(t, store, "store не должен быть nil")
	})
}

func TestGetAllUsers(t *testing.T) {
	store := stores.GetUserStore()

	t.Run("GetAllUsers возвращает не nil слайс", func(t *testing.T) {
		allUsers := store.GetAllUsers()
		assert.NotNil(t, allUsers, "GetAllUsers должен возвращать не nil слайс")
	})

	t.Run("GetAllUsers возвращает слайс", func(t *testing.T) {
		allUsers := store.GetAllUsers()
		assert.IsType(t, []models.User{}, allUsers, "GetAllUsers должен возвращать слайс User")
	})
}

func TestGetUserById(t *testing.T) {
	store := stores.GetUserStore()
	allUsers := store.GetAllUsers()

	if len(allUsers) > 0 {
		existingUser := allUsers[0]

		t.Run("GetUserById возвращает пользователя для существующего ID", func(t *testing.T) {
			user := store.GetUserById(existingUser.Id)
			assert.NotNil(t, user, "должен вернуть пользователя для существующего ID")
			assert.Equal(t, existingUser.Id, user.Id, "ID должны совпадать")
		})

		t.Run("GetUserById возвращает копию (изменения не влияют на store)", func(t *testing.T) {
			user1 := store.GetUserById(existingUser.Id)
			if user1 != nil {
				originalLogin := user1.Login
				user1.Login = "modified_test_login"

				user2 := store.GetUserById(existingUser.Id)
				assert.Equal(t, originalLogin, user2.Login, "изменение возвращенного объекта не должно влиять на store")
			}
		})
	}

	t.Run("GetUserById возвращает nil для несуществующего ID", func(t *testing.T) {
		user := store.GetUserById(models.UserId(999999999))
		assert.Nil(t, user, "должен вернуть nil для несуществующего ID")
	})
}

func TestGetUserByLogin(t *testing.T) {
	store := stores.GetUserStore()
	allUsers := store.GetAllUsers()

	if len(allUsers) > 0 {
		existingUser := allUsers[0]

		t.Run("GetUserByLogin возвращает пользователя для существующего login", func(t *testing.T) {
			user := store.GetUserByLogin(existingUser.Login)
			assert.NotNil(t, user, "должен вернуть пользователя для существующего login")
			assert.Equal(t, existingUser.Login, user.Login, "login должны совпадать")
		})

		t.Run("GetUserByLogin возвращает копию (изменения не влияют на store)", func(t *testing.T) {
			user1 := store.GetUserByLogin(existingUser.Login)
			if user1 != nil {
				originalId := user1.Id
				user1.Id = models.UserId(999999999)

				user2 := store.GetUserByLogin(existingUser.Login)
				assert.Equal(t, originalId, user2.Id, "изменение возвращенного объекта не должно влиять на store")
			}
		})
	}

	t.Run("GetUserByLogin возвращает nil для несуществующего login", func(t *testing.T) {
		user := store.GetUserByLogin("nonexistent_user_login_999999")
		assert.Nil(t, user, "должен вернуть nil для несуществующего login")
	})
}

func TestConcurrentReads(t *testing.T) {
	store := stores.GetUserStore()
	allUsers := store.GetAllUsers()

	if len(allUsers) == 0 {
		t.Skip("нет данных для тестирования конкурентности")
	}

	existingUserId := allUsers[0].Id

	t.Run("конкурентное чтение GetUserById безопасно", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				user := store.GetUserById(existingUserId)
				// Просто проверяем, что не паникует
				_ = user
			}()
		}
		wg.Wait()
	})

	t.Run("конкурентное чтение GetAllUsers безопасно", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				users := store.GetAllUsers()
				// Просто проверяем, что не паникует
				_ = users
			}()
		}
		wg.Wait()
	})

	t.Run("конкурентное чтение GetUserByLogin безопасно", func(t *testing.T) {
		existingLogin := allUsers[0].Login
		var wg sync.WaitGroup
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				user := store.GetUserByLogin(existingLogin)
				// Просто проверяем, что не паникует
				_ = user
			}()
		}
		wg.Wait()
	})
}

func TestDataImmutability(t *testing.T) {
	store := stores.GetUserStore()
	allUsers := store.GetAllUsers()

	if len(allUsers) == 0 {
		t.Skip("нет данных для тестирования неизменяемости")
	}

	existingUser := allUsers[0]

	t.Run("изменение данных через GetAllUsers не влияет на store", func(t *testing.T) {
		users1 := store.GetAllUsers()
		if len(users1) > 0 {
			originalLogin := users1[0].Login
			users1[0].Login = "modified_through_slice"

			users2 := store.GetAllUsers()
			assert.Equal(t, originalLogin, users2[0].Login, "изменение в возвращенном слайсе не должно влиять на store")
		}
	})

	t.Run("множественные изменения возвращенного объекта не влияют на store", func(t *testing.T) {
		user1 := store.GetUserById(existingUser.Id)
		if user1 != nil {
			originalLogin := user1.Login
			originalPassword := user1.Password

			user1.Login = "test_modified_1"
			user1.Password = "test_modified_pass_1"
			user1.Id = models.UserId(888888888)

			user2 := store.GetUserById(existingUser.Id)
			assert.Equal(t, originalLogin, user2.Login, "login в store не должен измениться")
			assert.Equal(t, originalPassword, user2.Password, "password в store не должен измениться")
			assert.Equal(t, existingUser.Id, user2.Id, "ID в store не должен измениться")
		}
	})
}
