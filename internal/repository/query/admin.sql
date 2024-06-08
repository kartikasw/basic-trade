-- name: CreateAdmin :one
INSERT INTO admins (
    name,
    email,
    password
) VALUES (
    $1, $2, $3
)
RETURNING id, uuid, name, email, password;

-- name: GetAdmin :one
SELECT id, uuid, name, email, password FROM admins
WHERE (sqlc.arg(login)::bool AND email = $1) OR (sqlc.arg(get)::bool AND uuid = $2)
LIMIT 1;

-- name: CheckProductFromAdmin :one
SELECT products.uuid 
FROM admins
JOIN products ON admins.id = products.admin_id
WHERE admins.uuid = sqlc.arg(admin_uuid)::uuid AND products.uuid = sqlc.arg(product_uuid)::uuid
LIMIT 1;

-- name: CheckVariantFromAdmin :one
SELECT variants.uuid
FROM admins
JOIN products ON admins.id = products.admin_id
JOIN variants ON products.id = variants.product_id
WHERE admins.uuid = sqlc.arg(admin_uuid)::uuid AND variants.uuid = sqlc.arg(variant_uuid)::uuid
LIMIT 1;