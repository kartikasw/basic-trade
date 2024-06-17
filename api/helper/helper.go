package apiHelper

import (
	"context"

	"github.com/gin-gonic/gin"
)

const ResChan = "response_channel"

type ResponseData struct {
	StatusCode int
	Error      error
	Message    string
	Data       any
}

func ResponseHandler(ctx *gin.Context, action func(context.Context, chan ResponseData)) {
	reqCtx := ctx.Request.Context()
	resChan := ctx.MustGet(ResChan)

	action(reqCtx, resChan.(chan ResponseData))
}
