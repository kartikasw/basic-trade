package entity

import (
	sqlc "basic-trade/internal/repository/sqlc"

	"github.com/google/uuid"
)

type Product struct {
	UUID      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	ImageURL  string    `json:"image_url"`
	AdminUUID uuid.UUID `json:"admin_id"`
}

func CreateProductToViewModel(product *sqlc.CreateProductRow, uuidAdm uuid.UUID) Product {
	return Product{
		UUID:      product.Uuid,
		Name:      product.Name,
		ImageURL:  product.ImageUrl,
		AdminUUID: uuidAdm,
	}
}

func GetProductToViewModel(product *sqlc.GetProductRow) Product {
	return Product{
		UUID:      product.Uuid,
		Name:      product.Name,
		ImageURL:  product.ImageUrl,
	}
}

func UpdateProductToViewModel(product *sqlc.UpdateAProductRow) Product {
	return Product{
		UUID:      product.Uuid,
		Name:      product.Name,
		ImageURL:  product.ImageUrl,
	}
}

func ListProductToViewModel(products []sqlc.ListProductsRow) []Product {
	list := make([]Product, len(products))

	for index, item := range products {
		list[index] = Product{
			UUID:     item.Uuid,
			Name:     item.Name,
			ImageURL: item.ImageUrl,
		}
	}

	return list
}
