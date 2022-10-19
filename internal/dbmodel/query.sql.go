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
set res_proxy = x.proxy,
    status    = 2 -- ProxieAssignedBotStatus
From (SELECT UNNEST($2::jsonb[]) as proxy,
             UNNEST($3::uuid[])      as id) x
where bot_accounts.id = x.id
  AND task_id = $1
`

type AssignProxiesToBotsForTaskParams struct {
	TaskID  uuid.UUID   `json:"task_id"`
	Proxies []string    `json:"proxies"`
	Ids     []uuid.UUID `json:"ids"`
}

func (q *Queries) AssignProxiesToBotsForTask(ctx context.Context, arg AssignProxiesToBotsForTaskParams) error {
	_, err := q.db.Exec(ctx, assignProxiesToBotsForTask, arg.TaskID, arg.Proxies, arg.Ids)
	return err
}

const createDraftTask = `-- name: CreateDraftTask :one
insert into tasks(manager_id, text_template, title, images, status, created_at)
VALUES ($1, $2, $3, $4, 1, now())
RETURNING id
`

type CreateDraftTaskParams struct {
	ManagerID    uuid.UUID `json:"manager_id"`
	TextTemplate string    `json:"text_template"`
	Title        string    `json:"title"`
	Images       [][]byte  `json:"images"`
}

func (q *Queries) CreateDraftTask(ctx context.Context, arg CreateDraftTaskParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createDraftTask,
		arg.ManagerID,
		arg.TextTemplate,
		arg.Title,
		arg.Images,
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
	Role         int16  `json:"role"`
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
where id in ($1::uuid[])
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
select id, task_id, username, password, user_agent, device_data, session, headers, res_proxy, work_proxy, status, posts_count, started_at, created_at, updated_at, deleted_at
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

const findProxiesForTask = `-- name: FindProxiesForTask :many

select id, task_id, assigned_to, host, port, login, pass, type
from proxies
where task_id = $1
`

// ProxieAssignedBotStatus
func (q *Queries) FindProxiesForTask(ctx context.Context, taskID uuid.UUID) ([]Proxy, error) {
	rows, err := q.db.Query(ctx, findProxiesForTask, taskID)
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

const findReadyBotsForTask = `-- name: FindReadyBotsForTask :many
select id, task_id, username, password, user_agent, device_data, session, headers, res_proxy, work_proxy, status, posts_count, started_at, created_at, updated_at, deleted_at
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
select id, manager_id, text_template, landing_accounts, account_profile_images, account_names, account_surnames, account_urls, images, status, title, bots_filename, cheap_proxies_filename, res_proxies_filename, targets_filename, created_at, started_at, stopped_at, updated_at, deleted_at
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
		&i.AccountSurnames,
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
	)
	return i, err
}

const findTasksByManagerID = `-- name: FindTasksByManagerID :many
select id, manager_id, text_template, landing_accounts, account_profile_images, account_names, account_surnames, account_urls, images, status, title, bots_filename, cheap_proxies_filename, res_proxies_filename, targets_filename, created_at, started_at, stopped_at, updated_at, deleted_at
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
			&i.AccountSurnames,
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
select id, task_id, username, user_id, status, created_at, updated_at
from target_users
where task_id = $1
  AND status = 0
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
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
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
select id, task_id, username, password, user_agent, device_data, session, headers, res_proxy, work_proxy, status, posts_count, started_at, created_at, updated_at, deleted_at
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
	)
	return i, err
}

const getTaskProgress = `-- name: GetTaskProgress :many
select username, posts_count, status
from bot_accounts
where task_id = $1
`

type GetTaskProgressRow struct {
	Username   string    `json:"username"`
	PostsCount int16     `json:"posts_count"`
	Status     botStatus `json:"status"`
}

func (q *Queries) GetTaskProgress(ctx context.Context, taskID uuid.UUID) ([]GetTaskProgressRow, error) {
	rows, err := q.db.Query(ctx, getTaskProgress, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTaskProgressRow
	for rows.Next() {
		var i GetTaskProgressRow
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

type SaveBotAccountsParams struct {
	TaskID     uuid.UUID              `json:"task_id"`
	Username   string                 `json:"username"`
	Password   string                 `json:"password"`
	UserAgent  string                 `json:"user_agent"`
	DeviceData headers.DeviceSettings `json:"device_data"`
	Session    headers.Session        `json:"session"`
	Headers    headers.Base           `json:"headers"`
	Status     botStatus              `json:"status"`
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
set status           = 2, --dbmodel.DataUploadedTaskStatus,
    bots_filename    = $2,
    res_proxies_filename = $3,
    cheap_proxies_filename = $4,
    targets_filename = $5
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
select (select count(*) from proxies p where p.task_id = $1)      as proxies_count,
       (select count(*) from bot_accounts b where b.task_id = $1) as bots_count,
       (select count(*) from target_users t where t.task_id = $1) as targets_count
`

type SelectCountsForTaskRow struct {
	ProxiesCount int64 `json:"proxies_count"`
	BotsCount    int64 `json:"bots_count"`
	TargetsCount int64 `json:"targets_count"`
}

func (q *Queries) SelectCountsForTask(ctx context.Context, taskID uuid.UUID) (SelectCountsForTaskRow, error) {
	row := q.db.QueryRow(ctx, selectCountsForTask, taskID)
	var i SelectCountsForTaskRow
	err := row.Scan(&i.ProxiesCount, &i.BotsCount, &i.TargetsCount)
	return i, err
}

const setBotPostsCount = `-- name: SetBotPostsCount :exec
update bot_accounts
set status      = 4, -- dbmodel.DoneBotStatus
    posts_count = $1
where id = $2
`

type SetBotPostsCountParams struct {
	PostsCount int16     `json:"posts_count"`
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
	Status targetStatus `json:"status"`
	Ids    []uuid.UUID  `json:"ids"`
}

func (q *Queries) SetTargetsStatus(ctx context.Context, arg SetTargetsStatusParams) error {
	_, err := q.db.Exec(ctx, setTargetsStatus, arg.Status, arg.Ids)
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

const updateTask = `-- name: UpdateTask :exec
update tasks
set text_template = $1,
    title         = $2,
    images         = $3,
    updated_at    = now()
where id = $4
`

type UpdateTaskParams struct {
	TextTemplate string    `json:"text_template"`
	Title        string    `json:"title"`
	Images       [][]byte  `json:"images"`
	ID           uuid.UUID `json:"id"`
}

func (q *Queries) UpdateTask(ctx context.Context, arg UpdateTaskParams) error {
	_, err := q.db.Exec(ctx, updateTask,
		arg.TextTemplate,
		arg.Title,
		arg.Images,
		arg.ID,
	)
	return err
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
