package handler

import (
	"bazaar/api/models"
	"context"
	"fmt"
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
	fmt.Println(totalPrice)

	if request.Status == "cancel" {
		totalPrice = 0
		return
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
				ProductID:              value.ProductID,
				StaffID:                sale.CashierID,
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

	salesResponse, err := h.storage.Sale().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})

	if err != nil {
		handleResponse(c, "error while getting sales list", http.StatusInternalServerError, err.Error())
		return
	}

	if salesResponse.Status == "succes" {

		cashierResponse, err := h.storage.Staff().Get(context.Background(), models.PrimaryKey{
			ID: sale.CashierID,
		})

		if err != nil {
			handleResponse(c, "error while get staff data", http.StatusInternalServerError, err.Error())
			return
		}

		cashierTarifResponse, err := h.storage.Tarif().Get(context.Background(), models.PrimaryKey{
			ID: cashierResponse.TarifID,
		})
		if err != nil {
			handleResponse(c, "error while getting tarif by id", http.StatusInternalServerError, err.Error())
			return
		}

		shopAssistantResponse, err := h.storage.Staff().Get(context.Background(), models.PrimaryKey{
			ID: sale.CashierID,
		})

		if err != nil {
			handleResponse(c, "error while get staff data", http.StatusInternalServerError, err.Error())
			return
		}

		shopAssistantTarifResponse, err := h.storage.Tarif().Get(context.Background(), models.PrimaryKey{
			ID: shopAssistantResponse.TarifID,
		})
		if err != nil {
			handleResponse(c, "error while getting tarif by id", http.StatusInternalServerError, err.Error())
			return
		}

		amount := 0.0

		if cashierTarifResponse.TarifType == "fixed" {

			if salesResponse.PaymentType == "card" {
				amount = (cashierTarifResponse.AmountForCard)
			} else {
				amount = (cashierTarifResponse.AmountForCash)
			}

		} else {

			if salesResponse.
				PaymentType == "card" {
				amount = (cashierTarifResponse.AmountForCard) * totalPrice
			} else {
				amount = (cashierTarifResponse.AmountForCash) * totalPrice
			}

		}

		reqToUpdate := models.UpdateStaffBalanceAndCreateTransaction{
			UpdateCashierBalance: models.StaffInfo{
				StaffID: cashierResponse.ID,
				Amount:  amount,
			},
			SaleID:          id,
			TransactionType: "topup",
			SourceType:      "sales",
			Amount:          salesResponse.Price,
			Description:     "qwerty",
		}

		if salesResponse.ShopAssistantID != "" {

			if shopAssistantTarifResponse.TarifType == "fixed" {

				if salesResponse.PaymentType == "card" {
					amount = (shopAssistantTarifResponse.AmountForCard)
				} else {
					amount = (shopAssistantTarifResponse.AmountForCash)
				}

			} else {

				if salesResponse.PaymentType == "card" {
					amount = (shopAssistantTarifResponse.AmountForCard) * totalPrice
				} else {
					amount = (shopAssistantTarifResponse.AmountForCash) * totalPrice

				}

			}

			reqToUpdate.UpdateShopAssistantBalance.StaffID = shopAssistantResponse.ID
			
			reqToUpdate.UpdateShopAssistantBalance.Amount = amount
		}

		err = h.storage.Transaction().UpdateStaffBalanceAndCreateTransaction(context.Background(), reqToUpdate)
		if err != nil {
			handleResponse(c, "error while update cashoier balance", http.StatusInternalServerError, err.Error())
			return
		}

	}
}

