package entity

import (
	sqlc "basic-trade/internal/repository/sqlc"

	"github.com/google/uuid"
)

type Variant struct {
	UUID        uuid.UUID `json:"uuid"`
	VariantName string    `json:"variant_name"`
	Quantity    int32     `json:"quantity"`
}

func CreateVariantToViewModel(variant *sqlc.CreateVariantRow) Variant {
	return Variant{
		UUID:        variant.Uuid,
		VariantName: variant.VariantName,
		Quantity:    variant.Quantity,
	}
}

func GetVariantToViewModel(variant *sqlc.GetVariantRow) Variant {
	return Variant{
		UUID:        variant.Uuid,
		VariantName: variant.VariantName,
		Quantity:    variant.Quantity,
	}
}

func UpdateVariantToViewModel(variant *sqlc.UpdateAVariantRow) Variant {
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
