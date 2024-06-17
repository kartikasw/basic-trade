package service

import (
	"basic-trade/internal/entity"
	"basic-trade/internal/repository"
	"context"
	"mime/multipart"

	"github.com/google/uuid"

	sqlc "basic-trade/internal/repository/sqlc"
)

type IProductService struct {
	productRepo repository.ProductRepository
	fileRepo    repository.FileRepository
}

type ProductService interface {
	CreateProduct(ctx context.Context, product entity.Product, admUUID uuid.UUID, image *multipart.FileHeader) (entity.Product, error)
	GetProduct(ctx context.Context, uuid uuid.UUID) (entity.ProductView, error)
	GetAllProducts(ctx context.Context, offset int32, limit int32) ([]entity.ProductViewList, error)
	SearchProducts(ctx context.Context, key string, offset int32, limit int32) ([]entity.ProductViewList, error)
	UpdateProduct(ctx context.Context, product entity.Product, admUUID uuid.UUID, image *multipart.FileHeader) (entity.Product, error)
	DeleteProduct(ctx context.Context, prdUUID uuid.UUID, admUUID uuid.UUID) error
}

func NewProductService(productRepo repository.ProductRepository, fileRepo repository.FileRepository) ProductService {
	return &IProductService{productRepo: productRepo, fileRepo: fileRepo}
}

func (s *IProductService) CreateProduct(
	ctx context.Context,
	product entity.Product,
	admUUID uuid.UUID,
	image *multipart.FileHeader,
) (entity.Product, error) {
	arg := sqlc.CreateProductParams{
		Name: product.Name,
	}

	result, err := s.productRepo.CreateProduct(ctx, arg, admUUID, func() (string, error) {
		imageURL, err := s.fileRepo.UploadImage(ctx, admUUID.String(), image)
		return imageURL, err
	})

	if err != nil {
		return entity.Product{}, err
	}

	return entity.CreateProductToViewModel(result), err
}

func (s *IProductService) GetProduct(ctx context.Context, uuid uuid.UUID) (entity.ProductView, error) {
	result, err := s.productRepo.GetProduct(ctx, uuid)
	if err != nil {
		return entity.ProductView{}, err
	}

	return entity.GetProductRowToViewModel(result), err
}

func (s *IProductService) GetAllProducts(ctx context.Context, offset int32, limit int32) ([]entity.ProductViewList, error) {
	arg := sqlc.ListProductsParams{
		Limit:  limit,
		Offset: offset,
	}
	result, err := s.productRepo.GetAllProducts(ctx, arg)
	if err != nil {
		return []entity.ProductViewList{}, err
	}

	return entity.ListProductsRowToViewModel(result), err
}

func (s *IProductService) SearchProducts(ctx context.Context, key string, offset int32, limit int32) ([]entity.ProductViewList, error) {
	arg := sqlc.ListProductsParams{
		Keyword: key,
		Limit:   limit,
		Offset:  offset,
	}
	result, err := s.productRepo.GetAllProducts(ctx, arg)
	if err != nil {
		return []entity.ProductViewList{}, err
	}

	return entity.ListProductsRowToViewModel(result), err
}

func (s *IProductService) UpdateProduct(
	ctx context.Context,
	product entity.Product,
	admUUID uuid.UUID,
	image *multipart.FileHeader,
) (entity.Product, error) {

	arg := sqlc.UpdateAProductParams{
		Uuid: product.UUID,
		Name: product.Name,
	}

	var uploadFunc func() (string, error)
	if image != nil {
		uploadFunc = func() (string, error) {
			if image != nil {
				imageURL, err := s.fileRepo.UploadImage(ctx, admUUID.String(), image)
				return imageURL, err
			}

			return "", nil
		}
	}

	result, err := s.productRepo.UpdateProduct(ctx, arg, uploadFunc)
	if err != nil {
		return entity.Product{}, err
	}

	return entity.UpdateProductToViewModel(result), err
}

func (s *IProductService) DeleteProduct(ctx context.Context, prdUUID uuid.UUID, admUUID uuid.UUID) error {
	err := s.productRepo.DeleteProduct(ctx, prdUUID, func() error {
		return s.fileRepo.DeleteImage(ctx, admUUID.String())
	})

	return err
}
