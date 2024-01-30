package models

type Category struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ParentID  string `json:"parent_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type CreateCategory struct {
	Name      string `json:"name"`
	ParentID  string `json:"parent_id"`
	CreatedAt string `json:"created_at"`
}

type UpdateCategory struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ParentID  string `json:"parent_id"`
	UpdatedAt string `json:"updated_at"`
}

type CategoriesResponse struct {
	Categories []Category `json:"categories"`
	Count      int        `json:"count"`
}
