package entity

import (
	sqlc "basic-trade/internal/repository/sqlc"
	"encoding/json"
	"math"

	"github.com/google/uuid"
)

type Pagination struct {
	LastPage int32 `json:"last_page"`
	Limit    int32 `json:"limit"`
	Offset   int32 `json:"offset"`
	Page     int32 `json:"page"`
	Total    int32 `json:"total"`
}

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
	RowID    int64        `json:"row_id"`
	UUID     uuid.UUID    `json:"uuid"`
	Name     string       `json:"name"`
	ImageURL string       `json:"image_url"`
	Variants []Variant    `json:"variants"`
	Admin    AdminProduct `json:"admin"`
}

type ProductPaginationView struct {
	Data       []ProductViewList `json:"products"`
	Pagination Pagination        `json:"pagination"`
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

func ListProductsRowToViewModel(products []sqlc.ListProductsRow, limit int32, offset int32, total int64) (ProductPaginationView, error) {
	list := make([]ProductViewList, len(products))

	for index, item := range products {
		var admin AdminProduct

		err := json.Unmarshal(item.Admin, &admin)
		if err != nil {
			return ProductPaginationView{}, err
		}

		list[index] = ProductViewList{
			RowID:    item.RowNumber,
			UUID:     item.Uuid,
			Name:     item.Name,
			ImageURL: item.ImageUrl,
			Variants: VariantInterfaceToEntityList(item.Variants.([]interface{})),
			Admin:    admin,
		}
	}

	return ListProductsViewPaginationToModel(list, limit, offset, total), nil
}

func ListProductsViewPaginationToModel(products []ProductViewList, limit int32, offset int32, total int64) ProductPaginationView {
	var pagination ProductPaginationView
	var productLen int32 = int32(len(products))

	pagination = ProductPaginationView{
		Data: products,
		Pagination: Pagination{
			Limit:    limit,
			Offset:   offset,
			Page:     offset + 1,
			LastPage: int32(math.Ceil(float64(total) / float64(limit))),
			Total:    productLen,
		},
	}

	return pagination
}
