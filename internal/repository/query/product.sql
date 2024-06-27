-- name: CreateProduct :one
INSERT INTO products (
    name,
    image_url,
    admin_id
) VALUES (
    $1, $2, $3
)
RETURNING uuid, name, image_url;

-- name: GetProductID :one
SELECT id FROM products
WHERE uuid = $1 
LIMIT 1;

-- name: GetProduct :one
SELECT uuid, name, image_url, variants FROM product__view
WHERE uuid = $1 
LIMIT 1;

-- name: GetProductForUpdate :one
SELECT uuid, name, image_url FROM products
WHERE uuid = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: UpdateAProduct :one
UPDATE products SET 
    name = CASE WHEN sqlc.arg(set_name)::bool
    THEN sqlc.arg(name)::text
    ELSE name
    END, 
    image_url = CASE WHEN sqlc.arg(set_image_url)::bool
    THEN sqlc.arg(image_url)::text
    ELSE image_url
    END
WHERE uuid = $1
RETURNING uuid, name, image_url;

-- name: ListProducts :many
SELECT 
    ROW_NUMBER() OVER (ORDER BY created_at DESC),
    uuid,
    name, 
    image_url, 
    variants
FROM product__view
WHERE sqlc.arg(keyword)::text = '' OR name_search @@ to_tsquery(sqlc.arg(keyword)::text)
LIMIT $1
OFFSET $2;

-- name: DeleteProduct :exec
DELETE FROM products WHERE uuid = $1;