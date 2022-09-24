// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package dbmodel

import (
	"time"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/headers"
)

type BotAccount struct {
	ID         uuid.UUID              `json:"id"`
	TaskID     uuid.UUID              `json:"task_id"`
	Username   string                 `json:"username"`
	Password   string                 `json:"password"`
	UserAgent  string                 `json:"user_agent"`
	DeviceData headers.DeviceSettings `json:"device_data"`
	Session    headers.Session        `json:"session"`
	Headers    headers.Base           `json:"headers"`
	ResProxy   *Proxy                 `json:"res_proxy"`
	WorkProxy  *Proxy                 `json:"work_proxy"`
	Status     int16                  `json:"status"`
	StartedAt  time.Time              `json:"started_at"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  *time.Time             `json:"updated_at"`
	DeletedAt  *time.Time             `json:"deleted_at"`
}

type Log struct {
	ID           uuid.UUID `json:"id"`
	BotID        uuid.UUID `json:"bot_id"`
	Request      string    `json:"request"`
	Response     string    `json:"response"`
	ResponseCode int32     `json:"response_code"`
	RequestTime  time.Time `json:"request_time"`
	ProxyUrl     *string   `json:"proxy_url"`
}

type Proxy struct {
	ID         uuid.UUID  `json:"id"`
	TaskID     uuid.UUID  `json:"task_id"`
	AssignedTo *uuid.UUID `json:"assigned_to"`
	Host       string     `json:"host"`
	Port       string     `json:"port"`
	Login      string     `json:"login"`
	Pass       string     `json:"pass"`
	Type       int16      `json:"type"`
}

type TargetUser struct {
	ID        uuid.UUID  `json:"id"`
	TaskID    uuid.UUID  `json:"task_id"`
	Username  string     `json:"username"`
	UserID    int64      `json:"user_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type TargetUsersToTask struct {
	TargetID   uuid.UUID  `json:"target_id"`
	TaskID     uuid.UUID  `json:"task_id"`
	NotifiedAt *time.Time `json:"notified_at"`
}

type Task struct {
	ID           uuid.UUID  `json:"id"`
	ManagerID    uuid.UUID  `json:"manager_id"`
	TextTemplate string     `json:"text_template"`
	Image        []byte     `json:"image"`
	Status       taskStatus `json:"status"`
	Title        string     `json:"title"`
	CreatedAt    time.Time  `json:"created_at"`
	StartedAt    *time.Time `json:"started_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}

type User struct {
	ID           uuid.UUID  `json:"id"`
	Login        string     `json:"login"`
	PasswordHash string     `json:"password_hash"`
	Role         int16      `json:"role"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}
