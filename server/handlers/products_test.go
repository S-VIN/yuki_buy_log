package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"yuki_buy_log/models"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestProductsHandler_GET(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)

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

	handler := ProductsHandler(deps)

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

	// Mock successful product insertion
	mock.ExpectQuery("INSERT INTO products").
		WithArgs("TestProduct", "500ml", "TestBrand", "test", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	handler := ProductsHandler(deps)

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

	handler := ProductsHandler(deps)

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

	handler := ProductsHandler(deps)

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

	handler := ProductsHandler(deps)

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