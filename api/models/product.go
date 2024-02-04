package models

import "time"

type Product struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Price      string    `json:"price"`
	Barcode    int64       `json:"barcode"`
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
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Price      string    `json:"price"`
	Barcode    int64       `json:"barcode"`
	CategoryID string    `json:"category_id"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ProductsResponse struct {
	Products []Product `json:"products"`
	Count    int       `json:"count"`
}
