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

-- name: CreateTask :exec
insert into tasks(manager_id, text_template, image, status, created_at)
VALUES ($1, $2, $3, $4, now());

-- name: StartTaskByID :one
update tasks
set status     = 3,
    started_at = now()
where id = $1
  AND status = 2 --
returning *;

-- name: GetBotByID :one
select *
from bot_accounts
where id = $1;


-- name: FindAccountsForTask :many
select *
from bot_accounts
where task_id = $1;

-- name: UpdateTaskStatus :exec
update tasks
set status     = $1,
    updated_at = now()
where id = $1;

-- name: SetAccountAsCompleted :exec
update bot_accounts
set status     = 4,
    updated_at = now()
where id = $1;