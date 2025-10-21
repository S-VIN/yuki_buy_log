package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGroupHandler_GET_WithGroup(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)

	// Mock getting user's group ID
	mock.ExpectQuery("SELECT id FROM groups WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Mock getting group members
	rows := sqlmock.NewRows([]string{"id", "user_id", "login"}).
		AddRow(1, 1, "user1").
		AddRow(1, 2, "user2")

	mock.ExpectQuery("SELECT g.id, g.user_id, u.login FROM groups g JOIN users u").
		WithArgs(1).
		WillReturnRows(rows)

	handler := GroupHandler(deps)

	req := httptest.NewRequest("GET", "/group", nil)
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

func TestGroupHandler_GET_NoGroup(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)

	// Mock getting user's group ID - no group found
	mock.ExpectQuery("SELECT id FROM groups WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	handler := GroupHandler(deps)

	req := httptest.NewRequest("GET", "/group", nil)
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

func TestGroupHandler_DELETE_Success(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)
	groupID := int64(1)

	// Mock getting user's group ID
	mock.ExpectQuery("SELECT id FROM groups WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(groupID))

	// Begin transaction
	mock.ExpectBegin()

	// Remove user from group
	mock.ExpectExec("DELETE FROM groups WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Check remaining members count (more than 1 remaining)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM groups WHERE id = \\$1").
		WithArgs(groupID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	// Commit transaction
	mock.ExpectCommit()

	handler := GroupHandler(deps)

	req := httptest.NewRequest("DELETE", "/group", nil)
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

func TestGroupHandler_DELETE_AutoDeleteGroup(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)
	groupID := int64(1)

	// Mock getting user's group ID
	mock.ExpectQuery("SELECT id FROM groups WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(groupID))

	// Begin transaction
	mock.ExpectBegin()

	// Remove user from group
	mock.ExpectExec("DELETE FROM groups WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Check remaining members count (only 1 remaining)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM groups WHERE id = \\$1").
		WithArgs(groupID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Delete the entire group
	mock.ExpectExec("DELETE FROM groups WHERE id = \\$1").
		WithArgs(groupID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Commit transaction
	mock.ExpectCommit()

	handler := GroupHandler(deps)

	req := httptest.NewRequest("DELETE", "/group", nil)
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

func TestGroupHandler_DELETE_NotInGroup(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)

	// Mock getting user's group ID - no group found
	mock.ExpectQuery("SELECT id FROM groups WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	handler := GroupHandler(deps)

	req := httptest.NewRequest("DELETE", "/group", nil)
	ctx := context.WithValue(req.Context(), UserIDKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database expectations were not met: %v", err)
	}
}

func TestGroupHandler_MethodNotAllowed(t *testing.T) {
	deps, _ := createTestDeps(t)
	defer deps.DB.Close()

	handler := GroupHandler(deps)

	req := httptest.NewRequest("POST", "/group", nil)
	ctx := context.WithValue(req.Context(), UserIDKey, int64(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}
