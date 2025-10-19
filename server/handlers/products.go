package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

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
	log.Printf("Fetching products for user ID: %d", uid)
	rows, err := deps.DB.Query(`SELECT id, name, volume, brand, default_tags, user_id FROM products WHERE user_id=$1`, uid)
	if err != nil {
		log.Printf("Failed to query products: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var p models.Product
		var defaultTagsStr string
		if err := rows.Scan(&p.Id, &p.Name, &p.Volume, &p.Brand, &defaultTagsStr, &p.UserId); err != nil {
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