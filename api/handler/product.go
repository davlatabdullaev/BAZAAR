package handler

import (
	"bazaar/api/models"
	"encoding/json"
	"errors"
	"net/http"
)

func (h Handler) Product(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateProduct(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		if _, ok := values["id"]; !ok {
			h.GetProductList(w)
		} else {
			h.GetProductByID(w, r)
		}
	case http.MethodPut:
		{
			h.UpdateProduct(w, r)
		}
	case http.MethodDelete:
		h.DeleteProduct(w, r)
	}
}

func (h Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	createProduct := models.CreateProduct{}

	if err := json.NewDecoder(r.Body).Decode(&createProduct); err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	id, err := h.storage.Product().Create(createProduct)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	product, err := h.storage.Branch().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusCreated, product)

}

func (h Handler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusInternalServerError, errors.New("id is required"))
		return
	}
	id := values["id"][0]
	var err error

	product, err := h.storage.Product().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, product)

}

func (h Handler) GetProductList(w http.ResponseWriter) {

	var (
		page, limit = 1, 50
		search      string
		err         error
	)

	response, err := h.storage.Product().GetList(models.GetListRequest{
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

func (h Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	updateProduct := models.UpdateProduct{}

	if err := json.NewDecoder(r.Body).Decode(&updateProduct); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Product().Update(updateProduct)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	product, err := h.storage.Product().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, product)

}

func (h Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]

	if err := h.storage.Product().Delete(id); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, "data succesfully deleted")

}
