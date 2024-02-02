package handler

import (
	"bazaar/api/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h Handler) CreateTransaction(c *gin.Context) {
	createTransaction := models.CreateTransaction{}

	if err := c.ShouldBindJSON(&createTransaction); err != nil {
		handleResponse(c, "error while reading body from client", http.StatusBadRequest, err)
	}

	id, err := h.storage.Transaction().Create(createTransaction)
	if err != nil {
		handleResponse(c, "error while create transaction", http.StatusInternalServerError, err)
		return
	}

	transaction, err := h.storage.Transaction().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while get transaction", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusCreated, transaction)

}

func (h Handler) GetTransactionByID(c *gin.Context) {
	var err error

	id := c.Param("id")

	transaction, err := h.storage.Transaction().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while get transaction by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, transaction)

}

func (h Handler) GetTransactionList(c *gin.Context) {

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

	response, err := h.storage.Transaction().GetList(models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	if err != nil {
		handleResponse(c, "error while get transaction", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, response)

}

func (h Handler) UpdateTransaction(c *gin.Context) {
	updateTransaction := models.UpdateTransaction{}

	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateTransaction.ID = uid

	if err := c.ShouldBindJSON(&updateTransaction); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Transaction().Update(updateTransaction)
	if err != nil {
		handleResponse(c, "error while updating transaction", http.StatusInternalServerError, err.Error())
		return
	}

	transaction, err := h.storage.Transaction().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while updating transaction", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, transaction)

}

func (h Handler) DeleteTransaction(c *gin.Context) {

	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.Transaction().Delete(id.String()); err != nil {
		handleResponse(c, "error while deleting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "data succesfully deleted")

}
