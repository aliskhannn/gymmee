-- +goose Up
-- +goose StatementBegin
CREATE TABLE plans
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id    INTEGER  NOT NULL,
    name       TEXT     NOT NULL,
    is_active  BOOLEAN  NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_plans_user_active ON plans (user_id, is_active);

CREATE TABLE plan_days
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    plan_id      INTEGER NOT NULL,
    day_of_week  INTEGER NOT NULL CHECK (day_of_week BETWEEN 1 AND 7),
    muscle_group TEXT    NOT NULL,
    FOREIGN KEY (plan_id) REFERENCES plans (id) ON DELETE CASCADE,
    UNIQUE (plan_id, day_of_week)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS plan_days;
DROP TABLE IF EXISTS plans;
-- +goose StatementEnd