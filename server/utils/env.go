package utils

import "os"


var DatabaseURL string

func init() {
	DatabaseURL = os.Getenv("DATABASE_URL")
	if DatabaseURL == "" {
		DatabaseURL = "postgres://postgres:postgres@localhost:5432/yuki_buy_log?sslmode=disable"
	}
}

