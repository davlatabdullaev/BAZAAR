package handler

import (
	"bazaar/api/models"
	"encoding/json"
	"errors"
	"net/http"
)

func (h Handler) Transaction(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateTransaction(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		if _, ok := values["id"]; !ok {
			h.GetTransactionList(w)
		} else {
			h.GetTransactionByID(w, r)
		}
	case http.MethodPut:
		{
			h.UpdateTransaction(w, r)
		}
	case http.MethodDelete:
		h.DeleteTarif(w, r)
	}
}

func (h Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	createTransaction := models.CreateTransaction{}

	if err := json.NewDecoder(r.Body).Decode(&createTransaction); err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	id, err := h.storage.Transaction().Create(createTransaction)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	transaction, err := h.storage.Transaction().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusCreated, transaction)

}

func (h Handler) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusInternalServerError, errors.New("id is required"))
		return
	}
	id := values["id"][0]
	var err error

	transaction, err := h.storage.Transaction().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, transaction)

}

func (h Handler) GetTransactionList(w http.ResponseWriter) {

	var (
		page, limit = 1, 50
		search      string
		err         error
	)

	response, err := h.storage.Transaction().GetList(models.GetListRequest{
		Page:  page,
		Limit: limit,
		Search: search,
	})

	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, response)

}

func (h Handler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	updateTransaction := models.UpdateTransaction{}

	if err := json.NewDecoder(r.Body).Decode(&updateTransaction); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Transaction().Update(updateTransaction)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	transaction, err := h.storage.Transaction().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, transaction)

}

func (h Handler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]

	if err := h.storage.Transaction().Delete(id); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, "data succesfully deleted")

}
