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
	GetAllVariants(ctx context.Context, key string, offset int32, limit int32) (entity.VariantPaginationView, error)
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

func (s *IVariantService) GetAllVariants(ctx context.Context, key string, offset int32, limit int32) (entity.VariantPaginationView, error) {
	arg := sqlc.ListVariantsParams{
		Keyword:   key,
		LimitVal:  limit,
		OffsetVal: offset,
	}

	result, total, err := s.variantRepo.GetAllVariants(ctx, arg)
	if err != nil {
		return entity.VariantPaginationView{}, err
	}

	return entity.ListVariantToViewModel(result, limit, offset, total), err
}

func (s *IVariantService) UpdateVariant(ctx context.Context, variant entity.Variant) (entity.Variant, error) {
	arg := sqlc.UpdateAVariantParams{
		Uuid:        variant.UUID,
		VariantName: variant.VariantName,
		Quantity:    variant.Quantity,
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
