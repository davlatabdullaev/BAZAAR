package models

type Product struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Price      string `json:"price"`
	Barcode    string `json:"barcode"`
	CategoryID string `json:"category_id"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	DeletedAt  string `json:"deleted_at"`
}

type CreateProduct struct {
	Name       string `json:"name"`
	Price      string `json:"price"`
	Barcode    string `json:"barcode"`
	CategoryID string `json:"category_id"`
}

type UpdateProduct struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Price      string `json:"price"`
	Barcode    string `json:"barcode"`
	CategoryID string `json:"category_id"`
	UpdatedAt  string `json:"updated_at"`
}

type ProductsResponse struct {
	Products []Product `json:"products"`
	Count    int       `json:"count"`
}
