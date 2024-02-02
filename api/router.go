package api

import (
	"bazaar/api/handler"
	"bazaar/storage"

	"github.com/gin-gonic/gin"
)

func New(store storage.IStorage) *gin.Engine {

	h := handler.New(store)

	r := gin.New()

	// BASKET

	r.POST("basket", h.CreateBasket)














	return r
}
