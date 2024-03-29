// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package dbmodel

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/headers"
)

type MediasKind string

const (
	MediasKindPhoto MediasKind = "photo"
	MediasKindReels MediasKind = "reels"
)

func (e *MediasKind) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = MediasKind(s)
	case string:
		*e = MediasKind(s)
	default:
		return fmt.Errorf("unsupported scan type for MediasKind: %T", src)
	}
	return nil
}

type NullMediasKind struct {
	MediasKind MediasKind
	Valid      bool // Valid is true if String is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullMediasKind) Scan(value interface{}) error {
	if value == nil {
		ns.MediasKind, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.MediasKind.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullMediasKind) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.MediasKind, nil
}

type PgqueueStatus string

const (
	PgqueueStatusNew            PgqueueStatus = "new"
	PgqueueStatusMustRetry      PgqueueStatus = "must_retry"
	PgqueueStatusNoAttemptsLeft PgqueueStatus = "no_attempts_left"
	PgqueueStatusCancelled      PgqueueStatus = "cancelled"
	PgqueueStatusSucceeded      PgqueueStatus = "succeeded"
)

func (e *PgqueueStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PgqueueStatus(s)
	case string:
		*e = PgqueueStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for PgqueueStatus: %T", src)
	}
	return nil
}

type NullPgqueueStatus struct {
	PgqueueStatus PgqueueStatus
	Valid         bool // Valid is true if String is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPgqueueStatus) Scan(value interface{}) error {
	if value == nil {
		ns.PgqueueStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PgqueueStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPgqueueStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.PgqueueStatus, nil
}

type TargetsInteraction string

const (
	TargetsInteractionNone            TargetsInteraction = "none"
	TargetsInteractionPostDescription TargetsInteraction = "post_description"
	TargetsInteractionPhotoTag        TargetsInteraction = "photo_tag"
)

func (e *TargetsInteraction) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TargetsInteraction(s)
	case string:
		*e = TargetsInteraction(s)
	default:
		return fmt.Errorf("unsupported scan type for TargetsInteraction: %T", src)
	}
	return nil
}

type NullTargetsInteraction struct {
	TargetsInteraction TargetsInteraction
	Valid              bool // Valid is true if String is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTargetsInteraction) Scan(value interface{}) error {
	if value == nil {
		ns.TargetsInteraction, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TargetsInteraction.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTargetsInteraction) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.TargetsInteraction, nil
}

type TargetsStatus string

const (
	TargetsStatusNew        TargetsStatus = "new"
	TargetsStatusInProgress TargetsStatus = "in_progress"
	TargetsStatusFailed     TargetsStatus = "failed"
	TargetsStatusNotified   TargetsStatus = "notified"
)

func (e *TargetsStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TargetsStatus(s)
	case string:
		*e = TargetsStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for TargetsStatus: %T", src)
	}
	return nil
}

type NullTargetsStatus struct {
	TargetsStatus TargetsStatus
	Valid         bool // Valid is true if String is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTargetsStatus) Scan(value interface{}) error {
	if value == nil {
		ns.TargetsStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TargetsStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTargetsStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.TargetsStatus, nil
}

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
	PostsCount int                    `json:"posts_count"`
	StartedAt  *time.Time             `json:"started_at"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  *time.Time             `json:"updated_at"`
	DeletedAt  *time.Time             `json:"deleted_at"`
	FileOrder  int32                  `json:"file_order"`
	InstID     int64                  `json:"inst_id"`
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

type Media struct {
	Kind      MediasKind `json:"kind"`
	InstID    string     `json:"inst_id"`
	BotID     uuid.UUID  `json:"bot_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Pk        int64      `json:"pk"`
	IsEdited  bool       `json:"is_edited"`
}

type Pgqueue struct {
	ID              int64         `json:"id"`
	Kind            int           `json:"kind"`
	Payload         []byte        `json:"payload"`
	ExternalKey     *string       `json:"external_key"`
	Status          PgqueueStatus `json:"status"`
	Messages        []string      `json:"messages"`
	AttemptsLeft    int           `json:"attempts_left"`
	AttemptsElapsed int           `json:"attempts_elapsed"`
	DelayedTill     time.Time     `json:"delayed_till"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
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
	ID              uuid.UUID          `json:"id"`
	TaskID          uuid.UUID          `json:"task_id"`
	Username        string             `json:"username"`
	UserID          int64              `json:"user_id"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       *time.Time         `json:"updated_at"`
	Status          TargetsStatus      `json:"status"`
	InteractionType TargetsInteraction `json:"interaction_type"`
	MediaFk         *int64             `json:"media_fk"`
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
	Type                  taskType   `json:"type"`
	VideoFilename         *string    `json:"video_filename"`
	PostsPerBot           int        `json:"posts_per_bot"`
	TargetsPerPost        int        `json:"targets_per_post"`
	PhotoTagsPostsPerBot  int        `json:"photo_tags_posts_per_bot"`
	PhotoTargetsPerPost   int        `json:"photo_targets_per_post"`
	FixedTag              *string    `json:"fixed_tag"`
	FixedPhotoTag         *int64     `json:"fixed_photo_tag"`
}

type User struct {
	ID           uuid.UUID  `json:"id"`
	Login        string     `json:"login"`
	PasswordHash string     `json:"password_hash"`
	Role         int        `json:"role"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}
