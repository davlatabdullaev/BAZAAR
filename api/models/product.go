package models

import "time"

type Product struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Price      float64    `json:"price"`
	Barcode    string    `json:"barcode"`
	CategoryID string    `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}

type CreateProduct struct {
	Name       string `json:"name"`
	Price      string `json:"price"`
	CategoryID string `json:"category_id"`
}

type UpdateProduct struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Price      string `json:"price"`
	CategoryID string `json:"category_id"`
}

type ProductsResponse struct {
	Products []Product `json:"products"`
	Count    int       `json:"count"`
}

type ProductGetListRequest struct {
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	Search  string `json:"Search"`
	Barcode int    `json:"barcode"`
}
