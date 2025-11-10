package database

import (
	"database/sql"
	"log"
	"os"
	"sync"
	"time"
	"yuki_buy_log/utils"
)

var db *sql.DB
var once sync.Once

func init() {
	// Skip database initialization during tests
	if os.Getenv("SKIP_DB_INIT") == "true" {
		log.Println("Skipping database initialization (SKIP_DB_INIT=true)")
		return
	}

	once.Do(func() {
		var err error
		log.Printf("Connecting to database with DSN: %s", utils.DatabaseURL)
		db, err = sql.Open("postgres", utils.DatabaseURL)
		if err != nil {
			log.Fatalf("Failed to open database connection: %v", err)
		}

		log.Println("Waiting for database connection...")
		for i := 0; i < 30; i++ {
			if err = db.Ping(); err == nil {
				log.Println("Database connection established successfully")
				break
			}
			log.Printf("Database connection attempt %d failed, retrying in 1 second...", i+1)
			time.Sleep(time.Second)
		}
		if err != nil {
			log.Fatalf("Failed to connect to database after 30 attempts: %v", err)
		}
	})
}

func Close() {
	if db != nil {
		db.Close()
	}
}
