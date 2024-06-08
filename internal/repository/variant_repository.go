package repository

import (
	sqlc "basic-trade/internal/repository/sqlc"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VariantRepository struct {
	store *sqlc.Store
}

type IVariantRepository interface {
	CreateVariant(arg sqlc.CreateVariantParams, uuidPrd uuid.UUID) (sqlc.CreateVariantRow, error)
	GetVariant(uuid uuid.UUID) (sqlc.GetVariantRow, error)
	GetAllVariants(arg sqlc.ListVariantsParams) ([]sqlc.ListVariantsRow, error)
	UpdateVariant(arg sqlc.UpdateAVariantParams) (sqlc.UpdateAVariantRow, error)
	DeleteVariant(uuid uuid.UUID) error
}

func NewVariantRepository(connPool *pgxpool.Pool) *VariantRepository {
	return &VariantRepository{store: sqlc.NewStore(connPool)}
}

func (r *VariantRepository) CreateVariant(arg sqlc.CreateVariantParams, uuidPrd uuid.UUID) (sqlc.CreateVariantRow, error) {
	ctx := context.Background()

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

func (r *VariantRepository) GetVariant(uuid uuid.UUID) (sqlc.GetVariantRow, error) {
	result, err := r.store.GetVariant(context.Background(), uuid)

	return result, err
}

func (r *VariantRepository) GetAllVariants(arg sqlc.ListVariantsParams) ([]sqlc.ListVariantsRow, error) {
	result, err := r.store.ListVariants(context.Background(), arg)

	return result, err
}

func (r *VariantRepository) UpdateVariant(arg sqlc.UpdateAVariantParams) (sqlc.UpdateAVariantRow, error) {
	var result sqlc.UpdateAVariantRow

	ctx := context.Background()

	err := r.store.ExecTx(ctx, func(q *sqlc.Queries) error {
		variant, err := q.GetVariantForUpdate(ctx, arg.Uuid)
		if err != nil {
			return err
		}

		if arg.VariantName != "" {
			arg.SetVariantName = true
			arg.VariantName = variant.VariantName
		}

		result, err = q.UpdateAVariant(ctx, arg)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func (r *VariantRepository) DeleteVariant(uuid uuid.UUID) error {
	ctx := context.Background()

	err := r.store.ExecTx(ctx, func(q *sqlc.Queries) error {
		err := q.DeleteVariant(ctx, uuid)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
