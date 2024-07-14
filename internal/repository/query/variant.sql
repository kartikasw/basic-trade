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
SELECT 
    ROW_NUMBER() OVER (ORDER BY created_at DESC),
    uuid, 
    variant_name, 
    quantity 
FROM variants
WHERE sqlc.arg(keyword)::text = ':*' OR (sqlc.arg(keyword)::text != ':*' AND variant_name_search @@ to_tsquery('simple', sqlc.arg(keyword)::text))
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: DeleteVariant :exec
DELETE FROM variants WHERE uuid = $1;