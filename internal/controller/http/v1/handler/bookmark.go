package handler

import (
	"strconv"

	"github.com/Akorm0181/yelp/config"
	"github.com/Akorm0181/yelp/internal/entity"
	"github.com/gin-gonic/gin"
)

// CreateBookmark godoc
// @Router /bookmark [post]
// @Summary Create a new bookmark
// @Description Create a new bookmark
// @Security BearerAuth
// @Tags bookmark
// @Accept  json
// @Produce  json
// @Param bookmark body entity.Bookmark true "Bookmark object"
// @Success 201 {object} entity.Bookmark
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) CreateBookmark(ctx *gin.Context) {
	var (
		body entity.Bookmark
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	body.UserID = ctx.GetHeader("sub")

	bookmark, err := h.UseCase.BookmarkRepo.Create(ctx, body)
	if h.HandleDbError(ctx, err, "Error creating bookmark") {
		return
	}

	ctx.JSON(201, bookmark)
}

// GetUser godoc
// @Router /bookmark/{id} [get]
// @Summary Get a bookmark by ID
// @Description Get a bookmark by ID
// @Security BearerAuth
// @Tags bookmark
// @Accept  json
// @Produce  json
// @Param id path string true "Bookmark ID"
// @Success 200 {object} entity.Bookmark
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetBookmark(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	bookmark, err := h.UseCase.BookmarkRepo.GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting bookmark") {
		return
	}

	ctx.JSON(200, bookmark)
}

// GetUsers godoc
// @Router /bookmark/list [get]
// @Summary Get a list of bookmarks
// @Description Get a list of bookmarks
// @Security BearerAuth
// @Tags bookmark
// @Accept  json
// @Produce  json
// @Param page query number true "page"
// @Param limit query number true "limit"
// @Param search query string false "search"
// @Success 200 {object} entity.BookmarksList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetBookmarks(ctx *gin.Context) {
	var (
		req entity.GetListFilter
	)

	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")
	search := ctx.DefaultQuery("search", "")

	req.Page, _ = strconv.Atoi(page)
	req.Limit, _ = strconv.Atoi(limit)
	req.Filters = append(req.Filters,
		entity.Filter{
			Column: "business_id",
			Type:   "search",
			Value:  search,
		},
	)

	req.OrderBy = append(req.OrderBy, entity.OrderBy{
		Column: "created_at",
		Order:  "desc",
	})

	users, err := h.UseCase.BookmarkRepo.GetList(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting bookmarks") {
		return
	}

	ctx.JSON(200, users)
}

// UpdateBookmark godoc
// @Router /bookmark [put]
// @Summary Update a bookmark
// @Description Update a bookmark
// @Security BearerAuth
// @Tags bookmark
// @Accept  json
// @Produce  json
// @Param bookmark body entity.Bookmark true "Bookmark object"
// @Success 200 {object} entity.Bookmark
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateBookmark(ctx *gin.Context) {
	var (
		body entity.Bookmark
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	bookmark, err := h.UseCase.BookmarkRepo.Update(ctx, body)
	if h.HandleDbError(ctx, err, "Error updating bookmark") {
		return
	}

	ctx.JSON(200, bookmark)
}

// DeleteUser godoc
// @Router /bookmark/{id} [delete]
// @Summary Delete a bookmark
// @Description Delete a bookmark
// @Security BearerAuth
// @Tags bookmark
// @Accept  json
// @Produce  json
// @Param id path string true "Bookmark ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeleteBookmark(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	if ctx.GetHeader("user_type") == "user" {
		req.ID = ctx.GetHeader("sub")
	}

	err := h.UseCase.BookmarkRepo.Delete(ctx, req)
	if h.HandleDbError(ctx, err, "Error deleting user") {
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Bookmark deleted successfully",
	})
}
