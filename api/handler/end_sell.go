package handler

import (
	"bazaar/api/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// EndSell godoc
// @Router           /end_sell/{id} [PUT]
// @Summary          end sell
// @Description      end sell
// @Tags             sell
// @Accept           json
// @Produce          json
// @Param            id path string true "sale_id"
// @Param            SaleRequest body models.SaleRequest true "sale request"
// @Succes           200 {object} models.Response
// @Failure          400 {object} models.Response
// @Failure          404 {object} models.Response
// @Failure          500 {object} models.Response
func (h Handler) EndSale(c *gin.Context) {

	var (
		totalPrice float64
		err        error
	)

	id := c.Param("id")

	request := models.SaleRequest{}

	if err = c.ShouldBindJSON(&request); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	baskets, err := h.storage.Basket().GetList(context.Background(), models.GetBasketsListRequest{
		Page:   1,
		Limit:  100,
		Search: id,
	})

	if err != nil {
		handleResponse(c, "error while getting baskets list", http.StatusInternalServerError, err.Error())
		return
	}

	selectedProducts := make(map[string]models.Basket)

	for _, basket := range baskets.Baskets {
		totalPrice += basket.Price
		selectedProducts[basket.ProductID] = basket
	}

	if request.Status == "cancel" {
		totalPrice = 0
	}

	saleID, err := h.storage.Sale().UpdateSalePrice(context.Background(), models.SaleRequest{
		ID:         id,
		TotalPrice: int(totalPrice),
		Status:     request.Status,
	})
	if err != nil {
		handleResponse(c, "error while updating sale price and status by id", http.StatusInternalServerError, err.Error())
		return
	}

	sale, err := h.storage.Sale().Get(context.Background(), models.PrimaryKey{
		ID: saleID,
	})
	if err != nil {
		handleResponse(c, "error while get sale by id", http.StatusInternalServerError, err.Error())
		return
	}

	if request.Status == "cancel" {
		_, err := h.storage.Transaction().Create(context.Background(), models.CreateTransaction{
			SaleID:          id,
			StaffID:         sale.CashierID,
			TransactionType: "withdraw",
			SourceType:      "sales",
			Amount:          totalPrice,
			Description:     "sale cancelled",
		})
		if err != nil {
			handleResponse(c, "error while creating transaction", http.StatusInternalServerError, err.Error())
			return
		}

		handleResponse(c, "succes", http.StatusOK, sale)

	}

	storageData, err := h.storage.Storage().GetList(context.Background(), models.GetListRequest{
		Page:  1,
		Limit: 100,
	})
	if err != nil {
		handleResponse(c, "error while getting storages list", http.StatusInternalServerError, err.Error())
		return
	}

	storageMap := make(map[string]models.Storage)
	for _, storage := range storageData.Storages {
		storageMap[storage.ID] = storage
	}

	for i, value := range storageMap {
		if value.ProductID == selectedProducts[value.ProductID].ProductID {
			_, err := h.storage.Storage().Update(context.Background(), models.UpdateStorage{
				ID:        i,
				ProductID: value.ProductID,
				BranchID:  value.BranchID,
				Count:     value.Count - selectedProducts[value.ProductID].Quantity,
			})

			if err != nil {
				handleResponse(c, "error while updating repositoryData prod quantities", http.StatusInternalServerError, err.Error())
				return
			}

			_, err = h.storage.StorageTransaction().Create(context.Background(), models.CreateStorageTransaction{
				StaffID:                sale.CashierID,
				ProductID:              value.ProductID,
				StorageTransactionType: "minus",
				Price:                  selectedProducts[value.ProductID].Price,
				Quantity:               float64(selectedProducts[value.ProductID].Quantity),
			})
			if err != nil {
				handleResponse(c, "error while creating storage data", http.StatusInternalServerError, err.Error())
				return
			}

		}
	}

	

}
