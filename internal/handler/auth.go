package handler

import (
	apiHelper "basic-trade/api/helper"
	"basic-trade/api/request"
	response "basic-trade/api/response"
	"basic-trade/common"
	"basic-trade/internal/entity"
	"basic-trade/internal/service"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	apiHelper.ResponseHandler(ctx, func(c context.Context, resChan chan apiHelper.ResponseData) {
		var req request.RegisterRequest
		if err := ctx.ShouldBind(&req); err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusBadRequest,
				Error:      common.ErrorValidation(err),
			}
		}

		arg := entity.Admin{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
		}

		result, err := h.authService.Register(c, arg)

		if err != nil {
			err := err
			if common.ErrorCode(err) == common.ErrUniqueViolation {
				err = errors.New("E-mail is already registered.")
			}
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusInternalServerError,
				Error:      err,
			}
		}

		resChan <- apiHelper.ResponseData{
			StatusCode: http.StatusCreated,
			Message:    "Registration completed successfully.",
			Data:       response.NewAdminResponse(result),
		}
	})
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	apiHelper.ResponseHandler(ctx, func(c context.Context, resChan chan apiHelper.ResponseData) {
		var req request.LoginRequest
		if err := ctx.ShouldBind(&req); err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusBadRequest,
				Error:      common.ErrorValidation(err),
			}
		}

		arg := entity.Admin{
			Email:    req.Email,
			Password: req.Password,
		}

		result, token, err := h.authService.Login(arg)
		if err != nil {
			var statusCode = http.StatusInternalServerError
			if errors.Is(err, common.ErrRecordNotFound) || common.ErrorCode(err) == fmt.Sprint(common.ErrCredentiials) {
				err = errors.New("E-mail or Password is incorrect")
				statusCode = http.StatusUnauthorized
			}
			resChan <- apiHelper.ResponseData{
				StatusCode: statusCode,
				Error:      err,
			}
		}

		resChan <- apiHelper.ResponseData{
			StatusCode: http.StatusOK,
			Message:    "Login successful.",
			Data:       response.NewLoginResponse(token, result),
		}
	})
}
