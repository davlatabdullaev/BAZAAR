package models

import "time"

type Tarif struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	TarifType     string    `json:"tarif_type"`
	AmountForCash float64   `json:"amount_for_cash"`
	AmountForCard float64   `json:"amount_for_card"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     time.Time `json:"deleted_at"`
}

type CreateTarif struct {
	Name          string    `json:"name"`
	TarifType     string    `json:"tarif_type"`
	AmountForCash float64   `json:"amount_for_cash"`
	AmountForCard float64   `json:"amount_for_card"`
}

type UpdateTarif struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	TarifType     string    `json:"tarif_type"`
	AmountForCash float64   `json:"amount_for_cash"`
	AmountForCard float64   `json:"amount_for_card"`
}

type TarifsResponse struct {
	Tarifs []Tarif `json:"tarifs"`
	Count  int     `json:"count"`
}
