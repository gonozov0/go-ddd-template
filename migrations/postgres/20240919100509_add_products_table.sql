-- +goose Up
-- +goose StatementBegin
CREATE TABLE products (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    status VARCHAR NOT NULL,
    image_filename VARCHAR NULL,
    image_url VARCHAR NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE products;
-- +goose StatementEnd
