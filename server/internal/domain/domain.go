package domain

import "time"

type (
	InviteId      int64
	GroupId       int64
	GroupMemberId int64
	UserId        int64
	ProductId     int64
	PurchaseId    int64
	ReceiptId     int64
)

type Product struct {
	Id          ProductId `json:"id"`
	Name        string    `json:"name"`
	Volume      string    `json:"volume"`
	Brand       string    `json:"brand"`
	DefaultTags []string  `json:"default_tags"`
	UserId      UserId    `json:"user_id"`
}

type Purchase struct {
	Id        PurchaseId `json:"id"`
	ProductId ProductId  `json:"product_id"`
	Quantity  int        `json:"quantity"`
	Price     int        `json:"price"`
	Date      time.Time  `json:"date"`
	Store     string     `json:"store"`
	Tags      []string   `json:"tags"`
	ReceiptId ReceiptId  `json:"receipt_id"`
	UserId    UserId     `json:"user_id"`
}

type User struct {
	Id       UserId `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
}

type Group struct {
	Id      GroupId       `json:"id"`
	Members []GroupMember `json:"members"`
}

type GroupMember struct {
	GroupId      GroupId `json:"group_id"`
	UserId       UserId  `json:"user_id"`
	Login        string  `json:"login"`
	MemberNumber int     `json:"member_number"`
}

type Invite struct {
	Id         InviteId  `json:"id"`
	FromUserId UserId    `json:"from_user_id"`
	ToUserId   UserId    `json:"to_user_id"`
	FromLogin  string    `json:"from_login"`
	ToLogin    string    `json:"to_login"`
	CreatedAt  time.Time `json:"created_at"`
}
