package repository

import (
	sqlc "basic-trade/internal/repository/sqlc"
	"fmt"

	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IProductRepository struct {
	store *sqlc.Store
}

type ProductRepository interface {
	CreateProduct(ctx context.Context, arg sqlc.CreateProductParams, admUUID uuid.UUID, uplImage func(prdUUID uuid.UUID) (string, error)) (sqlc.CreateProductRow, error)
	GetProduct(ctx context.Context, uuid uuid.UUID) (sqlc.GetProductRow, error)
	GetAllProducts(ctx context.Context, arg sqlc.ListProductsParams) ([]sqlc.ListProductsRow, error)
	UpdateProduct(ctx context.Context, arg sqlc.UpdateAProductParams, uplImage func() (string, error)) (sqlc.UpdateAProductRow, error)
	DeleteProduct(ctx context.Context, prdUUID uuid.UUID, delImage func() error) error
}

func NewProductRepository(connPool *pgxpool.Pool) ProductRepository {
	return &IProductRepository{store: sqlc.NewStore(connPool)}
}

func (r *IProductRepository) CreateProduct(
	ctx context.Context,
	arg sqlc.CreateProductParams,
	admUUID uuid.UUID,
	uplImage func(prdUUID uuid.UUID) (string, error),
) (sqlc.CreateProductRow, error) {
	var result sqlc.CreateProductRow

	err := r.store.ExecTx(ctx, func(q *sqlc.Queries) error {
		admArg := sqlc.GetAdminParams{
			Get:  true,
			Uuid: admUUID,
		}
		admin, err := r.store.GetAdmin(ctx, admArg)
		if err != nil {
			return err
		}

		arg.AdminID = admin.ID
		result, err = r.store.CreateProduct(ctx, arg)
		if err != nil {
			return err
		}

		imageURL, err := uplImage(result.Uuid)
		if err != nil {
			return err
		}

		updArg := sqlc.UpdateAProductParams{
			Uuid:        result.Uuid,
			SetImageUrl: true,
			ImageUrl:    imageURL,
		}
		updRes, err := r.store.UpdateAProduct(ctx, updArg)
		if err != nil {
			return err
		}

		result = sqlc.CreateProductRow{
			Uuid:     updRes.Uuid,
			Name:     updRes.Name,
			ImageUrl: updRes.ImageUrl,
		}

		return nil
	})

	return result, err
}

func (r *IProductRepository) GetProduct(ctx context.Context, uuid uuid.UUID) (sqlc.GetProductRow, error) {
	result, err := r.store.GetProduct(ctx, uuid)

	return result, err
}

func (r *IProductRepository) GetAllProducts(ctx context.Context, arg sqlc.ListProductsParams) ([]sqlc.ListProductsRow, error) {
	arg.Keyword = fmt.Sprintf("%s:*", arg.Keyword)
	result, err := r.store.ListProducts(ctx, arg)

	return result, err
}

func (r *IProductRepository) UpdateProduct(
	ctx context.Context,
	arg sqlc.UpdateAProductParams,
	uplImage func() (string, error),
) (sqlc.UpdateAProductRow, error) {
	var result sqlc.UpdateAProductRow

	err := r.store.ExecTx(ctx, func(q *sqlc.Queries) error {
		_, err := q.GetProductForUpdate(ctx, arg.Uuid)
		if err != nil {
			return err
		}

		if arg.Name != "" {
			arg.SetName = true
		}

		if uplImage != nil {
			imageURL, err := uplImage()
			if err != nil {
				return err
			}

			arg.SetImageUrl = true
			arg.ImageUrl = imageURL
		}

		result, err = q.UpdateAProduct(ctx, arg)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func (r *IProductRepository) DeleteProduct(ctx context.Context, prdUUID uuid.UUID, delImage func() error) error {
	err := r.store.ExecTx(ctx, func(q *sqlc.Queries) error {
		err := q.DeleteProduct(ctx, prdUUID)
		if err != nil {
			return err
		}

		err = delImage()
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
