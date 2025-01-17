// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: admin.sql

package repository

import (
	"context"

	"github.com/google/uuid"
)

const checkProductFromAdmin = `-- name: CheckProductFromAdmin :one
SELECT products.uuid 
FROM admins
JOIN products ON admins.id = products.admin_id
WHERE admins.uuid = $1::uuid AND products.uuid = $2::uuid
LIMIT 1
`

type CheckProductFromAdminParams struct {
	AdminUuid   uuid.UUID `json:"admin_uuid"`
	ProductUuid uuid.UUID `json:"product_uuid"`
}

func (q *Queries) CheckProductFromAdmin(ctx context.Context, arg CheckProductFromAdminParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, checkProductFromAdmin, arg.AdminUuid, arg.ProductUuid)
	var uuid uuid.UUID
	err := row.Scan(&uuid)
	return uuid, err
}

const checkVariantFromAdmin = `-- name: CheckVariantFromAdmin :one
SELECT variants.uuid
FROM admins
JOIN products ON admins.id = products.admin_id
JOIN variants ON products.id = variants.product_id
WHERE admins.uuid = $1::uuid AND variants.uuid = $2::uuid
LIMIT 1
`

type CheckVariantFromAdminParams struct {
	AdminUuid   uuid.UUID `json:"admin_uuid"`
	VariantUuid uuid.UUID `json:"variant_uuid"`
}

func (q *Queries) CheckVariantFromAdmin(ctx context.Context, arg CheckVariantFromAdminParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, checkVariantFromAdmin, arg.AdminUuid, arg.VariantUuid)
	var uuid uuid.UUID
	err := row.Scan(&uuid)
	return uuid, err
}

const createAdmin = `-- name: CreateAdmin :one
INSERT INTO admins (
    name,
    email,
    password
) VALUES (
    $1, $2, $3
)
RETURNING id, uuid, name, email, password
`

type CreateAdminParams struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateAdminRow struct {
	ID       int64     `json:"id"`
	Uuid     uuid.UUID `json:"uuid"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

func (q *Queries) CreateAdmin(ctx context.Context, arg CreateAdminParams) (CreateAdminRow, error) {
	row := q.db.QueryRow(ctx, createAdmin, arg.Name, arg.Email, arg.Password)
	var i CreateAdminRow
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.Name,
		&i.Email,
		&i.Password,
	)
	return i, err
}

const getAdmin = `-- name: GetAdmin :one
SELECT id, uuid, name, email, password FROM admins
WHERE ($3::bool AND email = $1) OR ($4::bool AND uuid = $2)
LIMIT 1
`

type GetAdminParams struct {
	Email string    `json:"email"`
	Uuid  uuid.UUID `json:"uuid"`
	Login bool      `json:"login"`
	Get   bool      `json:"get"`
}

type GetAdminRow struct {
	ID       int64     `json:"id"`
	Uuid     uuid.UUID `json:"uuid"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

func (q *Queries) GetAdmin(ctx context.Context, arg GetAdminParams) (GetAdminRow, error) {
	row := q.db.QueryRow(ctx, getAdmin,
		arg.Email,
		arg.Uuid,
		arg.Login,
		arg.Get,
	)
	var i GetAdminRow
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.Name,
		&i.Email,
		&i.Password,
	)
	return i, err
}
