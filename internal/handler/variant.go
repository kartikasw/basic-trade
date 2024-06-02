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
	VariantName string    `json:"variant_name" binding:"required,max=100"`
	Quantity    int32     `json:"quantity" binding:"required"`
	ProductID   uuid.UUID `json:"product_id" binding:"required,uuid4"`
}

type updateVariantRequest struct {
	VariantName string `json:"variant_name" binding:"max=100"`
	Quantity    int32  `json:"quantity"`
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

	result, err := h.variantService.CreateVariant(arg, req.ProductID)
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

	result, err := h.variantService.GetVariant(req.UUID)
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

	arg := entity.Variant{
		UUID:        idReq.UUID,
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

	err := h.variantService.DeleteVariant(req.UUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
	}

	ctx.JSON(http.StatusOK, nil)
}
