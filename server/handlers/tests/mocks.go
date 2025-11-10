package handlers_test

import (
	"time"
	"yuki_buy_log/models"
)

// MockUserStore - мок для UserStoreInterface
type MockUserStore struct {
	users              map[models.UserId]models.User
	addUserFunc        func(*models.User) error
	updateUserFunc     func(*models.User) error
	deleteUserFunc     func(models.UserId) error
}

func NewMockUserStore() *MockUserStore {
	return &MockUserStore{
		users: make(map[models.UserId]models.User),
	}
}

func (m *MockUserStore) GetUserById(id models.UserId) *models.User {
	if user, ok := m.users[id]; ok {
		return &user
	}
	return nil
}

func (m *MockUserStore) GetUserByLogin(login string) *models.User {
	for _, user := range m.users {
		if user.Login == login {
			return &user
		}
	}
	return nil
}

func (m *MockUserStore) AddUser(user *models.User) error {
	if m.addUserFunc != nil {
		return m.addUserFunc(user)
	}
	// Генерируем ID если не задан
	if user.Id == 0 {
		user.Id = models.UserId(len(m.users) + 1)
	}
	m.users[user.Id] = *user
	return nil
}

func (m *MockUserStore) UpdateUser(user *models.User) error {
	if m.updateUserFunc != nil {
		return m.updateUserFunc(user)
	}
	m.users[user.Id] = *user
	return nil
}

func (m *MockUserStore) DeleteUser(userId models.UserId) error {
	if m.deleteUserFunc != nil {
		return m.deleteUserFunc(userId)
	}
	delete(m.users, userId)
	return nil
}

func (m *MockUserStore) GetUsersByGroupId(groupId models.GroupId) []models.User {
	// Для простоты не реализуем в моке
	return nil
}

func (m *MockUserStore) GetAllUsers() []models.User {
	users := make([]models.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users
}

// MockProductStore - мок для ProductStoreInterface
type MockProductStore struct {
	products          map[models.ProductId]models.Product
	createProductFunc func(*models.Product) error
	updateProductFunc func(*models.Product) error
	deleteProductFunc func(models.ProductId, models.UserId) error
}

func NewMockProductStore() *MockProductStore {
	return &MockProductStore{
		products: make(map[models.ProductId]models.Product),
	}
}

func (m *MockProductStore) GetProductById(id models.ProductId) *models.Product {
	if product, ok := m.products[id]; ok {
		return &product
	}
	return nil
}

func (m *MockProductStore) GetProductsByUserId(userId models.UserId) []models.Product {
	var products []models.Product
	for _, product := range m.products {
		if product.UserId == userId {
			products = append(products, product)
		}
	}
	return products
}

func (m *MockProductStore) GetProductsByUserIds(userIds []models.UserId) []models.Product {
	userIdMap := make(map[models.UserId]bool)
	for _, userId := range userIds {
		userIdMap[userId] = true
	}

	var products []models.Product
	for _, product := range m.products {
		if userIdMap[product.UserId] {
			products = append(products, product)
		}
	}
	return products
}

func (m *MockProductStore) CreateProduct(product *models.Product) error {
	if m.createProductFunc != nil {
		return m.createProductFunc(product)
	}
	if product.Id == 0 {
		product.Id = models.ProductId(len(m.products) + 1)
	}
	m.products[product.Id] = *product
	return nil
}

func (m *MockProductStore) UpdateProduct(product *models.Product) error {
	if m.updateProductFunc != nil {
		return m.updateProductFunc(product)
	}
	m.products[product.Id] = *product
	return nil
}

func (m *MockProductStore) DeleteProduct(id models.ProductId, userId models.UserId) error {
	if m.deleteProductFunc != nil {
		return m.deleteProductFunc(id, userId)
	}
	if product, ok := m.products[id]; ok && product.UserId == userId {
		delete(m.products, id)
		return nil
	}
	return nil
}

// MockGroupStore - мок для GroupStoreInterface
type MockGroupStore struct {
	groups             map[models.GroupId]models.Group
	userGroups         map[models.UserId]models.GroupId
	createNewGroupFunc func(models.UserId) (*models.GroupId, error)
	addUserToGroupFunc func(models.GroupId, models.UserId) error
}

func NewMockGroupStore() *MockGroupStore {
	return &MockGroupStore{
		groups:     make(map[models.GroupId]models.Group),
		userGroups: make(map[models.UserId]models.GroupId),
	}
}

func (m *MockGroupStore) GetGroupById(id models.GroupId) *models.Group {
	if group, ok := m.groups[id]; ok {
		return &group
	}
	return nil
}

func (m *MockGroupStore) GetGroupIdByUserId(userId models.UserId) *models.GroupId {
	if groupId, ok := m.userGroups[userId]; ok {
		return &groupId
	}
	return nil
}

func (m *MockGroupStore) GetGroupByUserId(userId models.UserId) *models.Group {
	if groupId, ok := m.userGroups[userId]; ok {
		return m.GetGroupById(groupId)
	}
	return nil
}

func (m *MockGroupStore) GetGroupUserCount(groupId models.GroupId) int {
	if group, ok := m.groups[groupId]; ok {
		return len(group.Members)
	}
	return 0
}

func (m *MockGroupStore) IsUserInGroup(userId models.UserId) bool {
	_, ok := m.userGroups[userId]
	return ok
}

func (m *MockGroupStore) CreateNewGroup(userId models.UserId) (*models.GroupId, error) {
	if m.createNewGroupFunc != nil {
		return m.createNewGroupFunc(userId)
	}
	groupId := models.GroupId(len(m.groups) + 1)
	m.groups[groupId] = models.Group{
		Id: groupId,
		Members: []models.GroupMember{
			{GroupId: groupId, UserId: userId, MemberNumber: 0},
		},
	}
	m.userGroups[userId] = groupId
	return &groupId, nil
}

func (m *MockGroupStore) AddUserToGroup(groupId models.GroupId, userId models.UserId) error {
	if m.addUserToGroupFunc != nil {
		return m.addUserToGroupFunc(groupId, userId)
	}
	if group, ok := m.groups[groupId]; ok {
		group.Members = append(group.Members, models.GroupMember{
			GroupId:      groupId,
			UserId:       userId,
			MemberNumber: len(group.Members),
		})
		m.groups[groupId] = group
		m.userGroups[userId] = groupId
	}
	return nil
}

func (m *MockGroupStore) DeleteUserFromGroup(userId models.UserId) error {
	if groupId, ok := m.userGroups[userId]; ok {
		if group, ok := m.groups[groupId]; ok {
			var newMembers []models.GroupMember
			for _, member := range group.Members {
				if member.UserId != userId {
					newMembers = append(newMembers, member)
				}
			}
			group.Members = newMembers
			m.groups[groupId] = group
		}
		delete(m.userGroups, userId)
	}
	return nil
}

func (m *MockGroupStore) DeleteGroupById(id models.GroupId) error {
	if group, ok := m.groups[id]; ok {
		for _, member := range group.Members {
			delete(m.userGroups, member.UserId)
		}
		delete(m.groups, id)
	}
	return nil
}

// MockInviteStore - мок для InviteStoreInterface
type MockInviteStore struct {
	invites           []models.Invite
	createInviteFunc  func(*models.Invite) error
	deleteInviteFunc  func(models.InviteId) error
}

func NewMockInviteStore() *MockInviteStore {
	return &MockInviteStore{
		invites: make([]models.Invite, 0),
	}
}

func (m *MockInviteStore) GetInviteById(id models.InviteId) *models.Invite {
	for _, invite := range m.invites {
		if invite.Id == id {
			return &invite
		}
	}
	return nil
}

func (m *MockInviteStore) GetInvitesFromUser(fromUserId models.UserId) []models.Invite {
	var result []models.Invite
	for _, invite := range m.invites {
		if invite.FromUserId == fromUserId {
			result = append(result, invite)
		}
	}
	return result
}

func (m *MockInviteStore) GetInvitesToUser(toUserId models.UserId) []models.Invite {
	var result []models.Invite
	for _, invite := range m.invites {
		if invite.ToUserId == toUserId {
			result = append(result, invite)
		}
	}
	return result
}

func (m *MockInviteStore) GetInvite(fromUserId, toUserId models.UserId) *models.Invite {
	for _, invite := range m.invites {
		if invite.FromUserId == fromUserId && invite.ToUserId == toUserId {
			return &invite
		}
	}
	return nil
}

func (m *MockInviteStore) CreateInvite(invite *models.Invite) error {
	if m.createInviteFunc != nil {
		return m.createInviteFunc(invite)
	}
	if invite.Id == 0 {
		invite.Id = models.InviteId(len(m.invites) + 1)
	}
	m.invites = append(m.invites, *invite)
	return nil
}

func (m *MockInviteStore) DeleteInvite(id models.InviteId) error {
	if m.deleteInviteFunc != nil {
		return m.deleteInviteFunc(id)
	}
	var newInvites []models.Invite
	for _, invite := range m.invites {
		if invite.Id != id {
			newInvites = append(newInvites, invite)
		}
	}
	m.invites = newInvites
	return nil
}

func (m *MockInviteStore) DeleteInviteByUsers(fromUserId, toUserId models.UserId) error {
	var newInvites []models.Invite
	for _, invite := range m.invites {
		if !((invite.FromUserId == fromUserId && invite.ToUserId == toUserId) ||
			(invite.FromUserId == toUserId && invite.ToUserId == fromUserId)) {
			newInvites = append(newInvites, invite)
		}
	}
	m.invites = newInvites
	return nil
}

func (m *MockInviteStore) DeleteOldInvites(cutoffTime time.Time) (int64, error) {
	count := int64(0)
	var newInvites []models.Invite
	for _, invite := range m.invites {
		if invite.CreatedAt.Before(cutoffTime) {
			count++
		} else {
			newInvites = append(newInvites, invite)
		}
	}
	m.invites = newInvites
	return count, nil
}

// MockPurchaseStore - мок для PurchaseStoreInterface
type MockPurchaseStore struct {
	purchases           map[models.PurchaseId]models.Purchase
	createPurchaseFunc  func(*models.Purchase) error
	updatePurchaseFunc  func(*models.Purchase) error
	deletePurchaseFunc  func(models.PurchaseId, models.UserId) error
}

func NewMockPurchaseStore() *MockPurchaseStore {
	return &MockPurchaseStore{
		purchases: make(map[models.PurchaseId]models.Purchase),
	}
}

func (m *MockPurchaseStore) GetPurchaseById(id models.PurchaseId) *models.Purchase {
	if purchase, ok := m.purchases[id]; ok {
		return &purchase
	}
	return nil
}

func (m *MockPurchaseStore) GetPurchasesByUserId(userId models.UserId) []models.Purchase {
	var purchases []models.Purchase
	for _, purchase := range m.purchases {
		if purchase.UserId == userId {
			purchases = append(purchases, purchase)
		}
	}
	return purchases
}

func (m *MockPurchaseStore) GetPurchasesByUserIds(userIds []models.UserId) []models.Purchase {
	userIdMap := make(map[models.UserId]bool)
	for _, userId := range userIds {
		userIdMap[userId] = true
	}

	var purchases []models.Purchase
	for _, purchase := range m.purchases {
		if userIdMap[purchase.UserId] {
			purchases = append(purchases, purchase)
		}
	}
	return purchases
}

func (m *MockPurchaseStore) CreatePurchase(purchase *models.Purchase) error {
	if m.createPurchaseFunc != nil {
		return m.createPurchaseFunc(purchase)
	}
	if purchase.Id == 0 {
		purchase.Id = models.PurchaseId(len(m.purchases) + 1)
	}
	m.purchases[purchase.Id] = *purchase
	return nil
}

func (m *MockPurchaseStore) UpdatePurchase(purchase *models.Purchase) error {
	if m.updatePurchaseFunc != nil {
		return m.updatePurchaseFunc(purchase)
	}
	m.purchases[purchase.Id] = *purchase
	return nil
}

func (m *MockPurchaseStore) DeletePurchase(id models.PurchaseId, userId models.UserId) error {
	if m.deletePurchaseFunc != nil {
		return m.deletePurchaseFunc(id, userId)
	}
	if purchase, ok := m.purchases[id]; ok && purchase.UserId == userId {
		delete(m.purchases, id)
		return nil
	}
	return nil
}

// MockAuthenticator - мок для Authenticator interface
type MockAuthenticator struct {
	generateTokenFunc func(models.UserId) (string, error)
}

func NewMockAuthenticator() *MockAuthenticator {
	return &MockAuthenticator{}
}

func (m *MockAuthenticator) GenerateToken(userId models.UserId) (string, error) {
	if m.generateTokenFunc != nil {
		return m.generateTokenFunc(userId)
	}
	return "mock-token-123", nil
}
