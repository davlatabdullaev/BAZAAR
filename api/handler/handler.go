package handler

import (
	"bazaar/api/models"
	"bazaar/storage"
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler struct {
	storage storage.IStorage
}

func New(store storage.IStorage) Handler {
	return Handler{
		storage: store,
	}
}

func handleResponse(w http.ResponseWriter, statusCode int, data interface{}) {
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

	js, err := json.Marshal(response)
	if err != nil {
		fmt.Println("error while marshalling json", err.Error())
		return
	}

	w.WriteHeader(statusCode)
	w.Write(js)

}
