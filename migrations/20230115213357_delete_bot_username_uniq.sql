-- +goose Up
-- +goose StatementBegin


ALTER TABLE bot_accounts
    DROP CONSTRAINT bot_accounts_username_key;

CREATE UNIQUE INDEX bot_accounts_uniq_username_taskid_idx ON bot_accounts (task_id, username);

create index bots_task_id_idx on bot_accounts (task_id);
create index medias_bot_id_idx on medias (bot_id);
create index target_users_media_fk_idx on target_users(media_fk);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
