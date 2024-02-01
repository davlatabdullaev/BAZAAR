package models

type Basket struct {
	ID        string `json:"id"`
	SaleID    string `json:"sale_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Price     string `json:"price"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type CreateBasket struct {
	SaleID    string `json:"sale_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Price     string `json:"price"`
}

type UpdateBasket struct {
	ID        string `json:"id"`
	SaleID    string `json:"sale_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Price     string `json:"price"`
	UpdatedAt string `json:"updated_at"`
}

type BasketsResponse struct {
	Baskets []Basket `json:"baskets"`
	Count   int      `json:"count"`
}
