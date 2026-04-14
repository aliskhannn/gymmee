-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN barbell_weight REAL NOT NULL DEFAULT 20.0;
ALTER TABLE users
    ADD COLUMN available_plates TEXT NOT NULL DEFAULT '[25, 20, 15, 10, 5, 2.5, 1.25]';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN barbell_weight;
ALTER TABLE users
    DROP COLUMN available_plates;
-- +goose StatementEnd