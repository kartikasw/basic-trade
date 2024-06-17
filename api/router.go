package api

import (
	"basic-trade/api/middleware"
	"basic-trade/common"
	"basic-trade/internal/handler"
	"basic-trade/internal/repository"
	"basic-trade/pkg/config"
	"basic-trade/pkg/token"
	"basic-trade/pkg/validation"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	router         *gin.Engine
	jwtImpl        token.JWT
	authHandler    *handler.AuthHandler
	productHandler *handler.ProductHandler
	variantHandler *handler.VariantHandler
	authorization  middleware.AuthorizationMiddleware
}

func NewServer(
	cfg config.App,
	jwtImpl token.JWT,
	authHandler *handler.AuthHandler,
	productHandler *handler.ProductHandler,
	variantHandler *handler.VariantHandler,
	adminRepo repository.AdminRepository,
) *Server {
	server := &Server{
		jwtImpl:        jwtImpl,
		authHandler:    authHandler,
		productHandler: productHandler,
		variantHandler: variantHandler,
		authorization:  *middleware.NewAuthorizationMiddleware(adminRepo),
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validImage", validation.ValidImage)
		v.RegisterValidation("validUUID", validation.ValidUUID)

	}

	server.setupRouter(cfg)
	return server
}

func (server *Server) setupRouter(cfg config.App) {
	router := gin.Default()

	router.MaxMultipartMemory = common.MaxFileSize

	formRoutes := router.Group("/").Use(middleware.ContentTypeValidation(), middleware.Timeout(cfg.Timeout))
	{
		formRoutes.POST("/auth/register", server.authHandler.Register)
		formRoutes.POST("/auth/login", server.authHandler.Login)
	}

	authFormRoutes := router.Group("/").Use(
		middleware.ContentTypeValidation(),
		middleware.Authentication(server.jwtImpl),
		middleware.Timeout(cfg.Timeout),
	)
	{
		authFormRoutes.POST("/products", server.productHandler.CreateProduct)
		authFormRoutes.PUT("/products/:uuid", server.authorization.ProductAuthorization(), server.productHandler.UpdateProduct)
		authFormRoutes.POST("/variants", server.variantHandler.CreateVariant)
		authFormRoutes.PUT("/variants/:uuid", server.authorization.VariantAuthorization(), server.variantHandler.UpdateVariant)
	}

	timeout := router.Group("/").Use(middleware.Timeout(cfg.Timeout))
	{
		timeout.GET("/products", server.productHandler.GetAllProducts)
		timeout.GET("/products/search", server.productHandler.SearchProducts)
		timeout.GET("/products/:uuid", server.productHandler.GetProduct)
		timeout.GET("/variants", server.variantHandler.GetAllVariants)
		timeout.GET("/variants/search", server.variantHandler.SearchVariants)
		timeout.GET("/variants/:uuid", server.variantHandler.GetVariant)
	}

	authRoutes := router.Group("/").Use(middleware.Authentication(server.jwtImpl), middleware.Timeout(cfg.Timeout))
	{
		authRoutes.DELETE("/products/:uuid", server.authorization.ProductAuthorization(), server.productHandler.DeleteProduct)
		authRoutes.DELETE("/variants/:uuid", server.authorization.VariantAuthorization(), server.variantHandler.DeleteVariant)
	}

	server.router = router
}

func (server *Server) Start(port string) error {
	return server.router.Run(port)
}
