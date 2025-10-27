package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInviteHandler_GET(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)
	now := time.Now()

	// Mock getting incoming invites
	rows := sqlmock.NewRows([]string{"id", "from_user_id", "to_user_id", "from_login", "to_login", "created_at"}).
		AddRow(1, 2, 1, "sender", "receiver", now).
		AddRow(2, 3, 1, "sender2", "receiver", now)

	mock.ExpectQuery("SELECT i.id, i.from_user_id, i.to_user_id, u_from.login, u_to.login, i.created_at FROM invites i").
		WithArgs(userID).
		WillReturnRows(rows)

	handler := InviteHandler(deps)

	req := httptest.NewRequest("GET", "/invite", nil)
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

func TestInviteHandler_POST_NewInvite(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)
	targetUserID := int64(2)

	requestBody := map[string]string{
		"login": "target_user",
	}
	body, _ := json.Marshal(requestBody)

	// Get target user ID
	mock.ExpectQuery("SELECT id FROM users WHERE login = \\$1").
		WithArgs("target_user").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(targetUserID))

	// Check if users are already in groups (both not in groups)
	mock.ExpectQuery("SELECT id FROM groups WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectQuery("SELECT id FROM groups WHERE user_id = \\$1").
		WithArgs(targetUserID).
		WillReturnError(sql.ErrNoRows)

	// Begin transaction
	mock.ExpectBegin()

	// Check for reverse invite (not found)
	mock.ExpectQuery("SELECT id FROM invites WHERE from_user_id = \\$1 AND to_user_id = \\$2").
		WithArgs(targetUserID, userID).
		WillReturnError(sql.ErrNoRows)

	// Create invite
	mock.ExpectQuery("INSERT INTO invites \\(from_user_id, to_user_id\\) VALUES \\(\\$1, \\$2\\) RETURNING id").
		WithArgs(userID, targetUserID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Commit transaction
	mock.ExpectCommit()

	handler := InviteHandler(deps)

	req := httptest.NewRequest("POST", "/invite", bytes.NewReader(body))
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

func TestInviteHandler_POST_MutualInvite_NewGroup(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)
	targetUserID := int64(2)
	newGroupID := int64(1)

	requestBody := map[string]string{
		"login": "target_user",
	}
	body, _ := json.Marshal(requestBody)

	// Get target user ID
	mock.ExpectQuery("SELECT id FROM users WHERE login = \\$1").
		WithArgs("target_user").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(targetUserID))

	// Check if users are already in groups (both not in groups)
	mock.ExpectQuery("SELECT id FROM groups WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectQuery("SELECT id FROM groups WHERE user_id = \\$1").
		WithArgs(targetUserID).
		WillReturnError(sql.ErrNoRows)

	// Begin transaction
	mock.ExpectBegin()

	// Check for reverse invite (found - mutual invite!)
	mock.ExpectQuery("SELECT id FROM invites WHERE from_user_id = \\$1 AND to_user_id = \\$2").
		WithArgs(targetUserID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Delete both invites
	mock.ExpectExec("DELETE FROM invites WHERE \\(from_user_id = \\$1 AND to_user_id = \\$2\\) OR \\(from_user_id = \\$2 AND to_user_id = \\$1\\)").
		WithArgs(userID, targetUserID).
		WillReturnResult(sqlmock.NewResult(0, 2))

	// Create new group
	mock.ExpectQuery("INSERT INTO groups \\(user_id, member_number\\) VALUES \\(\\$1, \\$2\\) RETURNING id").
		WithArgs(userID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newGroupID))

	// Add second user to group
	mock.ExpectExec("INSERT INTO groups \\(id, user_id, member_number\\) VALUES \\(\\$1, \\$2, \\$3\\)").
		WithArgs(newGroupID, targetUserID, 2).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Commit transaction
	mock.ExpectCommit()

	handler := InviteHandler(deps)

	req := httptest.NewRequest("POST", "/invite", bytes.NewReader(body))
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

func TestInviteHandler_POST_MutualInvite_ExistingGroup(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)
	targetUserID := int64(2)
	existingGroupID := int64(1)

	requestBody := map[string]string{
		"login": "target_user",
	}
	body, _ := json.Marshal(requestBody)

	// Get target user ID
	mock.ExpectQuery("SELECT id FROM users WHERE login = \\$1").
		WithArgs("target_user").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(targetUserID))

	// Check if users are already in groups (user is in group, target is not)
	mock.ExpectQuery("SELECT id FROM groups WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(existingGroupID))

	mock.ExpectQuery("SELECT id FROM groups WHERE user_id = \\$1").
		WithArgs(targetUserID).
		WillReturnError(sql.ErrNoRows)

	// Check group size
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM groups WHERE id = \\$1").
		WithArgs(existingGroupID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

	// Begin transaction
	mock.ExpectBegin()

	// Check for reverse invite (found - mutual invite!)
	mock.ExpectQuery("SELECT id FROM invites WHERE from_user_id = \\$1 AND to_user_id = \\$2").
		WithArgs(targetUserID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Delete both invites
	mock.ExpectExec("DELETE FROM invites WHERE \\(from_user_id = \\$1 AND to_user_id = \\$2\\) OR \\(from_user_id = \\$2 AND to_user_id = \\$1\\)").
		WithArgs(userID, targetUserID).
		WillReturnResult(sqlmock.NewResult(0, 2))

	// Re-check groups within transaction
	mock.ExpectQuery("SELECT id FROM groups WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(existingGroupID))

	mock.ExpectQuery("SELECT id FROM groups WHERE user_id = \\$1").
		WithArgs(targetUserID).
		WillReturnError(sql.ErrNoRows)

	// Check group size again in transaction
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM groups WHERE id = \\$1").
		WithArgs(existingGroupID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

	// Add target user to existing group
	mock.ExpectExec("INSERT INTO groups \\(id, user_id, member_number\\) VALUES \\(\\$1, \\$2, \\$3\\)").
		WithArgs(existingGroupID, targetUserID, 4).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Commit transaction
	mock.ExpectCommit()

	handler := InviteHandler(deps)

	req := httptest.NewRequest("POST", "/invite", bytes.NewReader(body))
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

func TestInviteHandler_POST_TargetUserInGroup(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)
	targetUserID := int64(2)
	targetGroupID := int64(1)

	requestBody := map[string]string{
		"login": "target_user",
	}
	body, _ := json.Marshal(requestBody)

	// Get target user ID
	mock.ExpectQuery("SELECT id FROM users WHERE login = \\$1").
		WithArgs("target_user").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(targetUserID))

	// Check if users are already in groups (current user not in group, target is in group)
	mock.ExpectQuery("SELECT id FROM groups WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectQuery("SELECT id FROM groups WHERE user_id = \\$1").
		WithArgs(targetUserID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(targetGroupID))

	// Check target group size (not at maximum)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM groups WHERE id = \\$1").
		WithArgs(targetGroupID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

	// Begin transaction
	mock.ExpectBegin()

	// Check for reverse invite (not found)
	mock.ExpectQuery("SELECT id FROM invites WHERE from_user_id = \\$1 AND to_user_id = \\$2").
		WithArgs(targetUserID, userID).
		WillReturnError(sql.ErrNoRows)

	// Create invite
	mock.ExpectQuery("INSERT INTO invites \\(from_user_id, to_user_id\\) VALUES \\(\\$1, \\$2\\) RETURNING id").
		WithArgs(userID, targetUserID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Commit transaction
	mock.ExpectCommit()

	handler := InviteHandler(deps)

	req := httptest.NewRequest("POST", "/invite", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), UserIDKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Now this should succeed - we can invite users who are in groups
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database expectations were not met: %v", err)
	}
}

func TestInviteHandler_POST_BothUsersInDifferentGroups(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)
	targetUserID := int64(2)
	userGroupID := int64(1)
	targetGroupID := int64(2)

	requestBody := map[string]string{
		"login": "target_user",
	}
	body, _ := json.Marshal(requestBody)

	// Get target user ID
	mock.ExpectQuery("SELECT id FROM users WHERE login = \\$1").
		WithArgs("target_user").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(targetUserID))

	// Check if users are already in groups (both in different groups)
	mock.ExpectQuery("SELECT id FROM groups WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userGroupID))

	mock.ExpectQuery("SELECT id FROM groups WHERE user_id = \\$1").
		WithArgs(targetUserID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(targetGroupID))

	handler := InviteHandler(deps)

	req := httptest.NewRequest("POST", "/invite", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), UserIDKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Should fail - cannot invite users in different groups
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database expectations were not met: %v", err)
	}
}

func TestInviteHandler_POST_GroupSizeLimitReached(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)
	targetUserID := int64(2)
	existingGroupID := int64(1)

	requestBody := map[string]string{
		"login": "target_user",
	}
	body, _ := json.Marshal(requestBody)

	// Get target user ID
	mock.ExpectQuery("SELECT id FROM users WHERE login = \\$1").
		WithArgs("target_user").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(targetUserID))

	// Check if users are already in groups (user is in group, target is not)
	mock.ExpectQuery("SELECT id FROM groups WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(existingGroupID))

	mock.ExpectQuery("SELECT id FROM groups WHERE user_id = \\$1").
		WithArgs(targetUserID).
		WillReturnError(sql.ErrNoRows)

	// Check group size - already at maximum (5)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM groups WHERE id = \\$1").
		WithArgs(existingGroupID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	handler := InviteHandler(deps)

	req := httptest.NewRequest("POST", "/invite", bytes.NewReader(body))
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

func TestInviteHandler_POST_UserNotFound(t *testing.T) {
	deps, mock := createTestDeps(t)
	defer deps.DB.Close()

	userID := int64(1)

	requestBody := map[string]string{
		"login": "nonexistent_user",
	}
	body, _ := json.Marshal(requestBody)

	// Get target user ID - not found
	mock.ExpectQuery("SELECT id FROM users WHERE login = \\$1").
		WithArgs("nonexistent_user").
		WillReturnError(sql.ErrNoRows)

	handler := InviteHandler(deps)

	req := httptest.NewRequest("POST", "/invite", bytes.NewReader(body))
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

func TestInviteHandler_MethodNotAllowed(t *testing.T) {
	deps, _ := createTestDeps(t)
	defer deps.DB.Close()

	handler := InviteHandler(deps)

	req := httptest.NewRequest("PUT", "/invite", nil)
	ctx := context.WithValue(req.Context(), UserIDKey, int64(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}
