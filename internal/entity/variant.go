package entity

import (
	sqlc "basic-trade/internal/repository/sqlc"
	"encoding/json"
	"math"

	"github.com/google/uuid"
)

type Variant struct {
	UUID        uuid.UUID `json:"uuid"`
	VariantName string    `json:"variant_name"`
	Quantity    int32     `json:"quantity"`
}

type VariantViewList struct {
	RowID       int64     `json:"row_id"`
	UUID        uuid.UUID `json:"uuid"`
	VariantName string    `json:"variant_name"`
	Quantity    int32     `json:"quantity"`
}

type VariantPaginationView struct {
	Data       []VariantViewList `json:"variants"`
	Pagination Pagination        `json:"pagination"`
}

func VariantInterfaceToEntityList(items interface{}) []Variant {
	variants := make([]Variant, len(items.([]interface{})))

	for i, item := range items.([]interface{}) {
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

	return variants
}

func CreateVariantToViewModel(variant sqlc.CreateVariantRow) Variant {
	return Variant{
		UUID:        variant.Uuid,
		VariantName: variant.VariantName,
		Quantity:    variant.Quantity,
	}
}

func GetVariantToViewModel(variant sqlc.GetVariantRow) Variant {
	return Variant{
		UUID:        variant.Uuid,
		VariantName: variant.VariantName,
		Quantity:    variant.Quantity,
	}
}

func UpdateVariantToViewModel(variant sqlc.UpdateAVariantRow) Variant {
	return Variant{
		UUID:        variant.Uuid,
		VariantName: variant.VariantName,
		Quantity:    variant.Quantity,
	}
}

func ListVariantToViewModel(variants []sqlc.ListVariantsRow, limit int32, offset int32, total int64) VariantPaginationView {
	list := make([]VariantViewList, len(variants))

	for index, item := range variants {
		list[index] = VariantViewList{
			RowID:       item.RowNumber,
			UUID:        item.Uuid,
			VariantName: item.VariantName,
			Quantity:    item.Quantity,
		}
	}

	return ListVariantsViewPaginationToModel(list, limit, offset, total)
}

func ListVariantsViewPaginationToModel(variants []VariantViewList, limit int32, offset int32, total int64) VariantPaginationView {
	var pagination VariantPaginationView
	var variantLen int32 = int32(len(variants))

	pagination = VariantPaginationView{
		Data: variants,
		Pagination: Pagination{
			Limit:    limit,
			Offset:   offset,
			Page:     offset + 1,
			LastPage: int32(math.Ceil(float64(total) / float64(limit))),
			Total:    variantLen,
		},
	}

	return pagination
}
