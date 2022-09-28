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

const createDraftTask = `-- name: CreateDraftTask :one
insert into tasks(manager_id, text_template, title, image, status, created_at)
VALUES ($1, $2, $3, $4, 1, now())
RETURNING id
`

type CreateDraftTaskParams struct {
	ManagerID    uuid.UUID `json:"manager_id"`
	TextTemplate string    `json:"text_template"`
	Title        string    `json:"title"`
	Image        []byte    `json:"image"`
}

func (q *Queries) CreateDraftTask(ctx context.Context, arg CreateDraftTaskParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createDraftTask,
		arg.ManagerID,
		arg.TextTemplate,
		arg.Title,
		arg.Image,
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

const deleteUserByID = `-- name: DeleteUserByID :exec
UPDATE users
SET deleted_at= now()
where id = $1
`

func (q *Queries) DeleteUserByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteUserByID, id)
	return err
}

const findAccountsForTask = `-- name: FindAccountsForTask :many
select id, task_id, username, password, user_agent, device_data, session, headers, res_proxy, work_proxy, status, started_at, created_at, updated_at, deleted_at
from bot_accounts
where task_id = $1
`

func (q *Queries) FindAccountsForTask(ctx context.Context, taskID uuid.UUID) ([]BotAccount, error) {
	rows, err := q.db.Query(ctx, findAccountsForTask, taskID)
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

const findTaskByID = `-- name: FindTaskByID :one
select id, manager_id, text_template, image, status, title, created_at, started_at, updated_at, deleted_at
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
		&i.Image,
		&i.Status,
		&i.Title,
		&i.CreatedAt,
		&i.StartedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getBotByID = `-- name: GetBotByID :one
select id, task_id, username, password, user_agent, device_data, session, headers, res_proxy, work_proxy, status, started_at, created_at, updated_at, deleted_at
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
		&i.StartedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
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
	Type   int16     `json:"type"`
}

type SaveTargetUsersParams struct {
	TaskID   uuid.UUID `json:"task_id"`
	Username string    `json:"username"`
	UserID   int64     `json:"user_id"`
}

const selectNow = `-- name: SelectNow :exec
SELECT now()
`

func (q *Queries) SelectNow(ctx context.Context) error {
	_, err := q.db.Exec(ctx, selectNow)
	return err
}

const setAccountAsCompleted = `-- name: SetAccountAsCompleted :exec
update bot_accounts
set status     = 4,
    updated_at = now()
where id = $1
`

func (q *Queries) SetAccountAsCompleted(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, setAccountAsCompleted, id)
	return err
}

const startTaskByID = `-- name: StartTaskByID :one
update tasks
set status     = 3,
    started_at = now()
where id = $1
  AND status = 2 --
returning id, manager_id, text_template, image, status, title, created_at, started_at, updated_at, deleted_at
`

func (q *Queries) StartTaskByID(ctx context.Context, id uuid.UUID) (Task, error) {
	row := q.db.QueryRow(ctx, startTaskByID, id)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.ManagerID,
		&i.TextTemplate,
		&i.Image,
		&i.Status,
		&i.Title,
		&i.CreatedAt,
		&i.StartedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
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
