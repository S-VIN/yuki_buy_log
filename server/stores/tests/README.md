# Stores Tests

## Текущий статус

Тесты для stores находятся в `handlers/tests`, где используются моки для Store interfaces.

## Почему тесты не в этой директории?

Текущая архитектура stores использует:
- Singleton pattern (GetXStore())
- Приватные поля
- Прямые вызовы database функций

Это затрудняет юнит-тестирование stores изолированно.

## Рекомендации для рефакторинга

Для полноценного юнит-тестирования stores требуется:

1. **Убрать singleton pattern**
   ```go
   // Вместо:
   func GetUserStore() *UserStore { ... }

   // Использовать:
   func NewUserStore(db DatabaseInterface) *UserStore { ... }
   ```

2. **Dependency Injection для database**
   ```go
   type DatabaseOperations interface {
       AddUser(*models.User) error
       UpdateUser(*models.User) error
       // ...
   }

   type UserStore struct {
       db DatabaseOperations
       // ...
   }
   ```

3. **Избегать глобального состояния**
   - Stores должны создаваться через конструкторы
   - Database dependency должен инжектиться

## Текущие тесты

Функциональность stores тестируется через:
- `handlers/tests/products_test.go` - тесты Products и Login/Register handlers
- `handlers/tests/mocks.go` - моки для всех Store interfaces

Эти тесты покрывают основные сценарии использования stores через HTTP handlers.
