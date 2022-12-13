-- +goose Up
-- +goose StatementBegin

ALTER TABLE bot_accounts
    ADD COLUMN file_order integer not null,
    ADD CONSTRAINT uniq_file_order UNIQUE (task_id, file_order);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE bot_accounts
    DROP COLUMN file_order;

-- +goose StatementEnd
