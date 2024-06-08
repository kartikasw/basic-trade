package handler

import (
	response "basic-trade/api/response"
	"basic-trade/internal/entity"
	"basic-trade/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.IAuthService
}

func NewAuthHandler(authService service.IAuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type registerRequest struct {
	Name     string `form:"name" binding:"required,max=100"`
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required,min=8"`
}

type adminResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func newAdminResponse(admin entity.Admin) adminResponse {
	return adminResponse{
		Name:  admin.Name,
		Email: admin.Email,
	}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var req registerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}

	arg := entity.Admin{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := h.authService.Register(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
	}

	ctx.JSON(http.StatusOK, newAdminResponse(result))
}

type loginRequest struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type loginResponse struct {
	Token string        `json:"token"`
	Admin adminResponse `json:"admin"`
}

func newLoginResponse(token string, admin entity.Admin) loginResponse {
	return loginResponse{
		Token: token,
		Admin: newAdminResponse(admin),
	}
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}

	arg := entity.Admin{
		Email:    req.Email,
		Password: req.Password,
	}

	result, token, err := h.authService.Login(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
	}

	ctx.JSON(http.StatusOK, newLoginResponse(token, result))
}
