package handler

import (
	"bazaar/api/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h Handler) CreateStorageTransaction(c *gin.Context) {
	createStorageTransaction := models.CreateStorageTransaction{}

	if err := c.ShouldBindJSON(&createStorageTransaction); err != nil {
		handleResponse(c, "error while reading body from client", http.StatusBadRequest, err)
	}

	id, err := h.storage.StorageTransaction().Create(createStorageTransaction)
	if err != nil {
		handleResponse(c, "error while creating storage transaction", http.StatusInternalServerError, err)
		return
	}

	storageTransaction, err := h.storage.StorageTransaction().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while get storage transaction", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusCreated, storageTransaction)

}

func (h Handler) GetStorageTransactionByID(c *gin.Context) {
	var err error

	id := c.Param("id")

	storageTransaction, err := h.storage.StorageTransaction().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while get storage transaction by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, storageTransaction)

}

func (h Handler) GetStorageTransactionList(c *gin.Context) {

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

	response, err := h.storage.StorageTransaction().GetList(models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	if err != nil {
		handleResponse(c, "error while get storage transaction list", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, response)

}

func (h Handler) UpdateStorageTransaction(c *gin.Context) {
	updateStorageTransaction := models.UpdateStorageTransaction{}
	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateStorageTransaction.ID = uid

	if err := c.ShouldBindJSON(&updateStorageTransaction); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.StorageTransaction().Update(updateStorageTransaction)
	if err != nil {
		handleResponse(c, "error while reading body", http.StatusInternalServerError, err.Error())
		return
	}

	storageTransaction, err := h.storage.StorageTransaction().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while updating storage transaction", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, storageTransaction)

}

func (h Handler) DeleteStorageTransaction(c *gin.Context) {

	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.StorageTransaction().Delete(id.String()); err != nil {
		handleResponse(c, "error while deleting storage transaction by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "data succesfully deleted")

}
