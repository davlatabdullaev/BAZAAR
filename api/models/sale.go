package models

import "time"

type Sale struct {
	ID              string    `json:"id"`
	BranchID        string    `json:"branch_id"`
	ShopAssistantID string    `json:"shop_assistent_id"`
	CashierID       string    `json:"cashier_id"`
	PaymentType     string    `json:"payment_type"`
	Price           string    `json:"price"`
	Status          string    `json:"status"`
	ClientName      string    `json:"client_name"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       time.Time `json:"deleted_at"`
}

type CreateSale struct {
	BranchID        string `json:"branch_id"`
	ShopAssistantID string `json:"shop_assistent_id"`
	CashierID       string `json:"cashier_id"`
	PaymentType     string `json:"payment_type"`
	Price           string `json:"price"`
	Status          string `json:"status"`
	ClientName      string `json:"client_name"`
}

type UpdateSale struct {
	ID              string `json:"-"`
	BranchID        string `json:"branch_id"`
	ShopAssistantID string `json:"shop_assistent_id"`
	CashierID       string `json:"cashier_id"`
	PaymentType     string `json:"payment_type"`
	Price           string `json:"price"`
	Status          string `json:"status"`
	ClientName      string `json:"client_name"`
}

type SalesResponse struct {
	Sales []Sale `json:"sales"`
	Count int    `json:"count"`
}

type SaleRequest struct {
	ID         string `json:"id"`
	TotalPrice int    `json:"-"`
	Status     string `json:"status"`
}
