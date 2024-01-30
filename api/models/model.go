package models

type GetListRequest struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type PrimaryKey struct {
	ID string `json:"id"`
}

type Delete struct {
	DeletedAt string `json:"deleted_at"`
}
