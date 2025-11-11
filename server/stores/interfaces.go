package stores

import (
	"time"
	"yuki_buy_log/models"
)

// UserStoreInterface определяет методы для работы с пользователями
type UserStoreInterface interface {
	GetUserById(id models.UserId) *models.User
	GetUserByLogin(login string) *models.User
	AddUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(userId models.UserId) error
	GetUsersByGroupId(groupId models.GroupId) []models.User
	GetAllUsers() []models.User
}

// ProductStoreInterface определяет методы для работы с продуктами
type ProductStoreInterface interface {
	GetProductById(id models.ProductId) *models.Product
	GetProductsByUserId(userId models.UserId) []models.Product
	GetProductsByUserIds(userIds []models.UserId) []models.Product
	CreateProduct(product *models.Product) error
	UpdateProduct(product *models.Product) error
	DeleteProduct(id models.ProductId, userId models.UserId) error
}

// GroupStoreInterface определяет методы для работы с группами
type GroupStoreInterface interface {
	GetGroupById(id models.GroupId) *models.Group
	GetGroupIdByUserId(userId models.UserId) *models.GroupId
	GetGroupByUserId(userId models.UserId) *models.Group
	GetGroupUserCount(groupId models.GroupId) int
	IsUserInGroup(userId models.UserId) bool
	CreateNewGroup(userId models.UserId) (*models.GroupId, error)
	AddUserToGroup(groupId models.GroupId, userId models.UserId) error
	DeleteUserFromGroup(userId models.UserId) error
	DeleteGroupById(id models.GroupId) error
}

// InviteStoreInterface определяет методы для работы с приглашениями
type InviteStoreInterface interface {
	GetInviteById(id models.InviteId) *models.Invite
	GetInvitesFromUser(fromUserId models.UserId) []models.Invite
	GetInvitesToUser(toUserId models.UserId) []models.Invite
	GetInvite(fromUserId, toUserId models.UserId) *models.Invite
	CreateInvite(invite *models.Invite) error
	DeleteInvite(id models.InviteId) error
	DeleteInviteByUsers(fromUserId, toUserId models.UserId) error
	DeleteOldInvites(cutoffTime time.Time) (int64, error)
}

// PurchaseStoreInterface определяет методы для работы с покупками
type PurchaseStoreInterface interface {
	GetPurchaseById(id models.PurchaseId) *models.Purchase
	GetPurchasesByUserId(userId models.UserId) []models.Purchase
	GetPurchasesByUserIds(userIds []models.UserId) []models.Purchase
	CreatePurchase(purchase *models.Purchase) error
	UpdatePurchase(purchase *models.Purchase) error
	DeletePurchase(id models.PurchaseId, userId models.UserId) error
}
