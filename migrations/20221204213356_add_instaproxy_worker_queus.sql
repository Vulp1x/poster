-- +goose Up
-- +goose StatementBegin
CREATE TYPE pgqueue_status AS ENUM (
    'new',
    'must_retry',

    'no_attempts_left',

    'cancelled',
    'succeeded'
    );
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE pgqueue
(
    id               bigserial PRIMARY KEY,
    kind             smallint                                          NOT NULL,
    payload          bytea                                             NOT NULL,
    external_key     text,

    status           pgqueue_status           DEFAULT 'new'            NOT NULL,
    messages         text[]                   DEFAULT array []::text[] NOT NULL,

    attempts_left    smallint                                          NOT NULL,
    attempts_elapsed smallint                 DEFAULT 0                NOT NULL,
    delayed_till     timestamp with time zone                          NOT NULL,

    created_at       timestamp with time zone DEFAULT now()            NOT NULL,
    updated_at       timestamp with time zone DEFAULT now()            NOT NULL
) WITH (fillfactor = 80);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE UNIQUE INDEX pgqueue_idempotency_idx ON pgqueue (kind, external_key);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX pgqueue_open_tasks_idx ON pgqueue (kind, delayed_till)
    WHERE status IN ('new', 'must_retry');
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX pgqueue_broken_tasks_idx ON pgqueue (kind, created_at)
    WHERE status IN ('no_attempts_left');
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX pgqueue_terminal_tasks_idx ON pgqueue (kind, updated_at)
    WHERE status IN ('cancelled', 'succeeded');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE pgqueue_status;
DROP TABLE pgqueue;
-- +goose StatementEnd
