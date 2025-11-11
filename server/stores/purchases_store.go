package stores

import (
	"sync"
	"yuki_buy_log/database"
	"yuki_buy_log/models"
)

type PurchaseStore struct {
	db    database.Database
	data  map[models.PurchaseId]models.Purchase
	mutex sync.RWMutex
}

var (
	purchaseStoreInstance *PurchaseStore
	purchaseStoreOnce     sync.Once
)

func GetPurchaseStore(db database.Database) *PurchaseStore {
	purchaseStoreOnce.Do(func() {
		purchases, err := db.GetAllPurchases()
		if err != nil {
			purchases = []models.Purchase{}
		}

		// Преобразуем список покупок в map[PurchaseId]Purchase
		purchaseMap := make(map[models.PurchaseId]models.Purchase)
		for _, purchase := range purchases {
			purchaseMap[purchase.Id] = purchase
		}

		purchaseStoreInstance = &PurchaseStore{
			db:   db,
			data: purchaseMap,
		}
	})
	return purchaseStoreInstance
}

// GetPurchasesByUserIds возвращает покупки для списка пользователей (для группы)
func (s *PurchaseStore) GetPurchasesByUserIds(userIds []models.UserId) []models.Purchase {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Создаем map для быстрого поиска
	userIdMap := make(map[models.UserId]bool)
	for _, userId := range userIds {
		userIdMap[userId] = true
	}

	var purchases []models.Purchase
	for _, purchase := range s.data {
		if userIdMap[purchase.UserId] {
			purchases = append(purchases, purchase)
		}
	}
	return purchases
}

// GetPurchaseById возвращает покупку по ID
func (s *PurchaseStore) GetPurchaseById(id models.PurchaseId) *models.Purchase {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if purchase, ok := s.data[id]; ok {
		purchaseCopy := purchase
		return &purchaseCopy
	}
	return nil
}

// GetPurchasesByUserId возвращает все покупки пользователя
func (s *PurchaseStore) GetPurchasesByUserId(userId models.UserId) []models.Purchase {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var purchases []models.Purchase
	for _, purchase := range s.data {
		if purchase.UserId == userId {
			purchases = append(purchases, purchase)
		}
	}
	return purchases
}

// CreatePurchase добавляет новую покупку
func (s *PurchaseStore) CreatePurchase(purchase *models.Purchase) error {
	// Добавляем в БД
	err := s.db.AddPurchase(purchase)
	if err != nil {
		return err
	}

	// Обновляем локальный стор
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[purchase.Id] = *purchase
	return nil
}

// UpdatePurchase обновляет покупку
func (s *PurchaseStore) UpdatePurchase(purchase *models.Purchase) error {
	// Обновляем в БД
	err := s.db.UpdatePurchase(purchase)
	if err != nil {
		return err
	}

	// Обновляем локальный стор
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[purchase.Id] = *purchase
	return nil
}

// DeletePurchase удаляет покупку
func (s *PurchaseStore) DeletePurchase(purchaseId models.PurchaseId, userId models.UserId) error {
	// Удаляем из БД
	err := s.db.DeletePurchase(purchaseId, userId)
	if err != nil {
		return err
	}

	// Удаляем из локального стора
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.data, purchaseId)
	return nil
}