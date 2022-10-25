-- +goose Up
-- +goose StatementBegin

ALTER TABLE tasks
    ADD COLUMN follow_targets           bool not null default false,
    ADD COLUMN need_photo_tags          bool not null default false,
    ADD COLUMN per_post_sleep_seconds   int4 not null default 0,
    ADD COLUMN photo_tags_delay_seconds int4 not null default 0,
    ADD COLUMN posts_per_bot            int4 not null default 0,
    ADD COLUMN targets_per_post         int4 not null default 0;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE tasks
    DROP COLUMN follow_targets,
    DROP COLUMN need_photo_tags,
    DROP COLUMN per_post_sleep_seconds,
    DROP COLUMN photo_tags_delay_seconds,
    DROP COLUMN posts_per_bot,
    DROP COLUMN targets_per_post;

-- +goose StatementEnd
