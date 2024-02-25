package handler

import (
	"bazaar/api/models"
	"bazaar/pkg/logger"
	"bazaar/storage"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	storage storage.IStorage
	log     logger.ILogger
}

func New(store storage.IStorage, log logger.ILogger) Handler {
	return Handler{
		storage: store,
		log:     log,
	}
}

func handleResponse(c *gin.Context, log logger.ILogger, msg string, statusCode int, data interface{}) {
	response := models.Response{}

	switch code := statusCode; {
	case code < 400:
		response.Description = "OK"
		log.Info("~~~~> OK", logger.String("msg", msg), logger.Any("status", code))
	case code == 401:
		response.Description = "Unauthorized"
	case code < 500:
		response.Description = "bad request"
	default:
		response.Description = "internal server error"

	}

	response.StatusCode = statusCode
	response.Data = data

	c.JSON(response.StatusCode, response)

}
