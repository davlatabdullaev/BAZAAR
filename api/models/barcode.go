package models

type Barcode struct {
	SaleID  string `json:"sale_id"`
	Barcode int    `json:"barcode"`
	Count   int    `json:"count"`
}