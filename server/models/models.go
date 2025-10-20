package models

import "time"

type Product struct {
	Id          int64    `json:"id"`
	Name        string   `json:"name"`
	Volume      string   `json:"volume"`
	Brand       string   `json:"brand"`
	DefaultTags []string `json:"default_tags"`
	UserId      int64    `json:"user_id"`
}

type Purchase struct {
	Id        int64     `json:"id"`
	ProductId int64     `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     int       `json:"price"`
	Date      time.Time `json:"date"`
	Store     string    `json:"store"`
	Tags      []string  `json:"tags"`
	ReceiptId int64     `json:"receipt_id"`
	UserId    int64     `json:"user_id"`
}

type User struct {
	Id       int64  `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
}