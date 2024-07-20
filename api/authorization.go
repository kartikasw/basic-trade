package api

import (
	apiHelper "basic-trade/api/helper"
	"basic-trade/api/request"
	"basic-trade/internal/repository"
	"basic-trade/pkg/token"
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	contentType           = "Content-Type"
	multipartFormDataType = "multipart/form-data"
)

type AuthorizationMiddleware struct {
	adminRepo repository.AdminRepository
}

func NewAuthorizationMiddleware(adminRepo repository.AdminRepository) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{adminRepo: adminRepo}
}

// data is to determine wether to get param from the URL param or the body
func (a *AuthorizationMiddleware) ProductAuthorization(data ...bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiHelper.ResponseHandler(ctx, func(c context.Context, resChan chan apiHelper.ResponseData) {
			payload := ctx.MustGet(token.JWTClaim).(*token.Claim)

			var param string
			if len(data) != 0 && !data[0] {
				var req request.VariantProductIDRequest
				if err := ctx.ShouldBind(&req); err != nil {
					resChan <- apiHelper.ResponseData{
						StatusCode: http.StatusBadRequest,
						Error:      err,
					}
				}
				param = req.ProductID
			} else {
				param = ctx.Param("uuid")
			}

			productUUID, err := uuid.Parse(param)
			if err != nil {
				resChan <- apiHelper.ResponseData{
					StatusCode: http.StatusBadRequest,
					Error:      err,
				}
			}

			result := a.adminRepo.CheckProductFromAdmin(c, payload.UserID, productUUID)

			if !result {
				resChan <- apiHelper.ResponseData{
					StatusCode: http.StatusUnauthorized,
					Error:      errors.New("Product doesn't belong to the authorized admin"),
				}
			}

			ctx.Next()
		})
	}
}

func (a *AuthorizationMiddleware) VariantAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiHelper.ResponseHandler(ctx, func(c context.Context, resChan chan apiHelper.ResponseData) {
			payload := ctx.MustGet(token.JWTClaim).(*token.Claim)
			param := ctx.Param("uuid")

			variantUUID, err := uuid.Parse(param)
			if err != nil {
				resChan <- apiHelper.ResponseData{
					StatusCode: http.StatusBadRequest,
					Error:      err,
				}
			}

			result := a.adminRepo.CheckVariantFromAdmin(c, payload.UserID, variantUUID)
			if !result {
				resChan <- apiHelper.ResponseData{
					StatusCode: http.StatusUnauthorized,
					Error:      errors.New("Variant doesn't belong to the authorized admin"),
				}
			}

			ctx.Next()
		})
	}
}
