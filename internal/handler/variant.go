package handler

import (
	apiHelper "basic-trade/api/helper"
	"basic-trade/api/request"
	"basic-trade/common"
	"basic-trade/internal/entity"
	"basic-trade/internal/service"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VariantHandler struct {
	variantService service.VariantService
}

func NewVariantHandler(variantService service.VariantService) *VariantHandler {
	return &VariantHandler{variantService: variantService}
}

func (h *VariantHandler) CreateVariant(ctx *gin.Context) {
	apiHelper.ResponseHandler(ctx, func(c context.Context, resChan chan apiHelper.ResponseData) {
		var req request.CreateVariantRequest
		if err := ctx.ShouldBind(&req); err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusBadRequest,
				Error:      common.ErrorValidation(err),
			}
		}

		arg := entity.Variant{
			VariantName: req.VariantName,
			Quantity:    req.Quantity,
		}

		uuid, err := uuid.Parse(req.ProductID)
		if err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusBadRequest,
				Error:      err,
			}
		}

		result, err := h.variantService.CreateVariant(c, arg, uuid)
		if err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusInternalServerError,
				Error:      err,
			}
		}

		resChan <- apiHelper.ResponseData{
			StatusCode: http.StatusCreated,
			Message:    "Variant created successfully.",
			Data:       result,
		}
	})
}

func (h *VariantHandler) GetVariant(ctx *gin.Context) {
	apiHelper.ResponseHandler(ctx, func(c context.Context, resChan chan apiHelper.ResponseData) {
		var req request.GetDataByUUIDRequest
		if err := ctx.ShouldBindUri(&req); err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusBadRequest,
				Error:      common.ErrorValidation(err),
			}
		}

		uuid, err := uuid.Parse(req.UUID)
		if err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusBadRequest,
				Error:      err,
			}
		}

		result, err := h.variantService.GetVariant(c, uuid)
		if err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusInternalServerError,
				Error:      err,
			}
		}

		resChan <- apiHelper.ResponseData{
			StatusCode: http.StatusOK,
			Message:    "Variant retrieved successfully.",
			Data:       result,
		}
	})
}

func (h *VariantHandler) GetAllVariants(ctx *gin.Context) {
	apiHelper.ResponseHandler(ctx, func(c context.Context, resChan chan apiHelper.ResponseData) {
		var req request.PaginationRequest
		if err := ctx.ShouldBindQuery(&req); err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusBadRequest,
				Error:      common.ErrorValidation(err),
			}
		}

		result, err := h.variantService.GetAllVariants(c, req.Offset, req.Limit)
		if err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusInternalServerError,
				Error:      err,
			}
		}

		resChan <- apiHelper.ResponseData{
			StatusCode: http.StatusOK,
			Message:    "Variants retrieved successfully.",
			Data:       result,
		}
	})
}

func (h *VariantHandler) SearchVariants(ctx *gin.Context) {
	apiHelper.ResponseHandler(ctx, func(c context.Context, resChan chan apiHelper.ResponseData) {
		var req request.SearchRequest
		if err := ctx.ShouldBindQuery(&req); err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusBadRequest,
				Error:      common.ErrorValidation(err),
			}
		}

		result, err := h.variantService.SearchVariants(c, req.Keyword, req.Offset, req.Limit)
		if err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusInternalServerError,
				Error:      err,
			}
		}

		resChan <- apiHelper.ResponseData{
			StatusCode: http.StatusOK,
			Message:    "Variants retrieved successfully.",
			Data:       result,
		}
	})
}

func (h *VariantHandler) UpdateVariant(ctx *gin.Context) {
	apiHelper.ResponseHandler(ctx, func(c context.Context, resChan chan apiHelper.ResponseData) {
		var idReq request.GetDataByUUIDRequest
		if err := ctx.ShouldBindUri(&idReq); err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusBadRequest,
				Error:      common.ErrorValidation(err),
			}
		}

		var variantReq request.UpdateVariantRequest
		if err := ctx.ShouldBind(&variantReq); err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusBadRequest,
				Error:      common.ErrorValidation(err),
			}
		}

		uuid, err := uuid.Parse(idReq.UUID)
		if err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusBadRequest,
				Error:      err,
			}
		}

		arg := entity.Variant{
			UUID:        uuid,
			VariantName: variantReq.VariantName,
			Quantity:    variantReq.Quantity,
		}

		result, err := h.variantService.UpdateVariant(c, arg)
		if err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusInternalServerError,
				Error:      err,
			}
		}

		resChan <- apiHelper.ResponseData{
			StatusCode: http.StatusOK,
			Message:    "Variant updated successfully.",
			Data:       result,
		}
	})
}

func (h *VariantHandler) DeleteVariant(ctx *gin.Context) {
	apiHelper.ResponseHandler(ctx, func(c context.Context, resChan chan apiHelper.ResponseData) {
		var req request.GetDataByUUIDRequest
		if err := ctx.ShouldBindUri(&req); err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusBadRequest,
				Error:      common.ErrorValidation(err),
			}
		}

		uuid, err := uuid.Parse(req.UUID)
		if err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusBadRequest,
				Error:      err,
			}
		}

		err = h.variantService.DeleteVariant(c, uuid)
		if err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusNoContent,
				Error:      err,
			}
		}

		resChan <- apiHelper.ResponseData{
			StatusCode: http.StatusOK,
			Message:    "Variants deleted successfully.",
			Data:       nil,
		}
	})
}
