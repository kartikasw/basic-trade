package api

import (
	apiHelper "basic-trade/api/helper"
	"basic-trade/api/response"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// https://github.com/go-chi/chi/blob/master/middleware/timeout.go
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		resChan := make(chan apiHelper.ResponseData)

		go func() {
			c.Set(apiHelper.ResChan, resChan)
			c.Next()
			close(resChan)
		}()

		select {
		case res := <-resChan:
			if res.Error != nil && res.Error != context.DeadlineExceeded {
				c.AbortWithStatusJSON(res.StatusCode, response.ErrorResponse(res.Error))
				return
			}

			c.JSON(res.StatusCode, response.SuccessResponse(res.Message, res.Data))
			return
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				err := errors.New("Service is unavailable or timed out")
				c.AbortWithStatusJSON(http.StatusGatewayTimeout, response.ErrorResponse(err))
				return
			}
		}
	}
}
