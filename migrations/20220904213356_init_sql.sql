-- +goose Up
-- +goose StatementBegin
CREATE
    EXTENSION IF NOT EXISTS pgcrypto;


CREATE TABLE IF NOT EXISTS users
(
    id            uuid primary key         default gen_random_uuid(),
    login         text UNIQUE                            not null,
    password_hash text                                   not null,
    role          smallint                               not null,
    constraint valid_role check (role IN (0, 1) ), -- 0 is a manager, 1 is an admin

    created_at    TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    deleted_at    TIMESTAMP WITH TIME ZONE
);

CREATE TABLE tasks
(
    id               uuid                     not null primary key default gen_random_uuid(),
    manager_id       uuid                     not null references users,
    text_template    text                     not null,
    image            bytea                    not null,
    status           smallint                 not null,
    title            text                     not null,
    bots_filename    text, -- название файла с ботами
    proxies_filename text, -- название файла с прокси
    targets_filename text, -- название файла с получателями
    created_at       timestamp with time zone not null,
    started_at       timestamp with time zone,
    updated_at       timestamp with time zone,
    deleted_at       timestamp with time zone
);

CREATE TABLE bot_accounts
(
    id          uuid primary key     default gen_random_uuid(),
    task_id     uuid        not null references tasks,
    username    text UNIQUE not null,
    password    text        not null,
    user_agent  text        not null,
    device_data jsonb       not null,
    session     jsonb       not null,
    headers     jsonb       not null,
    res_proxy   jsonb, -- резидентские прокси
    work_proxy  jsonb,
    status      smallint    not null,
    posts_count smallint    not null,
    started_at  timestamp,
    created_at  timestamp   not null default now(),
    updated_at  timestamp,
    deleted_at  timestamp,

    constraint not_empty_device check ( device_data <> '[]' AND device_data <> '{}' ),
    constraint not_empty_headers check ( headers <> '[]' AND headers <> '{}' ),
    constraint not_empty_session check ( session <> '[]' AND session <> '{}' )
);

-- таблица с пользователями, которым будет показана реклама
create table target_users
(
    id         uuid primary key   default gen_random_uuid(),
    task_id    uuid      not null references tasks,
    username   text      not null,
    user_id    bigint    not null,
    status     smallint  not null default 0, -- 0 - не показывали рекламу, 1 - пытались показать рекламу, но не получилось, 2-показали рекламу
    created_at timestamp not null default now(),
    updated_at timestamp
);

create unique index target_users_uniq_idx on target_users (task_id, username, user_id);

create table target_users_to_tasks
(
    target_id   uuid not null references target_users,
    task_id     uuid not null references tasks,
    notified_at timestamp
);


create table proxies
(
    id          uuid primary key default gen_random_uuid(),
    task_id     uuid references tasks not null,
    assigned_to uuid references bot_accounts,
    host        text                  not null,
    port        integer               not null,
    login       text                  not null,
    pass        text                  not null,
    type        smallint              not null,-- 1 is for residential, 2 for usual
    unique (host, port)
);


create table logs
(
    id            uuid primary key default gen_random_uuid(),
    bot_id        uuid      not null references bot_accounts,
    operation     text      not null,
    request       jsonb     not null,
    response      jsonb     not null,
    response_code integer   not null,
    request_time  timestamp not null,
    proxy_url     text
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE tasks;
DROP TABLE users;
DROP TABLE bot_accounts;
drop table logs;

DROP EXTENSION pgcrypto

-- +goose StatementEnd
