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

// CreateStaff godoc
// @Router       /staff [POST]
// @Summary      Create a new staff
// @Description  Create a new staff
// @Tags         staff
// @Accept       json
// @Produce      json
// @Param        staff  body  models.CreateStaff  true  "staff data"
// @Success      201  {object}  models.Staff
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateStaff(c *gin.Context) {
	createStaff := models.CreateStaff{}

	if err := c.ShouldBindJSON(&createStaff); err != nil {
		handleResponse(c, h.log, "error while reading body from client", http.StatusBadRequest, err)
	}

	id, err := h.storage.Staff().Create(context.Background(), createStaff)
	if err != nil {
		handleResponse(c, h.log, "error while creating staff", http.StatusInternalServerError, err)
		return
	}

	staff, err := h.storage.Staff().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, h.log, "error while get product", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "", http.StatusCreated, staff)

}

// GetStaffByID godoc
// @Router       /staff/{id} [GET]
// @Summary      Get staff by id
// @Description  Get staff by id
// @Tags         staff
// @Accept       json
// @Produce      json
// @Param        id path string true "staff"
// @Success      200  {object}  models.Staff
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetStaffByID(c *gin.Context) {

	var err error

	uid := c.Param("id")

	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, h.log, "invalid uuid type ", http.StatusBadRequest, err.Error())
		return
	}

	staff, err := h.storage.Staff().Get(context.Background(), models.PrimaryKey{
		ID: id.String(),
	})
	if err != nil {
		handleResponse(c, h.log, "error while get staff by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, staff)

}

// GetStaffsList godoc
// @Router       /staff [GET]
// @Summary      Get staffs list
// @Description  Get staffs list
// @Tags         staff
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        search query string false "search"
// @Success      200  {object}  models.StaffsResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetStaffList(c *gin.Context) {

	var (
		page, limit int
		search      string
		err         error
	)

	pageStr := c.DefaultQuery("page", "1")
	page, err = strconv.Atoi(pageStr)
	if err != nil {
		handleResponse(c, h.log, "error while parsing page ", http.StatusBadRequest, err.Error())
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		handleResponse(c, h.log, "error while parsing limit", http.StatusBadRequest, err.Error())
		return
	}

	search = c.Query("search")

	response, err := h.storage.Staff().GetList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	if err != nil {
		handleResponse(c, h.log, "error while getting staff", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, response)

}

// UpdateStaff godoc
// @Router       /staff/{id} [PUT]
// @Summary      Update staff by id
// @Description  Update staff by id
// @Tags         staff
// @Accept       json
// @Produce      json
// @Param        id path string true "staff id"
// @Param        staff body models.UpdateStaff true "staff"
// @Success      200  {object}  models.Staff
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateStaff(c *gin.Context) {
	updateStaff := models.UpdateStaff{}

	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, h.log, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateStaff.ID = uid

	if err := c.ShouldBindJSON(&updateStaff); err != nil {
		handleResponse(c, h.log, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Staff().Update(context.Background(), updateStaff)
	if err != nil {
		handleResponse(c, h.log, "error while updating staff", http.StatusInternalServerError, err.Error())
		return
	}

	staff, err := h.storage.Staff().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, h.log, "error while get staff by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, staff)
}

// DeleteStaff godoc
// @Router       /staff/{id} [DELETE]
// @Summary      Delete Staff
// @Description  Delete Staff
// @Tags         staff
// @Accept       json
// @Produce      json
// @Param        id path string true "staff id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteStaff(c *gin.Context) {

	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, h.log, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.Staff().Delete(context.Background(), id.String()); err != nil {
		handleResponse(c, h.log, "error while deleting staff", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, "data succesfully deleted")

}
