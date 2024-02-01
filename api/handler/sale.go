package handler

import (
	"bazaar/api/models"
	"encoding/json"
	"errors"
	"net/http"
)

func (h Handler) Sale(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateSale(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		if _, ok := values["id"]; !ok {
			h.GetSaleList(w)
		} else {
			h.GetSaleByID(w, r)
		}
	case http.MethodPut:
		{
			h.UpdateSale(w, r)
		}
	case http.MethodDelete:
		h.DeleteSale(w, r)
	}
}

func (h Handler) CreateSale(w http.ResponseWriter, r *http.Request) {
	createSale := models.CreateSale{}

	if err := json.NewDecoder(r.Body).Decode(&createSale); err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	id, err := h.storage.Sale().Create(createSale)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	sale, err := h.storage.Sale().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusCreated, sale)

}

func (h Handler) GetSaleByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusInternalServerError, errors.New("id is required"))
		return
	}
	id := values["id"][0]
	var err error

	sale, err := h.storage.Sale().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, sale)

}

func (h Handler) GetSaleList(w http.ResponseWriter) {

	var (
		page, limit = 1, 50
		search      string
		err         error
	)

	response, err := h.storage.Sale().GetList(models.GetListRequest{
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

func (h Handler) UpdateSale(w http.ResponseWriter, r *http.Request) {
	updateSale := models.UpdateSale{}

	if err := json.NewDecoder(r.Body).Decode(&updateSale); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Sale().Update(updateSale)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sale, err := h.storage.Sale().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, sale)

}

func (h Handler) DeleteSale(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]

	if err := h.storage.Sale().Delete(id); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, "data succesfully deleted")

}
