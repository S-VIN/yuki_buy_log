package database

import (
	"database/sql"
	"log"
	"sync"
	"time"
	"yuki_buy_log/internal/utils"
)

type DatabaseManager struct {
	db *sql.DB
}

var (
	instance *DatabaseManager
	once     sync.Once
)

func GetDBManager() (*DatabaseManager, error) {
	once.Do(func() {
		instance = &DatabaseManager{}
		var err error
		log.Printf("Connecting to database with DSN: %s", utils.DatabaseURL)
		instance.db, err = sql.Open("postgres", utils.DatabaseURL)
		if err != nil {
			log.Fatalf("Failed to open database connection: %v", err)
		}

		log.Println("Waiting for database connection...")
		for i := 0; i < 30; i++ {
			if err = instance.db.Ping(); err == nil {
				log.Println("DatabaseManager connection established successfully")
				break
			}
			log.Printf("DatabaseManager connection attempt %d failed, retrying in 1 second...", i+1)
			time.Sleep(time.Second)
		}
		if err != nil {
			log.Fatalf("Failed to connect to database after 30 attempts: %v", err)
		}
	})
	return instance, nil
}

func (d *DatabaseManager) Close() {
	if d != nil && d.db != nil {
		d.db.Close()
	}
}
