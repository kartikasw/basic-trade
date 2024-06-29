package handler

import (
	apiHelper "basic-trade/api/helper"
	"basic-trade/api/request"
	"basic-trade/common"
	"basic-trade/internal/entity"
	"basic-trade/internal/service"
	"basic-trade/pkg/token"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) CreateProduct(ctx *gin.Context) {
	apiHelper.ResponseHandler(ctx, func(c context.Context, resChan chan apiHelper.ResponseData) {
		var req request.CreateProductRequest
		if err := ctx.ShouldBind(&req); err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusBadRequest,
				Error:      common.ErrorValidation(err),
			}
		}

		jwtPayload := ctx.MustGet(token.JWTClaim).(*token.Claim)

		arg := entity.Product{Name: req.Name}

		result, err := h.productService.CreateProduct(c, arg, jwtPayload.UserID, req.Image)
		if err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusInternalServerError,
				Error:      err,
			}
		}

		resChan <- apiHelper.ResponseData{
			StatusCode: http.StatusCreated,
			Message:    "Product created successfully.",
			Data:       result,
		}
	})
}

func (h *ProductHandler) GetProduct(ctx *gin.Context) {
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

		result, err := h.productService.GetProduct(c, uuid)
		if err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusInternalServerError,
				Error:      err,
			}
		}

		resChan <- apiHelper.ResponseData{
			StatusCode: http.StatusOK,
			Message:    "Product retrieved successfully.",
			Data:       result,
		}
	})
}

func (h *ProductHandler) GetAllProducts(ctx *gin.Context) {
	apiHelper.ResponseHandler(ctx, func(c context.Context, resChan chan apiHelper.ResponseData) {
		var req request.PaginationRequest
		if err := ctx.ShouldBindQuery(&req); err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusBadRequest,
				Error:      common.ErrorValidation(err),
			}
		}

		result, err := h.productService.GetAllProducts(c, req.Keyword, req.Offset, req.Limit)
		if err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusInternalServerError,
				Error:      err,
			}
		}

		resChan <- apiHelper.ResponseData{
			StatusCode: http.StatusOK,
			Message:    "Products retrieved successfully.",
			Data:       result,
		}
	})
}

func (h *ProductHandler) UpdateProduct(ctx *gin.Context) {
	apiHelper.ResponseHandler(ctx, func(c context.Context, resChan chan apiHelper.ResponseData) {
		var idReq request.GetDataByUUIDRequest

		if err := ctx.ShouldBindUri(&idReq); err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusBadRequest,
				Error:      common.ErrorValidation(err),
			}
		}

		var productReq request.UpdateProductRequest
		if err := ctx.ShouldBind(&productReq); err != nil {
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

		arg := entity.Product{
			UUID: uuid,
			Name: productReq.Name,
		}

		result, err := h.productService.UpdateProduct(c, arg, productReq.Image)
		if err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusInternalServerError,
				Error:      err,
			}
		}

		resChan <- apiHelper.ResponseData{
			StatusCode: http.StatusOK,
			Message:    "Product updated successfully.",
			Data:       result,
		}
	})
}

func (h *ProductHandler) DeleteProduct(ctx *gin.Context) {
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

		err = h.productService.DeleteProduct(c, uuid)
		if err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusInternalServerError,
				Error:      err,
			}
		}

		resChan <- apiHelper.ResponseData{
			StatusCode: http.StatusOK,
			Message:    "Product deleted successfully.",
			Data:       nil,
		}
	})
}
