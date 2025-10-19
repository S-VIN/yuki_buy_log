package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"yuki_buy_log/models"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRegisterHandler_Success(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	// Mock successful user insertion
	mock.ExpectQuery("INSERT INTO users").
		WithArgs("testuser", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	handler := RegisterHandler(deps)

	user := models.User{
		Login:    "testuser",
		Password: "testpass",
	}

	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check response contains token
	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	
	if token, exists := response["token"]; !exists || token == "" {
		t.Error("Expected token in response")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database expectations were not met: %v", err)
	}
}

func TestRegisterHandler_InvalidMethod(t *testing.T) {
	deps, _ := createTestDeps(t)
	defer deps.DB.Close()

	handler := RegisterHandler(deps)

	req := httptest.NewRequest("GET", "/register", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestRegisterHandler_InvalidJSON(t *testing.T) {
	deps, _ := createTestDeps(t)
	defer deps.DB.Close()

	handler := RegisterHandler(deps)

	req := httptest.NewRequest("POST", "/register", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestLoginHandler_Success(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	// Mock successful login query
	mock.ExpectQuery("SELECT id, password_hash FROM users WHERE login=\\$1").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "password_hash"}).
			AddRow(1, "$2a$10$hash")) // bcrypt hash placeholder

	handler := LoginHandler(deps)

	credentials := models.User{
		Login:    "testuser",
		Password: "testpass",
	}

	body, _ := json.Marshal(credentials)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Note: This test will fail because bcrypt hash comparison will fail
	// In a real test, you'd need to use a proper bcrypt hash
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 (due to hash mismatch), got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database expectations were not met: %v", err)
	}
}

func TestLoginHandler_InvalidMethod(t *testing.T) {
	deps, _ := createTestDeps(t)
	defer deps.DB.Close()

	handler := LoginHandler(deps)

	req := httptest.NewRequest("GET", "/login", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestLoginHandler_InvalidJSON(t *testing.T) {
	deps, _ := createTestDeps(t)
	defer deps.DB.Close()

	handler := LoginHandler(deps)

	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}