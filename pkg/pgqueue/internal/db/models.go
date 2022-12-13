package db

type PgqueueStatus string

const (
	PgqueueStatusNew            PgqueueStatus = "new"
	PgqueueStatusMustRetry      PgqueueStatus = "must_retry"
	PgqueueStatusNoAttemptsLeft PgqueueStatus = "no_attempts_left"
	PgqueueStatusCancelled      PgqueueStatus = "cancelled"
	PgqueueStatusSucceeded      PgqueueStatus = "succeeded"
)
