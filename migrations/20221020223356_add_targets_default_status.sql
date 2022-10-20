-- +goose Up
-- +goose StatementBegin

alter table target_users
    alter column status set default 1;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

alter table target_users
    alter column status drop default;

-- +goose StatementEnd
