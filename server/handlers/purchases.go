package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"yuki_buy_log/models"
	"yuki_buy_log/stores"
	"yuki_buy_log/validators"
)

func PurchasesHandler(auth Authenticator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Purchases handler called: %s %s", r.Method, r.URL.Path)
		switch r.Method {
		case http.MethodGet:
			getPurchases(w, r)
		case http.MethodPost:
			createPurchase(w, r)
		case http.MethodDelete:
			deletePurchase(w, r)
		default:
			log.Printf("Method not allowed for purchases: %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func getPurchases(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching purchases from store")
	user, err := getUser(r)
	if err != nil {
		log.Println("Unauthorized access attempt to purchases")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("Fetching purchases for user ID: %d and their group", user.Id)

	// Get group store and purchase store
	groupStore := stores.GetGroupStore()
	purchaseStore := stores.GetPurchaseStore()

	// Get all user IDs in the same group (including current user)
	// If user is not in a group, just return their own purchases
	var userIds []models.UserId
	group := groupStore.GetGroupByUserId(user.Id)

	if group.Members == nil {
		// User is not in a group, fetch only their purchases
		log.Printf("User %d is not in a group, fetching only their purchases", user.Id)
		userIds = []models.UserId{user.Id}
	} else {
		// User is in a group, get all group member IDs
		log.Printf("User %d is in a group with %d members, fetching purchases for all", user.Id, len(group.Members))
		userIds = make([]models.UserId, len(group.Members))
		for i, member := range group.Members {
			userIds[i] = member.UserId
		}
	}

	// Get purchases from store
	purchases := purchaseStore.GetPurchasesByUserIds(userIds)
	log.Printf("Successfully fetched %d purchases for user %d", len(purchases), user.Id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"purchases": purchases})
}

func createPurchase(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating new purchase")
	var p models.Purchase
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("Failed to decode purchase JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := getUser(r)
	if err != nil {
		log.Println("Unauthorized access attempt to create purchase")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	p.UserId = user.Id
	if err := validators.ValidatePurchase(&p); err != nil {
		log.Printf("Purchase validation failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Creating purchase for user ID: %d", user.Id)

	// Get purchase store and add purchase
	purchaseStore := stores.GetPurchaseStore()
	err = purchaseStore.AddPurchase(&p)
	if err != nil {
		log.Printf("Failed to insert purchase: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully created purchase with ID: %d", p.Id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func deletePurchase(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting purchase")
	user, err := getUser(r)
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

	// Get purchase store and delete purchase
	purchaseStore := stores.GetPurchaseStore()
	err = purchaseStore.DeletePurchase(models.PurchaseId(req.Id), user.Id)
	if err != nil {
		log.Printf("Failed to delete purchase: %v", err)
		http.Error(w, "purchase not found", http.StatusNotFound)
		return
	}

	log.Printf("Successfully deleted purchase with ID: %d", req.Id)
	w.WriteHeader(http.StatusNoContent)
}
