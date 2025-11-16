package database

import (
	"database/sql"
	"log"
	"sync"
	"time"
	"yuki_buy_log/models"
	"yuki_buy_log/utils"
)

type IDataBaseManager interface {
	// GetAllGroupMembers returns all group members from the database
	GetAllGroupMembers() (result []models.GroupMember, err error)
	GetGroupMembersByGroupId(id models.GroupId) (result []models.GroupMember, err error)
	DeleteGroupMembersByGroupId(id models.GroupId) error
	DeleteUserFromGroup(userId models.UserId) (err error)
	AddUserToGroup(groupId models.GroupId, userId models.UserId, memberNumber int) error
	UpdateGroupMember(groupMember *models.GroupMember) error
	CreateNewGroup(userId models.UserId) (groupId models.GroupId, err error)

	// GetIncomingInvites returns all incoming invites for a user
	GetIncomingInvites(userId models.UserId) ([]models.Invite, error)
	GetInvite(fromUserId, toUserId models.UserId) (models.Invite, error)
	GetInvitesFromUser(fromUserId models.InviteId) ([]models.Invite, error)
	GetInvitesToUser(toUserId models.UserId) ([]models.Invite, error)
	DeleteInvitesBetweenUsers(userId1, userId2 models.UserId) error
	GetAllInvites() ([]models.Invite, error)
	CreateInvite(fromUserId, toUserId models.UserId) (models.InviteId, error)
	DeleteOldInvites(cutoffTime time.Time) (int64, error)

	// GetAllProducts returns all products from the database
	GetAllProducts() ([]models.Product, error)
	GetProductById(id models.ProductId) (*models.Product, error)
	CreateProduct(product *models.Product) error
	UpdateProduct(product *models.Product) error
	DeleteProduct(id models.ProductId, userId models.UserId) error

	// GetAllPurchases returns all purchases from the database
	GetAllPurchases() ([]models.Purchase, error)
	GetPurchasesByUserIds(userIds []models.UserId) ([]models.Purchase, error)
	AddPurchase(purchase *models.Purchase) error
	DeletePurchase(purchaseId models.PurchaseId, userId models.UserId) error
	GetUserById(id *models.UserId) (user *models.User, err error)
	GetUserByLogin(login string) (user *models.User, err error)
	GetUsersByGroupId(id *models.GroupId) (users []models.User, err error)
	AddUser(user *models.User) (err error)
	UpdateUser(user *models.User) error
	DeleteUser(userId models.UserId) error
	GetAllUsers() ([]models.User, error)
}

type DatabaseManager struct {
	db *sql.DB
}

var (
	instance *DatabaseManager
	once     sync.Once
)

func NewDatabaseManager() (*DatabaseManager, error) {
	once.Do(func() {
		instance = &DatabaseManager{}
		var err error
		log.Printf("Connecting to database with DSN: %s", utils.DatabaseURL)
		instance.db, err = sql.Open("postgres", utils.DatabaseURL)
		if err != nil {
			log.Fatalf("Failed to open database connection: %v", err)
		}

		log.Println("Waiting for database connection...")
		for i := 0; i < 30; i++ {
			if err = instance.db.Ping(); err == nil {
				log.Println("DatabaseManager connection established successfully")
				break
			}
			log.Printf("DatabaseManager connection attempt %d failed, retrying in 1 second...", i+1)
			time.Sleep(time.Second)
		}
		if err != nil {
			log.Fatalf("Failed to connect to database after 30 attempts: %v", err)
		}
	})
	return instance, nil
}

func (d *DatabaseManager) Close() {
	if d != nil && d.db != nil {
		d.db.Close()
	}
}
