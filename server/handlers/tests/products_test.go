package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"yuki_buy_log/handlers"
	"yuki_buy_log/models"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestProductsHandler_GET(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)

	// Mock getUser call
	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).AddRow(userID, "user1", "hash1"))

	// Mock group query (user not in a group)
	mock.ExpectQuery("SELECT DISTINCT user_id FROM groups WHERE id =").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}))

	// Mock database query for products
	rows := sqlmock.NewRows([]string{"id", "name", "volume", "brand", "default_tags", "user_id"}).
		AddRow(1, "TestProduct", "500ml", "TestBrand", "tag1,tag2", 1)

	mock.ExpectQuery("SELECT id, name, volume, brand, default_tags, user_id FROM products WHERE user_id=\\$1").
		WithArgs(userID).
		WillReturnRows(rows)

	handler := handlers.ProductsHandler(deps)

	req := httptest.NewRequest("GET", "/products", nil)
	ctx := context.WithValue(req.Context(), UserIDKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database expectations were not met: %v", err)
	}
}

func TestProductsHandler_POST(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)

	// Mock getUser call
	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).AddRow(userID, "user1", "hash1"))

	// Mock successful product insertion
	mock.ExpectQuery("INSERT INTO products").
		WithArgs("TestProduct", "500ml", "TestBrand", "test", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	handler := handlers.ProductsHandler(deps)

	product := models.Product{
		Name:        "TestProduct",
		Volume:      "500ml",
		Brand:       "TestBrand",
		DefaultTags: []string{"test"},
	}

	body, _ := json.Marshal(product)
	req := httptest.NewRequest("POST", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), UserIDKey, int64(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database expectations were not met: %v", err)
	}
}

func TestProductsHandler_InvalidMethod(t *testing.T) {
	deps, _ := createTestDeps(t)
	defer deps.DB.Close()

	handler := handlers.ProductsHandler(deps)

	req := httptest.NewRequest("DELETE", "/products", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestProductsHandler_Unauthorized(t *testing.T) {
	deps, _ := createTestDeps(t)
	defer deps.DB.Close()

	handler := handlers.ProductsHandler(deps)

	// Request without user_id in context
	req := httptest.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestProductsHandler_GET_WithGroup(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)

	// Mock getUser call
	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).AddRow(userID, "user1", "hash1"))

	// Mock group query (user is in a group with 3 members)
	groupRows := sqlmock.NewRows([]string{"user_id"}).
		AddRow(1).
		AddRow(2).
		AddRow(3)

	mock.ExpectQuery("SELECT DISTINCT user_id FROM groups WHERE id =").
		WithArgs(userID).
		WillReturnRows(groupRows)

	// Mock database query for products from all group members
	productRows := sqlmock.NewRows([]string{"id", "name", "volume", "brand", "default_tags", "user_id"}).
		AddRow(1, "Product1", "500ml", "Brand1", "tag1", 1).
		AddRow(2, "Product2", "1L", "Brand2", "tag2", 2).
		AddRow(3, "Product3", "250ml", "Brand3", "tag3", 3)

	mock.ExpectQuery("SELECT id, name, volume, brand, default_tags, user_id FROM products WHERE user_id = ANY").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(productRows)

	handler := handlers.ProductsHandler(deps)

	req := httptest.NewRequest("GET", "/products", nil)
	ctx := context.WithValue(req.Context(), UserIDKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Parse response to verify we got products from all group members
	var response map[string][]models.Product
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	products := response["products"]
	if len(products) != 3 {
		t.Errorf("Expected 3 products from group members, got %d", len(products))
	}

	// Verify that products have correct user_ids
	expectedUserIDs := map[int64]bool{1: true, 2: true, 3: true}
	for _, p := range products {
		if !expectedUserIDs[p.UserId] {
			t.Errorf("Unexpected user_id %d in product", p.UserId)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database expectations were not met: %v", err)
	}
}

func TestProductsHandler_PUT_Success(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)

	// Mock getUser call
	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).AddRow(userID, "user1", "hash1"))

	// Mock successful product update
	result := sqlmock.NewResult(0, 1) // 1 row affected
	mock.ExpectExec("UPDATE products SET name=\\$1, volume=\\$2, brand=\\$3, default_tags=\\$4 WHERE id=\\$5 AND user_id=\\$6").
		WithArgs("UpdatedProduct", "1L", "UpdatedBrand", "tag1,tag2", 123, userID).
		WillReturnResult(result)

	handler := handlers.ProductsHandler(deps)

	product := models.Product{
		Id:          123,
		Name:        "UpdatedProduct",
		Volume:      "1L",
		Brand:       "UpdatedBrand",
		DefaultTags: []string{"tag1", "tag2"},
	}

	body, _ := json.Marshal(product)
	req := httptest.NewRequest("PUT", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), UserIDKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Verify response contains updated product
	var responseProduct models.Product
	if err := json.NewDecoder(w.Body).Decode(&responseProduct); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if responseProduct.Name != "UpdatedProduct" {
		t.Errorf("Expected product name 'UpdatedProduct', got '%s'", responseProduct.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database expectations were not met: %v", err)
	}
}

func TestProductsHandler_PUT_NotFound(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)

	// Mock getUser call
	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).AddRow(userID, "user1", "hash1"))

	// Mock update with no rows affected
	result := sqlmock.NewResult(0, 0) // 0 rows affected
	mock.ExpectExec("UPDATE products SET name=\\$1, volume=\\$2, brand=\\$3, default_tags=\\$4 WHERE id=\\$5 AND user_id=\\$6").
		WithArgs("UpdatedProduct", "1L", "UpdatedBrand", "tag1", 999, userID).
		WillReturnResult(result)

	handler := handlers.ProductsHandler(deps)

	product := models.Product{
		Id:          999,
		Name:        "UpdatedProduct",
		Volume:      "1L",
		Brand:       "UpdatedBrand",
		DefaultTags: []string{"tag1"},
	}

	body, _ := json.Marshal(product)
	req := httptest.NewRequest("PUT", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), UserIDKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database expectations were not met: %v", err)
	}
}

func TestProductsHandler_PUT_Unauthorized(t *testing.T) {
	deps, _ := createTestDeps(t)
	defer deps.DB.Close()

	handler := handlers.ProductsHandler(deps)

	product := models.Product{
		Id:          123,
		Name:        "UpdatedProduct",
		Volume:      "1L",
		Brand:       "UpdatedBrand",
		DefaultTags: []string{"tag1"},
	}

	body, _ := json.Marshal(product)
	req := httptest.NewRequest("PUT", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestProductsHandler_PUT_InvalidJSON(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)

	// Mock getUser call
	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).AddRow(userID, "user1", "hash1"))

	handler := handlers.ProductsHandler(deps)

	req := httptest.NewRequest("PUT", "/products", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), UserIDKey, int64(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestProductsHandler_PUT_MissingID(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)

	// Mock getUser call
	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).AddRow(userID, "user1", "hash1"))

	handler := handlers.ProductsHandler(deps)

	product := models.Product{
		Id:          0, // Missing/invalid ID
		Name:        "UpdatedProduct",
		Volume:      "1L",
		Brand:       "UpdatedBrand",
		DefaultTags: []string{"tag1"},
	}

	body, _ := json.Marshal(product)
	req := httptest.NewRequest("PUT", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), UserIDKey, int64(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestProductsHandler_PUT_ValidationError(t *testing.T) {
	deps, mock := createTestDepsWithFailingValidator(t)
	defer deps.DB.Close()

	userID := int64(1)

	// Mock getUser call
	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).AddRow(userID, "user1", "hash1"))

	handler := handlers.ProductsHandler(deps)

	product := models.Product{
		Id:          123,
		Name:        "UpdatedProduct",
		Volume:      "1L",
		Brand:       "UpdatedBrand",
		DefaultTags: []string{"tag1"},
	}

	body, _ := json.Marshal(product)
	req := httptest.NewRequest("PUT", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), UserIDKey, int64(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 (validation error), got %d", w.Code)
	}
}

func TestProductsHandler_PUT_DifferentUserProduct(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)

	// Mock getUser call
	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).AddRow(userID, "user1", "hash1"))

	// User 1 tries to update product owned by user 2
	// The query includes user_id check, so no rows will be affected
	result := sqlmock.NewResult(0, 0) // 0 rows affected
	mock.ExpectExec("UPDATE products SET name=\\$1, volume=\\$2, brand=\\$3, default_tags=\\$4 WHERE id=\\$5 AND user_id=\\$6").
		WithArgs("UpdatedProduct", "1L", "UpdatedBrand", "tag1", 456, userID).
		WillReturnResult(result)

	handler := handlers.ProductsHandler(deps)

	product := models.Product{
		Id:          456, // Product owned by different user
		Name:        "UpdatedProduct",
		Volume:      "1L",
		Brand:       "UpdatedBrand",
		DefaultTags: []string{"tag1"},
	}

	body, _ := json.Marshal(product)
	req := httptest.NewRequest("PUT", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), UserIDKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 (product not found for this user), got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database expectations were not met: %v", err)
	}
}

func TestProductsHandler_PUT_DatabaseError(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)

	// Mock getUser call
	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).AddRow(userID, "user1", "hash1"))

	// Mock database error during update
	mock.ExpectExec("UPDATE products SET name=\\$1, volume=\\$2, brand=\\$3, default_tags=\\$4 WHERE id=\\$5 AND user_id=\\$6").
		WithArgs("UpdatedProduct", "1L", "UpdatedBrand", "tag1", 123, userID).
		WillReturnError(sqlmock.ErrCancelled)

	handler := handlers.ProductsHandler(deps)

	product := models.Product{
		Id:          123,
		Name:        "UpdatedProduct",
		Volume:      "1L",
		Brand:       "UpdatedBrand",
		DefaultTags: []string{"tag1"},
	}

	body, _ := json.Marshal(product)
	req := httptest.NewRequest("PUT", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), UserIDKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database expectations were not met: %v", err)
	}
}

func TestProductsHandler_PUT_EmptyTags(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)

	// Mock getUser call
	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).AddRow(userID, "user1", "hash1"))

	// Mock successful product update with empty tags
	result := sqlmock.NewResult(0, 1) // 1 row affected
	mock.ExpectExec("UPDATE products SET name=\\$1, volume=\\$2, brand=\\$3, default_tags=\\$4 WHERE id=\\$5 AND user_id=\\$6").
		WithArgs("UpdatedProduct", "1L", "UpdatedBrand", "", 123, userID).
		WillReturnResult(result)

	handler := handlers.ProductsHandler(deps)

	product := models.Product{
		Id:          123,
		Name:        "UpdatedProduct",
		Volume:      "1L",
		Brand:       "UpdatedBrand",
		DefaultTags: []string{}, // Empty tags
	}

	body, _ := json.Marshal(product)
	req := httptest.NewRequest("PUT", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), UserIDKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database expectations were not met: %v", err)
	}
}

func TestProductsHandler_PUT_MultipleTags(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)

	// Mock getUser call
	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).AddRow(userID, "user1", "hash1"))

	// Mock successful product update with multiple tags
	result := sqlmock.NewResult(0, 1) // 1 row affected
	mock.ExpectExec("UPDATE products SET name=\\$1, volume=\\$2, brand=\\$3, default_tags=\\$4 WHERE id=\\$5 AND user_id=\\$6").
		WithArgs("UpdatedProduct", "1L", "UpdatedBrand", "tag1,tag2,tag3,tag4,tag5", 123, userID).
		WillReturnResult(result)

	handler := handlers.ProductsHandler(deps)

	product := models.Product{
		Id:          123,
		Name:        "UpdatedProduct",
		Volume:      "1L",
		Brand:       "UpdatedBrand",
		DefaultTags: []string{"tag1", "tag2", "tag3", "tag4", "tag5"},
	}

	body, _ := json.Marshal(product)
	req := httptest.NewRequest("PUT", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), UserIDKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database expectations were not met: %v", err)
	}
}

func TestProductsHandler_PUT_TooManyTags(t *testing.T) {
	deps, mock := createTestDepsWithFailingValidator(t)
	defer deps.DB.Close()

	userID := int64(1)

	// Mock getUser call
	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).AddRow(userID, "user1", "hash1"))

	handler := handlers.ProductsHandler(deps)

	// Create product with 11 tags (exceeds limit of 10)
	tags := make([]string, 11)
	for i := 0; i < 11; i++ {
		tags[i] = "tag"
	}

	product := models.Product{
		Id:          123,
		Name:        "UpdatedProduct",
		Volume:      "1L",
		Brand:       "UpdatedBrand",
		DefaultTags: tags,
	}

	body, _ := json.Marshal(product)
	req := httptest.NewRequest("PUT", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), UserIDKey, int64(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 (too many tags), got %d", w.Code)
	}
}