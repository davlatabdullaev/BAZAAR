package models

import "time"

type Sale struct {
	ID              string    `json:"id"`
	BranchID        string    `json:"branch_id"`
	ShopAssistantID string    `json:"shop_assistant_id"`
	CashierID       string    `json:"cashier_id"`
	PaymentType     string    `json:"payment_type"`
	Price           string    `json:"price"`
	Status          string    `json:"status"`
	ClientName      string    `json:"client_name"`
	CreatedAt       time.Time `json:"created_at"`
	Updated         time.Time `json:"updated_at"`
	DeletedAt       time.Time `json:"deleted_at"`
}

type CreateSale struct {
	BranchID        string `json:"branch_id"`
	ShopAssistantID string `json:"shop_assistant_id"`
	CashierID       string `json:"cashier_id"`
	PaymentType     string `json:"payment_type"`
	Price           string `json:"price"`
	Status          string `json:"status"`
	ClientName      string `json:"client_name"`
}

type UpdateSale struct {
	ID              string    `json:"id"`
	BranchID        string    `json:"branch_id"`
	ShopAssistantID string    `json:"shop_assistant_id"`
	CashierID       string    `json:"cashier_id"`
	PaymentType     string    `json:"payment_type"`
	Price           string    `json:"price"`
	Status          string    `json:"status"`
	ClientName      string    `json:"client_name"`
	Updated         time.Time `json:"updated_at"`
}

type SalesResponse struct {
	Sales []Sale `json:"sales"`
	Count int    `json:"count"`
}
