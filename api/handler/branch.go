package handler

import (
	"bazaar/api/models"
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h Handler) CreateBranch(c *gin.Context) {
	createBranch := models.CreateBranch{}

	if err := c.ShouldBindJSON(&createBranch); err != nil {
		handleResponse(c, "error while reading body from client", http.StatusBadRequest, err)
	}

	id, err := h.storage.Branch().Create(context.Background(), createBranch)
	if err != nil {
		handleResponse(c, "error while creating branch", http.StatusInternalServerError, err)
		return
	}

	branch, err := h.storage.Branch().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while get branch", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusCreated, branch)

}

func (h Handler) GetBranchByID(c *gin.Context) {

	var err error

	id := c.Param("id")

	branch, err := h.storage.Branch().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while get branch by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, branch)

}

func (h Handler) GetBranchList(c *gin.Context) {

	var (
		page, limit int
		search      string
		err         error
	)

	pageStr := c.DefaultQuery("page", "1")
	page, err = strconv.Atoi(pageStr)
	if err != nil {
		handleResponse(c, "error while parsing page ", http.StatusBadRequest, err.Error())
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		handleResponse(c, "error while parsing limit", http.StatusBadRequest, err.Error())
		return
	}

	search = c.Query("search")

	response, err := h.storage.Branch().GetList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	if err != nil {
		handleResponse(c, "error while getting branch", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, response)

}

func (h Handler) UpdateBranch(c *gin.Context) {
	updateBranch := models.UpdateBranch{}

	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateBranch.ID = uid

	if err := c.ShouldBindJSON(&updateBranch); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Branch().Update(context.Background(), updateBranch)
	if err != nil {
		handleResponse(c, "error while updating branch", http.StatusInternalServerError, err.Error())
		return
	}

	branch, err := h.storage.Branch().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while getting branch by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, branch)

}

func (h Handler) DeleteBranch(c *gin.Context) {

	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.Branch().Delete(context.Background(), id.String()); err != nil {
		handleResponse(c, "error while deleting branch by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "data succesfully deleted")

}
