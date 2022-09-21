CREATE USER docker WITH PASSWORD 'dN5mYdDVKbuyq6ry';
GRANT ALL PRIVILEGES ON DATABASE insta_poster TO docker;

CREATE TABLE users
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

INSERT INTO public.users (id, login, password_hash, role, created_at, updated_at, deleted_at)
VALUES ('b79d3bff-54c0-4f4a-b165-fa832e787648', 'admin', '$2a$10$HUX5NwM8LddgNXyRGQlyp.z7s2uTMnAkqWm5Wnt/eewwU0D7Ey2ry',
        0, '2022-09-20 07:55:33.613893 +00:00', '2022-09-20 07:55:33.613893 +00:00', null);
