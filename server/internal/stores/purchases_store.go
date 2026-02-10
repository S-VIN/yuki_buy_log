package stores

import (
	"sync"
	"yuki_buy_log/internal/database"
	"yuki_buy_log/internal/domain"
)

type PurchaseStore struct {
	data  map[domain.PurchaseId]domain.Purchase
	mutex sync.RWMutex
	db    database.DatabaseManager
}

var (
	purchaseStoreInstance *PurchaseStore
	purchaseStoreLock     sync.Once
)

func GetPurchaseStore() *PurchaseStore {
	purchaseStoreLock.Do(func() {
		var db, _ = database.GetDBManager()
		purchases, err := db.GetAllPurchases()
		if err != nil {
			purchases = []domain.Purchase{}
		}

		// Преобразуем список покупок в map[PurchaseId]Purchase
		purchaseMap := make(map[domain.PurchaseId]domain.Purchase)
		for _, purchase := range purchases {
			purchaseMap[purchase.Id] = purchase
		}

		purchaseStoreInstance = &PurchaseStore{
			data: purchaseMap,
			db:   *db,
		}
	})
	return purchaseStoreInstance
}

// GetPurchasesByUserIds возвращает покупки для списка пользователей (для группы)
func (s *PurchaseStore) GetPurchasesByUserIds(userIds []domain.UserId) []domain.Purchase {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Создаем map для быстрого поиска
	userIdMap := make(map[domain.UserId]bool)
	for _, userId := range userIds {
		userIdMap[userId] = true
	}

	var purchases []domain.Purchase
	for _, purchase := range s.data {
		if userIdMap[purchase.UserId] {
			purchases = append(purchases, purchase)
		}
	}
	return purchases
}

// AddPurchase добавляет новую покупку
func (s *PurchaseStore) AddPurchase(purchase *domain.Purchase) error {
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

// DeletePurchase удаляет покупку
func (s *PurchaseStore) DeletePurchase(purchaseId domain.PurchaseId, userId domain.UserId) error {
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
