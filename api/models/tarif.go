package models

type Tarif struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	TarifType     string  `json:"tarif_type"`
	AmountForCash float64 `json:"amount_for_cash"`
	AmountForCard float64 `json:"amount_for_card"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	DeletedAt     string  `json:"deleted_at"`
}

type CreateTarif struct {
	Name          string  `json:"name"`
	TarifType     string  `json:"tarif_type"`
	AmountForCash float64 `json:"amount_for_cash"`
	AmountForCard float64 `json:"amount_for_card"`
	CreatedAt     string  `json:"created_at"`
}

type UpdateTarif struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	TarifType     string  `json:"tarif_type"`
	AmountForCash float64 `json:"amount_for_cash"`
	AmountForCard float64 `json:"amount_for_card"`
	UpdatedAt     string  `json:"updated_at"`
}

type TarifsResponse struct {
	Tarifs []Tarif `json:"tarifs"`
	Count  int     `json:"count"`
}
