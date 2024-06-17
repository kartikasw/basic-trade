package entity

import (
	sqlc "basic-trade/internal/repository/sqlc"
	"encoding/json"

	"github.com/google/uuid"
)

type Variant struct {
	UUID        uuid.UUID `json:"uuid"`
	VariantName string    `json:"variant_name"`
	Quantity    int32     `json:"quantity"`
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

func ListVariantToViewModel(variants []sqlc.ListVariantsRow) []Variant {
	list := make([]Variant, len(variants))

	for index, item := range variants {
		list[index] = Variant{
			UUID:        item.Uuid,
			VariantName: item.VariantName,
			Quantity:    item.Quantity,
		}
	}

	return list
}
