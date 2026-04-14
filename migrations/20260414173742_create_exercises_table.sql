-- +goose Up
-- +goose StatementBegin
CREATE TABLE exercises
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id      INTEGER,
    name         TEXT NOT NULL,
    muscle_group TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS exercises;
-- +goose StatementEnd