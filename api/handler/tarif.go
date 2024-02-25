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

// CreateTarif godoc
// @Router       /tarif [POST]
// @Summary      Create a new tarif
// @Description  Create a new tarif
// @Tags         tarif
// @Accept       json
// @Produce      json
// @Param        tarif  body  models.CreateTarif  true  "tarif data"
// @Success      201  {object}  models.Tarif
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateTarif(c *gin.Context) {
	createTarif := models.CreateTarif{}

	if err := c.ShouldBindJSON(&createTarif); err != nil {
		handleResponse(c, h.log, "error while reading body from client", http.StatusBadRequest, err)
	}

	id, err := h.storage.Tarif().Create(context.Background(), createTarif)
	if err != nil {
		handleResponse(c, h.log, "error while creating tarif", http.StatusInternalServerError, err)
		return
	}

	tarif, err := h.storage.Tarif().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, h.log, "error while get tarif", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "", http.StatusCreated, tarif)

}

// GetTarifByID godoc
// @Router       /tarif/{id} [GET]
// @Summary      Get tarif by id
// @Description  Get tarif by id
// @Tags         tarif
// @Accept       json
// @Produce      json
// @Param        id path string true "tarif"
// @Success      200  {object}  models.Tarif
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetTarifByID(c *gin.Context) {
	var err error

	uid := c.Param("id")

	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, h.log, "invalid uuid type ", http.StatusBadRequest, err.Error())
		return
	}

	tarif, err := h.storage.Tarif().Get(context.Background(), models.PrimaryKey{
		ID: id.String(),
	})
	if err != nil {
		handleResponse(c, h.log, "error while get tarif by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, tarif)

}

// GetTarifsList godoc
// @Router       /tarif [GET]
// @Summary      Get tarifs list
// @Description  Get tarifs list
// @Tags         tarif
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        search query string false "search"
// @Success      200  {object}  models.TarifsResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetTarifList(c *gin.Context) {

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

	response, err := h.storage.Tarif().GetList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	if err != nil {
		handleResponse(c, h.log, "error while getting tarif", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, response)

}

// UpdateTarif godoc
// @Router       /tarif/{id} [PUT]
// @Summary      Update tarif by id
// @Description  Update tarif by id
// @Tags         tarif
// @Accept       json
// @Produce      json
// @Param        id path string true "tarif id"
// @Param        tarif body models.UpdateTarif true "tarif"
// @Success      200  {object}  models.Tarif
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateTarif(c *gin.Context) {
	updateTarif := models.UpdateTarif{}

	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, h.log, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateTarif.ID = uid

	if err := c.ShouldBindJSON(&updateTarif); err != nil {
		handleResponse(c, h.log, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Tarif().Update(context.Background(), updateTarif)
	if err != nil {
		handleResponse(c, h.log, "error while updating tarif", http.StatusInternalServerError, err.Error())
		return
	}

	tarif, err := h.storage.Tarif().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, h.log, "error while get tarif by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, tarif)

}

// DeleteTarif godoc
// @Router       /tarif/{id} [DELETE]
// @Summary      Delete Tarif
// @Description  Delete Tarif
// @Tags         tarif
// @Accept       json
// @Produce      json
// @Param        id path string true "tarif id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteTarif(c *gin.Context) {

	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, h.log, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.Tarif().Delete(context.Background(), id.String()); err != nil {
		handleResponse(c, h.log, "error while deleting tarif by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, "data succesfully deleted")

}
