package repository

import (
	"context"
	"testing"

	"basic-trade/common"
	sqlc "basic-trade/internal/repository/sqlc"
	"basic-trade/pkg/password"

	"github.com/stretchr/testify/require"
)

func createRandomAdmin(t *testing.T, ctx context.Context) sqlc.CreateAdminRow {
	hashPass, err := password.HashPassword(common.RandomString(8))
	require.NoError(t, err)

	arg := sqlc.CreateAdminParams{
		Name:     common.RandomName(),
		Email:    common.RandomEmail(),
		Password: hashPass,
	}

	admin, err := testAdminRepo.CreateAdmin(ctx, arg)

	require.NoError(t, err)
	require.Equal(t, arg.Name, admin.Name)
	require.Equal(t, arg.Email, admin.Email)
	require.Equal(t, arg.Email, admin.Email)

	return admin
}

func TestCreateAdmin(t *testing.T) {
	ctx := context.Background()
	defer tearDown(ctx)

	createRandomAdmin(t, ctx)
}

func TestGetAdmin(t *testing.T) {
	ctx := context.Background()
	defer tearDown(ctx)

	admin1 := createRandomAdmin(t, ctx)

	admin2, err := testAdminRepo.GetAdmin(ctx, admin1.Email)

	require.NoError(t, err)
	require.Equal(t, admin1.Name, admin2.Name)
	require.Equal(t, admin1.Email, admin2.Email)
	require.Equal(t, admin1.Password, admin2.Password)
}

func TestCheckProductFromAdmin(t *testing.T) {
	ctx := context.Background()
	defer tearDown(ctx)

	admin1 := createRandomAdmin(t, ctx)
	product1 := createRandomProduct(t, ctx, admin1.Uuid, false)
	admin2 := createRandomAdmin(t, ctx)

	result1 := testAdminRepo.CheckProductFromAdmin(ctx, admin1.Uuid, product1.Uuid)
	require.True(t, result1)

	result2 := testAdminRepo.CheckProductFromAdmin(ctx, admin2.Uuid, product1.Uuid)
	require.False(t, result2)
}

func TestCheckVariantFromAdmin(t *testing.T) {
	ctx := context.Background()
	defer tearDown(ctx)

	admin1 := createRandomAdmin(t, ctx)
	product := createRandomProduct(t, ctx, admin1.Uuid, false)
	variant := createRandomVariant(t, ctx, product.Uuid)
	admin2 := createRandomAdmin(t, ctx)

	result1 := testAdminRepo.CheckVariantFromAdmin(ctx, admin1.Uuid, variant.Uuid)
	require.True(t, result1)

	result2 := testAdminRepo.CheckVariantFromAdmin(ctx, admin2.Uuid, variant.Uuid)
	require.False(t, result2)
}
