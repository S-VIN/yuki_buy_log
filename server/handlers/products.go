package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/lib/pq"
	"yuki_buy_log/models"
)

func ProductsHandler(deps *Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Products handler called: %s %s", r.Method, r.URL.Path)
		switch r.Method {
		case http.MethodGet:
			getProducts(deps, w, r)
		case http.MethodPost:
			createProduct(deps, w, r)
		default:
			log.Printf("Method not allowed for products: %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func getProducts(deps *Dependencies, w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching products from database")
	uid, ok := userID(r)
	if !ok {
		log.Println("Unauthorized access attempt to products")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("Fetching products for user ID: %d and their group", uid)

	// Get all user IDs in the same group (including current user)
	// If user is not in a group, just return their own products
	var query string
	var args []interface{}

	// Try to get group members
	rows, err := deps.DB.Query(`
		SELECT DISTINCT user_id
		FROM groups
		WHERE id = (SELECT id FROM groups WHERE user_id = $1 LIMIT 1)
	`, uid)

	if err != nil {
		// User might not be in a group, just query their own products
		log.Printf("User %d is not in a group or error getting group members, fetching only their products", uid)
		query = `SELECT id, name, volume, brand, default_tags, user_id FROM products WHERE user_id=$1`
		args = []interface{}{uid}
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
			// User is not in a group, fetch only their products
			log.Printf("User %d is not in a group, fetching only their products", uid)
			query = `SELECT id, name, volume, brand, default_tags, user_id FROM products WHERE user_id=$1`
			args = []interface{}{uid}
		} else {
			// Build query with IN clause for all group members
			log.Printf("User %d is in a group with %d members, fetching products for all", uid, len(userIDs))
			query = `SELECT id, name, volume, brand, default_tags, user_id FROM products WHERE user_id = ANY($1)`
			args = []interface{}{pq.Array(userIDs)}
		}
	}

	productRows, err := deps.DB.Query(query, args...)
	if err != nil {
		log.Printf("Failed to query products: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer productRows.Close()

	products := []models.Product{}
	for productRows.Next() {
		var p models.Product
		var defaultTagsStr string
		if err := productRows.Scan(&p.Id, &p.Name, &p.Volume, &p.Brand, &defaultTagsStr, &p.UserId); err != nil {
			log.Printf("Failed to scan product row: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if defaultTagsStr != "" {
			p.DefaultTags = strings.Split(defaultTagsStr, ",")
		} else {
			p.DefaultTags = []string{}
		}
		products = append(products, p)
	}
	log.Printf("Successfully fetched %d products for user %d", len(products), uid)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"products": products})
}

func createProduct(deps *Dependencies, w http.ResponseWriter, r *http.Request) {
	log.Println("Creating new product")
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("Failed to decode product JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	uid, ok := userID(r)
	if !ok {
		log.Println("Unauthorized access attempt to create product")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	p.UserId = uid
	if err := deps.Validator.ValidateProduct(&p); err != nil {
		log.Printf("Product validation failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defaultTagsStr := strings.Join(p.DefaultTags, ",")
	log.Printf("Creating product for user ID: %d", uid)
	err := deps.DB.QueryRow(`INSERT INTO products (name, volume, brand, default_tags, user_id) VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		p.Name, p.Volume, p.Brand, defaultTagsStr, uid).Scan(&p.Id)
	if err != nil {
		log.Printf("Failed to insert product: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully created product with ID: %d for user %d", p.Id, uid)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}