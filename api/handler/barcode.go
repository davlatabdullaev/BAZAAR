package handler

import (
	"context"
	"net/http"
	"bazaar/api/models"

	"github.com/gin-gonic/gin"
)

// Barcode godoc
// @Router       /barcode [POST]
// @Summary      barcode
// @Description  barcode
// @Tags         barcode
// @Accept       json
// @Produce      json
// @Param		 info body models.Barcode true "info"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) Barcode(c *gin.Context) {
	info := models.Barcode{}
	if err := c.ShouldBindJSON(&info); err != nil {
		handleResponse(c, h.log, "error is while reading body", http.StatusBadRequest, err.Error())
		return
	}

	sale, err := h.storage.Sale().Get(context.Background(), models.PrimaryKey{
		ID: info.SaleID,
	})
	if err != nil {
		handleResponse(c, h.log, "error is getting sale by id", http.StatusInternalServerError, err.Error())
		return
	}

	if sale.Status == "success" {
		handleResponse(c, h.log, "sale ended", 300, "sale ended cannot add product")
		return
	}

	if sale.Status == "cancel" {
		handleResponse(c, h.log, "sale canceled", 300, "sale canceled cannot add product")
		return
	}

	products, err := h.storage.Product().GetList(context.Background(), models.ProductGetListRequest{
		Page:    1,
		Limit:   10,
		Barcode: info.Barcode,
	})

	if err != nil {
		handleResponse(c, h.log, "error is while getting product list by barcode", http.StatusInternalServerError, err.Error())
		return
	}

	var (
		prodID    string
		prodPrice int
	)
	for _, product := range products.Products {
		prodID = product.ID
		prodPrice = int(product.Price)
	}

	baskets, err := h.storage.Basket().GetList(context.Background(), models.GetBasketsListRequest{
		Page:   1,
		Limit:  10,
		Search: info.SaleID,
	})
	if err != nil {
		handleResponse(c, h.log, "error is while getting basket list", http.StatusInternalServerError, err.Error())
		return
	}

	var (
		basketsMap = make(map[string]models.Basket)
		totalPrice = 0
	)

	totalPrice = info.Count * prodPrice

	for _, basket := range baskets.Baskets {
		basketsMap[basket.ProductID] = basket
	}

	storage, err := h.storage.Storage().GetList(context.Background(), models.GetListRequest{
		Page:  1,
		Limit: 100,
	})
	if err != nil {
		handleResponse(c, h.log, "error is while getting repo list", http.StatusInternalServerError, err.Error())
		return
	}

	for _, r := range storage.Storages {
		if prodID == basketsMap[r.ProductID].ProductID {
			if r.Count < (basketsMap[r.ProductID].Quantity + info.Count) {
				handleResponse(c, h.log, "not enough product", 301, "not enough product")
				return
			}
		}

		if r.Count < info.Count {
			handleResponse(c, h.log, "not enough product", 300, "not enough product")
			return
		}
	}

	isTrue := false

	for _, value := range basketsMap {
		if prodID == value.ProductID {
			isTrue = true
			id, err := h.storage.Basket().Update(context.Background(), models.UpdateBasket{
				ID:        value.ID,
				SaleID:    value.SaleID,
				ProductID: prodID,
				Quantity:  value.Quantity + info.Count,
				Price:     value.Price + float64(totalPrice),
			})
			if err != nil {
				handleResponse(c, h.log, "error is while updating basket", 500, err.Error())
				return
			}
			updatedBasket, err := h.storage.Basket().Get(context.Background(), models.PrimaryKey{ID: id})
			if err != nil {
				handleResponse(c, h.log, "error is while getting basket", 500, err.Error())
				return
			}
			handleResponse(c, h.log, "updated", http.StatusOK, updatedBasket)
		}
	}

	if !isTrue {
		id, err := h.storage.Basket().Create(context.Background(), models.CreateBasket{
			SaleID:    info.SaleID,
			ProductID: prodID,
			Quantity:  info.Count,
			Price:     float64(totalPrice),
		})
		if err != nil {
			handleResponse(c, h.log, "error is while creating basket", 500, err.Error())
			return
		}
		createdBasket, err := h.storage.Basket().Get(context.Background(), models.PrimaryKey{ID: id})
		if err != nil {
			handleResponse(c, h.log, "error is while getting basket", 500, err.Error())
			return
		}
		handleResponse(c, h.log, "updated", http.StatusOK, createdBasket)
	}
}


/*



sale_id, barcode, count


barcode orqali prod get

count*price=batot

ba = cou

sal = said



*/
