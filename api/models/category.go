package models

import "time"

type Category struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	ParentID  string    `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type CreateCategory struct {
	Name     string `json:"name"`
	ParentID string `json:"parent_id"`
}

type UpdateCategory struct {
	ID        string    `json:"-"`
	Name      string    `json:"name"`
	ParentID  string    `json:"parent_id"`
}

type CategoriesResponse struct {
	Categories []Category `json:"categories"`
	Count      int        `json:"count"`
}
