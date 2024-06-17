package repository

import (
	"basic-trade/common"
	sqlc "basic-trade/internal/repository/sqlc"

	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomVariant(t *testing.T, prdUUID uuid.UUID) sqlc.CreateVariantRow {
	arg := sqlc.CreateVariantParams{
		VariantName: common.RandomName(),
		Quantity:    int32(common.RandomInt(1, 100)),
	}

	variant, err := testVariantRepo.CreateVariant(nil, arg, prdUUID)

	require.NoError(t, err)
	require.Equal(t, arg.VariantName, variant.VariantName)
	require.Equal(t, arg.Quantity, variant.Quantity)

	return variant
}

func TestCreateVariant(t *testing.T) {
	admin := createRandomAdmin(t)
	product := createRandomProduct(t, admin.Uuid, false)

	createRandomVariant(t, product.Uuid)
}

func TestGetVariant(t *testing.T) {
	admin := createRandomAdmin(t)
	product := createRandomProduct(t, admin.Uuid, false)

	variant1 := createRandomVariant(t, product.Uuid)

	variant2, err := testVariantRepo.GetVariant(nil, variant1.Uuid)

	require.NoError(t, err)
	require.Equal(t, variant1.VariantName, variant2.VariantName)
	require.Equal(t, variant1.Quantity, variant2.Quantity)
}

func TestGetAllVariants(t *testing.T) {
	admin := createRandomAdmin(t)
	product := createRandomProduct(t, admin.Uuid, false)

	for i := 1; i <= 10; i++ {
		createRandomVariant(t, product.Uuid)
	}

	arg := sqlc.ListVariantsParams{Limit: 100, Offset: 0}

	variants, err := testVariantRepo.GetAllVariants(nil, arg)

	require.NoError(t, err)
	require.Len(t, variants, 10)
}

func TestGetAllVariantsForSearch(t *testing.T) {}

func UpdateVariant(t *testing.T) {
	admin := createRandomAdmin(t)
	product := createRandomProduct(t, admin.Uuid, false)
	variant1 := createRandomVariant(t, product.Uuid)

	arg := sqlc.UpdateAVariantParams{
		Uuid:           variant1.Uuid,
		SetVariantName: true,
		VariantName:    common.RandomName(),
		Quantity:       int32(common.RandomInt(1, 20)),
	}

	variant2, err := testVariantRepo.UpdateVariant(nil, arg)

	require.NoError(t, err)
	require.Equal(t, variant1.Uuid, variant2.Uuid)
	require.Equal(t, arg.VariantName, variant2.VariantName)
	require.Equal(t, arg.Quantity, variant2.Quantity)
	require.NotEqual(t, variant1.VariantName, variant2.VariantName)
}

func DeleteVariant(t *testing.T) {
	admin := createRandomAdmin(t)
	product := createRandomProduct(t, admin.Uuid, false)
	variant := createRandomVariant(t, product.Uuid)

	err := testVariantRepo.DeleteVariant(nil, variant.Uuid)

	require.NoError(t, err)

	_, err = testVariantRepo.GetVariant(nil, variant.Uuid)

	require.Error(t, err)
}
