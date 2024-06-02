package api

import (
	"basic-trade/api/middleware"
	"basic-trade/internal/handler"
	"basic-trade/internal/repository"
	"basic-trade/pkg/token"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router         *gin.Engine
	tokenMaker     token.Maker
	authHandler    *handler.AuthHandler
	productHandler *handler.ProductHandler
	variantHandler *handler.VariantHandler
	authorization  middleware.AuthorizationMiddleware
}

func NewServer(
	tokenMaker token.Maker,
	authHandler *handler.AuthHandler,
	productHandler *handler.ProductHandler,
	variantHandler *handler.VariantHandler,
	adminRepo repository.IAdminRepository,
) *Server {
	server := &Server{
		tokenMaker:     tokenMaker,
		authHandler:    authHandler,
		productHandler: productHandler,
		variantHandler: variantHandler,
		authorization:  *middleware.NewAuthorizationMiddleware(adminRepo),
	}
	server.setupRouter()
	return server
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.MaxMultipartMemory = 5 << 20

	router.POST("/auth/register", server.authHandler.Register)
	router.POST("/auth/login", server.authHandler.Login)

	router.GET("/products", server.productHandler.GetAllProducts)
	router.GET("/products/search", server.productHandler.SearchProducts)
	router.GET("/products/:uuid", server.productHandler.GetProduct)

	router.GET("/variants", server.variantHandler.GetAllVariants)
	router.GET("/variants/search", server.variantHandler.SearchVariants)
	router.GET("/variants/:uuid", server.variantHandler.GetVariant)

	authRoutes := router.Group("/").Use(middleware.Authentication(server.tokenMaker))
	authRoutes.POST("/products", server.productHandler.CreateProduct)
	authRoutes.PUT("/products/:uuid", server.authorization.ProductAuthorization(), server.productHandler.UpdateProduct)
	authRoutes.DELETE("/products/:uuid", server.authorization.ProductAuthorization(), server.productHandler.DeleteProduct)

	authRoutes.POST("/variants", server.variantHandler.CreateVariant)
	authRoutes.PUT("/variants/:uuid", server.authorization.VariantAuthorization(), server.variantHandler.UpdateVariant)
	authRoutes.DELETE("/variants/:uuid", server.authorization.VariantAuthorization(), server.variantHandler.DeleteVariant)

	server.router = router
}

func (server *Server) Start(port string) error {
	return server.router.Run(port)
}
