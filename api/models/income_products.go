package models

import "time"

type IncomeProduct struct {
	ID        string    `json:"id"`
	IncomeID  string    `json:"income_id"`
	ProductID string    `json:"product_id"`
	Price     float64   `json:"price"`
	Count     int       `json:"count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type CreateIncomeProduct struct {
	IncomeID  string  `json:"income_id"`
	ProductID string  `json:"product_id"`
	Price     float64 `json:"price"`
	Count     int       `json:"count"`
}

type UpdateIncomeProduct struct {
	ID        string  `json:"-"`
	IncomeID  string  `json:"income_id"`
	ProductID string  `json:"product_id"`
	Price     float64 `json:"price"`
	Count     int       `json:"count"`
}

type IncomeProductsResponse struct {
	IncomeProducts []IncomeProduct `json:"income_products"`
	Count          int             `json:"count"`
}
