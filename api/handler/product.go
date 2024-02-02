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

func (h Handler) CreateProduct(c *gin.Context) {
	createProduct := models.CreateProduct{}

	if err := c.ShouldBindJSON(&createProduct); err != nil {
		handleResponse(c, "error while reading body from client", http.StatusBadRequest, err)
	}

	id, err := h.storage.Product().Create(context.Background(), createProduct)
	if err != nil {
		handleResponse(c, "error while creating product", http.StatusInternalServerError, err)
		return
	}

	product, err := h.storage.Branch().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while getting product", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusCreated, product)

}

func (h Handler) GetProductByID(c *gin.Context) {
	var err error

	id := c.Param("id")

	product, err := h.storage.Product().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while get product by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, product)

}

func (h Handler) GetProductList(c *gin.Context) {

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

	response, err := h.storage.Product().GetList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	if err != nil {
		handleResponse(c, "error while get product list", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, response)

}

func (h Handler) UpdateProduct(c *gin.Context) {
	updateProduct := models.UpdateProduct{}

	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateProduct.ID = uid

	if err := c.ShouldBindJSON(&updateProduct); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Product().Update(context.Background(), updateProduct)
	if err != nil {
		handleResponse(c, "error while updating product", http.StatusInternalServerError, err.Error())
		return
	}

	product, err := h.storage.Product().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while update product by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, product)

}

func (h Handler) DeleteProduct(c *gin.Context) {

	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.Product().Delete(context.Background(), id.String()); err != nil {
		handleResponse(c, "error while deleting product", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "data succesfully deleted")

}
