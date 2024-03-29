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

// CreateStorage godoc
// @Router       /storage [POST]
// @Summary      Create a new storage
// @Description  Create a new storage
// @Tags         storage
// @Accept       json
// @Produce      json
// @Param        storage  body  models.CreateStorage  true  "storage data"
// @Success      201  {object}  models.Storage
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateStorage(c *gin.Context) {
	createStorage := models.CreateStorage{}

	if err := c.ShouldBindJSON(&createStorage); err != nil {
		handleResponse(c, h.log, "error while reading body from client", http.StatusBadRequest, err)
	}

	id, err := h.storage.Storage().Create(context.Background(), createStorage)
	if err != nil {
		handleResponse(c, h.log, "error while creating storage", http.StatusInternalServerError, err)
		return
	}

	storage, err := h.storage.Storage().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, h.log, "error while get storage", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "", http.StatusCreated, storage)

}

// GetStorageByID godoc
// @Router       /storage/{id} [GET]
// @Summary      Get storage by id
// @Description  Get storage by id
// @Tags         storage
// @Accept       json
// @Produce      json
// @Param        id path string true "storage"
// @Success      200  {object}  models.Storage
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetStorageByID(c *gin.Context) {
	var err error

	uid := c.Param("id")

	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, h.log, "invalid uuid type ", http.StatusBadRequest, err.Error())
		return
	}

	storage, err := h.storage.Storage().Get(context.Background(), models.PrimaryKey{
		ID: id.String(),
	})
	if err != nil {
		handleResponse(c, h.log, "error while get storage by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, storage)

}

// GetStoragesList godoc
// @Router       /storage [GET]
// @Summary      Get storages list
// @Description  Get storages list
// @Tags         storage
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        search query string false "search"
// @Success      200  {object}  models.StoragesResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetStorageList(c *gin.Context) {

	var (
		page, limit = 1, 50
		search      string
		err         error
	)

	pageStr := c.DefaultQuery("page", "1")
	page, err = strconv.Atoi(pageStr)
	if err != nil {
		handleResponse(c, h.log, "error while parsing page ", http.StatusBadRequest, err.Error())
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		handleResponse(c, h.log, "error while parsing limit", http.StatusBadRequest, err.Error())
		return
	}

	search = c.Query("search")

	response, err := h.storage.Storage().GetList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	if err != nil {
		handleResponse(c, h.log, "error while getting storage", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, response)

}

// UpdateStorage godoc
// @Router       /storage/{id} [PUT]
// @Summary      Update storage by id
// @Description  Update storage by id
// @Tags         storage
// @Accept       json
// @Produce      json
// @Param        id path string true "storage id"
// @Param        storage body models.UpdateStorage true "storage data"
// @Success      200  {object}  models.Storage
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateStorage(c *gin.Context) {
	updateStorage := models.UpdateStorage{}

	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, h.log, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	if err := c.ShouldBindJSON(&updateStorage); err != nil {
		handleResponse(c, h.log, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	updateStorage.ID = uid

	id, err := h.storage.Storage().Update(context.Background(), updateStorage)
	if err != nil {
		handleResponse(c, h.log, "error while reading body", http.StatusInternalServerError, err.Error())
		return
	}

	storage, err := h.storage.Storage().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, h.log, "error while getting storage by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, storage)

}

// DeleteStorageTransaction godoc
// @Router       /storage/{id} [DELETE]
// @Summary      Delete Storage
// @Description  Delete Storage
// @Tags         storage
// @Accept       json
// @Produce      json
// @Param        id path string true "storage id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteStorage(c *gin.Context) {

	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, h.log, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.Storage().Delete(context.Background(), id.String()); err != nil {
		handleResponse(c, h.log, "error while deleting storage by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, "data succesfully deleted")

}
