package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"
	"yuki_buy_log/handlers"
	"yuki_buy_log/models"
)

// setupProductTest - вспомогательная функция для настройки тестового окружения
func setupProductTest() (*MockProductStore, *MockGroupStore, *MockUserStore, func()) {
	productStore := NewMockProductStore()
	groupStore := NewMockGroupStore()
	userStore := NewMockUserStore()

	// Сохраняем оригинальные фабрики
	origGetProductStore := handlers.GetProductStore
	origGetGroupStore := handlers.GetGroupStore
	origGetUserStore := handlers.GetUserStore

	// Подменяем фабрики на моки
	handlers.GetProductStore = func() handlers.ProductStoreInterface { return productStore }
	handlers.GetGroupStore = func() handlers.GroupStoreInterface { return groupStore }
	handlers.GetUserStore = func() handlers.UserStoreInterface { return userStore }

	// Возвращаем функцию для восстановления оригинальных фабрик
	cleanup := func() {
		handlers.GetProductStore = origGetProductStore
		handlers.GetGroupStore = origGetGroupStore
		handlers.GetUserStore = origGetUserStore
	}

	return productStore, groupStore, userStore, cleanup
}

// Тест: GET /products - успешное получение продуктов пользователя без группы
func TestGetProducts_UserNotInGroup_Success(t *testing.T) {
	productStore, _, userStore, cleanup := setupProductTest()
	defer cleanup()

	// Создаем тестового пользователя
	user := models.User{Id: 1, Login: "testuser"}
	userStore.AddUser(&user)

	// Создаем продукты для пользователя
	productStore.CreateProduct(&models.Product{Id: 1, Name: "Product1", Volume: "100ml", Brand: "Brand1", UserId: 1})
	productStore.CreateProduct(&models.Product{Id: 2, Name: "Product2", Volume: "200ml", Brand: "Brand2", UserId: 1})

	// Создаем HTTP запрос с контекстом пользователя
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	ctx := context.WithValue(req.Context(), "userId", models.UserId(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler := handlers.ProductsHandler(NewMockAuthenticator())
	handler.ServeHTTP(w, req)

	// Проверяем статус код
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Проверяем ответ
	var response map[string][]models.Product
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	products := response["products"]
	if len(products) != 2 {
		t.Errorf("Expected 2 products, got %d", len(products))
	}
}

// Тест: GET /products - успешное получение продуктов группы
func TestGetProducts_UserInGroup_Success(t *testing.T) {
	productStore, groupStore, userStore, cleanup := setupProductTest()
	defer cleanup()

	// Создаем пользователей
	user1 := models.User{Id: 1, Login: "user1"}
	user2 := models.User{Id: 2, Login: "user2"}
	userStore.AddUser(&user1)
	userStore.AddUser(&user2)

	// Создаем группу
	groupId, _ := groupStore.CreateNewGroup(1)
	groupStore.AddUserToGroup(*groupId, 2)

	// Создаем продукты для обоих пользователей
	productStore.CreateProduct(&models.Product{Id: 1, Name: "Product1", Volume: "100ml", Brand: "Brand1", UserId: 1})
	productStore.CreateProduct(&models.Product{Id: 2, Name: "Product2", Volume: "200ml", Brand: "Brand2", UserId: 2})

	// Создаем HTTP запрос
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	ctx := context.WithValue(req.Context(), "userId", models.UserId(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler := handlers.ProductsHandler(NewMockAuthenticator())
	handler.ServeHTTP(w, req)

	// Проверяем статус код
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Проверяем, что получены продукты обоих пользователей
	var response map[string][]models.Product
	json.NewDecoder(w.Body).Decode(&response)
	products := response["products"]

	if len(products) != 2 {
		t.Errorf("Expected 2 products from group, got %d", len(products))
	}
}

// Тест: GET /products - неавторизованный доступ
func TestGetProducts_Unauthorized(t *testing.T) {
	_, _, _, cleanup := setupProductTest()
	defer cleanup()

	// Создаем запрос без userId в контексте
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler := handlers.ProductsHandler(NewMockAuthenticator())
	handler.ServeHTTP(w, req)

	// Проверяем, что получили 401
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

// Тест: POST /products - успешное создание продукта
func TestCreateProduct_Success(t *testing.T) {
	productStore, _, userStore, cleanup := setupProductTest()
	defer cleanup()

	// Создаем пользователя
	user := models.User{Id: 1, Login: "testuser"}
	userStore.AddUser(&user)

	// Подготавливаем данные продукта
	newProduct := models.Product{
		Name:        "NewProduct",
		Volume:      "500ml",
		Brand:       "TestBrand",
		DefaultTags: []string{"tag1"},
	}

	body, _ := json.Marshal(newProduct)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "userId", models.UserId(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler := handlers.ProductsHandler(NewMockAuthenticator())
	handler.ServeHTTP(w, req)

	// Проверяем статус код
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Проверяем, что продукт был добавлен
	var createdProduct models.Product
	json.NewDecoder(w.Body).Decode(&createdProduct)

	if createdProduct.Name != "NewProduct" {
		t.Errorf("Expected product name 'NewProduct', got '%s'", createdProduct.Name)
	}

	if createdProduct.UserId != 1 {
		t.Errorf("Expected userId 1, got %d", createdProduct.UserId)
	}

	// Проверяем, что продукт в store
	if len(productStore.products) != 1 {
		t.Errorf("Expected 1 product in store, got %d", len(productStore.products))
	}
}

// Тест: POST /products - невалидные данные
func TestCreateProduct_InvalidData(t *testing.T) {
	_, _, userStore, cleanup := setupProductTest()
	defer cleanup()

	user := models.User{Id: 1, Login: "testuser"}
	userStore.AddUser(&user)

	// Продукт с невалидным именем (специальные символы)
	invalidProduct := models.Product{
		Name:   "Invalid@Name",
		Volume: "500ml",
		Brand:  "TestBrand",
	}

	body, _ := json.Marshal(invalidProduct)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "userId", models.UserId(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler := handlers.ProductsHandler(NewMockAuthenticator())
	handler.ServeHTTP(w, req)

	// Проверяем, что получили ошибку валидации
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// Тест: POST /products - неверный JSON
func TestCreateProduct_InvalidJSON(t *testing.T) {
	_, _, userStore, cleanup := setupProductTest()
	defer cleanup()

	user := models.User{Id: 1, Login: "testuser"}
	userStore.AddUser(&user)

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "userId", models.UserId(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler := handlers.ProductsHandler(NewMockAuthenticator())
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// Тест: POST /products - ошибка при сохранении
func TestCreateProduct_StoreError(t *testing.T) {
	productStore, _, userStore, cleanup := setupProductTest()
	defer cleanup()

	user := models.User{Id: 1, Login: "testuser"}
	userStore.AddUser(&user)

	// Настраиваем мок для возврата ошибки
	productStore.createProductFunc = func(p *models.Product) error {
		return errors.New("database error")
	}

	newProduct := models.Product{
		Name:   "Product",
		Volume: "500ml",
		Brand:  "Brand",
	}

	body, _ := json.Marshal(newProduct)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "userId", models.UserId(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler := handlers.ProductsHandler(NewMockAuthenticator())
	handler.ServeHTTP(w, req)

	// Проверяем, что получили ошибку сервера
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

// Тест: PUT /products - успешное обновление продукта
func TestUpdateProduct_Success(t *testing.T) {
	productStore, _, userStore, cleanup := setupProductTest()
	defer cleanup()

	user := models.User{Id: 1, Login: "testuser"}
	userStore.AddUser(&user)

	// Создаем продукт
	productStore.CreateProduct(&models.Product{
		Id:     1,
		Name:   "OldName",
		Volume: "100ml",
		Brand:  "OldBrand",
		UserId: 1,
	})

	// Обновляем продукт
	updatedProduct := models.Product{
		Id:     1,
		Name:   "NewName",
		Volume: "200ml",
		Brand:  "NewBrand",
	}

	body, _ := json.Marshal(updatedProduct)
	req := httptest.NewRequest(http.MethodPut, "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "userId", models.UserId(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler := handlers.ProductsHandler(NewMockAuthenticator())
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Проверяем, что продукт обновлен
	updated := productStore.GetProductById(1)
	if updated.Name != "NewName" {
		t.Errorf("Expected product name 'NewName', got '%s'", updated.Name)
	}
}

// Тест: PUT /products - отсутствует ID
func TestUpdateProduct_MissingId(t *testing.T) {
	_, _, userStore, cleanup := setupProductTest()
	defer cleanup()

	user := models.User{Id: 1, Login: "testuser"}
	userStore.AddUser(&user)

	// Продукт без ID
	product := models.Product{
		Name:   "Product",
		Volume: "100ml",
		Brand:  "Brand",
	}

	body, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPut, "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "userId", models.UserId(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler := handlers.ProductsHandler(NewMockAuthenticator())
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// Тест: метод не поддерживается
func TestProductsHandler_MethodNotAllowed(t *testing.T) {
	_, _, _, cleanup := setupProductTest()
	defer cleanup()

	req := httptest.NewRequest(http.MethodDelete, "/products", nil)
	w := httptest.NewRecorder()

	handler := handlers.ProductsHandler(NewMockAuthenticator())
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

// ===================
// LOGIN HANDLER TESTS
// ===================

// setupLoginTest - вспомогательная функция для настройки тестового окружения для login/register тестов
func setupLoginTest() (*MockUserStore, *MockAuthenticator, func()) {
	userStore := NewMockUserStore()
	auth := NewMockAuthenticator()

	origGetUserStore := handlers.GetUserStore
	handlers.GetUserStore = func() handlers.UserStoreInterface { return userStore }

	cleanup := func() {
		handlers.GetUserStore = origGetUserStore
	}

	return userStore, auth, cleanup
}

// Тест: POST /register - успешная регистрация нового пользователя
func TestRegister_Success(t *testing.T) {
	userStore, auth, cleanup := setupLoginTest()
	defer cleanup()

	newUser := models.User{Login: "newuser", Password: "password123"}
	body, _ := json.Marshal(newUser)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler := handlers.RegisterHandler(auth)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	if response["token"] == "" {
		t.Error("Expected token in response")
	}

	// Проверяем, что пользователь был добавлен в store
	if len(userStore.users) != 1 {
		t.Errorf("Expected 1 user in store, got %d", len(userStore.users))
	}
}

// Тест: POST /login - успешный вход
func TestLogin_Success(t *testing.T) {
	userStore, auth, cleanup := setupLoginTest()
	defer cleanup()

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := models.User{Id: 1, Login: "testuser", Password: string(hash)}
	userStore.AddUser(&user)

	credentials := models.User{Login: "testuser", Password: "password123"}
	body, _ := json.Marshal(credentials)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler := handlers.LoginHandler(auth)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// Тест: POST /login - пользователь не найден
func TestLogin_UserNotFound(t *testing.T) {
	_, auth, cleanup := setupLoginTest()
	defer cleanup()

	credentials := models.User{Login: "nonexistent", Password: "password123"}
	body, _ := json.Marshal(credentials)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler := handlers.LoginHandler(auth)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

// Тест: POST /login - неверный пароль
func TestLogin_InvalidPassword(t *testing.T) {
	userStore, auth, cleanup := setupLoginTest()
	defer cleanup()

	hash, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	user := models.User{Id: 1, Login: "testuser", Password: string(hash)}
	userStore.AddUser(&user)

	credentials := models.User{Login: "testuser", Password: "wrongpassword"}
	body, _ := json.Marshal(credentials)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler := handlers.LoginHandler(auth)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}
