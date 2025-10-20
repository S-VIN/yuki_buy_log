package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/lib/pq"
	"yuki_buy_log/models"
)

func PurchasesHandler(deps *Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	log.Printf("Purchases handler called: %s %s", r.Method, r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		getPurchases(deps, w, r)
	case http.MethodPost:
		createPurchase(deps, w, r)
	default:
		log.Printf("Method not allowed for purchases: %s", r.Method)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func getPurchases(deps *Dependencies, w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching purchases from database")
	uid, ok := userID(r)
	if !ok {
		log.Println("Unauthorized access attempt to purchases")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("Fetching purchases for user ID: %d", uid)
	rows, err := deps.DB.Query(`SELECT id, product_id, quantity, price, date, store, tags, receipt_id FROM purchases WHERE user_id=$1`, uid)
	if err != nil {
		log.Printf("Failed to query purchases: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	purchases := []models.Purchase{}
	for rows.Next() {
		var p models.Purchase
		if err := rows.Scan(&p.Id, &p.ProductId, &p.Quantity, &p.Price, &p.Date, &p.Store, pq.Array(&p.Tags), &p.ReceiptId); err != nil {
			log.Printf("Failed to scan purchase row: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		p.UserId = uid
		purchases = append(purchases, p)
	}
	log.Printf("Successfully fetched %d purchases for user %d", len(purchases), uid)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"purchases": purchases})
}

func createPurchase(deps *Dependencies, w http.ResponseWriter, r *http.Request) {
	log.Println("Creating new purchase")
	var p models.Purchase
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("Failed to decode purchase JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	uid, ok := userID(r)
	if !ok {
		log.Println("Unauthorized access attempt to create purchase")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	p.UserId = uid
	if err := deps.Validator.ValidatePurchase(&p); err != nil {
		log.Printf("Purchase validation failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Creating purchase for user ID: %d", uid)
	err := deps.DB.QueryRow(`INSERT INTO purchases (product_id, quantity, price, date, store, tags, receipt_id, user_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		p.ProductId, p.Quantity, p.Price, p.Date, p.Store, pq.Array(p.Tags), p.ReceiptId, uid).Scan(&p.Id)
	if err != nil {
		log.Printf("Failed to insert purchase: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully created purchase with ID: %d", p.Id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}