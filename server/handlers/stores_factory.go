package handlers

import (
	"yuki_buy_log/database"
	"yuki_buy_log/stores"
)

// Типы-алиасы для интерфейсов stores
type (
	UserStoreInterface     = stores.UserStoreInterface
	ProductStoreInterface  = stores.ProductStoreInterface
	GroupStoreInterface    = stores.GroupStoreInterface
	InviteStoreInterface   = stores.InviteStoreInterface
	PurchaseStoreInterface = stores.PurchaseStoreInterface
)

// Глобальный экземпляр database для использования в stores
var DB database.Database

// Переменные-фабрики для получения stores.
// По умолчанию используют глобальные singleton'ы,
// но в тестах можно подменить на моки.
var (
	GetUserStore     = func() UserStoreInterface { return stores.GetUserStore(DB) }
	GetProductStore  = func() ProductStoreInterface { return stores.GetProductStore(DB) }
	GetGroupStore    = func() GroupStoreInterface { return stores.GetGroupStore(DB) }
	GetInviteStore   = func() InviteStoreInterface { return stores.GetInviteStore(DB) }
	GetPurchaseStore = func() PurchaseStoreInterface { return stores.GetPurchaseStore(DB) }
)
