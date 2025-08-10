package main

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getProducts(w, r)
	case http.MethodPost:
		s.createProduct(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) purchasesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getPurchases(w, r)
	case http.MethodPost:
		s.createPurchase(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) getProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query(`SELECT id, name, volume, brand, category, description, creation_date FROM products`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		var created time.Time
		if err := rows.Scan(&p.Id, &p.Name, &p.Volume, &p.Brand, &p.Category, &p.Description, &created); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		p.CreationDate = created.Format("2006-01-02")
		products = append(products, p)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"products": products})
}

func (s *Server) createProduct(w http.ResponseWriter, r *http.Request) {
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p.CreationDate = ""
	if err := s.validator.ValidateProduct(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	created := time.Now()
	err := s.db.QueryRow(`INSERT INTO products (name, volume, brand, category, description, creation_date) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`,
		p.Name, p.Volume, p.Brand, p.Category, p.Description, created).Scan(&p.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (s *Server) getPurchases(w http.ResponseWriter, r *http.Request) {
	uid, ok := userID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	rows, err := s.db.Query(`SELECT id, product_id, quantity, price, date, store, receipt_id FROM purchases WHERE user_id=$1`, uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	purchases := []Purchase{}
	for rows.Next() {
		var p Purchase
		var d time.Time
		if err := rows.Scan(&p.Id, &p.ProductId, &p.Quantity, &p.Price, &d, &p.Store, &p.ReceiptId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		p.Date = d.Format("2006-01-02")
		purchases = append(purchases, p)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"purchases": purchases})
}

func (s *Server) createPurchase(w http.ResponseWriter, r *http.Request) {
	var p Purchase
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.validator.ValidatePurchase(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	d, _ := time.Parse("2006-01-02", p.Date)
	uid, ok := userID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	err := s.db.QueryRow(`INSERT INTO purchases (product_id, quantity, price, date, store, receipt_id, user_id) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		p.ProductId, p.Quantity, p.Price, d, p.Store, p.ReceiptId, uid).Scan(&p.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (s *Server) registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = s.db.QueryRow(`INSERT INTO users (login, password_hash) VALUES ($1,$2) RETURNING id`, u.Login, string(hash)).Scan(&u.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, err := s.auth.GenerateToken(u.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var creds User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var id int64
	var hash string
	err := s.db.QueryRow(`SELECT id, password_hash FROM users WHERE login=$1`, creds.Login).Scan(&id, &hash)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(creds.Password)) != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	token, err := s.auth.GenerateToken(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
