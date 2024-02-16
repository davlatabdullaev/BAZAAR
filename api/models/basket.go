package models

import "time"

type Basket struct {
	ID        string    `json:"id"`
	SaleID    string    `json:"sale_id"`
	ProductID string    `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type CreateBasket struct {
	SaleID    string  `json:"sale_id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type UpdateBasket struct {
	ID        string  `json:"id"`
	SaleID    string  `json:"sale_id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type BasketsResponse struct {
	Baskets []Basket `json:"baskets"`
	Count   int      `json:"count"`
}

type GetBasketsListRequest struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type UpdateBasketQuantity struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}
