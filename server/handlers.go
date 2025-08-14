package main

import (
	"database/sql"
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
	rows, err := s.db.Query(`SELECT id, name, volume, brand, category, description, creation_date, user_login FROM products`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		var created time.Time
		if err := rows.Scan(&p.Id, &p.Name, &p.Volume, &p.Brand, &p.Category, &p.Description, &created, &p.Login); err != nil {
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
	login, ok := userLogin(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	created := time.Now()
	err := s.db.QueryRow(`INSERT INTO products (name, volume, brand, category, description, creation_date, user_login) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		p.Name, p.Volume, p.Brand, p.Category, p.Description, created, login).Scan(&p.Id)
	p.Login = login
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (s *Server) getPurchases(w http.ResponseWriter, r *http.Request) {
	login, ok := userLogin(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var familyID int64
	err := s.db.QueryRow(`SELECT id FROM family WHERE user_login=$1`, login).Scan(&familyID)
	var rows *sql.Rows
	if err == sql.ErrNoRows {
		rows, err = s.db.Query(`SELECT id, product_id, quantity, price, date, store, receipt_id, user_login FROM purchases WHERE user_login=$1`, login)
	} else if err == nil {
		rows, err = s.db.Query(`SELECT id, product_id, quantity, price, date, store, receipt_id, user_login FROM purchases WHERE user_login IN (SELECT user_login FROM family WHERE id=$1)`, familyID)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	purchases := []Purchase{}
	for rows.Next() {
		var p Purchase
		var d time.Time
		if err := rows.Scan(&p.Id, &p.ProductId, &p.Quantity, &p.Price, &d, &p.Store, &p.ReceiptId, &p.Login); err != nil {
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
	login, ok := userLogin(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	err := s.db.QueryRow(`INSERT INTO purchases (product_id, quantity, price, date, store, receipt_id, user_login) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		p.ProductId, p.Quantity, p.Price, d, p.Store, p.ReceiptId, login).Scan(&p.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	p.Login = login
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
	if _, err = s.db.Exec(`INSERT INTO users (login, password_hash) VALUES ($1,$2)`, u.Login, string(hash)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, err := s.auth.GenerateToken(u.Login)
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
	var hash string
	err := s.db.QueryRow(`SELECT password_hash FROM users WHERE login=$1`, creds.Login).Scan(&hash)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(creds.Password)) != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	token, err := s.auth.GenerateToken(creds.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (s *Server) familyMembersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	login, ok := userLogin(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var familyID int64
	err := s.db.QueryRow(`SELECT id FROM family WHERE user_login=$1`, login).Scan(&familyID)
	if err == sql.ErrNoRows {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string][]string{"members": {}})
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rows, err := s.db.Query(`SELECT user_login FROM family WHERE id=$1`, familyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	members := []string{}
	for rows.Next() {
		var login string
		if err := rows.Scan(&login); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		members = append(members, login)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{"members": members})
}

func (s *Server) familyInviteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	login, ok := userLogin(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req InviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.validator.ValidateLogin(req.Login); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var tmp int64
	err := s.db.QueryRow(`SELECT 1 FROM users WHERE login=$1`, req.Login).Scan(&tmp)
	if err == sql.ErrNoRows {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if req.Login == login {
		http.Error(w, "cannot invite yourself", http.StatusBadRequest)
		return
	}
	var familyID int64
	err = s.db.QueryRow(`SELECT id FROM family WHERE user_login=$1`, login).Scan(&familyID)
	if err == sql.ErrNoRows {
		if err = s.db.QueryRow(`SELECT nextval('family_id_seq')`).Scan(&familyID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err = s.db.Exec(`INSERT INTO family (id, user_login) VALUES ($1,$2)`, familyID, login); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = s.db.QueryRow(`SELECT id FROM family WHERE user_login=$1`, req.Login).Scan(&tmp)
	if err != sql.ErrNoRows {
		http.Error(w, "user already in family", http.StatusBadRequest)
		return
	}
	err = s.db.QueryRow(`SELECT 1 FROM family_invitations WHERE family_id=$1 AND invitee_login=$2`, familyID, req.Login).Scan(&tmp)
	if err != sql.ErrNoRows {
		http.Error(w, "invitation already sent", http.StatusBadRequest)
		return
	}
	var count int
	if err = s.db.QueryRow(`SELECT COUNT(*) FROM family WHERE id=$1`, familyID).Scan(&count); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count >= 5 {
		http.Error(w, "family full", http.StatusBadRequest)
		return
	}
	if _, err = s.db.Exec(`INSERT INTO family_invitations (family_id, inviter_login, invitee_login) VALUES ($1,$2,$3)`, familyID, login, req.Login); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "invited"})
}

func (s *Server) familyRespondHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	login, ok := userLogin(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req InvitationResponse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.validator.ValidateLogin(req.Login); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var invitationID, familyID int64
	err := s.db.QueryRow(`SELECT id, family_id FROM family_invitations WHERE inviter_login=$1 AND invitee_login=$2`, req.Login, login).Scan(&invitationID, &familyID)
	if err != nil {
		http.Error(w, "invitation not found", http.StatusNotFound)
		return
	}
	if req.Accept {
		var tmp int64
		if err = s.db.QueryRow(`SELECT id FROM family WHERE user_login=$1`, login).Scan(&tmp); err != sql.ErrNoRows {
			http.Error(w, "already in family", http.StatusBadRequest)
			return
		}
		var count int
		if err = s.db.QueryRow(`SELECT COUNT(*) FROM family WHERE id=$1`, familyID).Scan(&count); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if count >= 5 {
			http.Error(w, "family full", http.StatusBadRequest)
			return
		}
		if _, err = s.db.Exec(`INSERT INTO family (id, user_login) VALUES ($1,$2)`, familyID, login); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if _, err = s.db.Exec(`DELETE FROM family_invitations WHERE id=$1`, invitationID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	status := "declined"
	if req.Accept {
		status = "accepted"
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}

func (s *Server) familyLeaveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	login, ok := userLogin(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var familyID int64
	err := s.db.QueryRow(`SELECT id FROM family WHERE user_login=$1`, login).Scan(&familyID)
	if err == sql.ErrNoRows {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err = s.db.Exec(`DELETE FROM family WHERE user_login=$1`, login); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var count int
	if err = s.db.QueryRow(`SELECT COUNT(*) FROM family WHERE id=$1`, familyID).Scan(&count); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count == 1 {
		if _, err = s.db.Exec(`DELETE FROM family WHERE id=$1`, familyID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
