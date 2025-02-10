package handler

import (
	"log"
	"strconv"
	"time"

	"github.com/Akorm0181/yelp/config"
	"github.com/Akorm0181/yelp/internal/entity"
	"github.com/Akorm0181/yelp/pkg/firebase"
	"github.com/gin-gonic/gin"
)

// CreateBusiness godoc
// @Router /business [post]
// @Summary Create a new business
// @Description Create a new business
// @Security BearerAuth
// @Tags business
// @Accept  json
// @Produce  json
// @Param business body entity.Business true "Business object"
// @Success 201 {object} entity.Business
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) CreateBusiness(ctx *gin.Context) {
	var (
		body entity.Business
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	body.OwnerID = ctx.GetHeader("sub")

	business, err := h.UseCase.BusinessRepo.Create(ctx, body)
	if h.HandleDbError(ctx, err, "Error creating business") {
		return
	}

	business.Attachments, err = h.UseCase.BusinessAttachmentRepo.MultipleUpsert(ctx, entity.BusinessAttachmentMultipleInsertRequest{
		BusinessId:  business.ID,
		Attachments: body.Attachments,
	})
	if h.HandleDbError(ctx, err, "Error creating business") {
		return
	}

	log.Println("Inserted attachments:", business.Attachments)

	ctx.JSON(201, business)
}

// GetBusiness godoc
// @Router /business/{id} [get]
// @Summary Get a business by ID
// @Description Get a business by ID
// @Security BearerAuth
// @Tags business
// @Accept  json
// @Produce  json
// @Param id path string true "Business ID"
// @Success 200 {object} entity.Business
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetBusiness(ctx *gin.Context) {
	var (
		req entity.BusinessSingleRequest
	)

	req.ID = ctx.Param("id")

	business, err := h.UseCase.BusinessRepo.GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting business") {
		return
	}

	businessAttachments, err := h.UseCase.BusinessAttachmentRepo.GetList(ctx,
		entity.GetListFilter{
			Filters: []entity.Filter{
				{
					Column: "business_id",
					Type:   "eq",
					Value:  req.ID,
				},
			},
			Page:  1,
			Limit: 10,
		},
	)

	if h.HandleDbError(ctx, err, "Error getting tweet attachments") {
		return
	}

	business.Attachments = businessAttachments.Items

	ctx.JSON(200, business)
}

// GetBusinesss godoc
// @Router /business/list [get]
// @Summary Get a list of businesses
// @Description Get a list of businesses
// @Security BearerAuth
// @Tags business
// @Accept  json
// @Produce  json
// @Param page query number true "page"
// @Param limit query number true "limit"
// @Param search query string false "search"
// @Success 200 {object} entity.BusinessList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetBusinesses(ctx *gin.Context) {
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
			Column: "name",
			Type:   "search",
			Value:  search,
		},
		entity.Filter{
			Column: "address",
			Type:   "search",
			Value:  search,
		},
		entity.Filter{
			Column: "description",
			Type:   "search",
			Value:  search,
		},
	)

	req.OrderBy = append(req.OrderBy, entity.OrderBy{
		Column: "created_at",
		Order:  "desc",
	})

	businesses, err := h.UseCase.BusinessRepo.GetList(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting businesses") {
		return
	}

	ctx.JSON(200, businesses)
}

// UpdateBusiness godoc
// @Router /business [put]
// @Summary Update a business
// @Description Update a business
// @Security BearerAuth
// @Tags business
// @Accept  json
// @Produce  json
// @Param business body entity.Business true "Business object"
// @Success 200 {object} entity.Business
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateBusiness(ctx *gin.Context) {
	var (
		body entity.Business
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	if ctx.GetHeader("sub") != body.OwnerID || ctx.GetHeader("user_type") != "admin" {
		h.ReturnError(ctx, config.ErrorForbidden, "Access denied, only owner or admin can update business", 403)
		return
	}

	business, err := h.UseCase.BusinessRepo.Update(ctx, body)
	if h.HandleDbError(ctx, err, "Error updating business") {
		return
	}

	business.Attachments, err = h.UseCase.BusinessAttachmentRepo.MultipleUpsert(ctx, entity.BusinessAttachmentMultipleInsertRequest{
		BusinessId:  business.ID,
		Attachments: body.Attachments,
	})
	if h.HandleDbError(ctx, err, "Error upserting business attachments") {
		return
	}

	ctx.JSON(200, business)
}

// DeleteBusiness godoc
// @Router /business/{id} [delete]
// @Summary Delete a business
// @Description Delete a business
// @Security BearerAuth
// @Tags business
// @Accept  json
// @Produce  json
// @Param id path string true "Business ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeleteBusiness(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	res, err := h.UseCase.BusinessRepo.GetSingle(ctx, entity.BusinessSingleRequest{ID: req.ID})
	if err != nil {
		h.ReturnError(ctx, config.ErrorNotFound, "Business not found", 404)
		return
	}

	if res.OwnerID != ctx.GetHeader("sub") && ctx.GetHeader("user_type") != "admin" {
		h.ReturnError(ctx, config.ErrorForbidden, "Access denied, only owner or admin can delete business", 403)
		return
	}

	err = h.UseCase.BusinessRepo.Delete(ctx, req)
	if h.HandleDbError(ctx, err, "Error deleting business") {
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Business deleted successfully",
	})
}

// UploadBusinessPic godoc
// @ID upload_business_pic_file
// @Router /business/upload/{id} [post]
// @Summary Upload Multiple Files
// @Description Upload Multiple Files
// @Security BearerAuth
// @Tags business
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Business attachment id"
// @Param file formData []file true "File to upload"
// @Success 200 {object} entity.MultipleFileUploadResponse "Success Request"
// @Failure 400 {object} entity.ErrorResponse "Bad Request"
// @Failure 500 {object} entity.ErrorResponse "Server error"
func (h *Handler) UploadBusinessPic(ctx *gin.Context) {
	var (
		param_id = ctx.Param("id")
	)

	log.Print(param_id)

	form, err := ctx.MultipartForm()
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid file upload request", 400)
		return
	}

	resp, err := firebase.UploadFiles(form)
	if h.HandleDbError(ctx, err, "Error uploading files") {
		return
	}

	log.Print(resp.Url[0].Url)
	_, err = h.UseCase.BusinessAttachmentRepo.Update(ctx, entity.BusinessAttachment{
		Id:        param_id,
		FilePath:  resp.Url[0].Url,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	})

	if h.HandleDbError(ctx, err, "Error updating business attachments") {
		return
	}

	ctx.JSON(200, resp)
}
