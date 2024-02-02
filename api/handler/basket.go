package handler

import (
	"bazaar/api/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h Handler) CreateBasket(c *gin.Context) {
	createBasket := models.CreateBasket{}

	if err := c.ShouldBindJSON(&createBasket); err != nil {
		handleResponse(c, "error while reading body from client", http.StatusBadRequest, err)
	}

	id, err := h.storage.Basket().Create(createBasket)
	if err != nil {
		handleResponse(c, "error while creating basket", http.StatusInternalServerError, err)
		return
	}

	basket, err := h.storage.Basket().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while get basket ", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusCreated, basket)

}

func (h Handler) GetBasketByID(c *gin.Context) {

	var err error

	id := c.Param("id")

	basket, err := h.storage.Basket().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while get basket by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, basket)

}

func (h Handler) GetBasketList(c *gin.Context) {

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

	response, err := h.storage.Basket().GetList(models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	if err != nil {
		handleResponse(c, "error while getting basket", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, response)

}

func (h Handler) UpdateBasket(c *gin.Context) {
	updateBasket := models.UpdateBasket{}

	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateBasket.ID = uid

	if err := c.ShouldBindJSON(&updateBasket); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Basket().Update(updateBasket)
	if err != nil {
		handleResponse(c, "error while updating basket", http.StatusInternalServerError, err.Error())
		return
	}

	basket, err := h.storage.Basket().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while getting basket by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, basket)

}

func (h Handler) DeleteBasket(c *gin.Context) {

	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.Basket().Delete(id.String()); err != nil {
		handleResponse(c, "error while deleting basket by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "data succesfully deleted")

}
