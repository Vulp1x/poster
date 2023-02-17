-- +goose Up
-- +goose StatementBegin


ALTER TABLE tasks
    ADD COLUMN fixed_tag       text,
    ADD COLUMN fixed_photo_tag bigint;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
