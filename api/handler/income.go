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

// CreateIncome godoc
// @Router       /income [POST]
// @Summary      Create a new income
// @Description  Create a new income
// @Tags         income
// @Accept       json
// @Produce      json
// @Param        income  body  models.CreateIncome  true  "income data"
// @Success      201  {object}  models.Income
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateIncome(c *gin.Context) {
	createIncome := models.CreateIncome{}

	if err := c.ShouldBindJSON(&createIncome); err != nil {
		handleResponse(c, "error while reading income body from client", http.StatusBadRequest, err)
	}

	id, err := h.storage.Income().Create(context.Background(), createIncome)
	if err != nil {
		handleResponse(c, "error while creating income", http.StatusInternalServerError, err)
		return
	}

	income, err := h.storage.Income().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while get income", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusCreated, income)

}

// GetIncomeByID godoc
// @Router       /income/{id} [GET]
// @Summary      Get income by id
// @Description  Get income by id
// @Tags         income
// @Accept       json
// @Produce      json
// @Param        id path string true "income"
// @Success      200  {object}  models.Income
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetIncomeByID(c *gin.Context) {

	var err error

	uid := c.Param("id")

	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "invalid uuid type ", http.StatusBadRequest, err.Error())
		return
	}

	income, err := h.storage.Income().Get(context.Background(), models.PrimaryKey{
		ID: id.String(),
	})
	if err != nil {
		handleResponse(c, "error while get income by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, income)

}

// GetIncomesList godoc
// @Router       /incomes [GET]
// @Summary      Get incomes list
// @Description  Get incomes list
// @Tags         income
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        search query string false "search"
// @Success      200  {object}  models.IncomesResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetIncomesList(c *gin.Context) {

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

	response, err := h.storage.Income().GetList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	if err != nil {
		handleResponse(c, "error while getting income", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, response)

}

// UpdateIncome godoc
// @Router       /income/{id} [PUT]
// @Summary      Update income by id
// @Description  Update income by id
// @Tags         income
// @Accept       json
// @Produce      json
// @Param        id path string true "income id"
// @Param        income body models.UpdateIncome true "income"
// @Success      200  {object}  models.Income
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateIncome(c *gin.Context) {
	updateIncome := models.UpdateIncome{}

	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateIncome.ID = uid

	if err := c.ShouldBindJSON(&updateIncome); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Income().Update(context.Background(), updateIncome)
	if err != nil {
		handleResponse(c, "error while updating income", http.StatusInternalServerError, err.Error())
		return
	}

	income, err := h.storage.Income().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while getting income by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, income)

}

// DeleteIncome godoc
// @Router       /income/{id} [DELETE]
// @Summary      Delete Income
// @Description  Delete Income
// @Tags         income
// @Accept       json
// @Produce      json
// @Param        id path string true "income id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteIncome(c *gin.Context) {

	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.Income().Delete(context.Background(), id.String()); err != nil {
		handleResponse(c, "error while deleting income by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "data succesfully deleted")

}
