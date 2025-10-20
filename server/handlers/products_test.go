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

	// Mock database query for products
	rows := sqlmock.NewRows([]string{"id", "name", "volume", "brand", "default_tags", "user_id"}).
		AddRow(1, "TestProduct", "500ml", "TestBrand", "tag1,tag2", 1)
	
	mock.ExpectQuery("SELECT id, name, volume, brand, default_tags, user_id FROM products WHERE user_id=\\$1").
		WithArgs(1).
		WillReturnRows(rows)

	handler := ProductsHandler(deps)

	req := httptest.NewRequest("GET", "/products", nil)
	ctx := context.WithValue(req.Context(), "user_id", int64(1))
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
	
	ctx := context.WithValue(req.Context(), "user_id", int64(1))
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