-- +goose Up
-- +goose StatementBegin

ALTER TABLE tasks
    ADD COLUMN type           smallint not null default 0,
    ADD COLUMN video_filename text;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE tasks
    DROP COLUMN type;

-- +goose StatementEnd
