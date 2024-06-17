package middleware

import (
	"basic-trade/api/response"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ContentTypeValidation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ct := ctx.Request.Header.Get(contentType)

		if !strings.HasPrefix(ct, multipartFormDataType) {
			err := errors.New("Content-Type is not supported")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(err))
			return
		}

		ctx.Next()
	}
}
