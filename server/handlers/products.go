package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"yuki_buy_log/models"

	"github.com/lib/pq"
)

func ProductsHandler(d *Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Products handler called: %s %s", r.Method, r.URL.Path)
		switch r.Method {
		case http.MethodGet:
			getProducts(d, w, r)
		case http.MethodPost:
			createProduct(d, w, r)
		case http.MethodPut:
			updateProduct(d, w, r)
		default:
			log.Printf("Method not allowed for products: %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func getProducts(d *Dependencies, w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching products from database")
	user, err := getUser(d, r)
	if err != nil {
		log.Println("Unauthorized access attempt to products")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("Fetching products for user ID: %d and their group", user.Id)

	// Get all user IDs in the same group (including current user)
	// If user is not in a group, just return their own products
	var query string
	var args []interface{}

	// Try to get group members
	rows, err := d.DB.Query(`
		SELECT DISTINCT user_id
		FROM groups
		WHERE id = (SELECT id FROM groups WHERE user_id = $1 LIMIT 1)
	`, user.Id)

	if err != nil {
		// User might not be in a group, just query their own products
		log.Printf("User %d is not in a group or error getting group members, fetching only their products", user.Id)
		query = `SELECT id, name, volume, brand, default_tags, user_id FROM products WHERE user_id=$1`
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
			// User is not in a group, fetch only their products
			log.Printf("User %d is not in a group, fetching only their products", user.Id)
			query = `SELECT id, name, volume, brand, default_tags, user_id FROM products WHERE user_id=$1`
			args = []interface{}{user.Id}
		} else {
			// Build query with IN clause for all group members
			log.Printf("User %d is in a group with %d members, fetching products for all", user.Id, len(userIDs))
			query = `SELECT id, name, volume, brand, default_tags, user_id FROM products WHERE user_id = ANY($1)`
			args = []interface{}{pq.Array(userIDs)}
		}
	}

	productRows, err := d.DB.Query(query, args...)
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
	log.Printf("Successfully fetched %d products for user %d", len(products), user.Id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"products": products})
}

func createProduct(d *Dependencies, w http.ResponseWriter, r *http.Request) {
	log.Println("Creating new product")
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("Failed to decode product JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := getUser(d, r)
	if err != nil {
		log.Println("Unauthorized access attempt to create product")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	p.UserId = user.Id
	if err := d.Validator.ValidateProduct(&p); err != nil {
		log.Printf("Product validation failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defaultTagsStr := strings.Join(p.DefaultTags, ",")
	log.Printf("Creating product for user ID: %d", user.Id)
	err = d.DB.QueryRow(`INSERT INTO products (name, volume, brand, default_tags, user_id) VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		p.Name, p.Volume, p.Brand, defaultTagsStr, user.Id).Scan(&p.Id)
	if err != nil {
		log.Printf("Failed to insert product: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully created product with ID: %d for user %d", p.Id, user.Id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func updateProduct(d *Dependencies, w http.ResponseWriter, r *http.Request) {
	log.Println("Updating product")
	user, err := getUser(d, r)
	if err != nil {
		log.Println("Unauthorized access attempt to update product")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("Failed to decode product JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if p.Id == 0 {
		log.Println("Missing id in request body")
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	p.UserId = user.Id
	if err := d.Validator.ValidateProduct(&p); err != nil {
		log.Printf("Product validation failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defaultTagsStr := strings.Join(p.DefaultTags, ",")
	log.Printf("Updating product ID: %d for user ID: %d", p.Id, user.Id)

	result, err := d.DB.Exec(`UPDATE products SET name=$1, volume=$2, brand=$3, default_tags=$4 WHERE id=$5 AND user_id=$6`,
		p.Name, p.Volume, p.Brand, defaultTagsStr, p.Id, user.Id)
	if err != nil {
		log.Printf("Failed to update product: %v", err)
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
		log.Printf("Product with ID %d not found for user %d", p.Id, user.Id)
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	log.Printf("Successfully updated product with ID: %d", p.Id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}
