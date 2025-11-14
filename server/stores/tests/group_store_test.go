package tests

import (
	"sync"
	"testing"
	"yuki_buy_log/models"
	"yuki_buy_log/stores"

	"github.com/stretchr/testify/assert"
)

// Тесты для GroupStore используют только публичный API

func TestGetGroupStore(t *testing.T) {
	t.Run("успешное создание store", func(t *testing.T) {
		store := stores.GetGroupStore()
		assert.NotNil(t, store, "store не должен быть nil")
	})
}

func TestGetGroupById(t *testing.T) {
	store := stores.GetGroupStore()

	t.Run("GetGroupById возвращает nil для несуществующего ID", func(t *testing.T) {
		group := store.GetGroupById(models.GroupId(999999999))
		assert.Nil(t, group, "должен вернуть nil для несуществующего ID")
	})

	t.Run("GetGroupById возвращает копию (изменения не влияют на store)", func(t *testing.T) {
		userStore := stores.GetUserStore()
		allUsers := userStore.GetAllUsers()

		if len(allUsers) == 0 {
			t.Skip("нет пользователей для тестирования")
		}

		// Находим пользователя, который в группе
		var groupId *models.GroupId
		for _, user := range allUsers {
			gid := store.GetGroupIdByUserId(user.Id)
			if gid != nil {
				groupId = gid
				break
			}
		}

		if groupId == nil {
			t.Skip("нет групп для тестирования")
		}

		group1 := store.GetGroupById(*groupId)
		assert.NotNil(t, group1, "группа должна существовать")

		if group1 != nil {
			originalId := group1.Id
			group1.Id = models.GroupId(999999999)

			group2 := store.GetGroupById(*groupId)
			assert.Equal(t, originalId, group2.Id, "изменение возвращенного объекта не должно влиять на store")
		}
	})
}

func TestGetGroupIdByUserId(t *testing.T) {
	store := stores.GetGroupStore()
	userStore := stores.GetUserStore()
	allUsers := userStore.GetAllUsers()

	if len(allUsers) == 0 {
		t.Skip("нет пользователей для тестирования")
	}

	t.Run("GetGroupIdByUserId возвращает nil для пользователя не в группе", func(t *testing.T) {
		groupId := store.GetGroupIdByUserId(models.UserId(999999999))
		assert.Nil(t, groupId, "должен вернуть nil для пользователя не в группе")
	})

	t.Run("GetGroupIdByUserId возвращает groupId для пользователя в группе", func(t *testing.T) {
		// Находим пользователя в группе
		for _, user := range allUsers {
			groupId := store.GetGroupIdByUserId(user.Id)
			if groupId != nil {
				assert.NotNil(t, groupId, "groupId не должен быть nil")
				assert.NotEqual(t, models.GroupId(0), *groupId, "groupId не должен быть 0")

				// Проверяем, что группа действительно существует
				group := store.GetGroupById(*groupId)
				assert.NotNil(t, group, "группа с таким ID должна существовать")
				return
			}
		}
		t.Skip("нет пользователей в группах для тестирования")
	})
}

func TestGetGroupByUserId(t *testing.T) {
	store := stores.GetGroupStore()
	userStore := stores.GetUserStore()
	allUsers := userStore.GetAllUsers()

	if len(allUsers) == 0 {
		t.Skip("нет пользователей для тестирования")
	}

	t.Run("GetGroupByUserId возвращает nil для пользователя не в группе", func(t *testing.T) {
		group := store.GetGroupByUserId(models.UserId(999999999))
		assert.Nil(t, group, "должен вернуть nil для пользователя не в группе")
	})

	t.Run("GetGroupByUserId возвращает группу для пользователя в группе", func(t *testing.T) {
		// Находим пользователя в группе
		for _, user := range allUsers {
			group := store.GetGroupByUserId(user.Id)
			if group != nil {
				assert.NotNil(t, group, "группа не должна быть nil")
				assert.NotEqual(t, models.GroupId(0), group.Id, "groupId не должен быть 0")
				assert.NotNil(t, group.Members, "members не должен быть nil")

				// Проверяем, что пользователь действительно в группе
				found := false
				for _, member := range group.Members {
					if member.UserId == user.Id {
						found = true
						break
					}
				}
				assert.True(t, found, "пользователь должен быть в возвращенной группе")
				return
			}
		}
		t.Skip("нет пользователей в группах для тестирования")
	})
}

func TestGetGroupUserCount(t *testing.T) {
	store := stores.GetGroupStore()
	userStore := stores.GetUserStore()
	allUsers := userStore.GetAllUsers()

	if len(allUsers) == 0 {
		t.Skip("нет пользователей для тестирования")
	}

	t.Run("GetGroupUserCount возвращает 0 для несуществующей группы", func(t *testing.T) {
		count := store.GetGroupUserCount(models.GroupId(999999999))
		assert.Equal(t, 0, count, "должен вернуть 0 для несуществующей группы")
	})

	t.Run("GetGroupUserCount возвращает корректное количество участников", func(t *testing.T) {
		// Находим пользователя в группе
		for _, user := range allUsers {
			groupId := store.GetGroupIdByUserId(user.Id)
			if groupId != nil {
				count := store.GetGroupUserCount(*groupId)
				assert.Greater(t, count, 0, "количество участников должно быть больше 0")

				// Проверяем, что количество совпадает с длиной Members
				group := store.GetGroupById(*groupId)
				if group != nil {
					assert.Equal(t, len(group.Members), count, "количество должно совпадать с длиной Members")
				}
				return
			}
		}
		t.Skip("нет групп для тестирования")
	})
}

func TestIsUserInGroup(t *testing.T) {
	store := stores.GetGroupStore()
	userStore := stores.GetUserStore()
	allUsers := userStore.GetAllUsers()

	if len(allUsers) == 0 {
		t.Skip("нет пользователей для тестирования")
	}

	t.Run("IsUserInGroup возвращает false для пользователя не в группе", func(t *testing.T) {
		inGroup := store.IsUserInGroup(models.UserId(999999999))
		assert.False(t, inGroup, "должен вернуть false для пользователя не в группе")
	})

	t.Run("IsUserInGroup возвращает true для пользователя в группе", func(t *testing.T) {
		// Находим пользователя в группе
		for _, user := range allUsers {
			groupId := store.GetGroupIdByUserId(user.Id)
			if groupId != nil {
				inGroup := store.IsUserInGroup(user.Id)
				assert.True(t, inGroup, "должен вернуть true для пользователя в группе")
				return
			}
		}
		t.Skip("нет пользователей в группах для тестирования")
	})

	t.Run("IsUserInGroup возвращает false для пользователя без группы", func(t *testing.T) {
		// Находим пользователя не в группе
		for _, user := range allUsers {
			groupId := store.GetGroupIdByUserId(user.Id)
			if groupId == nil {
				inGroup := store.IsUserInGroup(user.Id)
				assert.False(t, inGroup, "должен вернуть false для пользователя без группы")
				return
			}
		}
		t.Skip("все пользователи в группах")
	})
}

func TestGroupStoreConcurrency(t *testing.T) {
	store := stores.GetGroupStore()
	userStore := stores.GetUserStore()
	allUsers := userStore.GetAllUsers()

	if len(allUsers) == 0 {
		t.Skip("нет данных для тестирования конкурентности")
	}

	// Находим пользователя в группе
	var testUserId models.UserId
	var testGroupId models.GroupId
	for _, user := range allUsers {
		groupId := store.GetGroupIdByUserId(user.Id)
		if groupId != nil {
			testUserId = user.Id
			testGroupId = *groupId
			break
		}
	}

	if testGroupId == 0 {
		t.Skip("нет групп для тестирования конкурентности")
	}

	t.Run("конкурентное чтение GetGroupById безопасно", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				group := store.GetGroupById(testGroupId)
				_ = group
			}()
		}
		wg.Wait()
	})

	t.Run("конкурентное чтение GetGroupIdByUserId безопасно", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				groupId := store.GetGroupIdByUserId(testUserId)
				_ = groupId
			}()
		}
		wg.Wait()
	})

	t.Run("конкурентное чтение GetGroupByUserId безопасно", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				group := store.GetGroupByUserId(testUserId)
				_ = group
			}()
		}
		wg.Wait()
	})

	t.Run("конкурентное чтение IsUserInGroup безопасно", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				inGroup := store.IsUserInGroup(testUserId)
				_ = inGroup
			}()
		}
		wg.Wait()
	})

	t.Run("конкурентное чтение GetGroupUserCount безопасно", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				count := store.GetGroupUserCount(testGroupId)
				_ = count
			}()
		}
		wg.Wait()
	})
}

func TestGroupStoreDataImmutability(t *testing.T) {
	store := stores.GetGroupStore()
	userStore := stores.GetUserStore()
	allUsers := userStore.GetAllUsers()

	if len(allUsers) == 0 {
		t.Skip("нет данных для тестирования неизменяемости")
	}

	// Находим пользователя в группе
	var testGroupId models.GroupId
	for _, user := range allUsers {
		groupId := store.GetGroupIdByUserId(user.Id)
		if groupId != nil {
			testGroupId = *groupId
			break
		}
	}

	if testGroupId == 0 {
		t.Skip("нет групп для тестирования неизменяемости")
	}

	t.Run("изменение возвращенной группы не влияет на store", func(t *testing.T) {
		group1 := store.GetGroupById(testGroupId)
		assert.NotNil(t, group1)

		originalId := group1.Id
		originalMembersCount := len(group1.Members)

		// Изменяем данные
		group1.Id = models.GroupId(999999999)
		if len(group1.Members) > 0 {
			group1.Members[0].UserId = models.UserId(888888888)
		}

		// Получаем свежие данные
		group2 := store.GetGroupById(testGroupId)
		assert.Equal(t, originalId, group2.Id, "ID группы в store не должен измениться")
		assert.Equal(t, originalMembersCount, len(group2.Members), "количество участников не должно измениться")

		if len(group2.Members) > 0 && originalMembersCount > 0 {
			assert.NotEqual(t, models.UserId(888888888), group2.Members[0].UserId, "UserId участника не должен измениться")
		}
	})

	t.Run("изменение Members в возвращенной группе не влияет на store", func(t *testing.T) {
		group1 := store.GetGroupById(testGroupId)
		assert.NotNil(t, group1)

		if len(group1.Members) == 0 {
			t.Skip("нет участников для тестирования")
		}

		originalUserId := group1.Members[0].UserId

		// Изменяем слайс Members
		group1.Members = append(group1.Members, models.GroupMember{
			GroupId:      testGroupId,
			UserId:       models.UserId(777777777),
			MemberNumber: 999,
		})

		// Получаем свежие данные
		group2 := store.GetGroupById(testGroupId)

		// Проверяем, что новый участник не добавился
		found := false
		for _, member := range group2.Members {
			if member.UserId == models.UserId(777777777) {
				found = true
				break
			}
		}
		assert.False(t, found, "новый участник не должен быть в store")

		// Проверяем, что первый участник не изменился
		assert.Equal(t, originalUserId, group2.Members[0].UserId, "UserId первого участника не должен измениться")
	})
}
