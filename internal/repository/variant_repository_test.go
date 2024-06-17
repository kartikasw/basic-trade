package repository

import (
	"basic-trade/common"
	sqlc "basic-trade/internal/repository/sqlc"
	"context"

	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomVariant(t *testing.T, ctx context.Context, prdUUID uuid.UUID) sqlc.CreateVariantRow {
	arg := sqlc.CreateVariantParams{
		VariantName: common.RandomName(),
		Quantity:    int32(common.RandomInt(1, 100)),
	}

	variant, err := testVariantRepo.CreateVariant(ctx, arg, prdUUID)

	require.NoError(t, err)
	require.Equal(t, arg.VariantName, variant.VariantName)
	require.Equal(t, arg.Quantity, variant.Quantity)

	return variant
}

func TestCreateVariant(t *testing.T) {
	ctx := context.Background()
	defer tearDown(ctx)

	admin := createRandomAdmin(t, ctx)
	product := createRandomProduct(t, ctx, admin.Uuid, false)

	createRandomVariant(t, ctx, product.Uuid)
}

func TestGetVariant(t *testing.T) {
	ctx := context.Background()
	defer tearDown(ctx)

	admin := createRandomAdmin(t, ctx)
	product := createRandomProduct(t, ctx, admin.Uuid, false)

	variant1 := createRandomVariant(t, ctx, product.Uuid)

	variant2, err := testVariantRepo.GetVariant(ctx, variant1.Uuid)

	require.NoError(t, err)
	require.Equal(t, variant1.VariantName, variant2.VariantName)
	require.Equal(t, variant1.Quantity, variant2.Quantity)
}

func TestGetAllVariants(t *testing.T) {
	ctx := context.Background()
	defer tearDown(ctx)

	admin := createRandomAdmin(t, ctx)
	product := createRandomProduct(t, ctx, admin.Uuid, false)

	for i := 1; i <= 10; i++ {
		createRandomVariant(t, ctx, product.Uuid)
	}

	arg1 := sqlc.ListVariantsParams{Limit: 5, Offset: 0}
	variants1, err := testVariantRepo.GetAllVariants(ctx, arg1)

	require.NoError(t, err)
	require.Len(t, variants1, 5)
	require.Equal(t, variants1[0].RowNumber, int64(1))

	arg2 := sqlc.ListVariantsParams{Limit: 5, Offset: 5}

	variants2, err := testVariantRepo.GetAllVariants(ctx, arg2)

	require.NoError(t, err)
	require.Len(t, variants2, 5)
	require.Equal(t, variants2[0].RowNumber, int64(6))
}

func TestUpdateVariant(t *testing.T) {
	ctx := context.Background()
	defer tearDown(ctx)

	admin := createRandomAdmin(t, ctx)
	product := createRandomProduct(t, ctx, admin.Uuid, false)
	variant1 := createRandomVariant(t, ctx, product.Uuid)

	arg := sqlc.UpdateAVariantParams{
		Uuid:           variant1.Uuid,
		SetVariantName: true,
		VariantName:    common.RandomName(),
		Quantity:       int32(common.RandomInt(1, 20)),
	}

	variant2, err := testVariantRepo.UpdateVariant(ctx, arg)

	require.NoError(t, err)
	require.Equal(t, variant1.Uuid, variant2.Uuid)
	require.Equal(t, arg.VariantName, variant2.VariantName)
	require.Equal(t, arg.Quantity, variant2.Quantity)
	require.NotEqual(t, variant1.VariantName, variant2.VariantName)
}

func TestDeleteVariant(t *testing.T) {
	ctx := context.Background()
	defer tearDown(ctx)

	admin := createRandomAdmin(t, ctx)
	product := createRandomProduct(t, ctx, admin.Uuid, false)
	variant := createRandomVariant(t, ctx, product.Uuid)

	err := testVariantRepo.DeleteVariant(ctx, variant.Uuid)

	require.NoError(t, err)

	_, err = testVariantRepo.GetVariant(ctx, variant.Uuid)

	require.Error(t, err)
}
