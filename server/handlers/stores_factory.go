package handlers

import "yuki_buy_log/stores"

// Типы-алиасы для интерфейсов stores
type (
	UserStoreInterface     = stores.UserStoreInterface
	ProductStoreInterface  = stores.ProductStoreInterface
	GroupStoreInterface    = stores.GroupStoreInterface
	InviteStoreInterface   = stores.InviteStoreInterface
	PurchaseStoreInterface = stores.PurchaseStoreInterface
)

// Переменные-фабрики для получения stores.
// По умолчанию используют глобальные singleton'ы,
// но в тестах можно подменить на моки.
var (
	GetUserStore     = func() UserStoreInterface { return stores.GetUserStore() }
	GetProductStore  = func() ProductStoreInterface { return stores.GetProductStore() }
	GetGroupStore    = func() GroupStoreInterface { return stores.GetGroupStore() }
	GetInviteStore   = func() InviteStoreInterface { return stores.GetInviteStore() }
	GetPurchaseStore = func() PurchaseStoreInterface { return stores.GetPurchaseStore() }
)
