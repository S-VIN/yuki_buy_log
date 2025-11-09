package stores

import (
	"sync"
	"yuki_buy_log/database"
	"yuki_buy_log/models"
)

type ProductStore struct {
	data  map[models.ProductId]models.Product
	mutex sync.RWMutex
}

var (
	productStoreInstance *ProductStore
	productStoreOnce     sync.Once
)

func GetProductStore() *ProductStore {
	productStoreOnce.Do(func() {
		products, err := database.GetAllProducts()
		if err != nil {
			products = []models.Product{}
		}

		// Преобразуем список продуктов в map[ProductId]Product
		productMap := make(map[models.ProductId]models.Product)
		for _, product := range products {
			productMap[product.Id] = product
		}

		productStoreInstance = &ProductStore{
			data: productMap,
		}
	})
	return productStoreInstance
}

// GetProductById возвращает продукт по ID
func (s *ProductStore) GetProductById(id models.ProductId) *models.Product {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if product, ok := s.data[id]; ok {
		// Возвращаем копию, чтобы избежать модификации извне
		productCopy := product
		return &productCopy
	}
	return nil
}

// GetProductsByUserId возвращает все продукты пользователя
func (s *ProductStore) GetProductsByUserId(userId models.UserId) []models.Product {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var products []models.Product
	for _, product := range s.data {
		if product.UserId == userId {
			products = append(products, product)
		}
	}
	return products
}

// GetProductsByUserIds возвращает продукты для списка пользователей (для группы)
func (s *ProductStore) GetProductsByUserIds(userIds []models.UserId) []models.Product {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Создаем map для быстрого поиска
	userIdMap := make(map[models.UserId]bool)
	for _, userId := range userIds {
		userIdMap[userId] = true
	}

	var products []models.Product
	for _, product := range s.data {
		if userIdMap[product.UserId] {
			products = append(products, product)
		}
	}
	return products
}

// CreateProduct создает новый продукт
func (s *ProductStore) CreateProduct(product *models.Product) error {
	// Добавляем в БД
	err := database.CreateProduct(product)
	if err != nil {
		return err
	}

	// Обновляем локальный стор
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[product.Id] = *product
	return nil
}

// UpdateProduct обновляет данные продукта
func (s *ProductStore) UpdateProduct(product *models.Product) error {
	// Обновляем в БД
	err := database.UpdateProduct(product)
	if err != nil {
		return err
	}

	// Обновляем локальный стор
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[product.Id] = *product
	return nil
}

// DeleteProduct удаляет продукт
func (s *ProductStore) DeleteProduct(id models.ProductId, userId models.UserId) error {
	// Удаляем из БД
	err := database.DeleteProduct(id, userId)
	if err != nil {
		return err
	}

	// Удаляем из локального стора
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.data, id)
	return nil
}