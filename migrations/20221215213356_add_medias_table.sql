-- +goose Up
-- +goose StatementBegin
CREATE TYPE targets_status AS ENUM (
    'new',
    'in_progress',
    'failed',
    'notified'
    );

CREATE TYPE targets_interaction AS ENUM (
    'none',
    'post_description',
    'photo_tag'
    );

CREATE TYPE medias_kind AS ENUM (
    'photo',
    'reels'
    );


-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE medias
(
    id         bigserial PRIMARY KEY,
    kind       medias_kind                            NOT NULL,
    inst_id    text                                   not null,
    bot_id     uuid references bot_accounts           not null,

    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);

ALTER TABLE target_users
    ADD COLUMN media_fk         bigint references medias,
    DROP COLUMN status,
    ADD COLUMN status           targets_status      not null default 'new',
    ADD COLUMN interaction_type targets_interaction not null default 'none';

create index targets_user_id_idx
    on target_users (task_id, user_id);

ALTER TABLE tasks
    DROP COLUMN posts_per_bot,
    ADD COLUMN posts_per_bot            smallint default 0 not null,
    DROP COLUMN targets_per_post,
    ADD COLUMN targets_per_post         smallint default 0 not null,
    ADD COLUMN photo_tags_posts_per_bot smallint default 0 not null,
    ADD COLUMN photo_targets_per_post   smallint default 0 not null;

ALTER TABLE bot_accounts
    ADD COLUMN inst_id bigint not null;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE targets_status, targets_interaction, medias_kind;
DROP TABLE medias;
-- +goose StatementEnd
