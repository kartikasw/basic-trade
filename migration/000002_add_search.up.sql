ALTER TABLE products ADD COLUMN name_search tsvector GENERATED ALWAYS AS (to_tsvector('simple', name)) STORED;

ALTER TABLE variants ADD COLUMN variant_name_search tsvector GENERATED ALWAYS AS (to_tsvector('simple', variant_name)) STORED;

CREATE INDEX IF NOT EXISTS product__name_search__idx ON products USING GIN (name_search);

CREATE INDEX IF NOT EXISTS variant__name_search__idx ON variants USING GIN (variant_name_search);