package handler

import (
	"bazaar/api/models"
	"encoding/json"
	"errors"
	"net/http"
)

func (h Handler) Tarif(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateTarif(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		if _, ok := values["id"]; !ok {
			h.GetTarifList(w)
		} else {
			h.GetTarifByID(w, r)
		}
	case http.MethodPut:
		{
			h.UpdateTarif(w, r)
		}
	case http.MethodDelete:
		h.DeleteTarif(w, r)
	}
}

func (h Handler) CreateTarif(w http.ResponseWriter, r *http.Request) {
	createTarif := models.CreateTarif{}

	if err := json.NewDecoder(r.Body).Decode(&createTarif); err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	id, err := h.storage.Tarif().Create(createTarif)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	tarif, err := h.storage.Tarif().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusCreated, tarif)

}

func (h Handler) GetTarifByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusInternalServerError, errors.New("id is required"))
		return
	}
	id := values["id"][0]
	var err error

	storageTransaction, err := h.storage.Tarif().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, storageTransaction)

}

func (h Handler) GetTarifList(w http.ResponseWriter) {

	var (
		page, limit = 1, 50
		search      string
		err         error
	)

	response, err := h.storage.Tarif().GetList(models.GetListRequest{
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

func (h Handler) UpdateTarif(w http.ResponseWriter, r *http.Request) {
	updateTarif := models.UpdateTarif{}

	if err := json.NewDecoder(r.Body).Decode(&updateTarif); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Tarif().Update(updateTarif)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	tarif, err := h.storage.Tarif().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, tarif)

}

func (h Handler) DeleteTarif(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]

	if err := h.storage.Tarif().Delete(id); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, "data succesfully deleted")

}
