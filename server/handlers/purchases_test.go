package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"yuki_buy_log/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

func TestPurchasesHandler_GET(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)

	// Mock group query (user not in a group)
	mock.ExpectQuery("SELECT DISTINCT user_id FROM groups WHERE id =").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}))

	// Mock database query for purchases (now includes user_id)
	rows := sqlmock.NewRows([]string{"id", "product_id", "quantity", "price", "date", "store", "tags", "receipt_id", "user_id"}).
		AddRow(1, 1, 2, 1000, time.Now(), "TestStore", pq.Array([]string{"tag1", "tag2"}), 123, userID)

	mock.ExpectQuery("SELECT id, product_id, quantity, price, date, store, tags, receipt_id, user_id FROM purchases WHERE user_id=\\$1").
		WithArgs(userID).
		WillReturnRows(rows)

	handler := PurchasesHandler(deps)

	req := httptest.NewRequest("GET", "/purchases", nil)
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

func TestPurchasesHandler_POST(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	// Mock successful purchase insertion
	mock.ExpectQuery("INSERT INTO purchases").
		WithArgs(1, 2, 1000, sqlmock.AnyArg(), "TestStore", pq.Array([]string{"test"}), int64(123), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	handler := PurchasesHandler(deps)

	purchase := models.Purchase{
		ProductId: 1,
		Quantity:  2,
		Price:     1000,
		Date:      time.Now(),
		Store:     "TestStore",
		Tags:      []string{"test"},
		ReceiptId: 123,
	}

	body, _ := json.Marshal(purchase)
	req := httptest.NewRequest("POST", "/purchases", bytes.NewBuffer(body))
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

func TestPurchasesHandler_DELETE_Success(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	// Mock successful deletion
	result := sqlmock.NewResult(0, 1) // 1 row affected
	mock.ExpectExec("DELETE FROM purchases WHERE id=\\$1 AND user_id=\\$2").
		WithArgs(123, 1).
		WillReturnResult(result)

	handler := PurchasesHandler(deps)

	deleteReq := struct {
		Id int64 `json:"id"`
	}{Id: 123}

	body, _ := json.Marshal(deleteReq)
	req := httptest.NewRequest("DELETE", "/purchases", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), UserIDKey, int64(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database expectations were not met: %v", err)
	}
}

func TestPurchasesHandler_DELETE_NotFound(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	// Mock deletion with no rows affected
	result := sqlmock.NewResult(0, 0) // 0 rows affected
	mock.ExpectExec("DELETE FROM purchases WHERE id=\\$1 AND user_id=\\$2").
		WithArgs(999, 1).
		WillReturnResult(result)

	handler := PurchasesHandler(deps)

	deleteReq := struct {
		Id int64 `json:"id"`
	}{Id: 999}

	body, _ := json.Marshal(deleteReq)
	req := httptest.NewRequest("DELETE", "/purchases", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), UserIDKey, int64(1))
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

func TestPurchasesHandler_DELETE_Unauthorized(t *testing.T) {
	deps, _ := createTestDeps(t)
	defer deps.DB.Close()

	handler := PurchasesHandler(deps)

	deleteReq := struct {
		Id int64 `json:"id"`
	}{Id: 123}

	body, _ := json.Marshal(deleteReq)
	req := httptest.NewRequest("DELETE", "/purchases", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestPurchasesHandler_DELETE_InvalidJSON(t *testing.T) {
	deps, _ := createTestDeps(t)
	defer deps.DB.Close()

	handler := PurchasesHandler(deps)

	req := httptest.NewRequest("DELETE", "/purchases", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), UserIDKey, int64(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestPurchasesHandler_DELETE_MissingID(t *testing.T) {
	deps, _ := createTestDeps(t)
	defer deps.DB.Close()

	handler := PurchasesHandler(deps)

	deleteReq := struct {
		Id int64 `json:"id"`
	}{Id: 0}

	body, _ := json.Marshal(deleteReq)
	req := httptest.NewRequest("DELETE", "/purchases", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), UserIDKey, int64(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestPurchasesHandler_InvalidMethod(t *testing.T) {
	deps, _ := createTestDeps(t)
	defer deps.DB.Close()

	handler := PurchasesHandler(deps)

	req := httptest.NewRequest("PUT", "/purchases", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestPurchasesHandler_Unauthorized(t *testing.T) {
	deps, _ := createTestDeps(t)
	defer deps.DB.Close()

	handler := PurchasesHandler(deps)

	req := httptest.NewRequest("GET", "/purchases", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}