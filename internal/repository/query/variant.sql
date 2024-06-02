-- name: CreateVariant :one
INSERT INTO variants (
    variant_name,
    quantity,
    product_id
) VALUES (
    $1, $2, $3
)
RETURNING uuid, variant_name, quantity;

-- name: GetVariantForUpdate :one
SELECT uuid, variant_name, quantity FROM variants
WHERE uuid = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: UpdateAVariant :one
UPDATE variants SET 
    variant_name = CASE WHEN sqlc.arg(set_variant_name)::bool
    THEN sqlc.arg(variant_name)::text
    ELSE variant_name
    END,
    quantity = $2
WHERE uuid = $1
RETURNING uuid, variant_name, quantity, product_id;

-- name: GetVariant :one
SELECT uuid, variant_name, quantity, product_id FROM variants
WHERE uuid = $1 LIMIT 1;

-- name: ListVariants :many
SELECT uuid, variant_name, quantity FROM variants
WHERE sqlc.arg(search)::bool AND variant_name LIKE $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: DeleteVariant :exec
DELETE FROM variants WHERE uuid = $1;