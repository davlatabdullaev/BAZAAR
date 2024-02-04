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

func (h Handler) CreateStorage(c *gin.Context) {
	createStorage := models.CreateStorage{}

	if err := c.ShouldBindJSON(&createStorage); err != nil {
		handleResponse(c, "error while reading body from client", http.StatusBadRequest, err)
	}

	id, err := h.storage.Storage().Create(context.Background(), createStorage)
	if err != nil {
		handleResponse(c, "error while creating storage", http.StatusInternalServerError, err)
		return
	}

	storage, err := h.storage.Storage().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while get storage", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusCreated, storage)

}

func (h Handler) GetStorageByID(c *gin.Context) {
	var err error

	uid := c.Param("id")

	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "invalid uuid type ", http.StatusBadRequest, err.Error())
		return
	}

	storage, err := h.storage.Storage().Get(context.Background(), models.PrimaryKey{
		ID: id.String(),
	})
	if err != nil {
		handleResponse(c, "error while get storage by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, storage)

}

func (h Handler) GetStorageList(c *gin.Context) {

	var (
		page, limit = 1, 50
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

	response, err := h.storage.Storage().GetList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	if err != nil {
		handleResponse(c, "error while getting storage", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, response)

}

func (h Handler) UpdateStorage(c *gin.Context) {
	updateStorage := models.UpdateStorage{}

	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateStorage.ID = uid

	id, err := h.storage.Storage().Update(context.Background(), updateStorage)
	if err != nil {
		handleResponse(c, "error while reading body", http.StatusInternalServerError, err.Error())
		return
	}

	storage, err := h.storage.Sale().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while getting storage by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, storage)

}

func (h Handler) DeleteStorage(c *gin.Context) {

	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.Storage().Delete(context.Background(), id.String()); err != nil {
		handleResponse(c, "error while deleting storage by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "data succesfully deleted")

}
