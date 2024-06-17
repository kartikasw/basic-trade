package repository

import (
	"testing"

	"basic-trade/common"
	sqlc "basic-trade/internal/repository/sqlc"
	"basic-trade/pkg/password"

	"github.com/stretchr/testify/require"
)

func createRandomAdmin(t *testing.T) sqlc.CreateAdminRow {
	hashPass, err := password.HashPassword(common.RandomString(8))
	require.NoError(t, err)

	arg := sqlc.CreateAdminParams{
		Name:     common.RandomName(),
		Email:    common.RandomEmail(),
		Password: hashPass,
	}

	admin, err := testAdminRepo.CreateAdmin(nil, arg)

	require.NoError(t, err)
	require.Equal(t, arg.Name, admin.Name)
	require.Equal(t, arg.Email, admin.Email)
	require.Equal(t, arg.Email, admin.Email)

	return admin
}

func TextCreateAdmin(t *testing.T) {
	createRandomAdmin(t)
}

func TestGetAdmin(t *testing.T) {
	admin1 := createRandomAdmin(t)

	admin2, err := testAdminRepo.GetAdmin(admin1.Email)

	require.NoError(t, err)
	require.Equal(t, admin1.Name, admin2.Name)
	require.Equal(t, admin1.Email, admin2.Email)
	require.Equal(t, admin1.Password, admin2.Password)
}

func TestCheckProductFromAdmin(t *testing.T) {
	admin := createRandomAdmin(t)
	product := createRandomProduct(t, admin.Uuid, false)

	result := testAdminRepo.CheckProductFromAdmin(admin.Uuid, product.Uuid)
	require.True(t, result)
}

func TestCheckVariantFromAdmin(t *testing.T) {
	admin := createRandomAdmin(t)
	product := createRandomProduct(t, admin.Uuid, false)
	variant := createRandomVariant(t, product.Uuid)

	result := testAdminRepo.CheckVariantFromAdmin(admin.Uuid, variant.Uuid)
	require.True(t, result)
}
