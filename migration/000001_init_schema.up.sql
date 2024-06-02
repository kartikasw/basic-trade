BEGIN;

CREATE TABLE IF NOT EXISTS admins (
  id bigserial NOT NULL,
  uuid uuid NOT NULL UNIQUE DEFAULT gen_random_uuid(),
  name varchar(100) NOT NULL,
  email varchar(100) NOT NULL UNIQUE,
  password varchar NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now()),
  updated_at timestamptz,

  CONSTRAINT admin__pkey PRIMARY KEY (id)
);

CREATE TABLE products (
  id bigserial NOT NULL,
  uuid uuid NOT NULL UNIQUE DEFAULT gen_random_uuid(),
  name varchar(100) NOT NULL,
  image_url varchar NOT NULL,
  admin_id bigserial NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now()),
  updated_at timestamptz,

  CONSTRAINT product__pkey PRIMARY KEY (id),
  CONSTRAINT admin_id__fk FOREIGN KEY (admin_id) REFERENCES admins(id) ON DELETE CASCADE
);

CREATE TABLE variants (
  id bigserial NOT NULL,
  uuid uuid NOT NULL UNIQUE DEFAULT gen_random_uuid(),
  variant_name varchar(100) NOT NULL,
  quantity int NOT NULL,
  product_id bigserial NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now()),
  updated_at timestamptz,

  CONSTRAINT variant__pkey PRIMARY KEY (id),
  CONSTRAINT product_id__fk FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS admin__email__idx ON admins USING BTREE (email);

CREATE INDEX IF NOT EXISTS admin__uuid__idx ON admins USING BTREE (uuid);

CREATE INDEX IF NOT EXISTS product__admin_id__idx ON products USING BTREE (admin_id);

CREATE INDEX IF NOT EXISTS product__uuid__idx ON products USING BTREE (uuid);

CREATE INDEX IF NOT EXISTS variant__product_id__idx ON variants USING BTREE (product_id);

CREATE INDEX IF NOT EXISTS variant__uuid__idx ON variants USING BTREE (uuid);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_admins_updated_at
BEFORE UPDATE ON admins
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_products_updated_at
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_variants_updated_at
BEFORE UPDATE ON variants
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

COMMIT;