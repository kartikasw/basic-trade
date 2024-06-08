package handler

import (
	"basic-trade/api/middleware"
	"basic-trade/api/request"
	"basic-trade/api/response"
	"basic-trade/internal/entity"
	"basic-trade/internal/service"
	"basic-trade/pkg/token"
	"mime/multipart"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct {
	productService service.IProductService
	fileService    service.IFileService
}

func NewProductHandler(productService service.IProductService, cld *cloudinary.Cloudinary) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		fileService:    service.NewFileService(cld),
	}
}

type createProductRequest struct {
	Name  string                `form:"name" binding:"required,max=100"`
	Image *multipart.FileHeader `form:"image" binding:"required,validImage"`
}

type updateProductRequest struct {
	Name  string                `form:"name" binding:"max=100"`
	Image *multipart.FileHeader `form:"image" binding:"file,image"`
}

func (h *ProductHandler) CreateProduct(ctx *gin.Context) {
	var req createProductRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)

	imageURL, err := h.fileService.UploadImage(authPayload.UUID.String(), req.Image)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
	}

	arg := entity.Product{
		Name:     req.Name,
		ImageURL: imageURL,
	}

	result, err := h.productService.CreateProduct(arg, authPayload.UUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *ProductHandler) GetProduct(ctx *gin.Context) {
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

	result, err := h.productService.GetProduct(uuid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *ProductHandler) GetAllProducts(ctx *gin.Context) {
	var req request.PaginationRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}

	result, err := h.productService.GetAllProducts(req.Offset, req.Limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *ProductHandler) SearchProducts(ctx *gin.Context) {
	var req request.SearchRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}

	result, err := h.productService.SearchProducts(req.Keyword, req.Offset, req.Limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *ProductHandler) UpdateProduct(ctx *gin.Context) {
	var idReq request.GetDataByUUIDRequest

	if err := ctx.ShouldBindUri(&idReq); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}

	var productReq updateProductRequest
	if err := ctx.ShouldBindJSON(&productReq); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)

	var imageURL string
	if productReq.Image != nil {
		imageURL, _ = h.fileService.UploadImage(authPayload.UUID.String(), productReq.Image)
	}

	uuid, err := uuid.Parse(idReq.UUID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}

	arg := entity.Product{
		UUID:     uuid,
		Name:     productReq.Name,
		ImageURL: imageURL,
	}

	result, err := h.productService.UpdateProduct(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *ProductHandler) DeleteProduct(ctx *gin.Context) {
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

	err = h.productService.DeleteProduct(uuid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
	}

	ctx.JSON(http.StatusOK, nil)
}
