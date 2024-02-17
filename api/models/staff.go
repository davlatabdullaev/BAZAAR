package models

import "time"

type Staff struct {
	ID        string    `json:"id"`
	BranchID  string    `json:"branch_id"`
	TarifID   string    `json:"tarif_id"`
	TypeStaff string    `json:"type_staff"`
	Name      string    `json:"name"`
	Balance   float64   `json:"balance"`
	BirthDate string    `json:"birth_date"`
	Age       int       `json:"age"`
	Gender    string    `json:"gender"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type CreateStaff struct {
	BranchID  string  `json:"branch_id"`
	TarifID   string  `json:"tarif_id"`
	TypeStaff string  `json:"type_staff"`
	Name      string  `json:"name"`
	Balance   float64 `json:"balance"`
	BirthDate string  `json:"birth_date"`
	Gender    string  `json:"gender"`
	Login     string  `json:"login"`
	Password  string  `json:"password"`
}

type UpdateStaff struct {
	ID        string  `json:"id"`
	BranchID  string  `json:"branch_id"`
	TarifID   string  `json:"tarif_id"`
	TypeStaff string  `json:"type_staff"`
	Name      string  `json:"name"`
	Balance   float64 `json:"balance"`
	BirthDate string  `json:"birth_date"`
	Gender    string  `json:"gender"`
	Login     string  `json:"login"`
	Password  string  `json:"password"`
}

type StaffsResponse struct {
	Staffs []Staff `json:"staffs"`
	Count  int     `json:"count"`
}

type UpdateStaffBalance struct {
	ID      string  `json:"id"`
	Balance float64 `json:"balance"`
}
