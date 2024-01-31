package models

type Sale struct {
	ID              string `json:"id"`
	BranchID        string `json:"branch_id"`
	ShopAssistantID string `json:"shop_assistant_id"`
	CashierID       string `json:"cashier_id"`
	PaymentType     string `json:"payment_type"`
	Price           string `json:"price"`
	Status          string `json:"status"`
	ClientName      string `json:"client_name"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	DeletedAt       string `json:"deleted_at"`
}

type CreateSale struct {
	BranchID        string `json:"branch_id"`
	ShopAssistantID string `json:"shop_assistant_id"`
	CashierID       string `json:"cashier_id"`
	PaymentType     string `json:"payment_type"`
	Price           string `json:"price"`
	Status          string `json:"status"`
	ClientName      string `json:"client_name"`
}

type UpdateSale struct {
	BranchID        string `json:"branch_id"`
	ShopAssistantID string `json:"shop_assistant_id"`
	CashierID       string `json:"cashier_id"`
	PaymentType     string `json:"payment_type"`
	Price           string `json:"price"`
	Status          string `json:"status"`
	ClientName      string `json:"client_name"`
	UpdatedAt       string `json:"updated_at"`
}

type SalesResponse struct {
	Sales []Sale `json:"sales"`
	Count int    `json:"count"`
}
