package entity

import (
	sqlc "basic-trade/internal/repository/sqlc"
	"encoding/json"

	"github.com/google/uuid"
)

type Product struct {
	UUID     uuid.UUID `json:"uuid"`
	Name     string    `json:"name"`
	ImageURL string    `json:"image_url"`
	Variants []Variant `json:"variants"`
}

func CreateProductToViewModel(product *sqlc.CreateProductRow) Product {
	return Product{
		UUID:     product.Uuid,
		Name:     product.Name,
		ImageURL: product.ImageUrl,
	}
}

func ProductViewToViewModel(product *sqlc.ProductView) Product {
	variants := make([]Variant, len(product.Variants.([]interface{})))

	for i, item := range product.Variants.([]interface{}) {
		if m, ok := item.(map[string]interface{}); ok {
			var variant Variant
			byte, err := json.Marshal(m)
			if err != nil {
				break
			}

			err = json.Unmarshal(byte, &variant)
			if err != nil {
				break
			}

			variants[i] = variant
		}
	}

	return Product{
		UUID:     product.Uuid,
		Name:     product.Name,
		ImageURL: product.ImageUrl,
		Variants: variants,
	}
}

func UpdateProductToViewModel(product *sqlc.UpdateAProductRow) Product {
	return Product{
		UUID:     product.Uuid,
		Name:     product.Name,
		ImageURL: product.ImageUrl,
	}
}

func ListProductViewToViewModel(products []sqlc.ProductView) []Product {
	list := make([]Product, len(products))

	for index, item := range products {
		list[index] = ProductViewToViewModel(&item)
	}

	return list
}
