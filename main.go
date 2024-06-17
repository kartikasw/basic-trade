package main

import (
	"basic-trade/api"
	"basic-trade/internal/handler"
	"basic-trade/internal/repository"
	"basic-trade/internal/service"
	config "basic-trade/pkg/config"
	database "basic-trade/pkg/db"
	"basic-trade/pkg/token"
	"fmt"
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
)

func main() {
	cfg := config.LoadConfig("app.yaml")

	connPool, err := database.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("Init DB error: %v", err)
	}

	jwtImpl, err := token.NewJWT(cfg.Token)
	if err != nil {
		log.Fatalf("Couldn't create token maker: %v", err)
	}

	cld, err := cloudinary.NewFromParams(cfg.Cloudinary.Name, cfg.Cloudinary.ApiKey, cfg.Cloudinary.ApiSecret)
	cld.Config.URL.Secure = true
	if err != nil {
		log.Fatalf("Couldn't create cloudinary environment: %v", err)
	}
	fileRepo := repository.NewFileRepository(cld)

	adminRepo := repository.NewAdminRepository(connPool)
	authService := service.NewAuthService(adminRepo, jwtImpl)
	authHandler := handler.NewAuthHandler(authService)

	productRepo := repository.NewProductRepository(connPool)
	productService := service.NewProductService(productRepo, fileRepo)
	productHandler := handler.NewProductHandler(productService)

	variantRepo := repository.NewVariantRepository(connPool)
	variantService := service.NewVariantService(variantRepo)
	variantHandler := handler.NewVariantHandler(variantService)

	server := api.NewServer(cfg.App, jwtImpl, authHandler, productHandler, variantHandler, adminRepo)
	if err != nil {
		log.Fatal("Couldn't create server: ", err)
	}

	err = server.Start(fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port))
	if err != nil {
		log.Fatal("Couldn't start server: ", err)
	}

}
