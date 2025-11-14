package tests

import (
	"sync"
	"testing"
	"yuki_buy_log/models"
	"yuki_buy_log/stores"

	"github.com/stretchr/testify/assert"
)

// Тесты для PurchaseStore используют только публичный API

func TestGetPurchaseStore(t *testing.T) {
	t.Run("успешное создание store", func(t *testing.T) {
		store := stores.GetPurchaseStore()
		assert.NotNil(t, store, "store не должен быть nil")
	})
}

func TestGetPurchasesByUserIds(t *testing.T) {
	store := stores.GetPurchaseStore()
	userStore := stores.GetUserStore()
	allUsers := userStore.GetAllUsers()

	if len(allUsers) == 0 {
		t.Skip("нет пользователей для тестирования")
	}

	t.Run("GetPurchasesByUserIds возвращает слайс", func(t *testing.T) {
		userIds := []models.UserId{allUsers[0].Id}
		if len(allUsers) > 1 {
			userIds = append(userIds, allUsers[1].Id)
		}

		purchases := store.GetPurchasesByUserIds(userIds)
		assert.NotNil(t, purchases, "должен вернуть не nil слайс")
	})

	t.Run("GetPurchasesByUserIds возвращает пустой слайс для несуществующих userIds", func(t *testing.T) {
		userIds := []models.UserId{models.UserId(999999999), models.UserId(999999998)}
		purchases := store.GetPurchasesByUserIds(userIds)
		assert.NotNil(t, purchases, "должен вернуть не nil слайс")
		assert.Equal(t, 0, len(purchases), "должен вернуть пустой слайс")
	})

	t.Run("GetPurchasesByUserIds возвращает пустой слайс для пустого списка userIds", func(t *testing.T) {
		userIds := []models.UserId{}
		purchases := store.GetPurchasesByUserIds(userIds)
		assert.NotNil(t, purchases, "должен вернуть не nil слайс")
	})

	t.Run("GetPurchasesByUserIds корректно фильтрует по нескольким userId", func(t *testing.T) {
		if len(allUsers) < 2 {
			t.Skip("недостаточно пользователей для тестирования фильтрации")
		}

		userIds := []models.UserId{allUsers[0].Id, allUsers[1].Id}
		purchases := store.GetPurchasesByUserIds(userIds)

		// Проверяем, что все возвращенные покупки принадлежат указанным пользователям
		for _, purchase := range purchases {
			found := false
			for _, userId := range userIds {
				if purchase.UserId == userId {
					found = true
					break
				}
			}
			assert.True(t, found, "покупка должна принадлежать одному из указанных пользователей")
		}
	})

	t.Run("изменение возвращенного слайса не влияет на store", func(t *testing.T) {
		userIds := []models.UserId{allUsers[0].Id}
		purchases1 := store.GetPurchasesByUserIds(userIds)

		if len(purchases1) > 0 {
			originalQuantity := purchases1[0].Quantity
			purchases1[0].Quantity = 99999

			purchases2 := store.GetPurchasesByUserIds(userIds)
			if len(purchases2) > 0 {
				// Находим ту же покупку
				for _, p := range purchases2 {
					if p.Id == purchases1[0].Id {
						assert.Equal(t, originalQuantity, p.Quantity, "изменение в возвращенном слайсе не должно влиять на store")
						break
					}
				}
			}
		}
	})
}

func TestPurchaseStoreConcurrency(t *testing.T) {
	store := stores.GetPurchaseStore()
	userStore := stores.GetUserStore()
	allUsers := userStore.GetAllUsers()

	if len(allUsers) == 0 {
		t.Skip("нет данных для тестирования конкурентности")
	}

	t.Run("конкурентное чтение GetPurchasesByUserIds безопасно", func(t *testing.T) {
		userIds := []models.UserId{allUsers[0].Id}
		var wg sync.WaitGroup
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				purchases := store.GetPurchasesByUserIds(userIds)
				_ = purchases
			}()
		}
		wg.Wait()
	})

	t.Run("конкурентное чтение с разными userIds безопасно", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				userId := allUsers[index%len(allUsers)].Id
				purchases := store.GetPurchasesByUserIds([]models.UserId{userId})
				_ = purchases
			}(i)
		}
		wg.Wait()
	})
}

func TestPurchaseStoreDataImmutability(t *testing.T) {
	store := stores.GetPurchaseStore()
	userStore := stores.GetUserStore()
	allUsers := userStore.GetAllUsers()

	if len(allUsers) == 0 {
		t.Skip("нет данных для тестирования неизменяемости")
	}

	userIds := []models.UserId{allUsers[0].Id}
	purchases := store.GetPurchasesByUserIds(userIds)

	if len(purchases) == 0 {
		t.Skip("нет покупок для тестирования неизменяемости")
	}

	existingPurchase := purchases[0]

	t.Run("множественные изменения возвращенного объекта не влияют на store", func(t *testing.T) {
		purchases1 := store.GetPurchasesByUserIds(userIds)

		var originalPurchase models.Purchase
		purchaseIndex := -1
		for i, p := range purchases1 {
			if p.Id == existingPurchase.Id {
				originalPurchase = p
				purchaseIndex = i
				break
			}
		}

		if purchaseIndex == -1 {
			t.Skip("не найдена покупка для тестирования")
		}

		// Изменяем данные
		purchases1[purchaseIndex].Quantity = 99999
		purchases1[purchaseIndex].Price = 88888
		purchases1[purchaseIndex].Store = "modified_store"
		purchases1[purchaseIndex].ProductId = models.ProductId(777777)

		// Получаем свежие данные
		purchases2 := store.GetPurchasesByUserIds(userIds)
		for _, p := range purchases2 {
			if p.Id == existingPurchase.Id {
				assert.Equal(t, originalPurchase.Quantity, p.Quantity, "quantity в store не должен измениться")
				assert.Equal(t, originalPurchase.Price, p.Price, "price в store не должен измениться")
				assert.Equal(t, originalPurchase.Store, p.Store, "store в store не должен измениться")
				assert.Equal(t, originalPurchase.ProductId, p.ProductId, "productId в store не должен измениться")
				break
			}
		}
	})
}
