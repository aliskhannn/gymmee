-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    telegram_id   INTEGER  NOT NULL UNIQUE,
    username      TEXT,
    height        REAL CHECK (height > 0 OR height IS NULL),
    weight        REAL CHECK (weight > 0 OR weight IS NULL),
    target_weight REAL CHECK (target_weight > 0 OR target_weight IS NULL),
    created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd