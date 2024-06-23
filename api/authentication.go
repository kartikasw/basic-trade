package api

import (
	"basic-trade/api/response"
	"basic-trade/pkg/token"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authentication(jwt token.JWT) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(token.AuthorizationHeader)

		if authorizationHeader == "" {
			err := errors.New("Authorization header is not found")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("Authorization header is not in Bearer scheme")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse(err))
			return
		}

		authScheme := strings.ToLower(fields[0])
		if authScheme != token.BearerScheme {
			err := fmt.Errorf("Unsupported authorization scheme %s", authScheme)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse(err))
			return
		}

		accessToken := fields[1]
		claim, err := jwt.VerifyToken(accessToken, token.AccessTokenExpectation())
		if err != nil {
			err = fmt.Errorf("Couldn't verify token: %w", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse(err))
			return
		}

		ctx.Set(token.JWTClaim, claim)
		ctx.Next()
	}
}
