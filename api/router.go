package api

import (
	_ "bazaar/api/docs"
	"bazaar/api/handler"
	"bazaar/storage"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title           BAZAAR
// @version         1.0
// @description     An API for a store called BAZAAR
func New(store storage.IStorage) *gin.Engine {

	h := handler.New(store)

	r := gin.New()

	//BARCODE

	r.POST("barcode", h.Barcode)

	// BASKET

	r.POST("basket", h.CreateBasket)
	r.GET("basket/:id", h.GetBasketByID)
	r.GET("basket", h.GetBasketList)
	r.PUT("basket/:id", h.UpdateBasket)
	r.DELETE("basket/:id", h.DeleteBasket)

	// BRANCH

	r.POST("branch", h.CreateBranch)
	r.GET("branch/:id", h.GetBranchByID)
	r.GET("branch", h.GetBranchList)
	r.PUT("branch/:id", h.UpdateBranch)
	r.DELETE("branch/:id", h.DeleteBranch)

	// CATEGORY

	r.POST("category", h.CreateCategory)
	r.GET("category/:id", h.GetCategoryByID)
	r.GET("category", h.GetCategoryList)
	r.PUT("category/:id", h.UpdateCategory)
	r.DELETE("category/:id", h.DeleteCategory)

	// PRODUCT

	r.POST("product", h.CreateProduct)
	r.GET("product/:id", h.GetProductByID)
	r.GET("product", h.GetProductList)
	r.PUT("product/:id", h.UpdateProduct)
	r.DELETE("product/:id", h.DeleteProduct)

	// SALE

	r.POST("sale", h.CreateSale)
	r.GET("sale/:id", h.GetSaleByID)
	r.GET("sale", h.GetSaleList)
	r.PUT("sale/:id", h.UpdateSale)
	r.DELETE("sale/:id", h.DeleteSale)

	// STAFF

	r.POST("staff", h.CreateStaff)
	r.GET("staff/:id", h.GetStaffByID)
	r.GET("staff", h.GetStaffList)
	r.PUT("staff/:id", h.UpdateStaff)
	r.DELETE("staff/:id", h.DeleteStaff)

	// STORAGE-TRANSACTION

	r.POST("storage_transaction", h.CreateStorageTransaction)
	r.GET("storage_transaction/:id", h.GetStorageTransactionByID)
	r.GET("storage_transaction", h.GetStorageTransactionList)
	r.PUT("storage_transaction/:id", h.UpdateStorageTransaction)
	r.DELETE("storage_transaction/:id", h.DeleteStorageTransaction)

	// STORAGE

	r.POST("storage", h.CreateStorage)
	r.GET("storage/:id", h.GetStorageByID)
	r.GET("storage", h.GetStorageList)
	r.PUT("storage/:id", h.UpdateStorage)
	r.DELETE("storage/:id", h.DeleteStorage)

	// TARIF

	r.POST("tarif", h.CreateTarif)
	r.GET("tarif/:id", h.GetTarifByID)
	r.GET("tarif", h.GetTarifList)
	r.PUT("tarif/:id", h.UpdateTarif)
	r.DELETE("tarif/:id", h.DeleteTarif)

	// TRANSACTION

	r.POST("transaction", h.CreateTransaction)
	r.GET("transaction/:id", h.GetTransactionByID)
	r.GET("transaction", h.GetTransactionList)
	r.PUT("transaction/:id", h.UpdateTransaction)
	r.DELETE("transaction/:id", h.DeleteTransaction)

	// SELL

	r.POST("sell/", h.StartSell)
	r.PUT("end_sell/:id", h.EndSale)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
