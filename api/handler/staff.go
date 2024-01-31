package handler

import (
	"bazaar/api/models"
	"encoding/json"
	"errors"
	"net/http"
)

func (h Handler) Staff(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateStaff(w, r)
	case http.MethodGet:

		values := r.URL.Query()
		if _, ok := values["id"]; !ok {

			h.GetStaffList(w)
		} else {
			h.GetStaffByID(w, r)
		}
	case http.MethodPut:
		{
			h.UpdateStaff(w, r)
		}
	case http.MethodDelete:
		h.DeleteStaff(w, r)
	}
}

func (h Handler) CreateStaff(w http.ResponseWriter, r *http.Request) {
	createStaff := models.CreateStaff{}

	if err := json.NewDecoder(r.Body).Decode(&createStaff); err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	id, err := h.storage.Staff().Create(createStaff)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	staff, err := h.storage.Staff().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusCreated, staff)

}

func (h Handler) GetStaffByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusInternalServerError, errors.New("id is required"))
		return
	}

	id := values["id"][0]
	var err error

	staff, err := h.storage.Staff().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, staff)

}

func (h Handler) GetStaffList(w http.ResponseWriter) {

	var (
		page, limit = 1, 50
		search      string
		err         error
	)

	response, err := h.storage.Staff().GetList(models.GetListRequest{
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

func (h Handler) UpdateStaff(w http.ResponseWriter, r *http.Request) {
	updateStaff := models.UpdateStaff{}

	if err := json.NewDecoder(r.Body).Decode(&updateStaff); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Staff().Update(updateStaff)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	staff, err := h.storage.Staff().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, staff)
}

func (h Handler) DeleteStaff(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]

	if err := h.storage.Staff().Delete(id); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, "data succesfully deleted")

}
