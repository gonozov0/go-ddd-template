-- +goose Up
-- +goose StatementBegin

CREATE TABLE deliveries (
    id         UUID        PRIMARY KEY,
    order_id   UUID        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE deliveries;
-- +goose StatementEnd
