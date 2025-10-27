package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"yuki_buy_log/models"

	"github.com/lib/pq"
)

func PurchasesHandler(d *Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Purchases handler called: %s %s", r.Method, r.URL.Path)
		switch r.Method {
		case http.MethodGet:
			getPurchases(d, w, r)
		case http.MethodPost:
			createPurchase(d, w, r)
		case http.MethodDelete:
			deletePurchase(d, w, r)
		default:
			log.Printf("Method not allowed for purchases: %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func getPurchases(d *Dependencies, w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching purchases from database")
	user, err := getUser(d, r)
	if err != nil {
		log.Println("Unauthorized access attempt to purchases")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("Fetching purchases for user ID: %d and their group", user.Id)

	// Get all user IDs in the same group (including current user)
	// If user is not in a group, just return their own purchases
	var query string
	var args []interface{}

	// Try to get group members
	rows, err := d.DB.Query(`
		SELECT DISTINCT user_id
		FROM groups
		WHERE id = (SELECT id FROM groups WHERE user_id = $1 LIMIT 1)
	`, user.Id)

	if err != nil {
		// User might not be in a group, just query their own purchases
		log.Printf("User %d is not in a group or error getting group members, fetching only their purchases", user.Id)
		query = `SELECT id, product_id, quantity, price, date, store, tags, receipt_id, user_id FROM purchases WHERE user_id=$1`
		args = []interface{}{user.Id}
	} else {
		defer rows.Close()

		userIDs := []int64{}
		for rows.Next() {
			var userIDInGroup int64
			if err := rows.Scan(&userIDInGroup); err != nil {
				log.Printf("Failed to scan user ID: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			userIDs = append(userIDs, userIDInGroup)
		}

		if len(userIDs) == 0 {
			// User is not in a group, fetch only their purchases
			log.Printf("User %d is not in a group, fetching only their purchases", user.Id)
			query = `SELECT id, product_id, quantity, price, date, store, tags, receipt_id, user_id FROM purchases WHERE user_id=$1`
			args = []interface{}{user.Id}
		} else {
			// Build query with IN clause for all group members
			log.Printf("User %d is in a group with %d members, fetching purchases for all", user.Id, len(userIDs))
			query = `SELECT id, product_id, quantity, price, date, store, tags, receipt_id, user_id FROM purchases WHERE user_id = ANY($1)`
			args = []interface{}{pq.Array(userIDs)}
		}
	}

	purchaseRows, err := d.DB.Query(query, args...)
	if err != nil {
		log.Printf("Failed to query purchases: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer purchaseRows.Close()

	purchases := []models.Purchase{}
	for purchaseRows.Next() {
		var p models.Purchase
		if err := purchaseRows.Scan(&p.Id, &p.ProductId, &p.Quantity, &p.Price, &p.Date, &p.Store, pq.Array(&p.Tags), &p.ReceiptId, &p.UserId); err != nil {
			log.Printf("Failed to scan purchase row: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		purchases = append(purchases, p)
	}
	log.Printf("Successfully fetched %d purchases for user %d", len(purchases), user.Id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"purchases": purchases})
}

func createPurchase(d *Dependencies, w http.ResponseWriter, r *http.Request) {
	log.Println("Creating new purchase")
	var p models.Purchase
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("Failed to decode purchase JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := getUser(d, r)
	if err != nil {
		log.Println("Unauthorized access attempt to create purchase")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	p.UserId = user.Id
	if err := d.Validator.ValidatePurchase(&p); err != nil {
		log.Printf("Purchase validation failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Creating purchase for user ID: %d", user.Id)
	err = d.DB.QueryRow(`INSERT INTO purchases (product_id, quantity, price, date, store, tags, receipt_id, user_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		p.ProductId, p.Quantity, p.Price, p.Date, p.Store, pq.Array(p.Tags), p.ReceiptId, user.Id).Scan(&p.Id)
	if err != nil {
		log.Printf("Failed to insert purchase: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully created purchase with ID: %d", p.Id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func deletePurchase(d *Dependencies, w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting purchase")
	user, err := getUser(d, r)
	if err != nil {
		log.Println("Unauthorized access attempt to delete purchase")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Id int64 `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode delete request JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Id == 0 {
		log.Println("Missing id in request body")
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	log.Printf("Deleting purchase ID: %d for user ID: %d", req.Id, user.Id)
	result, err := d.DB.Exec(`DELETE FROM purchases WHERE id=$1 AND user_id=$2`, req.Id, user.Id)
	if err != nil {
		log.Printf("Failed to delete purchase: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to check rows affected: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		log.Printf("Purchase with ID %d not found for user %d", req.Id, user.Id)
		http.Error(w, "purchase not found", http.StatusNotFound)
		return
	}

	log.Printf("Successfully deleted purchase with ID: %d", req.Id)
	w.WriteHeader(http.StatusNoContent)
}
