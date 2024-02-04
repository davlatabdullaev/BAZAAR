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

func (h Handler) CreateStaff(c *gin.Context) {
	createStaff := models.CreateStaff{}

	if err := c.ShouldBindJSON(&createStaff); err != nil {
		handleResponse(c, "error while reading body from client", http.StatusBadRequest, err)
	}

	id, err := h.storage.Staff().Create(context.Background(), createStaff)
	if err != nil {
		handleResponse(c, "error while creating staff", http.StatusInternalServerError, err)
		return
	}

	staff, err := h.storage.Staff().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while get product", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusCreated, staff)

}

func (h Handler) GetStaffByID(c *gin.Context) {

	var err error

	uid := c.Param("id")

	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "invalid uuid type ", http.StatusBadRequest, err.Error())
		return
	}

	staff, err := h.storage.Staff().Get(context.Background(), models.PrimaryKey{
		ID: id.String(),
	})
	if err != nil {
		handleResponse(c, "error while get staff by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, staff)

}

func (h Handler) GetStaffList(c *gin.Context) {

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

	response, err := h.storage.Staff().GetList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	if err != nil {
		handleResponse(c, "error while getting staff", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, response)

}

func (h Handler) UpdateStaff(c *gin.Context) {
	updateStaff := models.UpdateStaff{}

	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateStaff.ID = uid

	if err := c.ShouldBindJSON(&updateStaff); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Staff().Update(context.Background(), updateStaff)
	if err != nil {
		handleResponse(c, "error while updating staff", http.StatusInternalServerError, err.Error())
		return
	}

	staff, err := h.storage.Staff().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while get staff by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, staff)
}

func (h Handler) DeleteStaff(c *gin.Context) {

	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.Staff().Delete(context.Background(), id.String()); err != nil {
		handleResponse(c, "error while deleting staff", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "data succesfully deleted")

}
