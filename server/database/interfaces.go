package database

import (
	"time"
	"yuki_buy_log/models"
)

// Database интерфейс для работы с базой данных
// Позволяет мокать функции database в тестах с помощью testify/mock
type Database interface {
	// User methods
	GetUserById(id *models.UserId) (*models.User, error)
	GetUserByLogin(login string) (*models.User, error)
	GetUsersByGroupId(id *models.GroupId) ([]models.User, error)
	AddUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(userId models.UserId) error
	GetAllUsers() ([]models.User, error)

	// Product methods
	GetAllProducts() ([]models.Product, error)
	GetProductById(id models.ProductId) (*models.Product, error)
	CreateProduct(product *models.Product) error
	UpdateProduct(product *models.Product) error
	DeleteProduct(id models.ProductId, userId models.UserId) error

	// Group methods
	GetAllGroupMembers() ([]models.GroupMember, error)
	GetGroupMembersByGroupId(id models.GroupId) ([]models.GroupMember, error)
	DeleteGroupMembersByGroupId(id models.GroupId) error
	DeleteUserFromGroup(userId models.UserId) error
	AddUserToGroup(groupId models.GroupId, userId models.UserId, memberNumber int) error
	UpdateGroupMember(groupMember *models.GroupMember) error
	CreateNewGroup(userId models.UserId) (models.GroupId, error)

	// Purchase methods
	GetAllPurchases() ([]models.Purchase, error)
	GetPurchasesByUserIds(userIds []models.UserId) ([]models.Purchase, error)
	AddPurchase(purchase *models.Purchase) error
	UpdatePurchase(purchase *models.Purchase) error
	DeletePurchase(purchaseId models.PurchaseId, userId models.UserId) error

	// Invite methods
	GetIncomingInvites(userId models.UserId) ([]models.Invite, error)
	GetInvite(fromUserId, toUserId models.UserId) (models.Invite, error)
	GetInvitesFromUser(fromUserId models.InviteId) ([]models.Invite, error)
	GetInvitesToUser(toUserId models.UserId) ([]models.Invite, error)
	DeleteInvitesBetweenUsers(userId1, userId2 models.UserId) error
	GetAllInvites() ([]models.Invite, error)
	CreateInvite(fromUserId, toUserId models.UserId) (models.InviteId, error)
	DeleteInvite(inviteId models.InviteId) error
	DeleteOldInvites(cutoffTime time.Time) (int64, error)

	// Connection management
	Close()
}
