package repository

import (
	sqlc "basic-trade/internal/repository/sqlc"

	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	store *sqlc.Store
}

type IProductRepository interface {
	CreateProduct(arg sqlc.CreateProductParams, uuidAdm uuid.UUID) (sqlc.CreateProductRow, error)
	GetProduct(uuid uuid.UUID) (sqlc.ProductView, error)
	GetAllProducts(arg sqlc.ListProductsParams) ([]sqlc.ProductView, error)
	UpdateProduct(arg sqlc.UpdateAProductParams) (sqlc.UpdateAProductRow, error)
	DeleteProduct(uuid uuid.UUID) error
}

func NewProductRepository(connPool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{store: sqlc.NewStore(connPool)}
}

func (r *ProductRepository) CreateProduct(arg sqlc.CreateProductParams, uuidAdm uuid.UUID) (sqlc.CreateProductRow, error) {
	ctx := context.Background()

	var result sqlc.CreateProductRow

	err := r.store.ExecTx(ctx, func(q *sqlc.Queries) error {
		admArg := sqlc.GetAdminParams{
			Get:  true,
			Uuid: uuidAdm,
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

		return nil
	})

	return result, err
}

func (r *ProductRepository) GetProduct(uuid uuid.UUID) (sqlc.ProductView, error) {
	result, err := r.store.GetProduct(context.Background(), uuid)

	return result, err
}

func (r *ProductRepository) GetAllProducts(arg sqlc.ListProductsParams) ([]sqlc.ProductView, error) {
	result, err := r.store.ListProducts(context.Background(), arg)

	return result, err
}

func (r *ProductRepository) UpdateProduct(arg sqlc.UpdateAProductParams) (sqlc.UpdateAProductRow, error) {
	var result sqlc.UpdateAProductRow

	ctx := context.Background()

	err := r.store.ExecTx(ctx, func(q *sqlc.Queries) error {
		product, err := q.GetProductForUpdate(ctx, arg.Uuid)
		if err != nil {
			return err
		}

		if arg.Name != "" {
			arg.SetName = true
			arg.Name = product.Name
		}

		if arg.ImageUrl != "" {
			arg.SetImageUrl = true
			arg.ImageUrl = product.ImageUrl
		}

		result, err = q.UpdateAProduct(ctx, arg)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func (r *ProductRepository) DeleteProduct(uuid uuid.UUID) error {
	ctx := context.Background()

	err := r.store.ExecTx(ctx, func(q *sqlc.Queries) error {
		err := q.DeleteProduct(ctx, uuid)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
