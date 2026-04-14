-- +goose Up
-- +goose StatementBegin
CREATE TABLE workout_sessions
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER  NOT NULL,
    plan_day_id INTEGER,
    started_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ended_at    DATETIME,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (plan_day_id) REFERENCES plan_days (id) ON DELETE SET NULL
);

CREATE TABLE workout_sets
(
    id                 INTEGER PRIMARY KEY AUTOINCREMENT,
    workout_session_id INTEGER  NOT NULL,
    exercise_id        INTEGER  NOT NULL,
    weight             REAL     NOT NULL CHECK (weight >= 0),
    reps               INTEGER  NOT NULL CHECK (reps > 0),
    created_at         DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workout_session_id) REFERENCES workout_sessions (id) ON DELETE CASCADE,
    FOREIGN KEY (exercise_id) REFERENCES exercises (id) ON DELETE RESTRICT
);

CREATE INDEX idx_workout_sets_history ON workout_sets (exercise_id, workout_session_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS workout_sets;
DROP TABLE IF EXISTS workout_sessions;
-- +goose StatementEnd