package validators

import (
	"testing"
	"time"
	
	"yuki_buy_log/models"
)

func TestValidateProduct(t *testing.T) {
	v := NewValidator()
	if err := v.ValidateProduct(&models.Product{Name: "Tea", Volume: "500ml", Brand: "Brand1", DefaultTags: []string{"tag1"}, UserId: 1}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := v.ValidateProduct(&models.Product{Name: "Tea1", Volume: "500ml", Brand: "Brand1", DefaultTags: []string{}, UserId: 1}); err == nil {
		t.Fatalf("expected error for invalid name")
	}
}

func TestValidatePurchase(t *testing.T) {
	v := NewValidator()
	validDate := time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)
	if err := v.ValidatePurchase(&models.Purchase{ProductId: 1, Quantity: 2, Price: 100, Date: validDate, Store: "Store", ReceiptId: 1, UserId: 1}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := v.ValidatePurchase(&models.Purchase{ProductId: 1, Quantity: 2, Price: 100, Date: time.Time{}, Store: "Store", ReceiptId: 1, UserId: 1}); err == nil {
		t.Fatalf("expected error for missing date")
	}
}