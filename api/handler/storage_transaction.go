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

// CreateStorageTransaction godoc
// @Router       /storage_transaction [POST]
// @Summary      Create a new storage_transaction
// @Description  Create a new storage_transaction
// @Tags         storage_transaction
// @Accept       json
// @Produce      json
// @Param        storage_transaction  body  models.CreateStorageTransaction  true  "storage transaction  data"
// @Success      201  {object}  models.StorageTransaction
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateStorageTransaction(c *gin.Context) {
	createStorageTransaction := models.CreateStorageTransaction{}

	if err := c.ShouldBindJSON(&createStorageTransaction); err != nil {
		handleResponse(c, h.log, "error while reading body from client", http.StatusBadRequest, err)
	}

	id, err := h.storage.StorageTransaction().Create(context.Background(), createStorageTransaction)
	if err != nil {
		handleResponse(c, h.log, "error while creating storage transaction", http.StatusInternalServerError, err)
		return
	}

	storageTransaction, err := h.storage.StorageTransaction().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, h.log, "error while get storage transaction", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "", http.StatusCreated, storageTransaction)

}

// GetStorageTransactionByID godoc
// @Router       /storage_transaction/{id} [GET]
// @Summary      Get storage transaction by id
// @Description  Get storage transaction by id
// @Tags         storage_transaction
// @Accept       json
// @Produce      json
// @Param        id path string true "storage transaction"
// @Success      200  {object}  models.StorageTransaction
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetStorageTransactionByID(c *gin.Context) {
	var err error

	uid := c.Param("id")

	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, h.log, "invalid uuid type ", http.StatusBadRequest, err.Error())
		return
	}

	storageTransaction, err := h.storage.StorageTransaction().Get(context.Background(), models.PrimaryKey{
		ID: id.String(),
	})
	if err != nil {
		handleResponse(c, h.log, "error while get storage transaction by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, storageTransaction)

}

// GetStorageTransactionsList godoc
// @Router       /storage_transaction [GET]
// @Summary      Get storage_transactions list
// @Description  Get storage_transactions list
// @Tags         storage_transaction
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        search query string false "search"
// @Success      200  {object}  models.StorageTransactionsResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetStorageTransactionList(c *gin.Context) {

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

	response, err := h.storage.StorageTransaction().GetList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	if err != nil {
		handleResponse(c, h.log, "error while get storage transaction list", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, response)

}

// UpdateStorageTransaction godoc
// @Router       /storage_transaction/{id} [PUT]
// @Summary      Update storage transaction by id
// @Description  Update storage transaction by id
// @Tags         storage_transaction
// @Accept       json
// @Produce      json
// @Param        id path string true "storage transaction id"
// @Param        storage_transaction body models.UpdateStorageTransaction true "storage_transaction"
// @Success      200  {object}  models.StorageTransaction
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateStorageTransaction(c *gin.Context) {
	updateStorageTransaction := models.UpdateStorageTransaction{}
	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, h.log, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateStorageTransaction.ID = uid

	if err := c.ShouldBindJSON(&updateStorageTransaction); err != nil {
		handleResponse(c, h.log, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.StorageTransaction().Update(context.Background(), updateStorageTransaction)
	if err != nil {
		handleResponse(c, h.log, "error while reading body", http.StatusInternalServerError, err.Error())
		return
	}

	storageTransaction, err := h.storage.StorageTransaction().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, h.log, "error while updating storage transaction", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, storageTransaction)

}

// DeleteStorageTransaction godoc
// @Router       /storage_transaction/{id} [DELETE]
// @Summary      Delete Storage Transaction
// @Description  Delete Storage Transaction
// @Tags         storage_transaction
// @Accept       json
// @Produce      json
// @Param        id path string true "storage transaction id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteStorageTransaction(c *gin.Context) {

	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, h.log, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.StorageTransaction().Delete(context.Background(), id.String()); err != nil {
		handleResponse(c, h.log, "error while deleting storage transaction by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, "data succesfully deleted")

}
