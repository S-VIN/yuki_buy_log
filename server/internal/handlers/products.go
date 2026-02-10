package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"yuki_buy_log/internal/domain"
	"yuki_buy_log/internal/stores"
	"yuki_buy_log/internal/validators"
)

func ProductsHandler(auth Authenticator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Products handler called: %s %s", r.Method, r.URL.Path)
		switch r.Method {
		case http.MethodGet:
			getProducts(w, r)
		case http.MethodPost:
			createProduct(w, r)
		case http.MethodPut:
			updateProduct(w, r)
		default:
			log.Printf("Method not allowed for products: %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching products from store")
	user, err := getUser(r)
	if err != nil {
		log.Println("Unauthorized access attempt to products")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("Fetching products for user ID: %d and their group", user.Id)

	// Get all user IDs in the same group (including current user)
	// If user is not in a group, just return their own products
	var products []domain.Product

	// Try to get group
	groupStore := stores.GetGroupStore()
	group := groupStore.GetGroupByUserId(user.Id)

	productStore := stores.GetProductStore()
	if group == nil {
		// User is not in a group, fetch only their products
		log.Printf("User %d is not in a group, fetching only their products", user.Id)
		products = productStore.GetProductsByUserId(user.Id)
	} else {
		// User is in a group, fetch products for all group members
		log.Printf("User %d is in a group with %d members, fetching products for all", user.Id, len(group.Members))

		userIds := make([]domain.UserId, len(group.Members))
		for i, member := range group.Members {
			userIds[i] = member.UserId
		}

		products = productStore.GetProductsByUserIds(userIds)
	}

	if products == nil {
		products = []domain.Product{}
	}

	log.Printf("Successfully fetched %d products for user %d", len(products), user.Id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"products": products})
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating new product")
	var p domain.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("Failed to decode product JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := getUser(r)
	if err != nil {
		log.Println("Unauthorized access attempt to create product")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	p.UserId = user.Id
	if err := validators.ValidateProduct(&p); err != nil {
		log.Printf("Product validation failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Creating product for user ID: %d", user.Id)
	productStore := stores.GetProductStore()
	err = productStore.CreateProduct(&p)
	if err != nil {
		log.Printf("Failed to create product: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully created product with ID: %d for user %d", p.Id, user.Id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating product")
	user, err := getUser(r)
	if err != nil {
		log.Println("Unauthorized access attempt to update product")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var p domain.Product
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
	if err := validators.ValidateProduct(&p); err != nil {
		log.Printf("Product validation failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Updating product ID: %d for user ID: %d", p.Id, user.Id)
	productStore := stores.GetProductStore()
	err = productStore.UpdateProduct(&p)
	if err != nil {
		log.Printf("Failed to update product: %v", err)
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "product not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully updated product with ID: %d", p.Id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}
