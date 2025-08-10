package main

// Product represents an item or service that can be purchased.
type Product struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	Volume       string `json:"volume"`
	Brand        string `json:"brand"`
	Category     string `json:"category"`
	Description  string `json:"description"`
	CreationDate string `json:"creation_date,omitempty"`
}

// Purchase represents a product bought at a specific price and time.
type Purchase struct {
	Id        int64  `json:"id"`
	ProductId int64  `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Price     int    `json:"price"`
	Date      string `json:"date"`
	Store     string `json:"store"`
	ReceiptId int64  `json:"receipt_id"`
}
