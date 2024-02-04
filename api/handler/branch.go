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

// CreateBranch godoc
// @Router       /branch [POST]
// @Summary      Create a new branch
// @Description  Create a new branch
// @Tags         branch
// @Accept       json
// @Produce      json
// @Param        branch  body  models.CreateBranch  true  "branch data"
// @Success      201  {object}  models.Branch
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateBranch(c *gin.Context) {
	createBranch := models.CreateBranch{}

	if err := c.ShouldBindJSON(&createBranch); err != nil {
		handleResponse(c, "error while reading body from client", http.StatusBadRequest, err)
	}

	id, err := h.storage.Branch().Create(context.Background(), createBranch)
	if err != nil {
		handleResponse(c, "error while creating branch", http.StatusInternalServerError, err)
		return
	}

	branch, err := h.storage.Branch().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while get branch", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusCreated, branch)

}

// GetBranchByID godoc
// @Router       /branch/{id} [GET]
// @Summary      Get branch by id
// @Description  Get branch by id
// @Tags         branch
// @Accept       json
// @Produce      json
// @Param        id path string true "branch"
// @Success      200  {object}  models.Branch
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetBranchByID(c *gin.Context) {

	var err error

	uid := c.Param("id")

	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "invalid uuid type ", http.StatusBadRequest, err.Error())
		return
	}

	branch, err := h.storage.Branch().Get(context.Background(), models.PrimaryKey{
		ID: id.String(),
	})
	if err != nil {
		handleResponse(c, "error while get branch by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, branch)

}

// GetBranchsList godoc
// @Router       /branch [GET]
// @Summary      Get branchs list
// @Description  Get branchs list
// @Tags         branch
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        search query string false "search"
// @Success      200  {object}  models.BranchsResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetBranchList(c *gin.Context) {

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

	response, err := h.storage.Branch().GetList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	if err != nil {
		handleResponse(c, "error while getting branch", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, response)

}

// UpdateBranch godoc
// @Router       /branch/{id} [PUT]
// @Summary      Update branch by id
// @Description  Update branch by id
// @Tags         branch
// @Accept       json
// @Produce      json
// @Param        id path string true "branch id"
// @Param        branch body models.UpdateBranch true "basket"
// @Success      200  {object}  models.Branch
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateBranch(c *gin.Context) {
	updateBranch := models.UpdateBranch{}

	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateBranch.ID = uid

	if err := c.ShouldBindJSON(&updateBranch); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Branch().Update(context.Background(), updateBranch)
	if err != nil {
		handleResponse(c, "error while updating branch", http.StatusInternalServerError, err.Error())
		return
	}

	branch, err := h.storage.Branch().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while getting branch by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, branch)

}

// DeleteBranch godoc
// @Router       /branch/{id} [DELETE]
// @Summary      Delete Branch
// @Description  Delete Branch
// @Tags         branch
// @Accept       json
// @Produce      json
// @Param        id path string true "branch id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteBranch(c *gin.Context) {

	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.Branch().Delete(context.Background(), id.String()); err != nil {
		handleResponse(c, "error while deleting branch by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "data succesfully deleted")

}
