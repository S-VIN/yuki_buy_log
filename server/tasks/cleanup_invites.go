package tasks

import (
	"log"
	"time"
	"yuki_buy_log/handlers"
)

// CleanupOldInvites removes invites that are older than 24 hours
func CleanupOldInvites() func() {
	return func() {
		cutoffTime := time.Now().Add(-24 * time.Hour)

		inviteStore := handlers.GetInviteStore()
		rowsAffected, err := inviteStore.DeleteOldInvites(cutoffTime)
		if err != nil {
			log.Printf("Failed to cleanup old invites: %v", err)
			return
		}

		if rowsAffected > 0 {
			log.Printf("Cleaned up %d old invite(s) created before %s", rowsAffected, cutoffTime.Format(time.RFC3339))
		}
	}
}
