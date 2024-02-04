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

// CreateProduct godoc
// @Router       /product [POST]
// @Summary      Create a new product
// @Description  Create a new product
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        product  body  models.CreateProduct  true  "product data"
// @Success      201  {object}  models.Product
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
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

	product, err := h.storage.Product().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while getting product", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusCreated, product)

}

// GetProductByID godoc
// @Router       /product/{id} [GET]
// @Summary      Get product by id
// @Description  Get product by id
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        id path string true "product"
// @Success      200  {object}  models.Product
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetProductByID(c *gin.Context) {
	var err error

	uid := c.Param("id")

	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "invalid uuid type ", http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.storage.Product().Get(context.Background(), models.PrimaryKey{
		ID: id.String(),
	})
	if err != nil {
		handleResponse(c, "error while get product by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, product)

}

// GetProductsList godoc
// @Router       /product [GET]
// @Summary      Get products list
// @Description  Get products list
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        search query string false "search"
// @Success      200  {object}  models.ProductsResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
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

// UpdateProduct godoc
// @Router       /product/{id} [PUT]
// @Summary      Update product by id
// @Description  Update product by id
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        id path string true "product id"
// @Param        product body models.UpdateProduct true "product"
// @Success      200  {object}  models.Product
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
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

// DeleteProduct godoc
// @Router       /product/{id} [DELETE]
// @Summary      Delete Product
// @Description  Delete Product
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        id path string true "product id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
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
