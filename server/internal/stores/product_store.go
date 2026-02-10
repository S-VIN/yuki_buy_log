package stores

import (
	"sync"
	"yuki_buy_log/internal/database"
	"yuki_buy_log/internal/domain"
)

type ProductStore struct {
	data  map[domain.ProductId]domain.Product
	mutex sync.RWMutex
	db    database.DatabaseManager
}

var (
	productStoreInstance *ProductStore
	productStoreLock     sync.Once
)

func GetProductStore() *ProductStore {
	productStoreLock.Do(func() {
		var db, _ = database.GetDBManager()
		products, err := db.GetAllProducts()
		if err != nil {
			products = []domain.Product{}
		}

		// Преобразуем список продуктов в map[ProductId]Product
		productMap := make(map[domain.ProductId]domain.Product)
		for _, product := range products {
			productMap[product.Id] = product
		}

		productStoreInstance = &ProductStore{
			data: productMap,
			db:   *db,
		}
	})
	return productStoreInstance
}

// GetProductById возвращает продукт по ID
func (s *ProductStore) GetProductById(id domain.ProductId) *domain.Product {
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
func (s *ProductStore) GetProductsByUserId(userId domain.UserId) []domain.Product {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var products []domain.Product
	for _, product := range s.data {
		if product.UserId == userId {
			products = append(products, product)
		}
	}
	return products
}

// GetProductsByUserIds возвращает продукты для списка пользователей (для группы)
func (s *ProductStore) GetProductsByUserIds(userIds []domain.UserId) []domain.Product {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Создаем map для быстрого поиска
	userIdMap := make(map[domain.UserId]bool)
	for _, userId := range userIds {
		userIdMap[userId] = true
	}

	var products []domain.Product
	for _, product := range s.data {
		if userIdMap[product.UserId] {
			products = append(products, product)
		}
	}
	return products
}

// CreateProduct создает новый продукт
func (s *ProductStore) CreateProduct(product *domain.Product) error {
	// Добавляем в БД
	err := s.db.CreateProduct(product)
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
func (s *ProductStore) UpdateProduct(product *domain.Product) error {
	// Обновляем в БД
	err := s.db.UpdateProduct(product)
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
func (s *ProductStore) DeleteProduct(id domain.ProductId, userId domain.UserId) error {
	// Удаляем из БД
	err := s.db.DeleteProduct(id, userId)
	if err != nil {
		return err
	}

	// Удаляем из локального стора
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.data, id)
	return nil
}
