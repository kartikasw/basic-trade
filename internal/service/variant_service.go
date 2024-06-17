package service

import (
	"basic-trade/internal/entity"
	"basic-trade/internal/repository"
	"context"

	"github.com/google/uuid"

	sqlc "basic-trade/internal/repository/sqlc"
)

type IVariantService struct {
	variantRepo repository.VariantRepository
}

type VariantService interface {
	CreateVariant(ctx context.Context, variant entity.Variant, uuidPrd uuid.UUID) (entity.Variant, error)
	GetVariant(ctx context.Context, uuid uuid.UUID) (entity.Variant, error)
	GetAllVariants(ctx context.Context, offset int32, limit int32) ([]entity.Variant, error)
	SearchVariants(ctx context.Context, key string, offset int32, limit int32) ([]entity.Variant, error)
	UpdateVariant(ctx context.Context, variant entity.Variant) (entity.Variant, error)
	DeleteVariant(ctx context.Context, uuid uuid.UUID) error
}

func NewVariantService(variantRepo repository.VariantRepository) VariantService {
	return &IVariantService{variantRepo: variantRepo}
}

func (s *IVariantService) CreateVariant(
	ctx context.Context,
	variant entity.Variant,
	uuidPrd uuid.UUID,
) (entity.Variant, error) {
	arg := sqlc.CreateVariantParams{
		VariantName: variant.VariantName,
		Quantity:    variant.Quantity,
	}

	result, err := s.variantRepo.CreateVariant(ctx, arg, uuidPrd)
	if err != nil {
		return entity.Variant{}, err
	}

	return entity.CreateVariantToViewModel(result), err
}

func (s *IVariantService) GetVariant(ctx context.Context, uuid uuid.UUID) (entity.Variant, error) {
	result, err := s.variantRepo.GetVariant(ctx, uuid)
	if err != nil {
		return entity.Variant{}, err
	}

	return entity.GetVariantToViewModel(result), err
}

func (s *IVariantService) GetAllVariants(ctx context.Context, offset int32, limit int32) ([]entity.Variant, error) {
	arg := sqlc.ListVariantsParams{
		Limit:  limit,
		Offset: offset,
	}

	result, err := s.variantRepo.GetAllVariants(ctx, arg)
	if err != nil {
		return []entity.Variant{}, err
	}

	return entity.ListVariantToViewModel(result), err
}

func (s *IVariantService) SearchVariants(ctx context.Context, key string, offset int32, limit int32) ([]entity.Variant, error) {
	arg := sqlc.ListVariantsParams{
		Keyword: key,
		Limit:   limit,
		Offset:  offset,
	}

	result, err := s.variantRepo.GetAllVariants(ctx, arg)
	if err != nil {
		return []entity.Variant{}, err
	}

	return entity.ListVariantToViewModel(result), err
}

func (s *IVariantService) UpdateVariant(ctx context.Context, product entity.Variant) (entity.Variant, error) {
	arg := sqlc.UpdateAVariantParams{
		Uuid:        product.UUID,
		VariantName: product.VariantName,
		Quantity:    product.Quantity,
	}

	result, err := s.variantRepo.UpdateVariant(ctx, arg)
	if err != nil {
		return entity.Variant{}, err
	}

	return entity.UpdateVariantToViewModel(result), err
}

func (s *IVariantService) DeleteVariant(ctx context.Context, uuid uuid.UUID) error {
	err := s.variantRepo.DeleteVariant(ctx, uuid)

	return err
}
