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

type UserId = int64
type User struct {
	Id       UserId `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
}

type GroupMember struct {
	GroupId      int64  `json:"group_id"`
	UserId       int64  `json:"user_id"`
	Login        string `json:"login"`
	MemberNumber int    `json:"member_number"`
}

type Invite struct {
	Id         int64     `json:"id"`
	FromUserId int64     `json:"from_user_id"`
	ToUserId   int64     `json:"to_user_id"`
	FromLogin  string    `json:"from_login"`
	ToLogin    string    `json:"to_login"`
	CreatedAt  time.Time `json:"created_at"`
}
