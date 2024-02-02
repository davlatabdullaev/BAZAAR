package handler

import (
	"bazaar/api/models"
	"bazaar/storage"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	storage storage.IStorage
}

func New(store storage.IStorage) Handler {
	return Handler{
		storage: store,
	}
}

func handleResponse(c *gin.Context, msg string, statusCode int, data interface{}) {
	response := models.Response{}

	switch code := statusCode; {
	case code < 400:
		response.Description = "succes"
	case code < 500:
		response.Description = "bad request"
	default:
		response.Description = "internal server error"

	}

	response.StatusCode = statusCode
	response.Data = data

	c.JSON(response.StatusCode, response)

}
