package stores_test

import (
	"testing"
)

// ПРИМЕЧАНИЕ:
// Текущая архитектура stores использует singleton pattern и приватные поля,
// что затрудняет юнит-тестирование. Для полноценных тестов требуется рефакторинг:
//
// 1. Сделать stores создаваемыми через конструкторы вместо GetXStore()
// 2. Использовать dependency injection для database операций
// 3. Избегать глобального состояния
//
// До рефакторинга, тестирование stores проводится через handlers тесты,
// где используются моки для stores interfaces.
// См. handlers/tests/products_test.go, purchases_test.go, group_test.go, invite_test.go

// Placeholder тест, чтобы пакет компилировался
func TestStoresPlaceholder(t *testing.T) {
	t.Skip("Stores тестируются через handlers tests с использованием моков")
}
