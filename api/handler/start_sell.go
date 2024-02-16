package handler

import (
	"bazaar/api/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StartSell godoc
// @Router       /sell [POST]
// @Summary      sell
// @Description  sell
// @Tags         sell
// @Accept       json
// @Produce      json
// @Param 		 sell body models.CreateSale false "sell"
// @Success      200  {object}  models.Sale
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) StartSell(c *gin.Context) {
	sell := models.CreateSale{}

	if err := c.ShouldBindJSON(&sell); err != nil {
		handleResponse(c, "error is while reading body", http.StatusBadRequest, err.Error())
		return
	}

	saleID, err := h.storage.Sale().Create(context.Background(), sell)
	if err != nil {
		handleResponse(c, "error is while creating sale", http.StatusInternalServerError, err.Error())
		return
	}

	sale, err := h.storage.Sale().Get(context.Background(), models.PrimaryKey{
		ID: saleID,
	})
	if err != nil {
		handleResponse(c, "error is while getting sale by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "success", http.StatusOK, sale)
}
