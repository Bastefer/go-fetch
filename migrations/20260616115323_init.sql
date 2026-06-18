-- +goose Up
CREATE TABLE "brands" (
    "id" SERIAL PRIMARY KEY,
    "name" text NOT NULL
);

CREATE UNIQUE INDEX brands_name_unique ON brands USING btree (name);

CREATE TABLE "categories" (
    "id" SERIAL PRIMARY KEY,
    "name" text NOT NULL
);

CREATE UNIQUE INDEX categories_name_unique ON categories USING btree (name);

CREATE TABLE "clients" (
    "id" integer PRIMARY KEY,
    "first_name" text NOT NULL,
    "last_name" text NOT NULL
);

CREATE TABLE "clients_products" (
    "id" SERIAL PRIMARY KEY,
    "client_id" integer NOT NULL,
    "product_id" integer NOT NULL
);
CREATE INDEX idx_clients_products_client_id ON "clients_products" (client_id);
CREATE INDEX idx_clients_products_product_id ON "clients_products" (product_id);

CREATE TABLE "products" (
    "id" integer PRIMARY KEY,
    "name" text NOT NULL,
    "brand_id" integer NOT NULL,
    "category_id" integer NOT NULL,
    "price" bigint NOT NULL,
    "stock" integer NOT NULL
);
CREATE INDEX idx_products_category_id ON "products" (category_id);
CREATE INDEX idx_products_brand_id ON "products" (brand_id);

ALTER TABLE ONLY "clients_products"
ADD CONSTRAINT "clients_products_client_fkey" FOREIGN KEY (client_id) REFERENCES clients (id) NOT DEFERRABLE;

ALTER TABLE ONLY "clients_products"
ADD CONSTRAINT "clients_products_product_fkey" FOREIGN KEY (product_id) REFERENCES products (id) NOT DEFERRABLE;

ALTER TABLE ONLY "products"
ADD CONSTRAINT "products_brand_fkey" FOREIGN KEY (brand_id) REFERENCES brands (id) NOT DEFERRABLE;

ALTER TABLE ONLY "products"
ADD CONSTRAINT "products_category_fkey" FOREIGN KEY (category_id) REFERENCES categories (id) NOT DEFERRABLE;

-- +goose Down
DROP TABLE IF EXISTS clients_products;

DROP TABLE IF EXISTS products;

DROP TABLE IF EXISTS clients;

DROP TABLE IF EXISTS categories;

DROP TABLE IF EXISTS brands;