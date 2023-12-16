-- +goose Up
-- +goose StatementBegin
CREATE TABLE products (
     id BIGINT NOT NULL AUTO_INCREMENT,
     title VARCHAR(100) NOT NULL,
     code VARCHAR(80) NOT NULL,
     description TEXT NOT NULL,
     price_in_cents BIGINT UNSIGNED NOT NULL,
     reference varchar(255) NOT NULL,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
     PRIMARY KEY (id),
     CONSTRAINT ukey_product_code UNIQUE (code)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
-- +goose StatementEnd



