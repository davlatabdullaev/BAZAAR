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

// CreateIncomeProduct godoc
// @Router       /income_product [POST]
// @Summary      Create a new income product
// @Description  Create a new income product
// @Tags         income_product
// @Accept       json
// @Produce      json
// @Param        income_product  body  models.CreateIncomeProduct  true  "income product data"
// @Success      201  {object}  models.IncomeProduct
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateIncomeProduct(c *gin.Context) {
	createIncomeProduct := models.CreateIncomeProduct{}

	if err := c.ShouldBindJSON(&createIncomeProduct); err != nil {
		handleResponse(c, "error while reading income product body from client", http.StatusBadRequest, err)
		return
	}

	incomeData, err := h.storage.Income().Get(context.Background(), models.PrimaryKey{
		ID: createIncomeProduct.IncomeID,
	})
	if err != nil {
		handleResponse(c, "error while search income for create income product", http.StatusInternalServerError, err)
		return
	}

	storageDataForProduct, err := h.storage.Storage().GetList(context.Background(), models.GetListRequest{
		Page:   1,
		Limit:  100,
		Search: createIncomeProduct.ProductID,
	})
	if err != nil {
		handleResponse(c, "error while search storage for create income product", http.StatusInternalServerError, err)
		return
	}

	storageDataForBranch, err := h.storage.Storage().GetList(context.Background(), models.GetListRequest{
		Page:   1,
		Limit:  100,
		Search: incomeData.BranchID,
	})
	if err != nil {
		handleResponse(c, "error while search storage for create income product", http.StatusInternalServerError, err)
		return
	}

	if createIncomeProduct.ProductID == storageDataForProduct.Storages[0].ProductID && incomeData.BranchID == storageDataForBranch.Storages[0].BranchID {

		if err := h.storage.Storage().UpdateCount(context.Background(), models.UpdateCount{
			ID:    storageDataForBranch.Storages[0].ID,
			Count: createIncomeProduct.Count,
		}); err != nil {
			handleResponse(c, "error while update storage count", http.StatusInternalServerError, err)
			return
		}

	} else {
		id, err := h.storage.IncomeProduct().Create(context.Background(), createIncomeProduct)
		if err != nil {
			handleResponse(c, "error while creating income product", http.StatusInternalServerError, err)
			return
		}

		incomeProduct, err := h.storage.IncomeProduct().Get(context.Background(), models.PrimaryKey{
			ID: id,
		})
		if err != nil {
			handleResponse(c, "error while get income product", http.StatusInternalServerError, err)
			return
		}

		handleResponse(c, "", http.StatusCreated, incomeProduct)
	}
}

// GetIncomeProductByID godoc
// @Router       /income_product/{id} [GET]
// @Summary      Get income product by id
// @Description  Get income product by id
// @Tags         income_product
// @Accept       json
// @Produce      json
// @Param        id path string true "income_product"
// @Success      200  {object}  models.IncomeProduct
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetIncomeProductByID(c *gin.Context) {

	var err error

	uid := c.Param("id")

	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "invalid uuid type ", http.StatusBadRequest, err.Error())
		return
	}

	incomeProduct, err := h.storage.IncomeProduct().Get(context.Background(), models.PrimaryKey{
		ID: id.String(),
	})
	if err != nil {
		handleResponse(c, "error while get income product by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, incomeProduct)

}

// GetIncomesList godoc
// @Router       /income_products [GET]
// @Summary      Get income products list
// @Description  Get income products list
// @Tags         income_product
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        search query string false "search"
// @Success      200  {object}  models.IncomeProductsResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetIncomeProductsList(c *gin.Context) {

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

	response, err := h.storage.IncomeProduct().GetList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	if err != nil {
		handleResponse(c, "error while getting income product", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, response)

}

// UpdateIncome godoc
// @Router       /income_product/{id} [PUT]
// @Summary      Update income product by id
// @Description  Update income product by id
// @Tags         income_product
// @Accept       json
// @Produce      json
// @Param        id path string true "income id"
// @Param        income_product body models.UpdateIncomeProduct true "income product"
// @Success      200  {object}  models.IncomeProduct
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateIncomeProduct(c *gin.Context) {
	updateIncomeProduct := models.UpdateIncomeProduct{}

	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateIncomeProduct.ID = uid

	if err := c.ShouldBindJSON(&updateIncomeProduct); err != nil {
		handleResponse(c, "error while reading income products body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.IncomeProduct().Update(context.Background(), updateIncomeProduct)
	if err != nil {
		handleResponse(c, "error while updating income product", http.StatusInternalServerError, err.Error())
		return
	}

	incomeProduct, err := h.storage.IncomeProduct().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while getting income by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, incomeProduct)

}

// DeleteIncomeProduct godoc
// @Router       /income_product/{id} [DELETE]
// @Summary      Delete Income
// @Description  Delete Income
// @Tags         income_product
// @Accept       json
// @Produce      json
// @Param        id path string true "income product id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteIncomeProduct(c *gin.Context) {

	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.IncomeProduct().Delete(context.Background(), id.String()); err != nil {
		handleResponse(c, "error while deleting income product by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "data succesfully deleted")

}
