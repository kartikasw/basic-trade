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
	fmt.Println("cfg: ", cfg)

	connPool, err := database.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("Init DB error: %v", err)
	}

	tokenMaker, err := token.NewJWTMaker(cfg.Token.SecretKey)
	if err != nil {
		log.Fatalf("Couldn't create token maker: %v", err)
	}

	adminRepo := repository.NewAdminRepository(connPool)
	authService := service.NewAuthService(adminRepo, tokenMaker, cfg.Token)
	authHandler := handler.NewAuthHandler(authService)

	cld, err := cloudinary.NewFromParams(cfg.Cloudinary.Name, cfg.Cloudinary.ApiKey, cfg.Cloudinary.ApiSecret)
	if err != nil {
		log.Fatalf("Couldn't create cloudinary environment: %v", err)
	}

	productRepo := repository.NewProductRepository(connPool)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService, cld)

	variantRepo := repository.NewVariantRepository(connPool)
	variantService := service.NewVariantService(variantRepo)
	variantHandler := handler.NewVariantHandler(variantService)

	server := api.NewServer(tokenMaker, authHandler, productHandler, variantHandler, adminRepo)
	if err != nil {
		log.Fatal("Couldn't create server: ", err)
	}

	err = server.Start(fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port))
	if err != nil {
		log.Fatal("Couldn't start server: ", err)
	}

}
