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

	// Mock getUser call
	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).AddRow(userID, "user1", "hash1"))

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

	userID := int64(1)

	// Mock getUser call
	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).AddRow(userID, "user1", "hash1"))

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

	userID := int64(1)

	// Mock getUser call
	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).AddRow(userID, "user1", "hash1"))

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

	userID := int64(1)

	// Mock getUser call
	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).AddRow(userID, "user1", "hash1"))

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
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)

	// Mock getUser call
	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).AddRow(userID, "user1", "hash1"))

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
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)

	// Mock getUser call
	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).AddRow(userID, "user1", "hash1"))

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

func TestPurchasesHandler_DELETE_DatabaseError(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)

	// Mock getUser call
	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).AddRow(userID, "user1", "hash1"))

	// Mock database error during deletion
	mock.ExpectExec("DELETE FROM purchases WHERE id=\\$1 AND user_id=\\$2").
		WithArgs(123, 1).
		WillReturnError(sqlmock.ErrCancelled)

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

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database expectations were not met: %v", err)
	}
}

func TestPurchasesHandler_DELETE_DifferentUserPurchase(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)

	// Mock getUser call
	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).AddRow(userID, "user1", "hash1"))

	// User 1 tries to delete purchase owned by user 2
	// The query includes user_id check, so no rows will be affected
	result := sqlmock.NewResult(0, 0) // 0 rows affected
	mock.ExpectExec("DELETE FROM purchases WHERE id=\\$1 AND user_id=\\$2").
		WithArgs(456, 1).
		WillReturnResult(result)

	handler := PurchasesHandler(deps)

	deleteReq := struct {
		Id int64 `json:"id"`
	}{Id: 456}

	body, _ := json.Marshal(deleteReq)
	req := httptest.NewRequest("DELETE", "/purchases", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), UserIDKey, int64(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 (purchase not found for this user), got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database expectations were not met: %v", err)
	}
}

func TestPurchasesHandler_DELETE_MultipleSuccessful(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)

	// First getUser call
	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).AddRow(userID, "user1", "hash1"))

	// First deletion
	result1 := sqlmock.NewResult(0, 1)
	mock.ExpectExec("DELETE FROM purchases WHERE id=\\$1 AND user_id=\\$2").
		WithArgs(100, userID).
		WillReturnResult(result1)

	// Second getUser call
	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).AddRow(userID, "user1", "hash1"))

	// Second deletion
	result2 := sqlmock.NewResult(0, 1)
	mock.ExpectExec("DELETE FROM purchases WHERE id=\\$1 AND user_id=\\$2").
		WithArgs(200, userID).
		WillReturnResult(result2)

	handler := PurchasesHandler(deps)

	// Delete first purchase
	deleteReq1 := struct {
		Id int64 `json:"id"`
	}{Id: 100}

	body1, _ := json.Marshal(deleteReq1)
	req1 := httptest.NewRequest("DELETE", "/purchases", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	ctx1 := context.WithValue(req1.Context(), UserIDKey, userID)
	req1 = req1.WithContext(ctx1)

	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, req1)

	if w1.Code != http.StatusNoContent {
		t.Errorf("First deletion: Expected status 204, got %d", w1.Code)
	}

	// Delete second purchase
	deleteReq2 := struct {
		Id int64 `json:"id"`
	}{Id: 200}

	body2, _ := json.Marshal(deleteReq2)
	req2 := httptest.NewRequest("DELETE", "/purchases", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	ctx2 := context.WithValue(req2.Context(), UserIDKey, userID)
	req2 = req2.WithContext(ctx2)

	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)

	if w2.Code != http.StatusNoContent {
		t.Errorf("Second deletion: Expected status 204, got %d", w2.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database expectations were not met: %v", err)
	}
}

func TestPurchasesHandler_GET_WithGroup(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)
	now := time.Now()

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

	// Mock database query for purchases from all group members
	purchaseRows := sqlmock.NewRows([]string{"id", "product_id", "quantity", "price", "date", "store", "tags", "receipt_id", "user_id"}).
		AddRow(1, 1, 2, 1000, now, "Store1", pq.Array([]string{"tag1"}), 100, 1).
		AddRow(2, 2, 3, 2000, now, "Store2", pq.Array([]string{"tag2"}), 200, 2).
		AddRow(3, 3, 1, 500, now, "Store3", pq.Array([]string{"tag3"}), 300, 3)

	mock.ExpectQuery("SELECT id, product_id, quantity, price, date, store, tags, receipt_id, user_id FROM purchases WHERE user_id = ANY").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(purchaseRows)

	handler := PurchasesHandler(deps)

	req := httptest.NewRequest("GET", "/purchases", nil)
	ctx := context.WithValue(req.Context(), UserIDKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Parse response to verify we got purchases from all group members
	var response map[string][]models.Purchase
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	purchases := response["purchases"]
	if len(purchases) != 3 {
		t.Errorf("Expected 3 purchases from group members, got %d", len(purchases))
	}

	// Verify that purchases have correct user_ids
	expectedUserIDs := map[int64]bool{1: true, 2: true, 3: true}
	for _, p := range purchases {
		if !expectedUserIDs[p.UserId] {
			t.Errorf("Unexpected user_id %d in purchase", p.UserId)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database expectations were not met: %v", err)
	}
}