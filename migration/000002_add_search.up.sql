ALTER TABLE products ADD COLUMN name_search tsvector GENERATED ALWAYS AS (to_tsvector('simple', name)) STORED;

ALTER TABLE variants ADD COLUMN variant_name_search tsvector GENERATED ALWAYS AS (to_tsvector('simple', variant_name)) STORED;

CREATE INDEX IF NOT EXISTS product__name_search__idx ON products USING GIN (name_search);

CREATE INDEX IF NOT EXISTS variant__name_search__idx ON variants USING GIN (variant_name_search);

CREATE VIEW product__view AS
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
  p.created_at
FROM products p
LEFT JOIN variants v ON p.id = v.product_id
GROUP BY p.id;