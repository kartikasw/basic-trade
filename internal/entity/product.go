package entity

import (
	sqlc "basic-trade/internal/repository/sqlc"

	"github.com/google/uuid"
)

type Product struct {
	UUID     uuid.UUID `json:"uuid"`
	Name     string    `json:"name"`
	ImageURL string    `json:"image_url"`
}

type ProductView struct {
	UUID     uuid.UUID `json:"uuid"`
	Name     string    `json:"name"`
	ImageURL string    `json:"image_url"`
	Variants []Variant `json:"variants"`
}

type ProductViewList struct {
	RowID    int64     `json:"row_id"`
	UUID     uuid.UUID `json:"uuid"`
	Name     string    `json:"name"`
	ImageURL string    `json:"image_url"`
	Variants []Variant `json:"variants"`
}

func CreateProductToViewModel(product sqlc.CreateProductRow) Product {
	return Product{
		UUID:     product.Uuid,
		Name:     product.Name,
		ImageURL: product.ImageUrl,
	}
}

func GetProductRowToViewModel(product sqlc.GetProductRow) ProductView {
	return ProductView{
		UUID:     product.Uuid,
		Name:     product.Name,
		ImageURL: product.ImageUrl,
		Variants: VariantInterfaceToEntityList(product.Variants),
	}
}

func UpdateProductToViewModel(product sqlc.UpdateAProductRow) Product {
	return Product{
		UUID:     product.Uuid,
		Name:     product.Name,
		ImageURL: product.ImageUrl,
	}
}

func ListProductsRowToViewModel(products []sqlc.ListProductsRow) []ProductViewList {
	list := make([]ProductViewList, len(products))

	for index, item := range products {
		list[index] = ProductViewList{
			RowID:    item.RowNumber,
			UUID:     item.Uuid,
			Name:     item.Name,
			ImageURL: item.ImageUrl,
			Variants: VariantInterfaceToEntityList(item.Variants.([]interface{})),
		}
	}

	return list
}
