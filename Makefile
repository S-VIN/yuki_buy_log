.PHONY: test test-verbose test-coverage

# Запуск всех тестов
test:
	cd server && SKIP_DB_INIT=true go test ./...

# Запуск тестов с подробным выводом
test-verbose:
	cd server && SKIP_DB_INIT=true go test -v ./...

# Запуск тестов с покрытием
test-coverage:
	cd server && SKIP_DB_INIT=true go test -coverprofile=coverage.out ./...
	cd server && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: server/coverage.html"
