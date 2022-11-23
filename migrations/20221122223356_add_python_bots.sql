-- +goose Up
-- +goose StatementBegin

CREATE TABLE python_bots
(
    session_id text  not null primary key,
    settings   jsonb not null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE python_bots;
-- +goose StatementEnd
