package handler

import (
	"basic-trade/api/request"
	response "basic-trade/api/response"
	"basic-trade/internal/entity"
	"basic-trade/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VariantHandler struct {
	variantService service.IVariantService
}

func NewVariantHandler(variantService service.IVariantService) *VariantHandler {
	return &VariantHandler{variantService: variantService}
}

type createVariantRequest struct {
	VariantName string `form:"variant_name" binding:"required,max=100"`
	Quantity    int32  `form:"quantity" binding:"required"`
	ProductID   string `form:"product_id" binding:"required,validUUID"`
}

type updateVariantRequest struct {
	VariantName string `form:"variant_name" binding:"max=100"`
	Quantity    int32  `form:"quantity" binding:"required"`
}

func (h *VariantHandler) CreateVariant(ctx *gin.Context) {
	var req createVariantRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}

	arg := entity.Variant{
		VariantName: req.VariantName,
		Quantity:    req.Quantity,
	}

	uuid, err := uuid.Parse(req.ProductID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}

	result, err := h.variantService.CreateVariant(arg, uuid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *VariantHandler) GetVariant(ctx *gin.Context) {
	var req request.GetDataByUUIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}

	uuid, err := uuid.Parse(req.UUID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}

	result, err := h.variantService.GetVariant(uuid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *VariantHandler) GetAllVariants(ctx *gin.Context) {
	var req request.PaginationRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}

	result, err := h.variantService.GetAllVariants(req.Offset, req.Limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *VariantHandler) SearchVariants(ctx *gin.Context) {
	var req request.SearchRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}

	result, err := h.variantService.SearchVariants(req.Keyword, req.Offset, req.Limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *VariantHandler) UpdateVariant(ctx *gin.Context) {
	var idReq request.GetDataByUUIDRequest
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}

	var variantReq updateVariantRequest
	if err := ctx.ShouldBind(&variantReq); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}

	uuid, err := uuid.Parse(idReq.UUID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}

	arg := entity.Variant{
		UUID:        uuid,
		VariantName: variantReq.VariantName,
		Quantity:    variantReq.Quantity,
	}

	result, err := h.variantService.UpdateVariant(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *VariantHandler) DeleteVariant(ctx *gin.Context) {
	var req request.GetDataByUUIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}

	uuid, err := uuid.Parse(req.UUID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}

	err = h.variantService.DeleteVariant(uuid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
	}

	ctx.JSON(http.StatusOK, nil)
}
