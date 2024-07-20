DROP VIEW IF EXISTS product__view;

DROP INDEX IF EXISTS product__name_search__idx;

DROP INDEX IF EXISTS variant__name_search__idx;

ALTER TABLE products DROP COLUMN IF EXISTS name_search;

ALTER TABLE variants DROP COLUMN IF EXISTS variant_name_search;