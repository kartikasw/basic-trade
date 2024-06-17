package repository

import (
	"basic-trade/common"
	sqlc "basic-trade/internal/repository/sqlc"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomProduct(t *testing.T, adminUUID uuid.UUID, withVariants bool) sqlc.CreateProductRow {
	arg := sqlc.CreateProductParams{
		Name:     common.RandomName(),
		ImageUrl: common.RandomString(6),
	}

	product, err := testProductRepo.CreateProduct(nil, arg, adminUUID, nil)

	require.NoError(t, err)
	require.Equal(t, arg.Name, product.Name)
	require.Equal(t, arg.ImageUrl, product.ImageUrl)

	if withVariants {
		n := int(common.RandomInt(1, 10))

		for i := 1; i <= n; i++ {
			createRandomVariant(t, product.Uuid)
		}
	}

	return product
}

func TestCreateProduct(t *testing.T) {
	admin := createRandomAdmin(t)

	createRandomProduct(t, admin.Uuid, false)
}

func TestGetProduct(t *testing.T) {
	admin := createRandomAdmin(t)
	product1 := createRandomProduct(t, admin.Uuid, false)

	product2, err := testProductRepo.GetProduct(nil, product1.Uuid)

	require.NoError(t, err)
	require.Equal(t, product1.Name, product2.Name)
	require.Equal(t, product1.ImageUrl, product2.ImageUrl)
	require.Equal(t, product1.Uuid, product2.Uuid)
	require.Equal(t, product1.Uuid, product2.Uuid)
	require.Len(t, product2.Variants.([]interface{}), 0)
}

func TestGetAllProducts(t *testing.T) {
	admin := createRandomAdmin(t)

	for i := 1; 1 <= 10; i++ {
		product := createRandomProduct(t, admin.Uuid, false)
		require.Equal(t, product.Uuid, admin.Uuid)
	}

	arg := sqlc.ListProductsParams{
		Limit:  10,
		Offset: 0,
	}

	products, err := testProductRepo.GetAllProducts(nil, arg)

	require.NoError(t, err)
	require.Len(t, products, 10)
}

func TestGetAllProductsForSearch(t *testing.T) {}

func TestUpdateProduct(t *testing.T) {
	admin := createRandomAdmin(t)
	product1 := createRandomProduct(t, admin.Uuid, false)

	arg := sqlc.UpdateAProductParams{
		Uuid:    product1.Uuid,
		SetName: true,
		Name:    common.RandomName(),
	}

	product2, err := testProductRepo.UpdateProduct(nil, arg, nil)

	require.NoError(t, err)
	require.Equal(t, product2.Uuid, product1.Uuid)
	require.Equal(t, product2.Name, arg.Name)
	require.NotEqual(t, product2.Name, product1.Name)
}

func TestDeleteProduct(t *testing.T) {
	admin := createRandomAdmin(t)
	product := createRandomProduct(t, admin.Uuid, true)

	err := testProductRepo.DeleteProduct(nil, product.Uuid, nil)

	require.NoError(t, err)

	_, err = testProductRepo.GetProduct(nil, product.Uuid)

	require.Error(t, err)

	arg := sqlc.ListVariantsParams{
		Limit:  100,
		Offset: 0,
	}

	variants, err := testVariantRepo.GetAllVariants(nil, arg)

	require.NoError(t, err)
	require.Len(t, variants, 0)
}
