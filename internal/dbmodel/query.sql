-- name: SelectNow :exec
SELECT now();

-- name: GetUserByID :one
SELECT *
from users
WHERE id = $1;

-- name: FindByLogin :one
SELECT *
from users
WHERE login = $1
  AND deleted_at IS NULL;

-- name: UpdateUserPassword :exec
UPDATE users
set password_hash = $1
where id = $2;

-- name: CreateUser :one
INSERT INTO users (login, password_hash, role, created_at)
VALUES ($1, $2, $3, now())
RETURNING *;

-- name: DeleteUserByID :exec
UPDATE users
SET deleted_at= now()
where id = $1;

-- name: CreateDraftTask :one
insert into tasks(manager_id, text_template, title, image, status, created_at)
VALUES ($1, $2, $3, $4, 1, now())
RETURNING id;

-- name: FindTaskByID :one
select *
from tasks
where id = $1;

-- name: FindTasksByManagerID :many
select *
from tasks
where manager_id = $1;

-- name: StartTaskByID :exec
update tasks
set status     = 4,
    started_at = now()
where id = $1
  AND status = 3;

-- name: GetBotByID :one
select *
from bot_accounts
where id = $1;

-- name: FindBotsForTask :many
select *
from bot_accounts
where task_id = $1;

-- name: FindReadyBotsForTask :many
select *
from bot_accounts
where task_id = $1
  and status = 2;
-- ProxieAssignedBotStatus

-- name: FindProxiesForTask :many
select *
from proxies
where task_id = $1;

-- name: FindUnprocessedTargetsForTask :many
select *
from target_users
where task_id = $1
  AND status = 0
limit $2;

-- name: UpdateTaskStatus :exec
update tasks
set status     = $1,
    updated_at = now()
where id = $2;

-- name: UpdateTask :exec
update tasks
set text_template = $1,
    title         = $2,
    image         = $3,
    updated_at    = now()
where id = $4;

-- name: SaveBotAccounts :copyfrom
insert into bot_accounts (task_id, username, password, user_agent, device_data, session, headers, status)
values ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: SaveProxies :copyfrom
insert into proxies (task_id, host, port, login, pass, type)
values ($1, $2, $3, $4, $5, $6);

-- name: SaveTargetUsers :copyfrom
insert into target_users (task_id, username, user_id)
values ($1, $2, $3);

-- name: ForceDeleteBotAccountsForTask :execrows
DELETE
FROM bot_accounts
where task_id = $1;

-- name: ForceDeleteProxiesForTask :execrows
DELETE
FROM proxies
where task_id = $1;

-- name: ForceDeleteTargetUsersForTask :execrows
DELETE
FROM target_users
where task_id = $1;

-- name: ForceDeleteTaskByID :exec
DELETE
FROM tasks
where id = $1;

-- name: AssignProxiesToBotsForTask :exec
UPDATE bot_accounts
set res_proxy = x.proxy,
    status    = 2 -- ProxieAssignedBotStatus
From (SELECT UNNEST(sqlc.arg(proxies)::jsonb[]) as proxy,
             UNNEST(sqlc.arg(ids)::uuid[])      as id) x
where bot_accounts.id = x.id
  AND task_id = $1;

-- name: AssignBotsToProxiesForTask :exec
UPDATE proxies
set assigned_to = x.bot_id
From (SELECT UNNEST(sqlc.arg(bot_ids)::uuid[]) as bot_id,
             UNNEST(sqlc.arg(ids)::uuid[])     as id) x
where proxies.id = x.id
  AND task_id = $1;

-- name: DeleteProxiesForTask :execrows
DELETE
FROM proxies
where id in ($1::uuid[])
RETURNING 1;

-- name: DeleteBotAccountsForTask :execrows
DELETE
FROM bot_accounts
where id = ANY ($1::uuid[])
RETURNING 1;

-- name: SelectCountsForTask :one
select (select count(*) from proxies p where p.task_id = $1)      as proxies_count,
       (select count(*) from bot_accounts b where b.task_id = $1) as bots_count,
       (select count(*) from target_users t where t.task_id = $1) as targets_count;

-- name: SaveUploadedDataToTask :exec
update tasks
set status           = 2, --dbmodel.DataUploadedTaskStatus,
    bots_filename    = $2,
    proxies_filename = $3,
    targets_filename = $4
where id = $1;

-- name: SetBotStatus :exec
update bot_accounts
set status = $1
where id = $2;

-- name: SetTargetsStatus :exec
update target_users
set status = $1
where id = ANY (sqlc.arg('ids')::uuid[]);