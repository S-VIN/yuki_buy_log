package main

import "testing"

func TestValidateProduct(t *testing.T) {
	v := NewValidator()
	if err := v.ValidateProduct(&Product{Name: "Tea", Volume: "500ml", Brand: "Brand1", Category: "Drink", Description: "Green"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := v.ValidateProduct(&Product{Name: "Tea1", Volume: "500ml", Brand: "Brand1", Category: "Drink"}); err == nil {
		t.Fatalf("expected error for invalid name")
	}
}

func TestValidatePurchase(t *testing.T) {
	v := NewValidator()
	if err := v.ValidatePurchase(&Purchase{ProductId: 1, Quantity: 2, Price: 100, Date: "2023-03-01", Store: "Store", ReceiptId: 1}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := v.ValidatePurchase(&Purchase{ProductId: 1, Quantity: 2, Price: 100, Date: "", Store: "Store", ReceiptId: 1}); err == nil {
		t.Fatalf("expected error for missing date")
	}
}
