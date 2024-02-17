package models

import "time"

type Transactions struct {
	ID              string    `json:"id"`
	SaleID          string    `json:"sale_id"`
	StaffID         string    `json:"staff_id"`
	TransactionType string    `json:"transaction_type"`
	SourceType      string    `json:"source_type"`
	Amount          float64   `json:"amount"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       time.Time `json:"deleted_at"`
}

type CreateTransactions struct {
	SaleID          string  `json:"sale_id"`
	StaffID         string  `json:"staff_id"`
	TransactionType string  `json:"transaction_type"`
	SourceType      string  `json:"source_type"`
	Amount          float64 `json:"amount"`
	Description     string  `json:"description"`
}

type UpdateTransactions struct {
	ID              string  `json:"id"`
	SaleID          string  `json:"sale_id"`
	StaffID         string  `json:"staff_id"`
	TransactionType string  `json:"transaction_type"`
	SourceType      string  `json:"source_type"`
	Amount          float64 `json:"amount"`
	Description     string  `json:"description"`
}

type GetListTransactionsRequest struct {
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
	FromAmount float64 `json:"from_amount"`
	ToAmount   float64 `json:"to_amount"`
}

type TransactionsResponse struct {
	Transactions []Transactions `json:"transaction"`
	Count        int            `json:"count"`
}

type UpdateStaffBalanceAndCreateTransaction struct {
	UpdateCashierBalance       StaffInfo
	UpdateShopAssistantBalance StaffInfo
	SaleID                     string
	StaffID                    string
	TransactionType            string
	SourceType                 string
	Amount                     string
	Description                string
}

type StaffInfo struct {
	StaffID string
	Amount  float64
}
