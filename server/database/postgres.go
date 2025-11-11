package database

import (
	"database/sql"
	"time"
	"yuki_buy_log/models"
)

// PostgresDB реализация интерфейса Database для PostgreSQL
type PostgresDB struct {
	db *sql.DB
}

// NewPostgresDB создает новый экземпляр PostgresDB
func NewPostgresDB() *PostgresDB {
	return &PostgresDB{db: db}
}

// User methods

func (p *PostgresDB) GetUserById(id *models.UserId) (*models.User, error) {
	return GetUserById(id)
}

func (p *PostgresDB) GetUserByLogin(login string) (*models.User, error) {
	return GetUserByLogin(login)
}

func (p *PostgresDB) GetUsersByGroupId(id *models.GroupId) ([]models.User, error) {
	return GetUsersByGroupId(id)
}

func (p *PostgresDB) AddUser(user *models.User) error {
	return AddUser(user)
}

func (p *PostgresDB) UpdateUser(user *models.User) error {
	return UpdateUser(user)
}

func (p *PostgresDB) DeleteUser(userId models.UserId) error {
	return DeleteUser(userId)
}

func (p *PostgresDB) GetAllUsers() ([]models.User, error) {
	return GetAllUsers()
}

// Product methods

func (p *PostgresDB) GetAllProducts() ([]models.Product, error) {
	return GetAllProducts()
}

func (p *PostgresDB) GetProductById(id models.ProductId) (*models.Product, error) {
	return GetProductById(id)
}

func (p *PostgresDB) CreateProduct(product *models.Product) error {
	return CreateProduct(product)
}

func (p *PostgresDB) UpdateProduct(product *models.Product) error {
	return UpdateProduct(product)
}

func (p *PostgresDB) DeleteProduct(id models.ProductId, userId models.UserId) error {
	return DeleteProduct(id, userId)
}

// Group methods

func (p *PostgresDB) GetAllGroupMembers() ([]models.GroupMember, error) {
	return GetAllGroupMembers()
}

func (p *PostgresDB) GetGroupMembersByGroupId(id models.GroupId) ([]models.GroupMember, error) {
	return GetGroupMembersByGroupId(id)
}

func (p *PostgresDB) DeleteGroupMembersByGroupId(id models.GroupId) error {
	return DeleteGroupMembersByGroupId(id)
}

func (p *PostgresDB) DeleteUserFromGroup(userId models.UserId) error {
	return DeleteUserFromGroup(userId)
}

func (p *PostgresDB) AddUserToGroup(groupId models.GroupId, userId models.UserId, memberNumber int) error {
	return AddUserToGroup(groupId, userId, memberNumber)
}

func (p *PostgresDB) UpdateGroupMember(groupMember *models.GroupMember) error {
	return UpdateGroupMember(groupMember)
}

func (p *PostgresDB) CreateNewGroup(userId models.UserId) (models.GroupId, error) {
	return CreateNewGroup(userId)
}

// Purchase methods

func (p *PostgresDB) GetAllPurchases() ([]models.Purchase, error) {
	return GetAllPurchases()
}

func (p *PostgresDB) GetPurchasesByUserIds(userIds []models.UserId) ([]models.Purchase, error) {
	return GetPurchasesByUserIds(userIds)
}

func (p *PostgresDB) AddPurchase(purchase *models.Purchase) error {
	return AddPurchase(purchase)
}

func (p *PostgresDB) UpdatePurchase(purchase *models.Purchase) error {
	return UpdatePurchase(purchase)
}

func (p *PostgresDB) DeletePurchase(purchaseId models.PurchaseId, userId models.UserId) error {
	return DeletePurchase(purchaseId, userId)
}

// Invite methods

func (p *PostgresDB) GetIncomingInvites(userId models.UserId) ([]models.Invite, error) {
	return GetIncomingInvites(userId)
}

func (p *PostgresDB) GetInvite(fromUserId, toUserId models.UserId) (models.Invite, error) {
	return GetInvite(fromUserId, toUserId)
}

func (p *PostgresDB) GetInvitesFromUser(fromUserId models.InviteId) ([]models.Invite, error) {
	return GetInvitesFromUser(fromUserId)
}

func (p *PostgresDB) GetInvitesToUser(toUserId models.UserId) ([]models.Invite, error) {
	return GetInvitesToUser(toUserId)
}

func (p *PostgresDB) DeleteInvitesBetweenUsers(userId1, userId2 models.UserId) error {
	return DeleteInvitesBetweenUsers(userId1, userId2)
}

func (p *PostgresDB) GetAllInvites() ([]models.Invite, error) {
	return GetAllInvites()
}

func (p *PostgresDB) CreateInvite(fromUserId, toUserId models.UserId) (models.InviteId, error) {
	return CreateInvite(fromUserId, toUserId)
}

func (p *PostgresDB) DeleteInvite(inviteId models.InviteId) error {
	return DeleteInvite(inviteId)
}

func (p *PostgresDB) DeleteOldInvites(cutoffTime time.Time) (int64, error) {
	return DeleteOldInvites(cutoffTime)
}

// Connection management

func (p *PostgresDB) Close() {
	Close()
}
