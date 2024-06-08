package service

import (
	"basic-trade/internal/entity"
	"basic-trade/internal/repository"
	"fmt"

	"github.com/google/uuid"

	sqlc "basic-trade/internal/repository/sqlc"
)

type ProductService struct {
	productRepo repository.IProductRepository
}

type IProductService interface {
	CreateProduct(product entity.Product, uuidAdm uuid.UUID) (entity.Product, error)
	GetProduct(uuid uuid.UUID) (entity.Product, error)
	GetAllProducts(offset int32, limit int32) ([]entity.Product, error)
	SearchProducts(key string, offset int32, limit int32) ([]entity.Product, error)
	UpdateProduct(product entity.Product) (entity.Product, error)
	DeleteProduct(uuid uuid.UUID) error
}

func NewProductService(productRepo repository.IProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) CreateProduct(product entity.Product, uuidAdm uuid.UUID) (entity.Product, error) {
	arg := sqlc.CreateProductParams{
		Name:     product.Name,
		ImageUrl: product.ImageURL,
	}

	result, err := s.productRepo.CreateProduct(arg, uuidAdm)
	if err != nil {
		return entity.Product{}, err
	}

	return entity.CreateProductToViewModel(&result), err
}

func (s *ProductService) GetProduct(uuid uuid.UUID) (entity.Product, error) {
	result, err := s.productRepo.GetProduct(uuid)
	if err != nil {
		return entity.Product{}, err
	}

	fmt.Println("GetProduct, result=", result)

	return entity.ProductViewToViewModel(&result), err
}

func (s *ProductService) GetAllProducts(offset int32, limit int32) ([]entity.Product, error) {
	arg := sqlc.ListProductsParams{
		Limit:  limit,
		Offset: offset,
	}
	result, err := s.productRepo.GetAllProducts(arg)
	if err != nil {
		return []entity.Product{}, err
	}

	return entity.ListProductViewToViewModel(result), err
}

func (s *ProductService) SearchProducts(key string, offset int32, limit int32) ([]entity.Product, error) {
	arg := sqlc.ListProductsParams{
		Keyword: key,
		Limit:  limit,
		Offset: offset,
	}
	result, err := s.productRepo.GetAllProducts(arg)
	if err != nil {
		return []entity.Product{}, err
	}

	return entity.ListProductViewToViewModel(result), err
}

func (s *ProductService) UpdateProduct(product entity.Product) (entity.Product, error) {
	arg := sqlc.UpdateAProductParams{
		Uuid:     product.UUID,
		Name:     product.Name,
		ImageUrl: product.ImageURL,
	}

	result, err := s.productRepo.UpdateProduct(arg)
	if err != nil {
		return entity.Product{}, err
	}

	return entity.UpdateProductToViewModel(&result), err
}

func (s *ProductService) DeleteProduct(uuid uuid.UUID) error {
	err := s.productRepo.DeleteProduct(uuid)

	return err
}
