-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    status VARCHAR(255) NOT NULL,
    product_ids UUID[] NOT NULL DEFAULT '{}'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd
