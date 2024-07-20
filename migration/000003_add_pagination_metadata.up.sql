CREATE VIEW product__view__admin AS
SELECT 
  p.uuid, 
  p.name, 
  p.image_url,
  p.name_search,
  COALESCE(
    json_agg(
      json_build_object(
        'uuid', v.uuid,
        'variant_name', v.variant_name,
        'quantity', v.quantity
      )
    ) FILTER (WHERE v.uuid IS NOT NULL), 
    '[]'
  ) AS variants,
  json_build_object(
    'uuid', a.uuid,
    'name', a.name,
    'email', a.email
  ) AS admin,
  p.created_at
FROM products p
LEFT JOIN variants v ON p.id = v.product_id
INNER JOIN admins a ON p.admin_id = a.id
GROUP BY p.id, a.id;