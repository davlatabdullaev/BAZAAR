package handler

import (
	"bazaar/api/models"
	"encoding/json"
	"errors"
	"net/http"
)

func (h Handler) StorageTransaction(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateStorageTransaction(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		if _, ok := values["id"]; !ok {
			h.GetStorageTransactionList(w)
		} else {
			h.GetStorageTransactionByID(w, r)
		}
	case http.MethodPut:
		{
			h.UpdateStorageTransaction(w, r)
		}
	case http.MethodDelete:
		h.DeleteStorageTransaction(w, r)
	}
}

func (h Handler) CreateStorageTransaction(w http.ResponseWriter, r *http.Request) {
	createStorageTransaction := models.CreateStorageTransaction{}

	if err := json.NewDecoder(r.Body).Decode(&createStorageTransaction); err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	id, err := h.storage.StorageTransaction().Create(createStorageTransaction)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	storageTransaction, err := h.storage.StorageTransaction().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusCreated, storageTransaction)

}

func (h Handler) GetStorageTransactionByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusInternalServerError, errors.New("id is required"))
		return
	}
	id := values["id"][0]
	var err error

	storageTransaction, err := h.storage.StorageTransaction().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, storageTransaction)

}

func (h Handler) GetStorageTransactionList(w http.ResponseWriter) {

	var (
		page, limit = 1, 50
		search      string
		err         error
	)

	response, err := h.storage.StorageTransaction().GetList(models.GetListRequest{
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

func (h Handler) UpdateStorageTransaction(w http.ResponseWriter, r *http.Request) {
	updateStorageTransaction := models.UpdateStorageTransaction{}

	if err := json.NewDecoder(r.Body).Decode(&updateStorageTransaction); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.StorageTransaction().Update(updateStorageTransaction)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	storageTransaction, err := h.storage.StorageTransaction().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, storageTransaction)

}

func (h Handler) DeleteStorageTransaction(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]

	if err := h.storage.StorageTransaction().Delete(id); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, "data succesfully deleted")

}
