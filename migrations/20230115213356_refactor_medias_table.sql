-- +goose Up
-- +goose StatementBegin


ALTER TABLE target_users
    DROP COLUMN media_fk;

ALTER TABLE medias
    DROP COLUMN id,
    ADD COLUMN pk        bigint primary key not null,
    ADD COLUMN is_edited bool default false not null;


ALTER TABLE target_users
    ADD COLUMN media_fk bigint references medias;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE medias
    DROP COLUMN is_edited;
-- +goose StatementEnd
