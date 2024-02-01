package handler

import (
	"bazaar/api/models"
	"encoding/json"
	"errors"
	"net/http"
)

func (h Handler) Basket(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateBasket(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		if _, ok := values["id"]; !ok {
			h.GetBasketList(w)
		} else {
			h.GetBasketByID(w, r)
		}
	case http.MethodPut:
		{
			h.UpdateBasket(w, r)
		}
	case http.MethodDelete:
		h.DeleteBasket(w, r)
	}
}

func (h Handler) CreateBasket(w http.ResponseWriter, r *http.Request) {
	createBasket := models.CreateBasket{}

	if err := json.NewDecoder(r.Body).Decode(&createBasket); err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	id, err := h.storage.Basket().Create(createBasket)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	basket, err := h.storage.Basket().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusCreated, basket)

}

func (h Handler) GetBasketByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusInternalServerError, errors.New("id is required"))
		return
	}
	id := values["id"][0]
	var err error

	basket, err := h.storage.Basket().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, basket)

}

func (h Handler) GetBasketList(w http.ResponseWriter) {

	var (
		page, limit = 1, 50
		search      string
		err         error
	)

	response, err := h.storage.Basket().GetList(models.GetListRequest{
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

func (h Handler) UpdateBasket(w http.ResponseWriter, r *http.Request) {
	updateBasket := models.UpdateBasket{}

	if err := json.NewDecoder(r.Body).Decode(&updateBasket); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Basket().Update(updateBasket)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	basket, err := h.storage.Basket().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, basket)

}

func (h Handler) DeleteBasket(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]

	if err := h.storage.Basket().Delete(id); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, "data succesfully deleted")

}
