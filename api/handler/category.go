package handler

import (
	"bazaar/api/models"
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateBasket godoc
// @Router       /category [POST]
// @Summary      Create a new category
// @Description  Create a new category
// @Tags         CATEGORY
// @Accept       json
// @Produce      json
// @Param        category  body  models.CreateCategory  true  "category data"
// @Success      201  {object}  models.Category.ID
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateCategory(c *gin.Context) {
	createCategory := models.CreateCategory{}

	if err := c.ShouldBindJSON(&createCategory); err != nil {
		handleResponse(c, "error while reading body from client", http.StatusBadRequest, err)
	}

	id, err := h.storage.Category().Create(context.Background(), createCategory)
	if err != nil {
		handleResponse(c, "error while creating category", http.StatusInternalServerError, err)
		return
	}

	category, err := h.storage.Category().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while get category", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusCreated, category)

}

func (h Handler) GetCategoryByID(c *gin.Context) {
	var err error

	id := c.Param("id")

	category, err := h.storage.Category().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while get category by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, category)

}

func (h Handler) GetCategoryList(c *gin.Context) {

	var (
		page, limit int
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

	response, err := h.storage.Category().GetList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	if err != nil {
		handleResponse(c, "error while getting category", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, response)

}

func (h Handler) UpdateCategory(c *gin.Context) {
	updateCategory := models.UpdateCategory{}

	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateCategory.ID = uid

	if err := c.ShouldBindJSON(&updateCategory); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Category().Update(context.Background(), updateCategory)
	if err != nil {
		handleResponse(c, "error while updating category", http.StatusInternalServerError, err.Error())
		return
	}

	category, err := h.storage.Category().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while get category by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, category)

}

func (h Handler) DeleteCategory(c *gin.Context) {

	uid := c.Param("id")

	if err := h.storage.Category().Delete(context.Background(), uid); err != nil {
		handleResponse(c, "error while deleting category by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "data succesfully deleted")

}
