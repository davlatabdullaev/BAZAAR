package models

import "time"

type Transaction struct {
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

type CreateTransaction struct {
	SaleID          string    `json:"sale_id"`
	StaffID         string    `json:"staff_id"`
	TransactionType string    `json:"transaction_type"`
	SourceType      string    `json:"source_type"`
	Amount          float64   `json:"amount"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
}

type UpdateTransaction struct {
	ID              string    `json:"id"`
	SaleID          string    `json:"sale_id"`
	StaffID         string    `json:"staff_id"`
	TransactionType string    `json:"transaction_type"`
	SourceType      string    `json:"source_type"`
	Amount          float64   `json:"amount"`
	Description     string    `json:"description"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type TransactionsResponse struct {
	Transactions []Transaction `json:"transaction"`
	Count        int           `json:"count"`
}
