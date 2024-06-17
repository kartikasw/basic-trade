package middleware

import (
	"basic-trade/api/response"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"golang.org/x/time/rate"
)

func RateLimiter() gin.HandlerFunc {
	// Implements a "token bucket" of size b,
	// initially full and refilled at the rate of r tokens per second
	limiter := rate.NewLimiter(1, 1000)
	return func(c *gin.Context) {

		if limiter.Allow() {
			c.Next()
		} else {
			err := errors.New("Request exceed")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, response.ErrorResponse(err))
			return
		}

	}
}
