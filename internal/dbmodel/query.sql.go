// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: query.sql

package dbmodel

import (
	"context"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/headers"
)

const addBotPost = `-- name: AddBotPost :one
update bot_accounts
set posts_count = 1 + $1
where id = $2
returning posts_count
`

type AddBotPostParams struct {
	PostsCount interface{} `json:"posts_count"`
	ID         uuid.UUID   `json:"id"`
}

func (q *Queries) AddBotPost(ctx context.Context, arg AddBotPostParams) (int, error) {
	row := q.db.QueryRow(ctx, addBotPost, arg.PostsCount, arg.ID)
	var posts_count int
	err := row.Scan(&posts_count)
	return posts_count, err
}

const assignBotsToProxiesForTask = `-- name: AssignBotsToProxiesForTask :exec
UPDATE proxies
set assigned_to = x.bot_id
From (SELECT UNNEST($2::uuid[]) as bot_id,
             UNNEST($3::uuid[])     as id) x
where proxies.id = x.id
  AND task_id = $1
`

type AssignBotsToProxiesForTaskParams struct {
	TaskID uuid.UUID   `json:"task_id"`
	BotIds []uuid.UUID `json:"bot_ids"`
	Ids    []uuid.UUID `json:"ids"`
}

func (q *Queries) AssignBotsToProxiesForTask(ctx context.Context, arg AssignBotsToProxiesForTaskParams) error {
	_, err := q.db.Exec(ctx, assignBotsToProxiesForTask, arg.TaskID, arg.BotIds, arg.Ids)
	return err
}

const assignProxiesToBotsForTask = `-- name: AssignProxiesToBotsForTask :exec
UPDATE bot_accounts
set res_proxy  = x.res_proxy,
    work_proxy = x.cheap_proxy,
    status     = 2 -- ProxieAssignedBotStatus
From (SELECT UNNEST($2::jsonb[]) as res_proxy,
             UNNEST($3::jsonb[])       as cheap_proxy,
             UNNEST($4::uuid[])                  as id) x
where bot_accounts.id = x.id
  AND task_id = $1
`

type AssignProxiesToBotsForTaskParams struct {
	TaskID             uuid.UUID   `json:"task_id"`
	ResidentialProxies []string    `json:"residential_proxies"`
	CheapProxies       []string    `json:"cheap_proxies"`
	Ids                []uuid.UUID `json:"ids"`
}

func (q *Queries) AssignProxiesToBotsForTask(ctx context.Context, arg AssignProxiesToBotsForTaskParams) error {
	_, err := q.db.Exec(ctx, assignProxiesToBotsForTask,
		arg.TaskID,
		arg.ResidentialProxies,
		arg.CheapProxies,
		arg.Ids,
	)
	return err
}

const countBotMedias = `-- name: CountBotMedias :one
select count(*)
from medias
where bot_id = $1
`

func (q *Queries) CountBotMedias(ctx context.Context, botID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countBotMedias, botID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createDraftTask = `-- name: CreateDraftTask :one
insert into tasks(manager_id, text_template, title, landing_accounts, images, account_names, account_last_names,
                  account_profile_images, account_urls, status, type, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 1, $10, now())
RETURNING id
`

type CreateDraftTaskParams struct {
	ManagerID            uuid.UUID `json:"manager_id"`
	TextTemplate         string    `json:"text_template"`
	Title                string    `json:"title"`
	LandingAccounts      []string  `json:"landing_accounts"`
	Images               [][]byte  `json:"images"`
	AccountNames         []string  `json:"account_names"`
	AccountLastNames     []string  `json:"account_last_names"`
	AccountProfileImages [][]byte  `json:"account_profile_images"`
	AccountUrls          []string  `json:"account_urls"`
	Type                 taskType  `json:"type"`
}

func (q *Queries) CreateDraftTask(ctx context.Context, arg CreateDraftTaskParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createDraftTask,
		arg.ManagerID,
		arg.TextTemplate,
		arg.Title,
		arg.LandingAccounts,
		arg.Images,
		arg.AccountNames,
		arg.AccountLastNames,
		arg.AccountProfileImages,
		arg.AccountUrls,
		arg.Type,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (login, password_hash, role, created_at)
VALUES ($1, $2, $3, now())
RETURNING id, login, password_hash, role, created_at, updated_at, deleted_at
`

type CreateUserParams struct {
	Login        string `json:"login"`
	PasswordHash string `json:"password_hash"`
	Role         int    `json:"role"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Login, arg.PasswordHash, arg.Role)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Login,
		&i.PasswordHash,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const deleteBotAccountsForTask = `-- name: DeleteBotAccountsForTask :execrows
DELETE
FROM bot_accounts
where id = ANY ($1::uuid[])
RETURNING 1
`

func (q *Queries) DeleteBotAccountsForTask(ctx context.Context, dollar_1 []uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deleteBotAccountsForTask, dollar_1)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const deleteProxiesForTask = `-- name: DeleteProxiesForTask :execrows
DELETE
FROM proxies
where id = ANY ($1::uuid[])
RETURNING 1
`

func (q *Queries) DeleteProxiesForTask(ctx context.Context, dollar_1 []uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deleteProxiesForTask, dollar_1)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const deleteUserByID = `-- name: DeleteUserByID :exec
UPDATE users
SET deleted_at= now()
where id = $1
`

func (q *Queries) DeleteUserByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteUserByID, id)
	return err
}

const findBotsForTask = `-- name: FindBotsForTask :many
select id, task_id, username, password, user_agent, device_data, session, headers, res_proxy, work_proxy, status, posts_count, started_at, created_at, updated_at, deleted_at, file_order, inst_id
from bot_accounts
where task_id = $1
`

func (q *Queries) FindBotsForTask(ctx context.Context, taskID uuid.UUID) ([]BotAccount, error) {
	rows, err := q.db.Query(ctx, findBotsForTask, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []BotAccount
	for rows.Next() {
		var i BotAccount
		if err := rows.Scan(
			&i.ID,
			&i.TaskID,
			&i.Username,
			&i.Password,
			&i.UserAgent,
			&i.DeviceData,
			&i.Session,
			&i.Headers,
			&i.ResProxy,
			&i.WorkProxy,
			&i.Status,
			&i.PostsCount,
			&i.StartedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.FileOrder,
			&i.InstID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findByLogin = `-- name: FindByLogin :one
SELECT id, login, password_hash, role, created_at, updated_at, deleted_at
from users
WHERE login = $1
  AND deleted_at IS NULL
`

func (q *Queries) FindByLogin(ctx context.Context, login string) (User, error) {
	row := q.db.QueryRow(ctx, findByLogin, login)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Login,
		&i.PasswordHash,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const findCheapProxiesForTask = `-- name: FindCheapProxiesForTask :many
select id, task_id, assigned_to, host, port, login, pass, type
from proxies
where task_id = $1
  and type = 2
`

func (q *Queries) FindCheapProxiesForTask(ctx context.Context, taskID uuid.UUID) ([]Proxy, error) {
	rows, err := q.db.Query(ctx, findCheapProxiesForTask, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Proxy
	for rows.Next() {
		var i Proxy
		if err := rows.Scan(
			&i.ID,
			&i.TaskID,
			&i.AssignedTo,
			&i.Host,
			&i.Port,
			&i.Login,
			&i.Pass,
			&i.Type,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findReadyBots = `-- name: FindReadyBots :many


select id, task_id, username, password, user_agent, device_data, session, headers, res_proxy, work_proxy, status, posts_count, started_at, created_at, updated_at, deleted_at, file_order, inst_id
from bot_accounts
where (res_proxy is not null or work_proxy is not null)
  and status >= 2
`

// -- name: GetTaskTargetsCount :many
// select status, count(*)
// from target_users
// where task_id = $1
// group by status;
func (q *Queries) FindReadyBots(ctx context.Context) ([]BotAccount, error) {
	rows, err := q.db.Query(ctx, findReadyBots)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []BotAccount
	for rows.Next() {
		var i BotAccount
		if err := rows.Scan(
			&i.ID,
			&i.TaskID,
			&i.Username,
			&i.Password,
			&i.UserAgent,
			&i.DeviceData,
			&i.Session,
			&i.Headers,
			&i.ResProxy,
			&i.WorkProxy,
			&i.Status,
			&i.PostsCount,
			&i.StartedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.FileOrder,
			&i.InstID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findReadyBotsForTask = `-- name: FindReadyBotsForTask :many
select id, task_id, username, password, user_agent, device_data, session, headers, res_proxy, work_proxy, status, posts_count, started_at, created_at, updated_at, deleted_at, file_order, inst_id
from bot_accounts
where task_id = $1
  and status = 2
`

func (q *Queries) FindReadyBotsForTask(ctx context.Context, taskID uuid.UUID) ([]BotAccount, error) {
	rows, err := q.db.Query(ctx, findReadyBotsForTask, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []BotAccount
	for rows.Next() {
		var i BotAccount
		if err := rows.Scan(
			&i.ID,
			&i.TaskID,
			&i.Username,
			&i.Password,
			&i.UserAgent,
			&i.DeviceData,
			&i.Session,
			&i.Headers,
			&i.ResProxy,
			&i.WorkProxy,
			&i.Status,
			&i.PostsCount,
			&i.StartedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.FileOrder,
			&i.InstID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findResidentialProxiesForTask = `-- name: FindResidentialProxiesForTask :many

select id, task_id, assigned_to, host, port, login, pass, type
from proxies
where task_id = $1
  and type = 1
`

// ProxieAssignedBotStatus
func (q *Queries) FindResidentialProxiesForTask(ctx context.Context, taskID uuid.UUID) ([]Proxy, error) {
	rows, err := q.db.Query(ctx, findResidentialProxiesForTask, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Proxy
	for rows.Next() {
		var i Proxy
		if err := rows.Scan(
			&i.ID,
			&i.TaskID,
			&i.AssignedTo,
			&i.Host,
			&i.Port,
			&i.Login,
			&i.Pass,
			&i.Type,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findTaskBotByUsername = `-- name: FindTaskBotByUsername :one
select id, task_id, username, password, user_agent, device_data, session, headers, res_proxy, work_proxy, status, posts_count, started_at, created_at, updated_at, deleted_at, file_order, inst_id
from bot_accounts
where task_id = $1
  and username = $2
  and status in (2, 5)
`

type FindTaskBotByUsernameParams struct {
	TaskID   uuid.UUID `json:"task_id"`
	Username string    `json:"username"`
}

func (q *Queries) FindTaskBotByUsername(ctx context.Context, arg FindTaskBotByUsernameParams) (BotAccount, error) {
	row := q.db.QueryRow(ctx, findTaskBotByUsername, arg.TaskID, arg.Username)
	var i BotAccount
	err := row.Scan(
		&i.ID,
		&i.TaskID,
		&i.Username,
		&i.Password,
		&i.UserAgent,
		&i.DeviceData,
		&i.Session,
		&i.Headers,
		&i.ResProxy,
		&i.WorkProxy,
		&i.Status,
		&i.PostsCount,
		&i.StartedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.FileOrder,
		&i.InstID,
	)
	return i, err
}

const findTaskBotsByUsername = `-- name: FindTaskBotsByUsername :many
select id, task_id, username, password, user_agent, device_data, session, headers, res_proxy, work_proxy, status, posts_count, started_at, created_at, updated_at, deleted_at, file_order, inst_id
from bot_accounts
where task_id = $1
  and username = any ($2::text[])
  and status in (2, 5)
ORDER BY file_order
`

type FindTaskBotsByUsernameParams struct {
	TaskID    uuid.UUID `json:"task_id"`
	Usernames []string  `json:"usernames"`
}

func (q *Queries) FindTaskBotsByUsername(ctx context.Context, arg FindTaskBotsByUsernameParams) ([]BotAccount, error) {
	rows, err := q.db.Query(ctx, findTaskBotsByUsername, arg.TaskID, arg.Usernames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []BotAccount
	for rows.Next() {
		var i BotAccount
		if err := rows.Scan(
			&i.ID,
			&i.TaskID,
			&i.Username,
			&i.Password,
			&i.UserAgent,
			&i.DeviceData,
			&i.Session,
			&i.Headers,
			&i.ResProxy,
			&i.WorkProxy,
			&i.Status,
			&i.PostsCount,
			&i.StartedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.FileOrder,
			&i.InstID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findTaskByID = `-- name: FindTaskByID :one
select id, manager_id, text_template, landing_accounts, account_profile_images, account_names, account_urls, images, status, title, bots_filename, cheap_proxies_filename, res_proxies_filename, targets_filename, created_at, started_at, stopped_at, updated_at, deleted_at, account_last_names, follow_targets, need_photo_tags, per_post_sleep_seconds, photo_tags_delay_seconds, type, video_filename, posts_per_bot, targets_per_post, photo_tags_posts_per_bot, photo_targets_per_post
from tasks
where id = $1
`

func (q *Queries) FindTaskByID(ctx context.Context, id uuid.UUID) (Task, error) {
	row := q.db.QueryRow(ctx, findTaskByID, id)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.ManagerID,
		&i.TextTemplate,
		&i.LandingAccounts,
		&i.AccountProfileImages,
		&i.AccountNames,
		&i.AccountUrls,
		&i.Images,
		&i.Status,
		&i.Title,
		&i.BotsFilename,
		&i.CheapProxiesFilename,
		&i.ResProxiesFilename,
		&i.TargetsFilename,
		&i.CreatedAt,
		&i.StartedAt,
		&i.StoppedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.AccountLastNames,
		&i.FollowTargets,
		&i.NeedPhotoTags,
		&i.PerPostSleepSeconds,
		&i.PhotoTagsDelaySeconds,
		&i.Type,
		&i.VideoFilename,
		&i.PostsPerBot,
		&i.TargetsPerPost,
		&i.PhotoTagsPostsPerBot,
		&i.PhotoTargetsPerPost,
	)
	return i, err
}

const findTasksByManagerID = `-- name: FindTasksByManagerID :many
select id, manager_id, text_template, landing_accounts, account_profile_images, account_names, account_urls, images, status, title, bots_filename, cheap_proxies_filename, res_proxies_filename, targets_filename, created_at, started_at, stopped_at, updated_at, deleted_at, account_last_names, follow_targets, need_photo_tags, per_post_sleep_seconds, photo_tags_delay_seconds, type, video_filename, posts_per_bot, targets_per_post, photo_tags_posts_per_bot, photo_targets_per_post
from tasks
where manager_id = $1
`

func (q *Queries) FindTasksByManagerID(ctx context.Context, managerID uuid.UUID) ([]Task, error) {
	rows, err := q.db.Query(ctx, findTasksByManagerID, managerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.ManagerID,
			&i.TextTemplate,
			&i.LandingAccounts,
			&i.AccountProfileImages,
			&i.AccountNames,
			&i.AccountUrls,
			&i.Images,
			&i.Status,
			&i.Title,
			&i.BotsFilename,
			&i.CheapProxiesFilename,
			&i.ResProxiesFilename,
			&i.TargetsFilename,
			&i.CreatedAt,
			&i.StartedAt,
			&i.StoppedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.AccountLastNames,
			&i.FollowTargets,
			&i.NeedPhotoTags,
			&i.PerPostSleepSeconds,
			&i.PhotoTagsDelaySeconds,
			&i.Type,
			&i.VideoFilename,
			&i.PostsPerBot,
			&i.TargetsPerPost,
			&i.PhotoTagsPostsPerBot,
			&i.PhotoTargetsPerPost,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findUnprocessedTargetsForTask = `-- name: FindUnprocessedTargetsForTask :many
select id, task_id, username, user_id, created_at, updated_at, media_fk, status, interaction_type
from target_users
where task_id = $1
  AND status = 'new'
limit $2
`

type FindUnprocessedTargetsForTaskParams struct {
	TaskID uuid.UUID `json:"task_id"`
	Limit  int32     `json:"limit"`
}

func (q *Queries) FindUnprocessedTargetsForTask(ctx context.Context, arg FindUnprocessedTargetsForTaskParams) ([]TargetUser, error) {
	rows, err := q.db.Query(ctx, findUnprocessedTargetsForTask, arg.TaskID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TargetUser
	for rows.Next() {
		var i TargetUser
		if err := rows.Scan(
			&i.ID,
			&i.TaskID,
			&i.Username,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.MediaFk,
			&i.Status,
			&i.InteractionType,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const forceDeleteBotAccountsForTask = `-- name: ForceDeleteBotAccountsForTask :execrows
DELETE
FROM bot_accounts
where task_id = $1
`

func (q *Queries) ForceDeleteBotAccountsForTask(ctx context.Context, taskID uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, forceDeleteBotAccountsForTask, taskID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const forceDeleteProxiesForTask = `-- name: ForceDeleteProxiesForTask :execrows
DELETE
FROM proxies
where task_id = $1
`

func (q *Queries) ForceDeleteProxiesForTask(ctx context.Context, taskID uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, forceDeleteProxiesForTask, taskID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const forceDeleteTargetUsersForTask = `-- name: ForceDeleteTargetUsersForTask :execrows
DELETE
FROM target_users
where task_id = $1
`

func (q *Queries) ForceDeleteTargetUsersForTask(ctx context.Context, taskID uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, forceDeleteTargetUsersForTask, taskID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const forceDeleteTaskByID = `-- name: ForceDeleteTaskByID :exec
DELETE
FROM tasks
where id = $1
`

func (q *Queries) ForceDeleteTaskByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, forceDeleteTaskByID, id)
	return err
}

const getBotByID = `-- name: GetBotByID :one
select id, task_id, username, password, user_agent, device_data, session, headers, res_proxy, work_proxy, status, posts_count, started_at, created_at, updated_at, deleted_at, file_order, inst_id
from bot_accounts
where id = $1
`

func (q *Queries) GetBotByID(ctx context.Context, id uuid.UUID) (BotAccount, error) {
	row := q.db.QueryRow(ctx, getBotByID, id)
	var i BotAccount
	err := row.Scan(
		&i.ID,
		&i.TaskID,
		&i.Username,
		&i.Password,
		&i.UserAgent,
		&i.DeviceData,
		&i.Session,
		&i.Headers,
		&i.ResProxy,
		&i.WorkProxy,
		&i.Status,
		&i.PostsCount,
		&i.StartedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.FileOrder,
		&i.InstID,
	)
	return i, err
}

const getBotsProgress = `-- name: GetBotsProgress :many
select username, posts_count, status
from bot_accounts
where task_id = $1
order by CASE
             WHEN $4::bool THEN posts_count
             else
                 CASE WHEN $5::bool THEN posts_count end END
LIMIT $2 OFFSET $3
`

type GetBotsProgressParams struct {
	TaskID    uuid.UUID `json:"task_id"`
	Limit     int32     `json:"limit"`
	Offset    int32     `json:"offset"`
	PostsAsc  bool      `json:"posts_asc"`
	PostsDesc bool      `json:"posts_desc"`
}

type GetBotsProgressRow struct {
	Username   string    `json:"username"`
	PostsCount int       `json:"posts_count"`
	Status     botStatus `json:"status"`
}

func (q *Queries) GetBotsProgress(ctx context.Context, arg GetBotsProgressParams) ([]GetBotsProgressRow, error) {
	rows, err := q.db.Query(ctx, getBotsProgress,
		arg.TaskID,
		arg.Limit,
		arg.Offset,
		arg.PostsAsc,
		arg.PostsDesc,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBotsProgressRow
	for rows.Next() {
		var i GetBotsProgressRow
		if err := rows.Scan(&i.Username, &i.PostsCount, &i.Status); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTaskTargetsCount = `-- name: GetTaskTargetsCount :one
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
          AND interaction_type = 'post_description')                                      as description_notified_targets
`

type GetTaskTargetsCountRow struct {
	UnusedTargets              int64 `json:"unused_targets"`
	FailedTargets              int64 `json:"failed_targets"`
	PhotoNotifiedTargets       int64 `json:"photo_notified_targets"`
	DescriptionNotifiedTargets int64 `json:"description_notified_targets"`
}

func (q *Queries) GetTaskTargetsCount(ctx context.Context, taskID uuid.UUID) (GetTaskTargetsCountRow, error) {
	row := q.db.QueryRow(ctx, getTaskTargetsCount, taskID)
	var i GetTaskTargetsCountRow
	err := row.Scan(
		&i.UnusedTargets,
		&i.FailedTargets,
		&i.PhotoNotifiedTargets,
		&i.DescriptionNotifiedTargets,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, login, password_hash, role, created_at, updated_at, deleted_at
from users
WHERE id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Login,
		&i.PasswordHash,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const markTargetsAsNotified = `-- name: MarkTargetsAsNotified :exec
update target_users
set media_fk= $1,
    status='notified',
    interaction_type= $2
where task_id = $3
  and user_id = ANY ($4::bigint[])
`

type MarkTargetsAsNotifiedParams struct {
	MediaFk         *int64             `json:"media_fk"`
	InteractionType TargetsInteraction `json:"interaction_type"`
	TaskID          uuid.UUID          `json:"task_id"`
	TargetIds       []int64            `json:"target_ids"`
}

func (q *Queries) MarkTargetsAsNotified(ctx context.Context, arg MarkTargetsAsNotifiedParams) error {
	_, err := q.db.Exec(ctx, markTargetsAsNotified,
		arg.MediaFk,
		arg.InteractionType,
		arg.TaskID,
		arg.TargetIds,
	)
	return err
}

type SaveBotAccountsParams struct {
	TaskID     uuid.UUID              `json:"task_id"`
	Username   string                 `json:"username"`
	Password   string                 `json:"password"`
	UserAgent  string                 `json:"user_agent"`
	DeviceData headers.DeviceSettings `json:"device_data"`
	Session    headers.Session        `json:"session"`
	Headers    headers.Base           `json:"headers"`
	Status     botStatus              `json:"status"`
	FileOrder  int32                  `json:"file_order"`
	InstID     int64                  `json:"inst_id"`
}

const savePostedMedia = `-- name: SavePostedMedia :one
insert into medias(kind, inst_id, bot_id, created_at)
VALUES ($1, $2, $3, now())
returning id, kind, inst_id, bot_id, created_at, updated_at
`

type SavePostedMediaParams struct {
	Kind   MediasKind `json:"kind"`
	InstID string     `json:"inst_id"`
	BotID  uuid.UUID  `json:"bot_id"`
}

func (q *Queries) SavePostedMedia(ctx context.Context, arg SavePostedMediaParams) (Media, error) {
	row := q.db.QueryRow(ctx, savePostedMedia, arg.Kind, arg.InstID, arg.BotID)
	var i Media
	err := row.Scan(
		&i.ID,
		&i.Kind,
		&i.InstID,
		&i.BotID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

type SaveProxiesParams struct {
	TaskID uuid.UUID `json:"task_id"`
	Host   string    `json:"host"`
	Port   int32     `json:"port"`
	Login  string    `json:"login"`
	Pass   string    `json:"pass"`
	Type   proxyType `json:"type"`
}

type SaveTargetUsersParams struct {
	TaskID   uuid.UUID `json:"task_id"`
	Username string    `json:"username"`
	UserID   int64     `json:"user_id"`
}

const saveUploadedDataToTask = `-- name: SaveUploadedDataToTask :exec
update tasks
set status                 = 2, --dbmodel.DataUploadedTaskStatus,
    bots_filename          = $2,
    res_proxies_filename   = $3,
    cheap_proxies_filename = $4,
    targets_filename       = $5
where id = $1
`

type SaveUploadedDataToTaskParams struct {
	ID                   uuid.UUID `json:"id"`
	BotsFilename         *string   `json:"bots_filename"`
	ResProxiesFilename   *string   `json:"res_proxies_filename"`
	CheapProxiesFilename *string   `json:"cheap_proxies_filename"`
	TargetsFilename      *string   `json:"targets_filename"`
}

func (q *Queries) SaveUploadedDataToTask(ctx context.Context, arg SaveUploadedDataToTaskParams) error {
	_, err := q.db.Exec(ctx, saveUploadedDataToTask,
		arg.ID,
		arg.BotsFilename,
		arg.ResProxiesFilename,
		arg.CheapProxiesFilename,
		arg.TargetsFilename,
	)
	return err
}

const selectCountsForTask = `-- name: SelectCountsForTask :one
select (select count(*) from proxies p where p.task_id = $1 and p.type = 1) as residential_proxies_count,
       (select count(*) from proxies p where p.task_id = $1 and p.type = 2) as cheap_proxies_count,
       (select count(*) from bot_accounts b where b.task_id = $1)           as bots_count,
       (select count(*) from target_users t where t.task_id = $1)           as targets_count
`

type SelectCountsForTaskRow struct {
	ResidentialProxiesCount int64 `json:"residential_proxies_count"`
	CheapProxiesCount       int64 `json:"cheap_proxies_count"`
	BotsCount               int64 `json:"bots_count"`
	TargetsCount            int64 `json:"targets_count"`
}

func (q *Queries) SelectCountsForTask(ctx context.Context, taskID uuid.UUID) (SelectCountsForTaskRow, error) {
	row := q.db.QueryRow(ctx, selectCountsForTask, taskID)
	var i SelectCountsForTaskRow
	err := row.Scan(
		&i.ResidentialProxiesCount,
		&i.CheapProxiesCount,
		&i.BotsCount,
		&i.TargetsCount,
	)
	return i, err
}

const setBotDoneStatus = `-- name: SetBotDoneStatus :exec
update bot_accounts
set status = 4 -- dbmodel.DoneBotStatus
where id = $1
`

func (q *Queries) SetBotDoneStatus(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, setBotDoneStatus, id)
	return err
}

const setBotPostsCount = `-- name: SetBotPostsCount :exec
update bot_accounts
set posts_count = $1
where id = $2
`

type SetBotPostsCountParams struct {
	PostsCount int       `json:"posts_count"`
	ID         uuid.UUID `json:"id"`
}

func (q *Queries) SetBotPostsCount(ctx context.Context, arg SetBotPostsCountParams) error {
	_, err := q.db.Exec(ctx, setBotPostsCount, arg.PostsCount, arg.ID)
	return err
}

const setBotStatus = `-- name: SetBotStatus :exec
update bot_accounts
set status = $1
where id = $2
`

type SetBotStatusParams struct {
	Status botStatus `json:"status"`
	ID     uuid.UUID `json:"id"`
}

func (q *Queries) SetBotStatus(ctx context.Context, arg SetBotStatusParams) error {
	_, err := q.db.Exec(ctx, setBotStatus, arg.Status, arg.ID)
	return err
}

const setTargetsStatus = `-- name: SetTargetsStatus :exec
update target_users
set status = $1
where id = ANY ($2::uuid[])
`

type SetTargetsStatusParams struct {
	Status TargetsStatus `json:"status"`
	Ids    []uuid.UUID   `json:"ids"`
}

func (q *Queries) SetTargetsStatus(ctx context.Context, arg SetTargetsStatusParams) error {
	_, err := q.db.Exec(ctx, setTargetsStatus, arg.Status, arg.Ids)
	return err
}

const setTaskVideoFilename = `-- name: SetTaskVideoFilename :exec
update tasks
set video_filename = $1
where id = $2
`

type SetTaskVideoFilenameParams struct {
	VideoFilename *string   `json:"video_filename"`
	ID            uuid.UUID `json:"id"`
}

func (q *Queries) SetTaskVideoFilename(ctx context.Context, arg SetTaskVideoFilenameParams) error {
	_, err := q.db.Exec(ctx, setTaskVideoFilename, arg.VideoFilename, arg.ID)
	return err
}

const startTaskByID = `-- name: StartTaskByID :exec
update tasks
set status     = 4,
    started_at = now()
where id = $1
  AND status = 3
`

func (q *Queries) StartTaskByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, startTaskByID, id)
	return err
}

const updateTask = `-- name: UpdateTask :one
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
    photo_targets_per_post   = $16,
    photo_tags_posts_per_bot = $17,
    updated_at               = now()
where id = $15
returning id, manager_id, text_template, landing_accounts, account_profile_images, account_names, account_urls, images, status, title, bots_filename, cheap_proxies_filename, res_proxies_filename, targets_filename, created_at, started_at, stopped_at, updated_at, deleted_at, account_last_names, follow_targets, need_photo_tags, per_post_sleep_seconds, photo_tags_delay_seconds, type, video_filename, posts_per_bot, targets_per_post, photo_tags_posts_per_bot, photo_targets_per_post
`

type UpdateTaskParams struct {
	TextTemplate          string    `json:"text_template"`
	Title                 string    `json:"title"`
	Images                [][]byte  `json:"images"`
	AccountNames          []string  `json:"account_names"`
	AccountLastNames      []string  `json:"account_last_names"`
	AccountUrls           []string  `json:"account_urls"`
	AccountProfileImages  [][]byte  `json:"account_profile_images"`
	LandingAccounts       []string  `json:"landing_accounts"`
	FollowTargets         bool      `json:"follow_targets"`
	NeedPhotoTags         bool      `json:"need_photo_tags"`
	PerPostSleepSeconds   int32     `json:"per_post_sleep_seconds"`
	PhotoTagsDelaySeconds int32     `json:"photo_tags_delay_seconds"`
	PostsPerBot           int       `json:"posts_per_bot"`
	TargetsPerPost        int       `json:"targets_per_post"`
	ID                    uuid.UUID `json:"id"`
	PhotoTargetsPerPost   int       `json:"photo_targets_per_post"`
	PhotoTagsPostsPerBot  int       `json:"photo_tags_posts_per_bot"`
}

func (q *Queries) UpdateTask(ctx context.Context, arg UpdateTaskParams) (Task, error) {
	row := q.db.QueryRow(ctx, updateTask,
		arg.TextTemplate,
		arg.Title,
		arg.Images,
		arg.AccountNames,
		arg.AccountLastNames,
		arg.AccountUrls,
		arg.AccountProfileImages,
		arg.LandingAccounts,
		arg.FollowTargets,
		arg.NeedPhotoTags,
		arg.PerPostSleepSeconds,
		arg.PhotoTagsDelaySeconds,
		arg.PostsPerBot,
		arg.TargetsPerPost,
		arg.ID,
		arg.PhotoTargetsPerPost,
		arg.PhotoTagsPostsPerBot,
	)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.ManagerID,
		&i.TextTemplate,
		&i.LandingAccounts,
		&i.AccountProfileImages,
		&i.AccountNames,
		&i.AccountUrls,
		&i.Images,
		&i.Status,
		&i.Title,
		&i.BotsFilename,
		&i.CheapProxiesFilename,
		&i.ResProxiesFilename,
		&i.TargetsFilename,
		&i.CreatedAt,
		&i.StartedAt,
		&i.StoppedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.AccountLastNames,
		&i.FollowTargets,
		&i.NeedPhotoTags,
		&i.PerPostSleepSeconds,
		&i.PhotoTagsDelaySeconds,
		&i.Type,
		&i.VideoFilename,
		&i.PostsPerBot,
		&i.TargetsPerPost,
		&i.PhotoTagsPostsPerBot,
		&i.PhotoTargetsPerPost,
	)
	return i, err
}

const updateTaskStatus = `-- name: UpdateTaskStatus :exec
update tasks
set status     = $1,
    updated_at = now()
where id = $2
`

type UpdateTaskStatusParams struct {
	Status taskStatus `json:"status"`
	ID     uuid.UUID  `json:"id"`
}

func (q *Queries) UpdateTaskStatus(ctx context.Context, arg UpdateTaskStatusParams) error {
	_, err := q.db.Exec(ctx, updateTaskStatus, arg.Status, arg.ID)
	return err
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE users
set password_hash = $1
where id = $2
`

type UpdateUserPasswordParams struct {
	PasswordHash string    `json:"password_hash"`
	ID           uuid.UUID `json:"id"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.Exec(ctx, updateUserPassword, arg.PasswordHash, arg.ID)
	return err
}
