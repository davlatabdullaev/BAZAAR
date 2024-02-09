package models

import "time"

type Basket struct {
	ID        string    `json:"id"`
	SaleID    string    `json:"sale_id"`
	ProductID string    `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     string    `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type CreateBasket struct {
	SaleID    string `json:"sale_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Price     string `json:"price"`
}

type UpdateBasket struct {
	ID        string    `json:"id"`
	SaleID    string    `json:"sale_id"`
	ProductID string    `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     string    `json:"price"`
}

type BasketsResponse struct {
	Baskets []Basket `json:"baskets"`
	Count   int      `json:"count"`
}
