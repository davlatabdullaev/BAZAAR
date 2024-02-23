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

// CreateBasket godoc
// @Router       /basket [POST]
// @Summary      Create a new basket
// @Description  Create a new basket
// @Tags         basket
// @Accept       json
// @Produce      json
// @Param        basket  body  models.CreateBasket  true  "basket data"
// @Success      201  {object}  models.Basket
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateBasket(c *gin.Context) {
	createBasket := models.CreateBasket{}

	if err := c.ShouldBindJSON(&createBasket); err != nil {
		handleResponse(c, "error while reading body from client", http.StatusBadRequest, err)
	}

	storage, err := h.storage.Storage().GetList(context.Background(), models.GetListRequest{
		Page:   1,
		Limit:  10,
		Search: createBasket.ProductID,
	})
	if err != nil {
		handleResponse(c, "error while searching product in storage", http.StatusInternalServerError, err)
	}

	if createBasket.Quantity <= storage.Storages[0].Count {

		basketForProductID, err := h.storage.Basket().GetList(context.Background(), models.GetBasketsListRequest{
			Page:   1,
			Limit:  100,
			Search: createBasket.ProductID,
		})
		if err != nil {
			handleResponse(c, "error while searching selected basket", http.StatusInternalServerError, "error searching selected basket")
		}

		basketForSaleID, err := h.storage.Basket().GetList(context.Background(), models.GetBasketsListRequest{
			Page:   1,
			Limit:  100,
			Search: createBasket.SaleID,
		})
		if err != nil {
			handleResponse(c, "error while searching selected basket", http.StatusInternalServerError, "error searching selected basket")
		}

		var foundProduct, foundSale bool

		for _, basket := range basketForProductID.Baskets {
			if basket.ProductID == createBasket.ProductID {
				foundProduct = true
				break
			}
		}

		for _, basket := range basketForSaleID.Baskets {
			if basket.SaleID == createBasket.SaleID {
				foundSale = true
				break
			}
		}

		if foundProduct && foundSale {

			h.storage.Basket().UpdateBasketQuantity(context.Background(), models.UpdateBasketQuantity{
				ID:       basketForProductID.Baskets[0].ID,
				Quantity: createBasket.Quantity,
			})

		} else {
			id, err := h.storage.Basket().Create(context.Background(), createBasket)
			if err != nil {
				handleResponse(c, "error while creating basket", http.StatusInternalServerError, err)
				return
			}

			basket, err := h.storage.Basket().Get(context.Background(), models.PrimaryKey{
				ID: id,
			})
			if err != nil {
				handleResponse(c, "error while get basket ", http.StatusInternalServerError, err)
				return
			}

			handleResponse(c, "", http.StatusCreated, basket)
		}

	} else {
		handleResponse(c, "", http.StatusOK, "the number of products is not enough")
	}

}

// GetBasketByID godoc
// @Router       /basket/{id} [GET]
// @Summary      Get basket by id
// @Description  Get basket by id
// @Tags         basket
// @Accept       json
// @Produce      json
// @Param        id path string true "basket"
// @Success      200  {object}  models.Basket
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetBasketByID(c *gin.Context) {

	var err error

	uid := c.Param("id")

	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "invalid uuid type ", http.StatusBadRequest, err.Error())
		return
	}

	basket, err := h.storage.Basket().Get(context.Background(), models.PrimaryKey{
		ID: id.String(),
	})
	if err != nil {
		handleResponse(c, "error while get basket by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, basket)

}

// GetBasketsList godoc
// @Router       /basket [GET]
// @Summary      Get baskets list
// @Description  Get baskets list
// @Tags         basket
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        search query string false "search"
// @Success      200  {object}  models.BasketsResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
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

	response, err := h.storage.Basket().GetList(context.Background(), models.GetBasketsListRequest{
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

// UpdateBasket godoc
// @Router       /basket/{id} [PUT]
// @Summary      Update basket by id
// @Description  Update basket by id
// @Tags         basket
// @Accept       json
// @Produce      json
// @Param        id path string true "basket id"
// @Param        basket body models.UpdateBasket true "basket"
// @Success      200  {object}  models.Basket
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
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

	id, err := h.storage.Basket().Update(context.Background(), updateBasket)
	if err != nil {
		handleResponse(c, "error while updating basket", http.StatusInternalServerError, err.Error())
		return
	}

	basket, err := h.storage.Basket().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while getting basket by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, basket)

}

// DeleteBasket godoc
// @Router       /basket/{id} [DELETE]
// @Summary      Delete Basket
// @Description  Delete Basket
// @Tags         basket
// @Accept       json
// @Produce      json
// @Param        id path string true "basket id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteBasket(c *gin.Context) {

	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.Basket().Delete(context.Background(), id.String()); err != nil {
		handleResponse(c, "error while deleting basket by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "data succesfully deleted")

}
