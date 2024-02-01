package handler

import (
	"bazaar/api/models"
	"encoding/json"
	"errors"
	"net/http"
)

func (h Handler) Branch(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateBranch(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		if _, ok := values["id"]; !ok {
			h.GetBranchList(w)
		} else {
			h.GetBranchByID(w, r)
		}
	case http.MethodPut:
		{
			h.UpdateBranch(w, r)
		}
	case http.MethodDelete:
		h.DeleteBranch(w, r)
	}
}

func (h Handler) CreateBranch(w http.ResponseWriter, r *http.Request) {
	createBranch := models.CreateBranch{}

	if err := json.NewDecoder(r.Body).Decode(&createBranch); err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	id, err := h.storage.Branch().Create(createBranch)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	branch, err := h.storage.Branch().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusCreated, branch)

}

func (h Handler) GetBranchByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusInternalServerError, errors.New("id is required"))
		return
	}
	id := values["id"][0]
	var err error

	basket, err := h.storage.Branch().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, basket)

}

func (h Handler) GetBranchList(w http.ResponseWriter) {

	var (
		page, limit = 1, 50
		search      string
		err         error
	)

	response, err := h.storage.Branch().GetList(models.GetListRequest{
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

func (h Handler) UpdateBranch(w http.ResponseWriter, r *http.Request) {
	updateBranch := models.UpdateBranch{}

	if err := json.NewDecoder(r.Body).Decode(&updateBranch); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Branch().Update(updateBranch)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	branch, err := h.storage.Branch().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, branch)

}

func (h Handler) DeleteBranch(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]

	if err := h.storage.Branch().Delete(id); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, "data succesfully deleted")

}
