package middleware

import (
	"basic-trade/api/response"
	"basic-trade/internal/repository"
	"basic-trade/pkg/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthorizationMiddleware struct {
	adminRepo repository.IAdminRepository
}

func NewAuthorizationMiddleware(adminRepo repository.IAdminRepository) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{adminRepo: adminRepo}
}

func (a *AuthorizationMiddleware) ProductAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
		param := ctx.Param("uuid")

		productUUID, err := uuid.Parse(param)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse(err))
			return
		}

		result := a.adminRepo.CheckProductFromAdmin(authPayload.UUID, productUUID)

		if !result {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, nil)
			return
		}

		ctx.Next()
	}
}

func (a *AuthorizationMiddleware) VariantAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
		param := ctx.Param("uuid")

		variantUUID, err := uuid.Parse(param)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse(err))
			return
		}

		result := a.adminRepo.CheckVariantFromAdmin(authPayload.UUID, variantUUID)
		if !result {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, nil)
			return
		}

		ctx.Next()
	}
}
