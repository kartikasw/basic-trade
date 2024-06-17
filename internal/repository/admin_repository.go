package repository

import (
	sqlc "basic-trade/internal/repository/sqlc"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IAdminRepository struct {
	store *sqlc.Store
}

type AdminRepository interface {
	CreateAdmin(ctx context.Context, admin sqlc.CreateAdminParams) (sqlc.CreateAdminRow, error)
	GetAdmin(email string) (sqlc.GetAdminRow, error)
	CheckProductFromAdmin(uuiAdm uuid.UUID, uuidPrd uuid.UUID) bool
	CheckVariantFromAdmin(uuiAdm uuid.UUID, uuidVrt uuid.UUID) bool
}

func NewAdminRepository(connPool *pgxpool.Pool) AdminRepository {
	return &IAdminRepository{store: sqlc.NewStore(connPool)}
}

func (r *IAdminRepository) CreateAdmin(ctx context.Context, arg sqlc.CreateAdminParams) (sqlc.CreateAdminRow, error) {
	var result sqlc.CreateAdminRow

	err := r.store.ExecTx(ctx, func(q *sqlc.Queries) error {
		var err error
		result, err = r.store.CreateAdmin(ctx, arg)

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func (r *IAdminRepository) GetAdmin(email string) (sqlc.GetAdminRow, error) {
	arg := sqlc.GetAdminParams{
		Login: true,
		Email: email,
	}
	result, err := r.store.GetAdmin(context.Background(), arg)

	return result, err
}

func (r *IAdminRepository) CheckProductFromAdmin(uuiAdm uuid.UUID, uuidPrd uuid.UUID) bool {
	arg := sqlc.CheckProductFromAdminParams{
		AdminUuid:   uuiAdm,
		ProductUuid: uuidPrd,
	}

	result, err := r.store.CheckProductFromAdmin(context.Background(), arg)
	if err != nil {
		return false
	}

	if result != uuidPrd {
		return false
	}

	return true
}

func (r *IAdminRepository) CheckVariantFromAdmin(uuiAdm uuid.UUID, uuidVrt uuid.UUID) bool {
	arg := sqlc.CheckVariantFromAdminParams{
		AdminUuid:   uuiAdm,
		VariantUuid: uuidVrt,
	}

	result, err := r.store.CheckVariantFromAdmin(context.Background(), arg)
	if err != nil {
		return false
	}

	if result != uuidVrt {
		return false
	}

	return true
}
