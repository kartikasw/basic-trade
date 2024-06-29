package repository

import (
	"basic-trade/common"
	sqlc "basic-trade/internal/repository/sqlc"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomProduct(t *testing.T, ctx context.Context, adminUUID uuid.UUID, withVariants bool) sqlc.CreateProductRow {
	arg := sqlc.CreateProductParams{
		Name: common.RandomName(),
	}

	product, err := testProductRepo.CreateProduct(ctx, arg, adminUUID, func(prdUUID uuid.UUID) (string, error) {
		return common.RandomString(6), nil
	})

	require.NoError(t, err)
	require.Equal(t, arg.Name, product.Name)

	if withVariants {
		n := int(common.RandomInt(1, 10))

		for i := 1; i <= n; i++ {
			createRandomVariant(t, ctx, product.Uuid)
		}
	}

	return product
}

func TestCreateProduct(t *testing.T) {
	ctx := context.Background()
	defer tearDown(ctx)

	admin := createRandomAdmin(t, ctx)

	createRandomProduct(t, ctx, admin.Uuid, false)
}

func TestGetProduct(t *testing.T) {
	ctx := context.Background()
	defer tearDown(ctx)

	admin := createRandomAdmin(t, ctx)
	product1 := createRandomProduct(t, ctx, admin.Uuid, false)

	product2, err := testProductRepo.GetProduct(ctx, product1.Uuid)

	require.NoError(t, err)
	require.Equal(t, product1.Name, product2.Name)
	require.Equal(t, product1.ImageUrl, product2.ImageUrl)
	require.Equal(t, product1.Uuid, product2.Uuid)
	require.Equal(t, product1.Uuid, product2.Uuid)
	require.Len(t, product2.Variants.([]interface{}), 0)
}

func TestGetAllProducts(t *testing.T) {
	ctx := context.Background()
	defer tearDown(ctx)

	admin := createRandomAdmin(t, ctx)

	n := 10
	for i := 1; i <= n; i++ {
		createRandomProduct(t, ctx, admin.Uuid, false)
	}

	arg1 := sqlc.ListProductsParams{Limit: 5, Offset: 0}
	products1, err := testProductRepo.GetAllProducts(ctx, arg1)

	require.NoError(t, err)
	require.Len(t, products1, 5)
	require.Equal(t, products1[0].RowNumber, int64(1))

	arg2 := sqlc.ListProductsParams{Limit: 5, Offset: 5}
	products2, err := testProductRepo.GetAllProducts(ctx, arg2)

	require.NoError(t, err)
	require.Len(t, products2, 5)
	require.Equal(t, products2[0].RowNumber, int64(6))
}

func TestUpdateProduct(t *testing.T) {
	ctx := context.Background()
	defer tearDown(ctx)

	admin := createRandomAdmin(t, ctx)
	product1 := createRandomProduct(t, ctx, admin.Uuid, false)

	arg := sqlc.UpdateAProductParams{
		Uuid:    product1.Uuid,
		SetName: true,
		Name:    common.RandomName(),
	}

	product2, err := testProductRepo.UpdateProduct(ctx, arg, func() (string, error) {
		return common.RandomString(10), nil
	})

	require.NoError(t, err)
	require.Equal(t, product2.Uuid, product1.Uuid)
	require.Equal(t, product2.Name, arg.Name)
	require.NotEqual(t, product2.Name, product1.Name)
	require.NotEqual(t, product2.ImageUrl, product1.ImageUrl)
}

func TestDeleteProduct(t *testing.T) {
	ctx := context.Background()
	defer tearDown(ctx)

	admin := createRandomAdmin(t, ctx)
	product := createRandomProduct(t, ctx, admin.Uuid, true)

	err := testProductRepo.DeleteProduct(ctx, product.Uuid, func() error {
		return errors.New("error")
	})

	require.Error(t, err)

	get, err := testProductRepo.GetProduct(ctx, product.Uuid)
	require.Equal(t, get.Uuid, product.Uuid)

	err = testProductRepo.DeleteProduct(ctx, product.Uuid, func() error {
		return nil
	})

	require.NoError(t, err)

	_, err = testProductRepo.GetProduct(ctx, product.Uuid)
	require.Error(t, err)

	arg := sqlc.ListVariantsParams{Limit: 100, Offset: 0}

	variants, err := testVariantRepo.GetAllVariants(ctx, arg)

	require.NoError(t, err)
	require.Len(t, variants, 0)
}
