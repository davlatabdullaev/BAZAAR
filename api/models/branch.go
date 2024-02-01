package models

type Branch struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type CreateBranch struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type UpdateBranch struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	UpdatedAt string `json:"updated_at"`
}

type BranchsResponse struct {
	Branchs []Branch `json:"branchs"`
	Count   int      `json:"count"`
}
