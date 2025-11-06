package validators_test

import (
	"testing"
	"time"
	
	"yuki_buy_log/models"
	"yuki_buy_log/validators"
)

func TestValidateProduct(t *testing.T) {
	v := validators.NewValidator()
	// Valid product with English letters
	if err := v.ValidateProduct(&models.Product{Name: "Tea", Volume: "500ml", Brand: "Brand1", DefaultTags: []string{"tag1"}, UserId: 1}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Valid product with digits (now allowed)
	if err := v.ValidateProduct(&models.Product{Name: "Tea1", Volume: "500ml", Brand: "Brand1", DefaultTags: []string{}, UserId: 1}); err != nil {
		t.Fatalf("unexpected error for name with digits: %v", err)
	}
	// Valid product with Russian letters
	if err := v.ValidateProduct(&models.Product{Name: "Чай", Volume: "500ml", Brand: "Бренд 123", DefaultTags: []string{"тег1"}, UserId: 1}); err != nil {
		t.Fatalf("unexpected error for Russian letters: %v", err)
	}
	// Invalid product with special characters
	if err := v.ValidateProduct(&models.Product{Name: "Tea@123", Volume: "500ml", Brand: "Brand1", DefaultTags: []string{}, UserId: 1}); err == nil {
		t.Fatalf("expected error for name with special characters")
	}
	// Invalid product with empty name
	if err := v.ValidateProduct(&models.Product{Name: "", Volume: "500ml", Brand: "Brand1", DefaultTags: []string{}, UserId: 1}); err == nil {
		t.Fatalf("expected error for empty name")
	}
}

func TestValidatePurchase(t *testing.T) {
	v := validators.NewValidator()
	validDate := time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)
	// Valid purchase with English letters
	if err := v.ValidatePurchase(&models.Purchase{ProductId: 1, Quantity: 2, Price: 100, Date: validDate, Store: "Store", ReceiptId: 1, UserId: 1}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Valid purchase with Russian letters and digits
	if err := v.ValidatePurchase(&models.Purchase{ProductId: 1, Quantity: 2, Price: 100, Date: validDate, Store: "Магазин 5", Tags: []string{"тег1", "sale"}, ReceiptId: 1, UserId: 1}); err != nil {
		t.Fatalf("unexpected error for Russian letters: %v", err)
	}
	// Invalid purchase with special characters in store name
	if err := v.ValidatePurchase(&models.Purchase{ProductId: 1, Quantity: 2, Price: 100, Date: validDate, Store: "Store@123", ReceiptId: 1, UserId: 1}); err == nil {
		t.Fatalf("expected error for store with special characters")
	}
	// Invalid purchase with missing date
	if err := v.ValidatePurchase(&models.Purchase{ProductId: 1, Quantity: 2, Price: 100, Date: time.Time{}, Store: "Store", ReceiptId: 1, UserId: 1}); err == nil {
		t.Fatalf("expected error for missing date")
	}
	// Invalid purchase with special characters in tag
	if err := v.ValidatePurchase(&models.Purchase{ProductId: 1, Quantity: 2, Price: 100, Date: validDate, Store: "Store", Tags: []string{"tag@#$"}, ReceiptId: 1, UserId: 1}); err == nil {
		t.Fatalf("expected error for tag with special characters")
	}
}