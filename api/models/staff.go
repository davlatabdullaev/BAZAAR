package models

type Staff struct {
	ID        string `json:"id"`
	BranchID  string `json:"branch_id"`
	TarifID   string `json:"tarif_id"`
	TypeStaff string `json:"type_staff"`
	Name      string `json:"name"`
	Balance   string `json:"balance"`
	BirthDate string `json:"birth_date"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	Login     string `json:"login"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type CreateStaff struct {
	BranchID  string `json:"branch_id"`
	TarifID   string `json:"tarif_id"`
	TypeStaff string `json:"type_staff"`
	Name      string `json:"name"`
	Balance   string `json:"balance"`
	BirthDate string `json:"birth_date"`
	Gender    string `json:"gender"`
	Login     string `json:"login"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

type UpdateStaff struct {
	ID        string `json:"id"`
	BranchID  string `json:"branch_id"`
	TarifID   string `json:"tarif_id"`
	TypeStaff string `json:"type_staff"`
	Name      string `json:"name"`
	Balance   string `json:"balance"`
	BirthDate string `json:"birth_date"`
	Gender    string `json:"gender"`
	Login     string `json:"login"`
	Password  string `json:"password"`
	UpdatedAt string `json:"updated_at"`
}

type StaffsResponse struct {
	Staffs []Staff `json:"staffs"`
	Count  int     `json:"count"`
}
