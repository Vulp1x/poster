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
	Status     botStatus              `json:"status"`
	PostsCount int16                  `json:"posts_count"`
	StartedAt  *time.Time             `json:"started_at"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  *time.Time             `json:"updated_at"`
	DeletedAt  *time.Time             `json:"deleted_at"`
}

type Log struct {
	ID           uuid.UUID `json:"id"`
	BotID        uuid.UUID `json:"bot_id"`
	Operation    string    `json:"operation"`
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
	Port       int32      `json:"port"`
	Login      string     `json:"login"`
	Pass       string     `json:"pass"`
	Type       proxyType  `json:"type"`
}

type PythonBot struct {
	SessionID string `json:"session_id"`
	Settings  string `json:"settings"`
}

type TargetUser struct {
	ID        uuid.UUID    `json:"id"`
	TaskID    uuid.UUID    `json:"task_id"`
	Username  string       `json:"username"`
	UserID    int64        `json:"user_id"`
	Status    targetStatus `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt *time.Time   `json:"updated_at"`
}

type TargetUsersToTask struct {
	TargetID   uuid.UUID  `json:"target_id"`
	TaskID     uuid.UUID  `json:"task_id"`
	NotifiedAt *time.Time `json:"notified_at"`
}

type Task struct {
	ID                    uuid.UUID  `json:"id"`
	ManagerID             uuid.UUID  `json:"manager_id"`
	TextTemplate          string     `json:"text_template"`
	LandingAccounts       []string   `json:"landing_accounts"`
	AccountProfileImages  [][]byte   `json:"account_profile_images"`
	AccountNames          []string   `json:"account_names"`
	AccountUrls           []string   `json:"account_urls"`
	Images                [][]byte   `json:"images"`
	Status                taskStatus `json:"status"`
	Title                 string     `json:"title"`
	BotsFilename          *string    `json:"bots_filename"`
	CheapProxiesFilename  *string    `json:"cheap_proxies_filename"`
	ResProxiesFilename    *string    `json:"res_proxies_filename"`
	TargetsFilename       *string    `json:"targets_filename"`
	CreatedAt             time.Time  `json:"created_at"`
	StartedAt             *time.Time `json:"started_at"`
	StoppedAt             *time.Time `json:"stopped_at"`
	UpdatedAt             *time.Time `json:"updated_at"`
	DeletedAt             *time.Time `json:"deleted_at"`
	AccountLastNames      []string   `json:"account_last_names"`
	FollowTargets         bool       `json:"follow_targets"`
	NeedPhotoTags         bool       `json:"need_photo_tags"`
	PerPostSleepSeconds   int32      `json:"per_post_sleep_seconds"`
	PhotoTagsDelaySeconds int32      `json:"photo_tags_delay_seconds"`
	PostsPerBot           int32      `json:"posts_per_bot"`
	TargetsPerPost        int32      `json:"targets_per_post"`
	Type                  taskType   `json:"type"`
	VideoFilename         *string    `json:"video_filename"`
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
