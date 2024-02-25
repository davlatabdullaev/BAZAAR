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

// CreateSale godoc
// @Router       /sale [POST]
// @Summary      Create a new sale
// @Description  Create a new sale
// @Tags         sale
// @Accept       json
// @Produce      json
// @Param        sale body  models.CreateSale  true  "sale data"
// @Success      201  {object}  models.Sale
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateSale(c *gin.Context) {
	createSale := models.CreateSale{}

	if err := c.ShouldBindJSON(&createSale); err != nil {
		handleResponse(c, h.log, "error while reading body from client", http.StatusBadRequest, err)
	}

	id, err := h.storage.Sale().Create(context.Background(), createSale)
	if err != nil {
		handleResponse(c, h.log, "error while creating sale", http.StatusInternalServerError, err)
		return
	}

	sale, err := h.storage.Sale().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, h.log, "error while get sale", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "", http.StatusCreated, sale)

}

// GetSaleByID godoc
// @Router       /sale/{id} [GET]
// @Summary      Get sale by id
// @Description  Get sale by id
// @Tags         sale
// @Accept       json
// @Produce      json
// @Param        id path string true "sale"
// @Success      200  {object}  models.Sale
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetSaleByID(c *gin.Context) {
	var err error

	uid := c.Param("id")

	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, h.log, "invalid uuid type ", http.StatusBadRequest, err.Error())
		return
	}

	sale, err := h.storage.Sale().Get(context.Background(), models.PrimaryKey{
		ID: id.String(),
	})
	if err != nil {
		handleResponse(c, h.log, "error while get sale by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, sale)

}

// GetSalesList godoc
// @Router       /sale [GET]
// @Summary      Get sales list
// @Description  Get sales list
// @Tags         sale
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        search query string false "search"
// @Success      200  {object}  models.SalesResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetSaleList(c *gin.Context) {

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

	response, err := h.storage.Sale().GetList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	if err != nil {
		handleResponse(c, h.log, "error while getting sale", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, response)

}

// UpdateSale godoc
// @Router       /sale/{id} [PUT]
// @Summary      Update sale by id
// @Description  Update sale by id
// @Tags         sale
// @Accept       json
// @Produce      json
// @Param        id path string true "sale id"
// @Param        sale body models.UpdateSale true "sale"
// @Success      200  {object}  models.Sale
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateSale(c *gin.Context) {
	updateSale := models.UpdateSale{}

	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, h.log, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateSale.ID = uid

	if err := c.ShouldBindJSON(&updateSale); err != nil {
		handleResponse(c, h.log, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Sale().Update(context.Background(), updateSale)
	if err != nil {
		handleResponse(c, h.log, "error while updating sale", http.StatusInternalServerError, err.Error())
		return
	}

	sale, err := h.storage.Sale().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, h.log, "error while get sale by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, sale)

}

// DeleteSale godoc
// @Router       /sale/{id} [DELETE]
// @Summary      Delete Sale
// @Description  Delete Sale
// @Tags         sale
// @Accept       json
// @Produce      json
// @Param        id path string true "sale id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteSale(c *gin.Context) {

	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, h.log, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.Sale().Delete(context.Background(), id.String()); err != nil {
		handleResponse(c, h.log, "error while deleting sale by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, "data succesfully deleted")

}
