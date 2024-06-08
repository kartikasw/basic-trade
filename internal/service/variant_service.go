package service

import (
	"basic-trade/internal/entity"
	"basic-trade/internal/repository"

	"github.com/google/uuid"

	sqlc "basic-trade/internal/repository/sqlc"
)

type VariantService struct {
	variantRepo repository.IVariantRepository
}

type IVariantService interface {
	CreateVariant(variant entity.Variant, uuidPrd uuid.UUID) (entity.Variant, error)
	GetVariant(uuid uuid.UUID) (entity.Variant, error)
	GetAllVariants(offset int32, limit int32) ([]entity.Variant, error)
	SearchVariants(key string, offset int32, limit int32) ([]entity.Variant, error)
	UpdateVariant(variant entity.Variant) (entity.Variant, error)
	DeleteVariant(uuid uuid.UUID) error
}

func NewVariantService(variantRepo repository.IVariantRepository) *VariantService {
	return &VariantService{variantRepo: variantRepo}
}

func (s *VariantService) CreateVariant(variant entity.Variant, uuidPrd uuid.UUID) (entity.Variant, error) {
	arg := sqlc.CreateVariantParams{
		VariantName: variant.VariantName,
		Quantity:    variant.Quantity,
	}

	result, err := s.variantRepo.CreateVariant(arg, uuidPrd)
	if err != nil {
		return entity.Variant{}, err
	}

	return entity.CreateVariantToViewModel(&result), err
}

func (s *VariantService) GetVariant(uuid uuid.UUID) (entity.Variant, error) {
	result, err := s.variantRepo.GetVariant(uuid)
	if err != nil {
		return entity.Variant{}, err
	}

	return entity.GetVariantToViewModel(&result), err
}

func (s *VariantService) GetAllVariants(offset int32, limit int32) ([]entity.Variant, error) {
	arg := sqlc.ListVariantsParams{
		Limit:  limit,
		Offset: offset,
	}

	result, err := s.variantRepo.GetAllVariants(arg)
	if err != nil {
		return []entity.Variant{}, err
	}

	return entity.ListVariantToViewModel(result), err
}

func (s *VariantService) SearchVariants(key string, offset int32, limit int32) ([]entity.Variant, error) {
	arg := sqlc.ListVariantsParams{
		Keyword: key,
		Limit:   limit,
		Offset:  offset,
	}

	result, err := s.variantRepo.GetAllVariants(arg)
	if err != nil {
		return []entity.Variant{}, err
	}

	return entity.ListVariantToViewModel(result), err
}

func (s *VariantService) UpdateVariant(product entity.Variant) (entity.Variant, error) {
	arg := sqlc.UpdateAVariantParams{
		Uuid:        product.UUID,
		VariantName: product.VariantName,
		Quantity:    product.Quantity,
	}

	result, err := s.variantRepo.UpdateVariant(arg)
	if err != nil {
		return entity.Variant{}, err
	}

	return entity.UpdateVariantToViewModel(&result), err
}

func (s *VariantService) DeleteVariant(uuid uuid.UUID) error {
	err := s.variantRepo.DeleteVariant(uuid)

	return err
}
