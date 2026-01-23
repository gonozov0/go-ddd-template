-- +goose Up
-- +goose StatementBegin
CREATE TOPIC `order_processing` (
    CONSUMER order_processing_consumer WITH (important = true)
) WITH (
    retention_period = Interval('P7D')

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TOPIC `order_processing`;
-- +goose StatementEnd
