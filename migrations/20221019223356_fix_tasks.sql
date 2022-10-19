-- +goose Up
-- +goose StatementBegin

ALTER TABLE tasks
    DROP COLUMN account_last_names,
    ADD COLUMN account_last_names text[];

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

SELECT 'skipped down query';

-- +goose StatementEnd
