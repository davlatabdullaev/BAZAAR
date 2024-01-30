package handler

import (
	"bazaar/api/models"
	"encoding/json"
	"errors"
	"net/http"
)

func (h Handler) Category(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateCategory(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		if _, ok := values["id"]; !ok {
			h.GetCategoryList(w)
		} else {
			h.GetCategoryByID(w, r)
		}
	case http.MethodPut:
		{
			h.UpdateCategory(w, r)
		}
	case http.MethodDelete:
		h.DeleteCategory(w, r)
	}
}

func (h Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	createCategory := models.CreateCategory{}

	if err := json.NewDecoder(r.Body).Decode(&createCategory); err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	id, err := h.storage.Category().Create(createCategory)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	car, err := h.storage.Category().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusCreated, car)

}

func (h Handler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusInternalServerError, errors.New("id is required"))
		return
	}
	id := values["id"][0]
	var err error

	category, err := h.storage.Category().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, category)

}

func (h Handler) GetCategoryList(w http.ResponseWriter) {

	var (
		page, limit = 1, 50
		search      string
		err         error
	)

	response, err := h.storage.Category().GetList(models.GetListRequest{
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

func (h Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	updateCategory := models.UpdateCategory{}

	if err := json.NewDecoder(r.Body).Decode(&updateCategory); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Category().Update(updateCategory)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	category, err := h.storage.Category().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, category)

}

func (h Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]

	if err := h.storage.Category().Delete(id); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, "data succesfully deleted")

}
