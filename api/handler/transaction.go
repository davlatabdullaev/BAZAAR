package handler

import (
	"bazaar/api/models"
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateTransaction godoc
// @Router       /transaction [POST]
// @Summary      Create a new transaction
// @Description  Create a new transaction
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        transaction  body  models.CreateTransaction  true  "transaction data"
// @Success      201  {object}  models.Transaction
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateTransaction(c *gin.Context) {
	createTransaction := models.CreateTransaction{}

	if err := c.ShouldBindJSON(&createTransaction); err != nil {
		handleResponse(c, "error while reading body from client", http.StatusBadRequest, err)
	}

	id, err := h.storage.Transaction().Create(context.Background(), createTransaction)
	if err != nil {
		handleResponse(c, "error while create transaction", http.StatusInternalServerError, err)
		return
	}

	transaction, err := h.storage.Transaction().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while get transaction", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusCreated, transaction)

}

// GetTransactionByID godoc
// @Router       /transaction/{id} [GET]
// @Summary      Get transaction by id
// @Description  Get transaction by id
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        id path string true "transaction"
// @Success      200  {object}  models.Transaction
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetTransactionByID(c *gin.Context) {
	var err error

	uid := c.Param("id")

	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "invalid uuid type ", http.StatusBadRequest, err.Error())
		return
	}

	transaction, err := h.storage.Transaction().Get(context.Background(), models.PrimaryKey{
		ID: id.String(),
	})
	if err != nil {
		handleResponse(c, "error while get transaction by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, transaction)

}

// GetTransactionsList godoc
// @Router       /transaction [GET]
// @Summary      Get transactions list
// @Description  Get transactions list
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param		 from-amount query string false "from-amount"
// @Param		 to-amount query string false "to-amount"
// @Success      200  {object}  models.TransactionsResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetTransactionList(c *gin.Context) {

	var (
		page, limit int
		err         error
		fromAmount  float64
		toAmount    float64
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

	toAmountStr := c.DefaultQuery("to-amount", fmt.Sprintf("%f", math.MaxFloat64))
	toAmount, err = strconv.ParseFloat(toAmountStr, 64)
	if err != nil {
		handleResponse(c, "error is while converting to amount", http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.storage.Transaction().GetList(context.Background(), models.GetListTransactionsRequest{
		Page:   page,
		Limit:  limit,
		FromAmount: fromAmount,
		ToAmount: toAmount,
	})

	if err != nil {
		handleResponse(c, "error while get transaction", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, response)

}

// UpdateTransaction godoc
// @Router       /transaction/{id} [PUT]
// @Summary      Update transaction by id
// @Description  Update transaction by id
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        id path string true "transaction id"
// @Param        transaction body models.UpdateTransaction true "transaction"
// @Success      200  {object}  models.Tarif
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
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

	id, err := h.storage.Transaction().Update(context.Background(), updateTransaction)
	if err != nil {
		handleResponse(c, "error while updating transaction", http.StatusInternalServerError, err.Error())
		return
	}

	transaction, err := h.storage.Transaction().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while updating transaction", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, transaction)

}

// DeleteTransaction godoc
// @Router       /transaction/{id} [DELETE]
// @Summary      Delete Transaction
// @Description  Delete Transaction
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        id path string true "transaction id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteTransaction(c *gin.Context) {

	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.Transaction().Delete(context.Background(), id.String()); err != nil {
		handleResponse(c, "error while deleting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "data succesfully deleted")

}
