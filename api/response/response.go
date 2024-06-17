package response

import (
	"github.com/gin-gonic/gin"
)

func ErrorResponse(err error) gin.H {
	return gin.H{
		"status":  "error",
		"message": err.Error(),
	}
}

func SuccessResponse(message string, data any) gin.H {
	result := gin.H{
		"result":  "success",
		"message": message,
	}

	if data != nil {
		result["data"] = data
	}

	return result
}
