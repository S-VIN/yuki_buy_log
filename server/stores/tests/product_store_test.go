package tests

import (
	"sync"
	"testing"
	"yuki_buy_log/models"
	"yuki_buy_log/stores"

	"github.com/stretchr/testify/assert"
)

// Тесты для ProductStore используют только публичный API

func TestGetProductStore(t *testing.T) {
	t.Run("успешное создание store", func(t *testing.T) {
		store := stores.GetProductStore()
		assert.NotNil(t, store, "store не должен быть nil")
	})
}

func TestGetProductById(t *testing.T) {
	store := stores.GetProductStore()

	t.Run("GetProductById возвращает nil для несуществующего ID", func(t *testing.T) {
		product := store.GetProductById(models.ProductId(999999999))
		assert.Nil(t, product, "должен вернуть nil для несуществующего ID")
	})

	t.Run("GetProductById возвращает копию (изменения не влияют на store)", func(t *testing.T) {
		// Получаем все продукты конкретного пользователя
		userStore := stores.GetUserStore()
		allUsers := userStore.GetAllUsers()

		if len(allUsers) == 0 {
			t.Skip("нет пользователей для тестирования")
		}

		products := store.GetProductsByUserId(allUsers[0].Id)
		if len(products) == 0 {
			t.Skip("нет продуктов для тестирования")
		}

		existingProduct := products[0]
		product1 := store.GetProductById(existingProduct.Id)
		if product1 != nil {
			originalName := product1.Name
			product1.Name = "modified_test_name"

			product2 := store.GetProductById(existingProduct.Id)
			assert.Equal(t, originalName, product2.Name, "изменение возвращенного объекта не должно влиять на store")
		}
	})
}

func TestGetProductsByUserId(t *testing.T) {
	store := stores.GetProductStore()
	userStore := stores.GetUserStore()
	allUsers := userStore.GetAllUsers()

	if len(allUsers) == 0 {
		t.Skip("нет пользователей для тестирования")
	}

	t.Run("GetProductsByUserId возвращает слайс", func(t *testing.T) {
		products := store.GetProductsByUserId(allUsers[0].Id)
		assert.NotNil(t, products, "должен вернуть не nil слайс")
	})

	t.Run("GetProductsByUserId возвращает пустой слайс для несуществующего userId", func(t *testing.T) {
		products := store.GetProductsByUserId(models.UserId(999999999))
		assert.NotNil(t, products, "должен вернуть не nil слайс")
		assert.Equal(t, 0, len(products), "должен вернуть пустой слайс")
	})

	t.Run("изменение возвращенного слайса не влияет на store", func(t *testing.T) {
		products1 := store.GetProductsByUserId(allUsers[0].Id)
		if len(products1) > 0 {
			originalName := products1[0].Name
			products1[0].Name = "modified_through_slice"

			product := store.GetProductById(products1[0].Id)
			if product != nil {
				assert.Equal(t, originalName, product.Name, "изменение в возвращенном слайсе не должно влиять на store")
			}
		}
	})
}

func TestGetProductsByUserIds(t *testing.T) {
	store := stores.GetProductStore()
	userStore := stores.GetUserStore()
	allUsers := userStore.GetAllUsers()

	if len(allUsers) == 0 {
		t.Skip("нет пользователей для тестирования")
	}

	t.Run("GetProductsByUserIds возвращает слайс", func(t *testing.T) {
		userIds := []models.UserId{allUsers[0].Id}
		if len(allUsers) > 1 {
			userIds = append(userIds, allUsers[1].Id)
		}

		products := store.GetProductsByUserIds(userIds)
		assert.NotNil(t, products, "должен вернуть не nil слайс")
	})

	t.Run("GetProductsByUserIds возвращает пустой слайс для несуществующих userIds", func(t *testing.T) {
		userIds := []models.UserId{models.UserId(999999999), models.UserId(999999998)}
		products := store.GetProductsByUserIds(userIds)
		assert.NotNil(t, products, "должен вернуть не nil слайс")
		assert.Equal(t, 0, len(products), "должен вернуть пустой слайс")
	})

	t.Run("GetProductsByUserIds возвращает пустой слайс для пустого списка userIds", func(t *testing.T) {
		userIds := []models.UserId{}
		products := store.GetProductsByUserIds(userIds)
		assert.NotNil(t, products, "должен вернуть не nil слайс")
	})
}

func TestProductStoreConcurrency(t *testing.T) {
	store := stores.GetProductStore()
	userStore := stores.GetUserStore()
	allUsers := userStore.GetAllUsers()

	if len(allUsers) == 0 {
		t.Skip("нет данных для тестирования конкурентности")
	}

	t.Run("конкурентное чтение GetProductsByUserId безопасно", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				products := store.GetProductsByUserId(allUsers[0].Id)
				_ = products
			}()
		}
		wg.Wait()
	})

	t.Run("конкурентное чтение GetProductsByUserIds безопасно", func(t *testing.T) {
		userIds := []models.UserId{allUsers[0].Id}
		var wg sync.WaitGroup
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				products := store.GetProductsByUserIds(userIds)
				_ = products
			}()
		}
		wg.Wait()
	})
}

func TestProductStoreDataImmutability(t *testing.T) {
	store := stores.GetProductStore()
	userStore := stores.GetUserStore()
	allUsers := userStore.GetAllUsers()

	if len(allUsers) == 0 {
		t.Skip("нет данных для тестирования неизменяемости")
	}

	products := store.GetProductsByUserId(allUsers[0].Id)
	if len(products) == 0 {
		t.Skip("нет продуктов для тестирования неизменяемости")
	}

	existingProduct := products[0]

	t.Run("множественные изменения возвращенного объекта не влияют на store", func(t *testing.T) {
		product1 := store.GetProductById(existingProduct.Id)
		if product1 != nil {
			originalName := product1.Name
			originalVolume := product1.Volume
			originalBrand := product1.Brand

			product1.Name = "test_modified_name"
			product1.Volume = "test_modified_volume"
			product1.Brand = "test_modified_brand"
			product1.Id = models.ProductId(888888888)

			product2 := store.GetProductById(existingProduct.Id)
			assert.Equal(t, originalName, product2.Name, "name в store не должен измениться")
			assert.Equal(t, originalVolume, product2.Volume, "volume в store не должен измениться")
			assert.Equal(t, originalBrand, product2.Brand, "brand в store не должен измениться")
			assert.Equal(t, existingProduct.Id, product2.Id, "ID в store не должен измениться")
		}
	})
}
