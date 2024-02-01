package handler

import (
	"bazaar/api/models"
	"encoding/json"
	"errors"
	"net/http"
)

func (h Handler) Storage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateStorage(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		if _, ok := values["id"]; !ok {
			h.GetStorageList(w)
		} else {
			h.GetStorageByID(w, r)
		}
	case http.MethodPut:
		{
			h.UpdateStorage(w, r)
		}
	case http.MethodDelete:
		h.DeleteStorage(w, r)
	}
}

func (h Handler) CreateStorage(w http.ResponseWriter, r *http.Request) {
	createStorage := models.CreateStorage{}

	if err := json.NewDecoder(r.Body).Decode(&createStorage); err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	id, err := h.storage.Storage().Create(createStorage)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	storage, err := h.storage.Storage().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusCreated, storage)

}

func (h Handler) GetStorageByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusInternalServerError, errors.New("id is required"))
		return
	}
	id := values["id"][0]
	var err error

	sale, err := h.storage.Storage().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, sale)

}

func (h Handler) GetStorageList(w http.ResponseWriter) {

	var (
		page, limit = 1, 50
		search      string
		err         error
	)

	response, err := h.storage.Storage().GetList(models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, response)

}

func (h Handler) UpdateStorage(w http.ResponseWriter, r *http.Request) {
	updateStorage := models.UpdateStorage{}

	if err := json.NewDecoder(r.Body).Decode(&updateStorage); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Storage().Update(updateStorage)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	storage, err := h.storage.Sale().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, storage)

}

func (h Handler) DeleteStorage(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]

	if err := h.storage.Storage().Delete(id); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, "data succesfully deleted")

}
