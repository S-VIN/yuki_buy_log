package handlers

import "yuki_buy_log/stores"

// Переменные-фабрики для получения stores.
// По умолчанию используют глобальные singleton'ы,
// но в тестах можно подменить на моки.
var (
	getUserStore     = func() stores.UserStoreInterface { return stores.GetUserStore() }
	getProductStore  = func() stores.ProductStoreInterface { return stores.GetProductStore() }
	getGroupStore    = func() stores.GroupStoreInterface { return stores.GetGroupStore() }
	getInviteStore   = func() stores.InviteStoreInterface { return stores.GetInviteStore() }
	getPurchaseStore = func() stores.PurchaseStoreInterface { return stores.GetPurchaseStore() }
)
