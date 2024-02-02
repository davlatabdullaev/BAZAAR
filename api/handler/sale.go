package handler

import (
	"bazaar/api/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h Handler) CreateSale(c *gin.Context) {
	createSale := models.CreateSale{}

	if err := c.ShouldBindJSON(&createSale); err != nil {
		handleResponse(c, "error while reading body from client", http.StatusBadRequest, err)
	}

	id, err := h.storage.Sale().Create(createSale)
	if err != nil {
		handleResponse(c, "error while creating sale", http.StatusInternalServerError, err)
		return
	}

	sale, err := h.storage.Sale().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while get sale", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusCreated, sale)

}

func (h Handler) GetSaleByID(c *gin.Context) {
	var err error

	id := c.Param("id")

	sale, err := h.storage.Sale().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while get sale by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, sale)

}

func (h Handler) GetSaleList(c *gin.Context) {

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

	response, err := h.storage.Sale().GetList(models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	if err != nil {
		handleResponse(c, "error while getting sale", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, response)

}

func (h Handler) UpdateSale(c *gin.Context) {
	updateSale := models.UpdateSale{}

	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateSale.ID = uid

	if err := c.ShouldBindJSON(&updateSale); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Sale().Update(updateSale)
	if err != nil {
		handleResponse(c, "error while updating sale", http.StatusInternalServerError, err.Error())
		return
	}

	sale, err := h.storage.Sale().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while get sale by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, sale)

}

func (h Handler) DeleteSale(c *gin.Context) {

	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.Sale().Delete(id.String()); err != nil {
		handleResponse(c, "error while deleting sale by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "data succesfully deleted")

}
