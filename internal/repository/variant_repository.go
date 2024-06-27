package repository

import (
	sqlc "basic-trade/internal/repository/sqlc"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IVariantRepository struct {
	store *sqlc.Store
}

type VariantRepository interface {
	CreateVariant(ctx context.Context, arg sqlc.CreateVariantParams, uuidPrd uuid.UUID) (sqlc.CreateVariantRow, error)
	GetVariant(ctx context.Context, uuid uuid.UUID) (sqlc.GetVariantRow, error)
	GetAllVariants(ctx context.Context, arg sqlc.ListVariantsParams) ([]sqlc.ListVariantsRow, error)
	UpdateVariant(ctx context.Context, arg sqlc.UpdateAVariantParams) (sqlc.UpdateAVariantRow, error)
	DeleteVariant(ctx context.Context, uuid uuid.UUID) error
}

func NewVariantRepository(connPool *pgxpool.Pool) VariantRepository {
	return &IVariantRepository{store: sqlc.NewStore(connPool)}
}

func (r *IVariantRepository) CreateVariant(
	ctx context.Context,
	arg sqlc.CreateVariantParams,
	uuidPrd uuid.UUID,
) (sqlc.CreateVariantRow, error) {
	var result sqlc.CreateVariantRow

	err := r.store.ExecTx(ctx, func(q *sqlc.Queries) error {
		prdID, err := r.store.GetProductID(ctx, uuidPrd)
		if err != nil {
			return err
		}

		arg.ProductID = prdID
		result, err = r.store.CreateVariant(context.Background(), arg)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func (r *IVariantRepository) GetVariant(ctx context.Context, uuid uuid.UUID) (sqlc.GetVariantRow, error) {
	result, err := r.store.GetVariant(ctx, uuid)

	return result, err
}

func (r *IVariantRepository) GetAllVariants(ctx context.Context, arg sqlc.ListVariantsParams) ([]sqlc.ListVariantsRow, error) {
	arg.Keyword = fmt.Sprintf("%s:*", arg.Keyword)
	result, err := r.store.ListVariants(ctx, arg)

	return result, err
}

func (r *IVariantRepository) UpdateVariant(ctx context.Context, arg sqlc.UpdateAVariantParams) (sqlc.UpdateAVariantRow, error) {
	var result sqlc.UpdateAVariantRow

	err := r.store.ExecTx(ctx, func(q *sqlc.Queries) error {
		_, err := q.GetVariantForUpdate(ctx, arg.Uuid)
		if err != nil {
			return err
		}

		if arg.VariantName != "" {
			arg.SetVariantName = true
		}

		result, err = q.UpdateAVariant(ctx, arg)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func (r *IVariantRepository) DeleteVariant(ctx context.Context, uuid uuid.UUID) error {
	err := r.store.ExecTx(ctx, func(q *sqlc.Queries) error {
		err := q.DeleteVariant(ctx, uuid)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
