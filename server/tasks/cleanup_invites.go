package tasks

import (
	"database/sql"
	"log"
	"time"
)

// CleanupOldInvites removes invites that are older than 24 hours
func CleanupOldInvites(db *sql.DB) func() {
	return func() {
		cutoffTime := time.Now().Add(-24 * time.Hour)

		result, err := db.Exec(`DELETE FROM invites WHERE created_at < $1`, cutoffTime)
		if err != nil {
			log.Printf("Failed to cleanup old invites: %v", err)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Printf("Failed to get rows affected: %v", err)
			return
		}

		if rowsAffected > 0 {
			log.Printf("Cleaned up %d old invite(s) created before %s", rowsAffected, cutoffTime.Format(time.RFC3339))
		}
	}
}
