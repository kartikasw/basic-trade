package entity

import (
	sqlc "basic-trade/internal/repository/sqlc"

	"github.com/google/uuid"
)

type Admin struct {
	UUID     uuid.UUID `json:"uuid"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

func GetAdminToViewModel(admin sqlc.GetAdminRow) Admin {
	return Admin{
		UUID:     admin.Uuid,
		Name:     admin.Name,
		Email:    admin.Email,
		Password: admin.Password,
	}
}

func CreateAdminToViewModel(admin sqlc.CreateAdminRow) Admin {
	return Admin{
		UUID:     admin.Uuid,
		Name:     admin.Name,
		Email:    admin.Email,
		Password: admin.Password,
	}
}
