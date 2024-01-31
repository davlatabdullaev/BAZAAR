package models

type Storage struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	BranchID  string `json:"branch_id"`
	Count     int    `json:"count"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type CreateStorage struct {
	ProductID string `json:"product_id"`
	BranchID  string `json:"branch_id"`
	Count     int    `json:"count"`
}

type UpdateStorage struct {
	ProductID string `json:"product_id"`
	BranchID  string `json:"branch_id"`
	Count     int    `json:"count"`
	UpdatedAt string `json:"updated_at"`
}

type StoragesResponse struct {
	Storages []Storage `json:"storages"`
	Count    int       `json:"count"`
}
