package tests

import (
	"sync"
	"testing"
	"time"
	"yuki_buy_log/models"
	"yuki_buy_log/stores"

	"github.com/stretchr/testify/assert"
)

// Тесты для InviteStore используют только публичный API

func TestGetInviteStore(t *testing.T) {
	t.Run("успешное создание store", func(t *testing.T) {
		store := stores.GetInviteStore()
		assert.NotNil(t, store, "store не должен быть nil")
	})
}

func TestGetInviteById(t *testing.T) {
	store := stores.GetInviteStore()

	t.Run("GetInviteById возвращает nil для несуществующего ID", func(t *testing.T) {
		invite := store.GetInviteById(models.InviteId(999999999))
		assert.Nil(t, invite, "должен вернуть nil для несуществующего ID")
	})

	t.Run("GetInviteById возвращает ссылку (не копию)", func(t *testing.T) {
		userStore := stores.GetUserStore()
		allUsers := userStore.GetAllUsers()

		if len(allUsers) < 2 {
			t.Skip("недостаточно пользователей для тестирования")
		}

		// Получаем все инвайты от первого пользователя
		invites := store.GetInvitesFromUser(allUsers[0].Id)
		if len(invites) == 0 {
			t.Skip("нет инвайтов для тестирования")
		}

		existingInvite := invites[0]
		invite := store.GetInviteById(existingInvite.Id)
		assert.NotNil(t, invite, "должен вернуть инвайт для существующего ID")
		if invite != nil {
			assert.Equal(t, existingInvite.Id, invite.Id, "ID должен совпадать")
		}
	})
}

func TestGetInvitesFromUser(t *testing.T) {
	store := stores.GetInviteStore()
	userStore := stores.GetUserStore()
	allUsers := userStore.GetAllUsers()

	if len(allUsers) == 0 {
		t.Skip("нет пользователей для тестирования")
	}

	t.Run("GetInvitesFromUser возвращает слайс", func(t *testing.T) {
		invites := store.GetInvitesFromUser(allUsers[0].Id)
		assert.NotNil(t, invites, "должен вернуть не nil слайс")
	})

	t.Run("GetInvitesFromUser возвращает пустой слайс для несуществующего userId", func(t *testing.T) {
		invites := store.GetInvitesFromUser(models.UserId(999999999))
		assert.NotNil(t, invites, "должен вернуть не nil слайс")
		assert.Equal(t, 0, len(invites), "должен вернуть пустой слайс")
	})

	t.Run("GetInvitesFromUser возвращает только инвайты от указанного пользователя", func(t *testing.T) {
		invites := store.GetInvitesFromUser(allUsers[0].Id)
		for _, invite := range invites {
			assert.Equal(t, allUsers[0].Id, invite.FromUserId, "все инвайты должны быть от указанного пользователя")
		}
	})
}

func TestGetInvitesToUser(t *testing.T) {
	store := stores.GetInviteStore()
	userStore := stores.GetUserStore()
	allUsers := userStore.GetAllUsers()

	if len(allUsers) == 0 {
		t.Skip("нет пользователей для тестирования")
	}

	t.Run("GetInvitesToUser возвращает слайс", func(t *testing.T) {
		invites := store.GetInvitesToUser(allUsers[0].Id)
		assert.NotNil(t, invites, "должен вернуть не nil слайс")
	})

	t.Run("GetInvitesToUser возвращает пустой слайс для несуществующего userId", func(t *testing.T) {
		invites := store.GetInvitesToUser(models.UserId(999999999))
		assert.NotNil(t, invites, "должен вернуть не nil слайс")
		assert.Equal(t, 0, len(invites), "должен вернуть пустой слайс")
	})

	t.Run("GetInvitesToUser возвращает только инвайты для указанного пользователя", func(t *testing.T) {
		invites := store.GetInvitesToUser(allUsers[0].Id)
		for _, invite := range invites {
			assert.Equal(t, allUsers[0].Id, invite.ToUserId, "все инвайты должны быть для указанного пользователя")
		}
	})
}

func TestGetInvite(t *testing.T) {
	store := stores.GetInviteStore()
	userStore := stores.GetUserStore()
	allUsers := userStore.GetAllUsers()

	if len(allUsers) < 2 {
		t.Skip("недостаточно пользователей для тестирования")
	}

	t.Run("GetInvite возвращает nil для несуществующей пары пользователей", func(t *testing.T) {
		invite := store.GetInvite(models.UserId(999999999), models.UserId(999999998))
		assert.Nil(t, invite, "должен вернуть nil для несуществующей пары")
	})

	t.Run("GetInvite возвращает инвайт для существующей пары", func(t *testing.T) {
		// Получаем все инвайты от первого пользователя
		invites := store.GetInvitesFromUser(allUsers[0].Id)
		if len(invites) == 0 {
			t.Skip("нет инвайтов для тестирования")
		}

		existingInvite := invites[0]
		invite := store.GetInvite(existingInvite.FromUserId, existingInvite.ToUserId)

		if invite != nil {
			assert.Equal(t, existingInvite.FromUserId, invite.FromUserId, "FromUserId должен совпадать")
			assert.Equal(t, existingInvite.ToUserId, invite.ToUserId, "ToUserId должен совпадать")
		}
	})
}

func TestAddInvite(t *testing.T) {
	store := stores.GetInviteStore()
	userStore := stores.GetUserStore()
	allUsers := userStore.GetAllUsers()

	if len(allUsers) < 2 {
		t.Skip("недостаточно пользователей для тестирования")
	}

	t.Run("AddInvite добавляет инвайт с автоинкрементным ID", func(t *testing.T) {
		// Создаем новый инвайт
		newInvite := models.Invite{
			FromUserId: allUsers[0].Id,
			ToUserId:   allUsers[1].Id,
			CreatedAt:  time.Now(),
		}

		// Сохраняем текущее количество инвайтов от пользователя
		invitesBefore := store.GetInvitesFromUser(allUsers[0].Id)
		countBefore := len(invitesBefore)

		// Добавляем инвайт
		err := store.AddInvite(newInvite)
		assert.NoError(t, err, "не должно быть ошибки при добавлении инвайта")

		// Проверяем, что инвайт добавлен
		invitesAfter := store.GetInvitesFromUser(allUsers[0].Id)
		assert.Equal(t, countBefore+1, len(invitesAfter), "количество инвайтов должно увеличиться на 1")

		// Проверяем, что у нового инвайта есть ID
		addedInvite := invitesAfter[len(invitesAfter)-1]
		assert.NotEqual(t, models.InviteId(0), addedInvite.Id, "ID нового инвайта не должен быть 0")
	})
}

func TestInviteStoreConcurrency(t *testing.T) {
	store := stores.GetInviteStore()
	userStore := stores.GetUserStore()
	allUsers := userStore.GetAllUsers()

	if len(allUsers) == 0 {
		t.Skip("нет данных для тестирования конкурентности")
	}

	t.Run("конкурентное чтение GetInvitesFromUser безопасно", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				invites := store.GetInvitesFromUser(allUsers[0].Id)
				_ = invites
			}()
		}
		wg.Wait()
	})

	t.Run("конкурентное чтение GetInvitesToUser безопасно", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				invites := store.GetInvitesToUser(allUsers[0].Id)
				_ = invites
			}()
		}
		wg.Wait()
	})

	t.Run("конкурентное чтение GetInviteById безопасно", func(t *testing.T) {
		invites := store.GetInvitesFromUser(allUsers[0].Id)
		if len(invites) == 0 {
			t.Skip("нет инвайтов для тестирования конкурентности")
		}

		inviteId := invites[0].Id
		var wg sync.WaitGroup
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				invite := store.GetInviteById(inviteId)
				_ = invite
			}()
		}
		wg.Wait()
	})
}

func TestInviteStoreDataImmutability(t *testing.T) {
	store := stores.GetInviteStore()
	userStore := stores.GetUserStore()
	allUsers := userStore.GetAllUsers()

	if len(allUsers) == 0 {
		t.Skip("нет данных для тестирования неизменяемости")
	}

	invites := store.GetInvitesFromUser(allUsers[0].Id)
	if len(invites) == 0 {
		t.Skip("нет инвайтов для тестирования неизменяемости")
	}

	t.Run("изменение возвращенного слайса не влияет на store", func(t *testing.T) {
		invites1 := store.GetInvitesFromUser(allUsers[0].Id)
		if len(invites1) > 0 {
			originalId := invites1[0].Id
			invites1[0].Id = models.InviteId(999999999)

			invites2 := store.GetInvitesFromUser(allUsers[0].Id)
			if len(invites2) > 0 {
				// Проверяем, что первый инвайт имеет оригинальный ID
				assert.Equal(t, originalId, invites2[0].Id, "ID инвайта в store не должен измениться")
			}
		}
	})
}
