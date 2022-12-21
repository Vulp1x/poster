-- +goose Up
-- +goose StatementBegin
ALTER TABLE proxies
    DROP CONSTRAINT IF EXISTS proxies_host_port_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table proxies
    ADD CONSTRAINT proxies_host_port_key unique (host, port);
-- +goose StatementEnd
