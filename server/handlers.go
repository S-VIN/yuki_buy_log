package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) productsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Products handler called: %s %s", r.Method, r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		s.getProducts(w, r)
	case http.MethodPost:
		s.createProduct(w, r)
	default:
		log.Printf("Method not allowed for products: %s", r.Method)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) purchasesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Purchases handler called: %s %s", r.Method, r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		s.getPurchases(w, r)
	case http.MethodPost:
		s.createPurchase(w, r)
	default:
		log.Printf("Method not allowed for purchases: %s", r.Method)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) getProducts(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching products from database")
	rows, err := s.db.Query(`SELECT id, name, volume, brand, category, description, creation_date FROM products`)
	if err != nil {
		log.Printf("Failed to query products: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		var created time.Time
		if err := rows.Scan(&p.Id, &p.Name, &p.Volume, &p.Brand, &p.Category, &p.Description, &created); err != nil {
			log.Printf("Failed to scan product row: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		p.CreationDate = created.Format("2006-01-02")
		products = append(products, p)
	}
	log.Printf("Successfully fetched %d products", len(products))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"products": products})
}

func (s *Server) createProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating new product")
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("Failed to decode product JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p.CreationDate = ""
	if err := s.validator.ValidateProduct(&p); err != nil {
		log.Printf("Product validation failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	created := time.Now()
	err := s.db.QueryRow(`INSERT INTO products (name, volume, brand, category, description, creation_date) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`,
		p.Name, p.Volume, p.Brand, p.Category, p.Description, created).Scan(&p.Id)
	if err != nil {
		log.Printf("Failed to insert product: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully created product with ID: %d", p.Id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (s *Server) getPurchases(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching purchases from database")
	uid, ok := userID(r)
	if !ok {
		log.Println("Unauthorized access attempt to purchases")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("Fetching purchases for user ID: %d", uid)
	rows, err := s.db.Query(`SELECT id, product_id, quantity, price, date, store, receipt_id FROM purchases WHERE user_id=$1`, uid)
	if err != nil {
		log.Printf("Failed to query purchases: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	purchases := []Purchase{}
	for rows.Next() {
		var p Purchase
		var d time.Time
		if err := rows.Scan(&p.Id, &p.ProductId, &p.Quantity, &p.Price, &d, &p.Store, &p.ReceiptId); err != nil {
			log.Printf("Failed to scan purchase row: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		p.Date = d.Format("2006-01-02")
		purchases = append(purchases, p)
	}
	log.Printf("Successfully fetched %d purchases for user %d", len(purchases), uid)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"purchases": purchases})
}

func (s *Server) createPurchase(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating new purchase")
	var p Purchase
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("Failed to decode purchase JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.validator.ValidatePurchase(&p); err != nil {
		log.Printf("Purchase validation failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	d, _ := time.Parse("2006-01-02", p.Date)
	uid, ok := userID(r)
	if !ok {
		log.Println("Unauthorized access attempt to create purchase")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("Creating purchase for user ID: %d", uid)
	err := s.db.QueryRow(`INSERT INTO purchases (product_id, quantity, price, date, store, receipt_id, user_id) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		p.ProductId, p.Quantity, p.Price, d, p.Store, p.ReceiptId, uid).Scan(&p.Id)
	if err != nil {
		log.Printf("Failed to insert purchase: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully created purchase with ID: %d", p.Id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (s *Server) registerHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Register handler called: %s %s", r.Method, r.URL.Path)
	if r.Method != http.MethodPost {
		log.Printf("Method not allowed for register: %s", r.Method)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Printf("Failed to decode user JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Registering new user: %s", u.Login)
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = s.db.QueryRow(`INSERT INTO users (login, password_hash) VALUES ($1,$2) RETURNING id`, u.Login, string(hash)).Scan(&u.Id)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, err := s.auth.GenerateToken(u.Id)
	if err != nil {
		log.Printf("Failed to generate token for user %d: %v", u.Id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully registered user %s with ID: %d", u.Login, u.Id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Login handler called: %s %s", r.Method, r.URL.Path)
	if r.Method != http.MethodPost {
		log.Printf("Method not allowed for login: %s", r.Method)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var creds User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		log.Printf("Failed to decode login JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Login attempt for user: %s", creds.Login)
	var id int64
	var hash string
	err := s.db.QueryRow(`SELECT id, password_hash FROM users WHERE login=$1`, creds.Login).Scan(&id, &hash)
	if err != nil {
		log.Printf("User not found or database error for login %s: %v", creds.Login, err)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(creds.Password)) != nil {
		log.Printf("Invalid password for user: %s", creds.Login)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	token, err := s.auth.GenerateToken(id)
	if err != nil {
		log.Printf("Failed to generate token for user %d: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully logged in user %s with ID: %d", creds.Login, id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
