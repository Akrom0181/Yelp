package handler

import (
	"strconv"

	"github.com/Akorm0181/yelp/config"
	"github.com/Akorm0181/yelp/internal/entity"
	"github.com/gin-gonic/gin"
)

// CreatePromotion godoc
// @Router /promotion [post]
// @Summary Create a new promotion
// @Description Create a new promotion
// @Security BearerAuth
// @Tags promotion
// @Accept  json
// @Produce  json
// @Param promotion body entity.Promotion true "Promotion object"
// @Success 201 {object} entity.Promotion
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) CreatePromotion(ctx *gin.Context) {
	var (
		body entity.Promotion
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	body.UserID = ctx.GetHeader("sub")

	if body.ExpiresAt.Before(body.StartedAt) {
		h.ReturnError(ctx, config.ErrorBadRequest, "End date must be after start date", 400)
		return
	}

	if body.DiscountPercentage < 0 && body.DiscountPercentage >= 100 {
		h.ReturnError(ctx, config.ErrorBadRequest, "Discount percentage must be between 0 and 100", 400)
		return
	}

	promotion, err := h.UseCase.PromotionRepo.Create(ctx, body)
	if h.HandleDbError(ctx, err, "Error creating promotion") {
		return
	}

	ctx.JSON(201, promotion)
}

// GetPromotions godoc
// @Router /promotion/{id} [get]
// @Summary Get a promotion by ID
// @Description Get a promotion by ID
// @Security BearerAuth
// @Tags promotion
// @Accept  json
// @Produce  json
// @Param id path string true "Promotion ID"
// @Success 200 {object} entity.Promotion
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetPromotion(ctx *gin.Context) {
	var (
		req entity.PromotionSingleRequest
	)

	req.ID = ctx.Param("id")

	promotion, err := h.UseCase.PromotionRepo.GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting promotion") {
		return
	}

	ctx.JSON(200, promotion)
}

// GetPromotions godoc
// @Router /promotion/list [get]
// @Summary Get a list of promotions
// @Description Get a list of promotions
// @Security BearerAuth
// @Tags promotion
// @Accept  json
// @Produce  json
// @Param page query number true "page"
// @Param limit query number true "limit"
// @Param search query string false "search"
// @Success 200 {object} entity.PromotionGetList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetPromotions(ctx *gin.Context) {
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
			Column: "description",
			Type:   "search",
			Value:  search,
		},
		entity.Filter{
			Column: "title",
			Type:   "search",
			Value:  search,
		},
	)

	req.OrderBy = append(req.OrderBy, entity.OrderBy{
		Column: "created_at",
		Order:  "desc",
	})

	promotions, err := h.UseCase.PromotionRepo.GetList(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting promotions") {
		return
	}

	ctx.JSON(200, promotions)
}

// DeletePromotion godoc
// @Router /promotion/{id} [delete]
// @Summary Delete a promotion
// @Description Delete a promotion
// @Security BearerAuth
// @Tags promotion
// @Accept  json
// @Produce  json
// @Param id path string true "Promotion ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeletePromotion(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	if ctx.GetHeader("user_type") == "user" {
		req.ID = ctx.GetHeader("sub")
	}

	err := h.UseCase.UserRepo.Delete(ctx, req)
	if h.HandleDbError(ctx, err, "Error deleting promotion") {
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Promotion deleted successfully",
	})
}
