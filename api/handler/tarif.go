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

func (h Handler) CreateTarif(c *gin.Context) {
	createTarif := models.CreateTarif{}

	if err := c.ShouldBindJSON(&createTarif); err != nil {
		handleResponse(c, "error while reading body from client", http.StatusBadRequest, err)
	}

	id, err := h.storage.Tarif().Create(context.Background(), createTarif)
	if err != nil {
		handleResponse(c, "error while creating tarif", http.StatusInternalServerError, err)
		return
	}

	tarif, err := h.storage.Tarif().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while get tarif", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusCreated, tarif)

}

func (h Handler) GetTarifByID(c *gin.Context) {
	var err error

	uid := c.Param("id")

	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "invalid uuid type ", http.StatusBadRequest, err.Error())
		return
	}

	tarif, err := h.storage.Tarif().Get(context.Background(), models.PrimaryKey{
		ID: id.String(),
	})
	if err != nil {
		handleResponse(c, "error while get tarif by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, tarif)

}

func (h Handler) GetTarifList(c *gin.Context) {

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

	response, err := h.storage.Tarif().GetList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	if err != nil {
		handleResponse(c, "error while getting tarif", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, response)

}

func (h Handler) UpdateTarif(c *gin.Context) {
	updateTarif := models.UpdateTarif{}

	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateTarif.ID = uid

	if err := c.ShouldBindJSON(&updateTarif); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Tarif().Update(context.Background(), updateTarif)
	if err != nil {
		handleResponse(c, "error while updating tarif", http.StatusInternalServerError, err.Error())
		return
	}

	tarif, err := h.storage.Tarif().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while get tarif by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, tarif)

}

func (h Handler) DeleteTarif(c *gin.Context) {

	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.Tarif().Delete(context.Background(), id.String()); err != nil {
		handleResponse(c, "error while deleting tarif by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "data succesfully deleted")

}
