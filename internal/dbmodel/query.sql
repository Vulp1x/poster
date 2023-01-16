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
insert into tasks(manager_id, text_template, title, landing_accounts, images, account_names, account_last_names,
                  account_profile_images, account_urls, status, type, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 1, $10, now())
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

-- name: FindResidentialProxiesForTask :many
select *
from proxies
where task_id = $1
  and type = 1;

-- name: FindCheapProxiesForTask :many
select *
from proxies
where task_id = $1
  and type = 2;


-- name: FindUnprocessedTargetsForTask :many
select *
from target_users
where task_id = $1
  AND status = 'new'
limit $2;

-- name: UpdateTaskStatus :exec
update tasks
set status     = $1,
    updated_at = now()
where id = $2;

-- name: UpdateTask :one
update tasks
set text_template            = $1,
    title                    = $2,
    images                   = $3,
    account_names            = $4,
    account_last_names       = $5,
    account_urls             = $6,
    account_profile_images   = $7,
    landing_accounts         = $8,
    follow_targets           = $9,
    need_photo_tags          = $10,
    per_post_sleep_seconds   =$11,
    photo_tags_delay_seconds = $12,
    posts_per_bot            = $13,
    targets_per_post         = $14,
    photo_targets_per_post   = @photo_targets_per_post,
    photo_tags_posts_per_bot = @photo_tags_posts_per_bot,
    updated_at               = now()
where id = $15
returning *;

-- name: SaveBotAccounts :copyfrom
insert into bot_accounts (task_id, username, password, user_agent, device_data, session, headers, status, file_order,
                          inst_id)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

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
set res_proxy  = x.res_proxy,
    work_proxy = x.cheap_proxy,
    status     = 2 -- ProxieAssignedBotStatus
From (SELECT UNNEST(sqlc.arg(residential_proxies)::jsonb[]) as res_proxy,
             UNNEST(sqlc.arg(cheap_proxies)::jsonb[])       as cheap_proxy,
             UNNEST(sqlc.arg(ids)::uuid[])                  as id) x
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
where id = ANY ($1::uuid[])
RETURNING 1;

-- name: DeleteBotAccountsForTask :execrows
DELETE
FROM bot_accounts
where id = ANY ($1::uuid[])
RETURNING 1;

-- name: SelectCountsForTask :one
select (select count(*) from proxies p where p.task_id = $1 and p.type = 1) as residential_proxies_count,
       (select count(*) from proxies p where p.task_id = $1 and p.type = 2) as cheap_proxies_count,
       (select count(*) from bot_accounts b where b.task_id = $1)           as bots_count,
       (select count(*) from target_users t where t.task_id = $1)           as targets_count;

-- name: SaveUploadedDataToTask :exec
update tasks
set status                 = 2, --dbmodel.DataUploadedTaskStatus,
    bots_filename          = $2,
    res_proxies_filename   = $3,
    cheap_proxies_filename = $4,
    targets_filename       = $5
where id = $1;

-- name: SetBotStatus :exec
update bot_accounts
set status = $1
where id = $2;

-- name: SetBotPostsCount :exec
update bot_accounts
set posts_count = $1
where id = $2;

-- name: SetBotDoneStatus :exec
update bot_accounts
set status = 4 -- dbmodel.DoneBotStatus
where id = @id;


-- name: SetTargetsStatus :exec
update target_users
set status = $1
where id = ANY (sqlc.arg('ids')::uuid[]);

-- name: GetBotsProgress :many
select username, posts_count, status
from bot_accounts
where task_id = $1
order by CASE
             WHEN @posts_asc::bool THEN posts_count
             else
                 CASE WHEN @posts_desc::bool THEN posts_count end END
LIMIT $2 OFFSET $3;

-- name: GetTaskTargetsCount :one
select (select count(*) from target_users t where t.task_id = $1 and t.status = 'new')    as unused_targets,
       (select count(*) from target_users t where t.task_id = $1 and t.status = 'failed') as failed_targets,
       (select count(*)
        from target_users t
        where t.task_id = $1
          and t.status = 'notified'
          AND interaction_type = 'photo_tag')                                             as photo_notified_targets,
       (select count(*)
        from target_users t
        where t.task_id = $1
          and t.status = 'notified'
          AND interaction_type = 'post_description')                                      as description_notified_targets;

-- name: GetTaskBotsCount :one
select count(*)
from bot_accounts
where task_id = $1;

-- name: FindReadyBots :many
select *
from bot_accounts
where (res_proxy is not null or work_proxy is not null)
  and status >= 2;

-- name: SetTaskVideoFilename :exec
update tasks
set video_filename = @video_filename
where id = @id;

-- name: FindTaskBotsByUsername :many
select *
from bot_accounts
where task_id = @task_id
  and username = any (sqlc.arg('usernames')::text[])
  and status in (2, 5)
ORDER BY file_order;


-- name: FindTaskBotByUsername :one
select *
from bot_accounts
where task_id = @task_id
  and username = @username
  and status in (2, 5);

-- name: SavePostedMedia :one
insert into medias(kind, inst_id, bot_id, created_at)
VALUES (@kind, @inst_id, @bot_id, now())
returning *;

-- name: MarkTargetsAsNotified :exec
update target_users
set media_fk= @media_fk,
    status='notified',
    interaction_type= @interaction_type
where task_id = @task_id
  and user_id = ANY (sqlc.arg('target_ids')::bigint[]);

-- name: AddBotPost :one
update bot_accounts
set posts_count = 1 + @posts_count
where id = @id
returning posts_count;

-- name: CountBotMedias :one
select count(*)
from medias
where bot_id = @bot_id;