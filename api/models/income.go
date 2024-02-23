package models

import "time"

type Income struct {
	ID        string    `json:"id"`
	BranchID  string    `json:"branch_id"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type CreateIncome struct {
	BranchID string  `json:"branch_id"`
	Price    float64 `json:"price"`
}

type UpdateIncome struct {
	ID       string  `json:"-"`
	BranchID string  `json:"branch_id"`
	Price    float64 `json:"price"`
}

type IncomesResponse struct {
	Incomes []Income `json:"incomes"`
	Count   int      `json:"count"`
}
