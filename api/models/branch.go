package models

import "time"

type Branch struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt string    `json:"deleted_at"`
}

type CreateBranch struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type UpdateBranch struct {
	ID        string    `json:"-"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
}

type BranchsResponse struct {
	Branchs []Branch `json:"branchs"`
	Count   int      `json:"count"`
}
