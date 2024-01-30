package models

type StorageTransaction struct {
	ID                     string  `json:"id"`
	StaffID                string  `json:"staff_id"`
	ProductID              string  `json:"product_id"`
	StorageTransactionType string  `json:"storage_transaction_type"`
	Price                  float64 `json:"price"`
	Quantity               float64 `json:"quantity"`
	CreatedAt              string  `json:"created_at"`
	UpdatedAt              string  `json:"updated_at"`
	DeletedAt              string  `json:"deleted_at"`
}

type CreateStorageTransaction struct {
	StaffID                string  `json:"staff_id"`
	ProductID              string  `json:"product_id"`
	StorageTransactionType string  `json:"storage_transaction_type"`
	Price                  float64 `json:"price"`
	Quantity               float64 `json:"quantity"`
	CreatedAt              string  `json:"created_at"`
}

type UpdateStorageTransaction struct {
	ID                     string  `json:"id"`
	StaffID                string  `json:"staff_id"`
	ProductID              string  `json:"product_id"`
	StorageTransactionType string  `json:"storage_transaction_type"`
	Price                  float64 `json:"price"`
	Quantity               float64 `json:"quantity"`
	UpdatedAt              string  `json:"updated_at"`
}

type StorageTransactionsResponse struct {
	StorageTransactions []StorageTransaction `json:"storage_transactions"`
	Count               int                  `json:"count"`
}
