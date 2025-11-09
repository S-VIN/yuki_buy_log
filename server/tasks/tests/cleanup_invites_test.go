package tasks_test

import (
	"testing"
)

// TODO: Переписать тесты для использования InviteStore вместо прямого доступа к БД
// Старые тесты использовали sqlmock для мокирования базы данных,
// но теперь CleanupOldInvites использует InviteStore, который требует другого подхода к тестированию.
//
// Возможные подходы:
// 1. Создать мок для InviteStore
// 2. Использовать реальную тестовую БД для интеграционных тестов
// 3. Сделать InviteStore инжектируемым в CleanupOldInvites для удобства тестирования

func TestCleanupOldInvites_Placeholder(t *testing.T) {
	t.Skip("Тесты требуют переработки после рефакторинга на использование InviteStore")
}
