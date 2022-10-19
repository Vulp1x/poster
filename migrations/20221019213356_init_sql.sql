-- +goose Up
-- +goose StatementBegin

ALTER TABLE tasks
    DROP COLUMN account_surnames,
    ADD COLUMN account_last_names text,
    ALTER COLUMN account_names drop not null,
    ALTER COLUMN account_profile_images drop not null,
    ALTER COLUMN account_urls drop not null;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

SELECT 'skipped down query';

-- +goose StatementEnd
