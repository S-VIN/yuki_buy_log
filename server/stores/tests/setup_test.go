package stores_test

import (
	"os"
)

func init() {
	// Устанавливаем флаг для пропуска инициализации БД в тестах
	os.Setenv("SKIP_DB_INIT", "true")
}
