package api

import (
	"bazaar/api/handler"
	"net/http"
)

func New(h handler.Handler) {

	http.HandleFunc("/category", h.Category)
	http.HandleFunc("/staff", h.Staff)
	http.HandleFunc("/storageTransaction", h.StorageTransaction)
	http.HandleFunc("/tarif", h.Tarif)
	http.HandleFunc("/transaction", h.Transaction)
}
